package report

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/models"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/util/snowflakenodes"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/stringutil"
)

// PushReport creates a new Report object with the given executorID,
// victimID, reason, attachmentID, and typ. The report is saved to the database
// using the passed db databse rpovider and an embed is created with the attachment
// url assembled with publicAddr as image endpoint root. This embed is then sent to
// the specified mod log channel for this guild, if existent.
func PushReport(s *discordgo.Session, db database.Database, publicAddr,
	guildID, executorID, victimID, reason, attachmentID string, typ models.Type) (*models.Report, error) {

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

	err := db.AddReport(rep)
	if err != nil {
		return nil, err
	}

	if modlogChan, err := db.GetGuildModLog(guildID); err == nil {
		s.ChannelMessageSendEmbed(modlogChan, rep.AsEmbed(publicAddr))
	}

	dmChan, err := s.UserChannelCreate(victimID)
	if err == nil {
		s.ChannelMessageSendEmbed(dmChan.ID, rep.AsEmbed(publicAddr))
	}

	return rep, nil
}

// PushKick is shorthand for PushReport as member kick action and also
// kicks the member from the guild with the given reason and case ID
// for the audit log.
func PushKick(s *discordgo.Session, db database.Database, publicAddr, guildID,
	executorID, victimID, reason, attachment string) (*models.Report, error) {

	const typ = 0

	rep, err := PushReport(s, db, publicAddr, guildID, executorID, victimID, reason, attachment, typ)
	if err != nil {
		return nil, err
	}

	if err = s.GuildMemberDeleteWithReason(guildID, victimID, fmt.Sprintf(`[CASE %s] %s`, rep.ID, reason)); err != nil {
		db.DeleteReport(rep.ID)
		return nil, err
	}

	return rep, nil
}

// PushBan is shorthand for PushReport as member ban action and also
// bans the member from the guild with the given reason and case ID
// for the audit log.
func PushBan(s *discordgo.Session, db database.Database, publicAddr, guildID,
	executorID, victimID, reason, attachment string) (*models.Report, error) {

	const typ = 1

	rep, err := PushReport(s, db, publicAddr, guildID, executorID, victimID, reason, attachment, typ)
	if err != nil {
		return nil, err
	}

	if err = s.GuildBanCreateWithReason(guildID, victimID, fmt.Sprintf(`[CASE %s] %s`, rep.ID, reason), 7); err != nil {
		db.DeleteReport(rep.ID)
		return nil, err
	}

	return rep, nil
}

// PushMute is shorthand for PushReport as member mute action and also
// adds the mute role to the specified victim.
func PushMute(s *discordgo.Session, db database.Database, publicAddr, guildID,
	executorID, victimID, reason, attachment, muteRoleID string) (*models.Report, error) {

	const typ = 2

	if reason == "" {
		reason = "no reason specified"
	}

	rep, err := PushReport(s, db, publicAddr, guildID, executorID, victimID, reason, attachment, typ)
	if err != nil {
		return nil, err
	}

	err = s.GuildMemberRoleAdd(guildID, victimID, muteRoleID)
	if err != nil {
		db.DeleteReport(rep.ID)
		return nil, err
	}

	return rep, nil
}

// RevokeMute removes the mute role of the specified victim and sends
// an unmute embed to the users DMs and to the mod log channel.
func RevokeMute(s *discordgo.Session, db database.Database, publicAddr, guildID,
	executorID, victimID, reason, muteRoleID string) (emb *discordgo.MessageEmbed, err error) {

	err = s.GuildMemberRoleRemove(guildID, victimID, muteRoleID)
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

	if modlogChan, err := db.GetGuildModLog(guildID); err == nil {
		s.ChannelMessageSendEmbed(modlogChan, emb)
	}

	dmChan, err := s.UserChannelCreate(victimID)
	if err == nil {
		s.ChannelMessageSendEmbed(dmChan.ID, emb)
	}

	return
}

func RevokeReport(rep *models.Report, executorID, reason,
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
