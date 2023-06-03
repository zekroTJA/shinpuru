package modnot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/pkg/discordutil"
)

// Send embed messages into the mod notification channel
// specified for the given guildID.
func Send(
	db database.Database,
	s discordutil.ISession,
	guildID string,
	embed *discordgo.MessageEmbed,
) error {
	chanID, err := db.GetGuildModNot(guildID)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return err
	}
	if chanID == "" {
		return nil
	}

	_, err = s.ChannelMessageSendEmbed(chanID, embed)

	return err
}
