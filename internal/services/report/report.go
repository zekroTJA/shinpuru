package report

import (
	"errors"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/config"
	"github.com/zekroTJA/shinpuru/internal/models"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/util/snowflakenodes"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/discordutil"
	"github.com/zekroTJA/shinpuru/pkg/roleutil"
	"github.com/zekroTJA/shinpuru/pkg/stringutil"
)

var (
	errRoleDiff = errors.New("You can only ban or kick members with lower permissions than yours.")
)

type ReportService struct {
	s   *discordgo.Session
	db  database.Database
	cfg *config.Config
}

func New(container di.Container) *ReportService {
	return &ReportService{
		s:   container.Get(static.DiDiscordSession).(*discordgo.Session),
		db:  container.Get(static.DiDatabase).(database.Database),
		cfg: container.Get(static.DiConfig).(*config.Config),
	}
}

// PushReport creates a new Report object with the given executorID,
// victimID, reason, attachmentID, and typ. The report is saved to the database
// using the passed db databse rpovider and an embed is created with the attachment
// url assembled with publicAddr as image endpoint root. This embed is then sent to
// the specified mod log channel for this guild, if existent.
func (r *ReportService) PushReport(guildID, executorID, victimID, reason, attachmentID string, typ models.Type) (*models.Report, error) {
	repID := snowflakenodes.NodesReport[typ].Generate()

	rep := &models.Report{
		ID:            repID,
		Type:          typ,
		GuildID:       guildID,
		ExecutorID:    executorID,
		VictimID:      victimID,
		Msg:           reason,
		AttachmehtURL: attachmentID,
	}

	err := r.db.AddReport(rep)
	if err != nil {
		return nil, err
	}

	if modlogChan, err := r.db.GetGuildModLog(guildID); err == nil {
		r.s.ChannelMessageSendEmbed(modlogChan, rep.AsEmbed(r.cfg.WebServer.PublicAddr))
	}

	dmChan, err := r.s.UserChannelCreate(victimID)
	if err == nil {
		r.s.ChannelMessageSendEmbed(dmChan.ID, rep.AsEmbed(r.cfg.WebServer.PublicAddr))
	}

	return rep, nil
}

// PushKick is shorthand for PushReport as member kick action and also
// kicks the member from the guild with the given reason and case ID
// for the audit log.
func (r *ReportService) PushKick(guildID, executorID, victimID, reason, attachment string) (*models.Report, error) {
	const typ = 0

	guild, err := discordutil.GetGuild(r.s, guildID)
	if err != nil {
		return nil, err
	}

	victim, err := discordutil.GetMember(r.s, guildID, victimID)
	if err != nil {
		return nil, err
	}

	executor, err := discordutil.GetMember(r.s, guildID, executorID)
	if err != nil {
		return nil, err
	}

	if roleutil.PositionDiff(victim, executor, guild) >= 0 {
		return nil, errRoleDiff
	}

	rep, err := r.PushReport(guildID, executorID, victimID, reason, attachment, typ)
	if err != nil {
		return nil, err
	}

	if err = r.s.GuildMemberDeleteWithReason(guildID, victimID, fmt.Sprintf(`[CASE %s] %s`, rep.ID, reason)); err != nil {
		r.db.DeleteReport(rep.ID)
		return nil, err
	}

	return rep, nil
}

// PushBan is shorthand for PushReport as member ban action and also
// bans the member from the guild with the given reason and case ID
// for the audit log.
func (r *ReportService) PushBan(guildID, executorID, victimID, reason, attachment string) (*models.Report, error) {
	const typ = 1

	guild, err := discordutil.GetGuild(r.s, guildID)
	if err != nil {
		return nil, err
	}

	victim, err := discordutil.GetMember(r.s, guildID, victimID)
	if err != nil {
		return nil, err
	}

	executor, err := discordutil.GetMember(r.s, guildID, executorID)
	if err != nil {
		return nil, err
	}

	if roleutil.PositionDiff(victim, executor, guild) >= 0 {
		return nil, errRoleDiff
	}

	rep, err := r.PushReport(guildID, executorID, victimID, reason, attachment, typ)
	if err != nil {
		return nil, err
	}

	if err = r.s.GuildBanCreateWithReason(guildID, victimID, fmt.Sprintf(`[CASE %s] %s`, rep.ID, reason), 7); err != nil {
		r.db.DeleteReport(rep.ID)
		return nil, err
	}

	return rep, nil
}

// PushMute is shorthand for PushReport as member mute action and also
// adds the mute role to the specified victim.
func (r *ReportService) PushMute(guildID, executorID, victimID, reason, attachment, muteRoleID string) (*models.Report, error) {
	const typ = 2

	if reason == "" {
		reason = "no reason specified"
	}

	rep, err := r.PushReport(guildID, executorID, victimID, reason, attachment, typ)
	if err != nil {
		return nil, err
	}

	err = r.s.GuildMemberRoleAdd(guildID, victimID, muteRoleID)
	if err != nil {
		r.db.DeleteReport(rep.ID)
		return nil, err
	}

	return rep, nil
}

// RevokeMute removes the mute role of the specified victim and sends
// an unmute embed to the users DMs and to the mod log channel.
func (r *ReportService) RevokeMute(guildID, executorID, victimID, reason, muteRoleID string) (emb *discordgo.MessageEmbed, err error) {
	err = r.s.GuildMemberRoleRemove(guildID, victimID, muteRoleID)
	if err != nil {
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
