package modnot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/pkg/discordutil"
)

func Send(
	db database.Database,
	s discordutil.ISession,
	guildID string,
	embed *discordgo.MessageEmbed,
) error {
	chanID, err := db.GetGuildModNot(guildID)
	if err != nil {
		return err
	}

	_, err = s.ChannelMessageSendEmbed(chanID, embed)

	return err
}
