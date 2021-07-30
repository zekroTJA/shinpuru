package report

import (
	"errors"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/sarulabs/di/v2"
	"github.com/sirupsen/logrus"
	"github.com/zekroTJA/shinpuru/internal/config"
	"github.com/zekroTJA/shinpuru/internal/models"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/util/snowflakenodes"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/multierror"
	"github.com/zekroTJA/shinpuru/pkg/roleutil"
	"github.com/zekroTJA/shinpuru/pkg/stringutil"
	"github.com/zekrotja/dgrs"
)

var (
	ErrRoleDiff = errors.New("You can only ban or kick members with lower permissions than yours.")
)

type ReportService struct {
	s   *discordgo.Session
	db  database.Database
	cfg *config.Config
	st  *dgrs.State
}

type ReportError struct {
	error
	*models.Report
}

func New(container di.Container) *ReportService {
	return &ReportService{
		s:   container.Get(static.DiDiscordSession).(*discordgo.Session),
		db:  container.Get(static.DiDatabase).(database.Database),
		cfg: container.Get(static.DiConfig).(*config.Config),
		st:  container.Get(static.DiState).(*dgrs.State),
	}
}

// PushReport creates a new Report object with the given executorID,
// victimID, reason, attachmentID, and typ. The report is saved to the database
// using the passed db databse rpovider and an embed is created with the attachment
// url assembled with publicAddr as image endpoint root. This embed is then sent to
// the specified mod log channel for this guild, if existent.
func (r *ReportService) PushReport(rep *models.Report) (*models.Report, error) {
	repID := snowflakenodes.NodesReport[rep.Type].Generate()

	rep.ID = repID

	err := r.db.AddReport(rep)
	if err != nil {
		return nil, err
	}

	if modlogChan, err := r.db.GetGuildModLog(rep.GuildID); err == nil {
		r.s.ChannelMessageSendEmbed(modlogChan, rep.AsEmbed(r.cfg.WebServer.PublicAddr))
	}

	dmChan, err := r.s.UserChannelCreate(rep.VictimID)
	if err == nil {
		r.s.ChannelMessageSendEmbed(dmChan.ID, rep.AsEmbed(r.cfg.WebServer.PublicAddr))
	}

	return rep, nil
}

// PushKick is shorthand for PushReport as member kick action and also
// kicks the member from the guild with the given reason and case ID
// for the audit log.
func (r *ReportService) PushKick(rep *models.Report) (*models.Report, error) {
	const typ = 0
	rep.ID = typ

	guild, err := r.st.Guild(rep.GuildID, true)
	if err != nil {
		return nil, err
	}

	victim, err := r.st.Member(rep.GuildID, rep.VictimID)
	if err != nil {
		return nil, err
	}

	executor, err := r.st.Member(rep.GuildID, rep.ExecutorID)
	if err != nil {
		return nil, err
	}

	if roleutil.PositionDiff(victim, executor, guild) >= 0 {
		return nil, ErrRoleDiff
	}

	rep, err = r.PushReport(rep)
	if err != nil {
		return nil, err
	}

	if err = r.s.GuildMemberDeleteWithReason(rep.GuildID, rep.VictimID, fmt.Sprintf(`[CASE %s] %s`, rep.ID, rep.Msg)); err != nil {
		r.db.DeleteReport(rep.ID)
		return nil, err
	}

	return rep, nil
}

// PushBan is shorthand for PushReport as member ban action and also
// bans the member from the guild with the given reason and case ID
// for the audit log.
func (r *ReportService) PushBan(rep *models.Report) (*models.Report, error) {
	const typ = 1
	rep.Type = typ

	if err := checkTimeout(rep.Timeout); err != nil {
		return nil, err
	}

	guild, err := r.st.Guild(rep.GuildID, true)
	if err != nil {
		return nil, err
	}

	var victim *discordgo.Member
	if !rep.Anonymous {
		victim, err = r.st.Member(rep.GuildID, rep.VictimID)
		if err != nil {
			return nil, err
		}
	}

	executor, err := r.st.Member(rep.GuildID, rep.ExecutorID)
	if err != nil {
		return nil, err
	}

	if !rep.Anonymous {
		if roleutil.PositionDiff(victim, executor, guild) >= 0 {
			return nil, ErrRoleDiff
		}
	}

	rep, err = r.PushReport(rep)
	if err != nil {
		return nil, err
	}

	if err = r.s.GuildBanCreateWithReason(rep.GuildID, rep.VictimID, fmt.Sprintf(`[CASE %s] %s`, rep.ID, rep.Msg), 7); err != nil {
		r.db.DeleteReport(rep.ID)
		return nil, err
	}

	return rep, nil
}

// PushMute is shorthand for PushReport as member mute action and also
// adds the mute role to the specified victim.
func (r *ReportService) PushMute(rep *models.Report, muteRoleID string) (*models.Report, error) {
	const typ = 2
	rep.Type = typ

	if err := checkTimeout(rep.Timeout); err != nil {
		return nil, err
	}

	guild, err := r.st.Guild(rep.GuildID, true)
	if err != nil {
		return nil, err
	}

	victim, err := r.st.Member(rep.GuildID, rep.VictimID)
	if err != nil {
		return nil, err
	}

	executor, err := r.st.Member(rep.GuildID, rep.ExecutorID)
	if err != nil {
		return nil, err
	}

	if roleutil.PositionDiff(victim, executor, guild) >= 0 {
		return nil, ErrRoleDiff
	}

	if rep.Msg == "" {
		rep.Msg = "no reason specified"
	}

	rep, err = r.PushReport(rep)
	if err != nil {
		return nil, err
	}

	err = r.s.GuildMemberRoleAdd(rep.GuildID, rep.VictimID, muteRoleID)
	if err != nil {
		r.db.DeleteReport(rep.ID)
		return nil, err
	}

	return rep, nil
}

// RevokeMute removes the mute role of the specified victim and sends
// an unmute embed to the users DMs and to the mod log channel.
func (r *ReportService) RevokeMute(guildID, executorID, victimID, reason, muteRoleID string) (emb *discordgo.MessageEmbed, err error) {
	guild, err := r.st.Guild(guildID, true)
	if err != nil {
		return nil, err
	}

	victim, err := r.st.Member(guildID, victimID)
	if err != nil {
		return nil, err
	}

	executor, err := r.st.Member(guildID, executorID)
	if err != nil {
		return nil, err
	}

	if roleutil.PositionDiff(victim, executor, guild) >= 0 {
		return nil, ErrRoleDiff
	}

	err = r.s.GuildMemberRoleRemove(guildID, victimID, muteRoleID)
	if err != nil {
		return
	}

	if err = r.ExpireLastReport(guildID, victimID, int(models.TypeMute)); err != nil {
		return
	}

	repType := stringutil.IndexOf("MUTE", models.ReportTypes)
	repID := snowflakenodes.NodesReport[repType].Generate()

	emb = &discordgo.MessageEmbed{
		Title: "Case " + repID.String(),
		Color: models.ReportColors[repType],
		Fields: []*discordgo.MessageEmbedField{
			{
				Inline: true,
				Name:   "Executor",
				Value:  fmt.Sprintf("<@%s>", executorID),
			},
			{
				Inline: true,
				Name:   "Victim",
				Value:  fmt.Sprintf("<@%s>", victimID),
			},
			{
				Name:  "Type",
				Value: "UNMUTE",
			},
			{
				Name:  "Description",
				Value: "MANUAL UNMUTE",
			},
		},
		Timestamp: time.Unix(repID.Time()/1000, 0).Format("2006-01-02T15:04:05.000Z"),
	}

	if modlogChan, err := r.db.GetGuildModLog(guildID); err == nil {
		r.s.ChannelMessageSendEmbed(modlogChan, emb)
	}

	dmChan, err := r.s.UserChannelCreate(victimID)
	if err == nil {
		r.s.ChannelMessageSendEmbed(dmChan.ID, emb)
	}

	return
}

func (r *ReportService) RevokeReport(rep *models.Report, executorID, reason,
	wsPublicAddr string, db database.Database,
	s *discordgo.Session) (*discordgo.MessageEmbed, error) {

	err := db.DeleteReport(rep.ID)
	if err != nil {
		return nil, err
	}

	repRevEmb := &discordgo.MessageEmbed{
		Color:       static.ReportRevokedColor,
		Title:       "REPORT REVOCATION",
		Description: "Revoked reports are deleted from the database and no more visible in any commands.",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Revoke Executor",
				Value: fmt.Sprintf("<@%s>", executorID),
			},
			{
				Name:  "Revocation Reason",
				Value: reason,
			},
			rep.AsEmbedField(wsPublicAddr),
		},
	}

	if modlogChan, err := db.GetGuildModLog(rep.GuildID); err == nil {
		s.ChannelMessageSendEmbed(modlogChan, repRevEmb)
	}
	dmChan, err := s.UserChannelCreate(rep.VictimID)
	if err == nil {
		s.ChannelMessageSendEmbed(dmChan.ID, repRevEmb)
	}

	return repRevEmb, nil
}

func (r *ReportService) ExpireLastReport(guildID, victimID string, typ int) (err error) {
	reps, err := r.db.GetReportsFiltered(guildID, victimID, typ, 0, 1)
	if err != nil {
		return
	}
	if len(reps) > 0 && reps[0].Timeout != nil {
		err = r.db.ExpireReports(reps[0].ID.String())
	}
	return
}

func (r *ReportService) ExpireExpiredReports() (mErr *multierror.MultiError) {
	mErr = multierror.New()

	reps, err := r.db.GetExpiredReports()
	mErr.Append(err)
	logrus.WithField("n", len(reps)).Debug("REPORTS :: start expiring cleanup ...")

	expIDs := make([]string, 0, len(reps))
	for _, rep := range reps {
		logrus.WithFields(logrus.Fields{
			"id":  rep.ID,
			"typ": rep.Type,
		}).Debug("REPORTS :: expiring report")
		err = r.revokeReportOnExpiration(rep)
		if err == nil {
			expIDs = append(expIDs, rep.ID.String())
		} else {
			mErr.Append(&ReportError{
				error:  err,
				Report: rep,
			})
		}
	}

	mErr.Append(
		r.db.ExpireReports(expIDs...))

	return
}

func (r *ReportService) revokeReportOnExpiration(rep *models.Report) (err error) {
	switch rep.Type {
	case models.TypeBan:
		err = r.s.GuildBanDelete(rep.GuildID, rep.VictimID)
	case models.TypeMute:
		var rid string
		rid, err = r.db.GetGuildMuteRole(rep.GuildID)
		if err != nil {
			return
		}
		_, err = r.RevokeMute(rep.GuildID, rep.ExecutorID, rep.VictimID, "Automatic Timeout", rid)
	}
	return
}

func checkTimeout(t *time.Time) (err error) {
	if t != nil && t.Before(time.Now()) {
		err = errors.New("timeout must be in the future")
	}
	return
}
