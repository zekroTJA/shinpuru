package util

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/util/static"
)

type EmbedMessage struct {
	*discordgo.Message

	s   *discordgo.Session
	err error
}

func (emb *EmbedMessage) DeleteAfter(d time.Duration) *EmbedMessage {
	if emb.Message != nil {
		time.AfterFunc(d, func() {
			emb.err = emb.s.ChannelMessageDelete(emb.ChannelID, emb.ID)
		})
	}
	return emb
}

func (emb *EmbedMessage) Error() error {
	return emb.err
}

func SendEmbedError(s *discordgo.Session, chanID, content string, title ...string) *EmbedMessage {
	emb := &discordgo.MessageEmbed{
		Description: content,
		Color:       static.ColorEmbedError,
		Title:       "Error",
	}

	if len(title) > 0 {
		emb.Title = title[0]
	}

	return sendEmbedRaw(s, chanID, emb)
}

func SendEmbed(s *discordgo.Session, chanID, content string, title string, color int) *EmbedMessage {
	emb := &discordgo.MessageEmbed{
		Description: content,
		Color:       color,
	}

	emb.Title = title
	if color == 0 {
		emb.Color = static.ColorEmbedError
	}

	return sendEmbedRaw(s, chanID, emb)
}

func sendEmbedRaw(s *discordgo.Session, chanID string, emb *discordgo.MessageEmbed) *EmbedMessage {
	msg, err := s.ChannelMessageSendEmbed(chanID, emb)

	return &EmbedMessage{msg, s, err}
}
