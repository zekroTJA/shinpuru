package util

import (
	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/util/static"
)

func SendEmbedError(s *discordgo.Session, chanID, content string, title ...string) (*discordgo.Message, error) {
	emb := &discordgo.MessageEmbed{
		Description: content,
		Color:       static.ColorEmbedError,
	}
	if len(title) > 0 {
		emb.Title = title[0]
	}
	return s.ChannelMessageSendEmbed(chanID, emb)
}

func SendEmbed(s *discordgo.Session, chanID, content string, title string, color int) (*discordgo.Message, error) {
	emb := &discordgo.MessageEmbed{
		Description: content,
		Color:       color,
	}
	emb.Title = title
	if color == 0 {
		emb.Color = static.ColorEmbedError
	}
	return s.ChannelMessageSendEmbed(chanID, emb)
}
