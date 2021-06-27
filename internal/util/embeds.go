package util

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/util/static"
)

// EmbedMessage extends a discordgo.MessageEmbedMessage
// with extra utilities.
type EmbedMessage struct {
	*discordgo.Message

	s   *discordgo.Session
	err error
}

// DeleteAfter deletes the message after the specified
// duration when the message still exists and sets occured
// error to the EmbedMessage.
func (emb *EmbedMessage) DeleteAfter(d time.Duration) *EmbedMessage {
	if emb.Message != nil {
		time.AfterFunc(d, func() {
			emb.err = emb.s.ChannelMessageDelete(emb.ChannelID, emb.ID)
		})
	}
	return emb
}

// Error returns the embedded error.
func (emb *EmbedMessage) Error() error {
	return emb.err
}

// Edit updates the current embed message with
// the given content replacing the internal message
// and error of this embed instance.
func (emb *EmbedMessage) Edit(content string, title string, color int) *EmbedMessage {
	newEmb := &discordgo.MessageEmbed{
		Description: content,
		Color:       color,
	}

	newEmb.Title = title
	if color == 0 {
		newEmb.Color = static.ColorEmbedDefault
	}

	return emb.EditRaw(newEmb)
}

// EditRaw updates the current embed message with
// the given raw embed replacing the internal message
// and error of this embed instance.
func (emb *EmbedMessage) EditRaw(newEmb *discordgo.MessageEmbed) *EmbedMessage {
	emb.Message, emb.err = emb.s.ChannelMessageEditEmbed(emb.ChannelID, emb.ID, newEmb)
	return emb
}

// SendEmbed creates an discordgo.MessageEmbed from the passed
// content, title and color and sends it to the specified channel.
//
// If color == 0, static.ColorEmbedDefault will be set as color.
//
// Occured errors are set to the internal error.
func SendEmbed(s *discordgo.Session, chanID, content string, title string, color int) *EmbedMessage {
	emb := &discordgo.MessageEmbed{
		Description: content,
		Color:       color,
	}

	emb.Title = title
	if color == 0 {
		emb.Color = static.ColorEmbedDefault
	}

	return SendEmbedRaw(s, chanID, emb)
}

// SendEmbedError is shorthand for SendEmbed with
// static.ColorEmbedError as color and title "Error"
// if no title was passed.
func SendEmbedError(s *discordgo.Session, chanID, content string, title ...string) *EmbedMessage {
	emb := &discordgo.MessageEmbed{
		Description: content,
		Color:       static.ColorEmbedError,
		Title:       "Error",
	}

	if len(title) > 0 {
		emb.Title = title[0]
	}

	return SendEmbedRaw(s, chanID, emb)
}

// SendEmbedRaw sends the passed emb to the passed
// channel and sets occured errors to the internal
// error.
func SendEmbedRaw(s *discordgo.Session, chanID string, emb *discordgo.MessageEmbed) *EmbedMessage {
	msg, err := s.ChannelMessageSendEmbed(chanID, emb)

	return &EmbedMessage{msg, s, err}
}
