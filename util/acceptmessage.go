package util

import (
	"github.com/bwmarrin/discordgo"
)

const (
	acceptMessageEmoteAccept  = "✅"
	acceptMessageEmoteDecline = "❌"
)

type AcceptMessage struct {
	*discordgo.Message
	Session     *discordgo.Session
	Embed       *discordgo.MessageEmbed
	AcceptFunc  func(*discordgo.Message) error
	DeclineFUnc func(*discordgo.Message) error
	embedUnsub  func()
}

func (am *AcceptMessage) Send(chanID string) (*AcceptMessage, error) {
	msg, err := am.Session.ChannelMessageSendEmbed(chanID, am.Embed)
	if err != nil {
		return nil, err
	}
	am.Message = msg
	am.embedUnsub = am.Session.AddHandler(func(s *discordgo.Session, e *discordgo.MessageReactionAdd) {

	})
	return am, nil
}
