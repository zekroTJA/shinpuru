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
	Session        *discordgo.Session
	Embed          *discordgo.MessageEmbed
	UserID         string
	DeleteMsgAfter bool
	AcceptFunc     func(*discordgo.Message)
	DeclineFunc    func(*discordgo.Message)
	eventUnsub     func()
}

func (am *AcceptMessage) Send(chanID string) (*AcceptMessage, error) {
	msg, err := am.Session.ChannelMessageSendEmbed(chanID, am.Embed)
	if err != nil {
		return nil, err
	}
	am.Message = msg
	err = am.Session.MessageReactionAdd(chanID, msg.ID, acceptMessageEmoteAccept)
	err = am.Session.MessageReactionAdd(chanID, msg.ID, acceptMessageEmoteDecline)
	if err != nil {
		return nil, err
	}
	am.eventUnsub = am.Session.AddHandler(func(s *discordgo.Session, e *discordgo.MessageReactionAdd) {
		if e.MessageID != msg.ID || e.UserID == s.State.User.ID || (am.UserID != "" && am.UserID != e.UserID) {
			return
		}
		if e.Emoji.Name != acceptMessageEmoteAccept && e.Emoji.Name != acceptMessageEmoteDecline {
			return
		}
		switch e.Emoji.Name {
		case acceptMessageEmoteAccept:
			if am.AcceptFunc != nil {
				am.AcceptFunc(msg)
			}
		case acceptMessageEmoteDecline:
			if am.DeclineFunc != nil {
				am.DeclineFunc(msg)
			}
		}
		am.eventUnsub()
		if am.DeleteMsgAfter {
			am.Session.ChannelMessageDelete(chanID, msg.ID)
		} else {
			am.Session.MessageReactionsRemoveAll(chanID, msg.ID)
		}
	})
	return am, nil
}
