// Package acceptmsg provides a message model for
// discordgo which can be accepted or declined
// via message reactions.
package acceptmsg

import (
	"errors"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekrotja/ken"
)

const (
	acceptMessageEmoteAccept  = "✅"
	acceptMessageEmoteDecline = "❌"
)

var (
	ErrTimeout = errors.New("timed out")
)

type ActionHandler func(*discordgo.Message) error

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
	cErr           chan error
	timeout        time.Duration

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

func (am *AcceptMessage) WithTimeout(t time.Duration) *AcceptMessage {
	am.timeout = t
	return am
}

// Error blocks until either one of the action functions was
// called or until the accept message timed out. Then, it
// returns an error or nil.
func (am *AcceptMessage) Error() error {
	return <-am.cErr
}

type senderFunc func(emb *discordgo.MessageEmbed) (*discordgo.Message, error)

// Send pushes the accept message into the specified
// channel and sets up listener handlers for reactions.
func (am *AcceptMessage) send(sender senderFunc) (*AcceptMessage, error) {
	if am.Session == nil {
		return nil, errors.New("session is not defined")
	}
	if am.Embed == nil {
		return nil, errors.New("embed not defined")
	}

	if am.timeout <= 0 {
		am.timeout = 1 * time.Minute
	}

	if am.AcceptFunc == nil {
		am.AcceptFunc = func(m *discordgo.Message) error {
			return nil
		}
	}

	if am.DeclineFunc == nil {
		am.DeclineFunc = func(m *discordgo.Message) error {
			return nil
		}
	}

	am.cErr = make(chan error, 1)

	msg, err := sender(am.Embed)
	if err != nil {
		return nil, err
	}
	am.Message = msg
	err = am.Session.MessageReactionAdd(msg.ChannelID, msg.ID, acceptMessageEmoteAccept)
	err = am.Session.MessageReactionAdd(msg.ChannelID, msg.ID, acceptMessageEmoteDecline)
	if err != nil {
		return nil, err
	}

	go func() {
		time.Sleep(am.timeout)
		am.cErr <- ErrTimeout

		if am.eventUnsub != nil {
			am.eventUnsub()
		}
	}()

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
			am.cErr <- am.AcceptFunc(msg)
		case acceptMessageEmoteDecline:
			am.cErr <- am.DeclineFunc(msg)
		}
		am.eventUnsub()
		if am.DeleteMsgAfter {
			am.Session.ChannelMessageDelete(msg.ChannelID, msg.ID)
		} else {
			am.Session.MessageReactionsRemoveAll(msg.ChannelID, msg.ID)
		}
	})

	return am, nil
}

// Send pushes the accept message into the specified
// channel and sets up listener handlers for reactions.
func (am *AcceptMessage) Send(chanID string) (*AcceptMessage, error) {
	return am.send(func(emb *discordgo.MessageEmbed) (*discordgo.Message, error) {
		return am.Session.ChannelMessageSendEmbed(chanID, am.Embed)
	})
}

// AsFollowUp pushes the accept messages as follow up
// message to the command context and sets up listener
// handlers for reactions.
func (am *AcceptMessage) AsFollowUp(ctx ken.Context) (*AcceptMessage, error) {
	return am.send(func(emb *discordgo.MessageEmbed) (*discordgo.Message, error) {
		fum := ctx.FollowUpEmbed(am.Embed).Send()
		return fum.Message, fum.Error
	})
}
