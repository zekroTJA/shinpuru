// Package acceptmsg provides a message model for
// discordgo which can be accepted or declined
// via message reactions.
package acceptmsg

import (
	"errors"

	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekrotja/discordgo"
)

const (
	acceptMessageEmoteAccept  = "✅"
	acceptMessageEmoteDecline = "❌"
)

type ActionHandler func(*discordgo.Message)

// AcceptMessage extends discordgo.Message to build
// and send an AcceptMessage.
type AcceptMessage struct {
	*discordgo.Message
	Session        *discordgo.Session
	Embed          *discordgo.MessageEmbed
	UserID         string
	DeleteMsgAfter bool
	AcceptFunc     ActionHandler
	DeclineFunc    ActionHandler

	eventUnsub func()
}

// New creates an empty instance of AcceptMessage.
func New() *AcceptMessage {
	return new(AcceptMessage)
}

// WithSession sets the discordgo.Session instance.
func (am *AcceptMessage) WithSession(s *discordgo.Session) *AcceptMessage {
	am.Session = s
	return am
}

// WithEmbed sets the Embed instance to be set.
func (am *AcceptMessage) WithEmbed(e *discordgo.MessageEmbed) *AcceptMessage {
	am.Embed = e
	return am
}

// WithContent creates an embed with default color and
// specified content as description and sets it as
// embed instance.
func (am *AcceptMessage) WithContent(cont string) *AcceptMessage {
	am.Embed = &discordgo.MessageEmbed{
		Color:       static.ColorEmbedDefault,
		Description: cont,
	}
	return am
}

// LockOnUser specifies, that only reaction inputs from
// the defined user are accepted.
func (am *AcceptMessage) LockOnUser(userID string) *AcceptMessage {
	am.UserID = userID
	return am
}

// DeleteAfterAnser enables that the whole accept
// embed message is being deleted after users
// answer.
func (am *AcceptMessage) DeleteAfterAnswer() *AcceptMessage {
	am.DeleteMsgAfter = true
	return am
}

// DoOnAccept specifies the action handler executed
// on acception.
func (am *AcceptMessage) DoOnAccept(onAccept ActionHandler) *AcceptMessage {
	am.AcceptFunc = onAccept
	return am
}

// DoOnDecline specifies the action handler executed
// on decline.
func (am *AcceptMessage) DoOnDecline(onDecline ActionHandler) *AcceptMessage {
	am.DeclineFunc = onDecline
	return am
}

// Send pushes the accept message into the specified
// channel and sets up listener handlers for reactions.
func (am *AcceptMessage) Send(chanID string) (*AcceptMessage, error) {
	if am.Session == nil {
		return nil, errors.New("session is not defined")
	}
	if am.Embed == nil {
		return nil, errors.New("embed not defined")
	}

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
		if e.MessageID != msg.ID {
			return
		}

		if e.UserID != am.Session.State.User.ID {
			am.Session.MessageReactionRemove(am.ChannelID, am.ID, e.Emoji.Name, e.UserID)
		}

		if e.UserID == s.State.User.ID || (am.UserID != "" && am.UserID != e.UserID) {
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
