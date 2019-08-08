package shared

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/core"
	"github.com/zekroTJA/shinpuru/internal/util"
)

func PushReport(s *discordgo.Session, db core.Database, guildID, executorID, victimID, reason, attachment string, typ int) (*util.Report, error) {
	repID := util.NodesReport[typ].Generate()

	rep := &util.Report{
		ID:            repID,
		Type:          typ,
		GuildID:       guildID,
		ExecutorID:    executorID,
		VictimID:      victimID,
		Msg:           reason,
		AttachmehtURL: attachment,
	}

	err := db.AddReport(rep)
	if err != nil {
		return nil, err
	}

	if modlogChan, err := db.GetGuildModLog(guildID); err == nil {
		s.ChannelMessageSendEmbed(modlogChan, rep.AsEmbed())
	}

	dmChan, err := s.UserChannelCreate(victimID)
	if err == nil {
		s.ChannelMessageSendEmbed(dmChan.ID, rep.AsEmbed())
	}

	return rep, nil
}

func PushKick(s *discordgo.Session, db core.Database, guildID, executorID, victimID, reason, attachment string) (*util.Report, error) {
	const typ = 0

	rep, err := PushReport(s, db, guildID, executorID, victimID, reason, attachment, typ)
	if err != nil {
		return nil, err
	}

	if err = s.GuildMemberDeleteWithReason(guildID, victimID, fmt.Sprintf(`[CASE %s] %s`, rep.ID, reason)); err != nil {
		db.DeleteReport(rep.ID)
		return nil, err
	}

	return rep, nil
}

func PushBan(s *discordgo.Session, db core.Database, guildID, executorID, victimID, reason, attachment string) (*util.Report, error) {
	const typ = 1

	rep, err := PushReport(s, db, guildID, executorID, victimID, reason, attachment, typ)
	if err != nil {
		return nil, err
	}

	if err = s.GuildBanCreateWithReason(guildID, victimID, fmt.Sprintf(`[CASE %s] %s`, rep.ID, reason), 7); err != nil {
		db.DeleteReport(rep.ID)
		return nil, err
	}

	return rep, nil
}

func PushMute(s *discordgo.Session, db core.Database, guildID, executorID, victimID, reason, attachment, muteRoleID string) (*util.Report, error) {
	const typ = 2

	if reason == "" {
		reason = "no reason specified"
	}

	rep, err := PushReport(s, db, guildID, executorID, victimID, reason, attachment, typ)
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
