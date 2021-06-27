// Package msgcollector provides functionalities to
// collect messages in a channel in conect of a single
// command request.
package msgcollector

import (
	"errors"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Options for MessageCollector
// initialization.
type Options struct {
	MaxMessages        int
	MaxMatches         int
	Timeout            time.Duration
	DeleteMatchesAfter bool
}

// MessageCollector provides functionalities to collect
// sent messages in a channel and matching them by a
// specified filter function.
type MessageCollector struct {
	collectedMessages []*discordgo.Message
	collectedMatches  []*discordgo.Message
	closed            bool
	channelID         string
	session           *discordgo.Session
	options           *Options
	filter            func(*discordgo.Message) bool
	eventUnsub        func()
	onCollected       func(*discordgo.Message, *MessageCollector)
	onMatched         func(*discordgo.Message, *MessageCollector)
	onClosed          func(string, *MessageCollector)
}

// New creates a new instance of MessageCollector using the passed
// Session s, channelID, filter function and options.
func New(s *discordgo.Session, channelID string, filter func(*discordgo.Message) bool, options *Options) (*MessageCollector, error) {
	if s == nil {
		return nil, errors.New("session is not defined")
	}
	if options == nil {
		options = new(Options)
	}

	mc := &MessageCollector{
		session:           s,
		channelID:         channelID,
		filter:            filter,
		options:           options,
		collectedMessages: make([]*discordgo.Message, 0),
	}

	mc.eventUnsub = mc.session.AddHandler(mc.messageCreateHandler)

	if mc.options.Timeout != 0*time.Second {
		time.AfterFunc(mc.options.Timeout, func() {
			if mc != nil {
				mc.Close("timeout")
			}
		})
	}
	return mc, nil
}

// Close unregisters the message create event listener,
// sets the status of the message collector to closed
// and removed all collected messages if set in options.
func (mc *MessageCollector) Close(reason string) {
	if mc.closed {
		return
	}
	if mc.eventUnsub != nil {
		mc.eventUnsub()
	}
	if mc.onClosed != nil {
		mc.onClosed(reason, mc)
	}
	mc.closed = true
	if mc.options.DeleteMatchesAfter {
		for _, msg := range mc.collectedMatches {
			mc.session.ChannelMessageDelete(mc.channelID, msg.ID)
		}
	}
}

// OnCollected sets a handler function which is called
// everytime any message was collected in the specified
// channel.
func (mc *MessageCollector) OnColelcted(handler func(*discordgo.Message, *MessageCollector)) {
	mc.onCollected = handler
}

// OnMatch sets a handler function which is called
// everytime a message was collected which went
// positively through the specified match filter.
func (mc *MessageCollector) OnMatched(handler func(*discordgo.Message, *MessageCollector)) {
	mc.onMatched = handler
}

// OnClosed is called when the message collector
// was closed actively or automatically getting
// passed the reason why the message collector
// was closed.
func (mc *MessageCollector) OnClosed(handler func(string, *MessageCollector)) {
	mc.onClosed = handler
}

// IsClosed returns the current closed status of
// the message collector.
func (mc *MessageCollector) IsClosed() bool {
	return mc.closed
}

// CollectedMatches returns the array of collected
// messages which went positively through the specified
// filter function.
func (mc *MessageCollector) CollectedMatches() []*discordgo.Message {
	return mc.collectedMatches
}

// CollectedMessages returns the array of collected
// messages.
func (mc *MessageCollector) CollectedMessages() []*discordgo.Message {
	return mc.collectedMessages
}

// messageCreateHandler is the event handler set to
// the discordgo session to process incomming messages.
func (mc *MessageCollector) messageCreateHandler(s *discordgo.Session, msg *discordgo.MessageCreate) {
	if msg.ChannelID != mc.channelID {
		return
	}

	mc.collectedMessages = append(mc.collectedMessages, msg.Message)

	if mc.onCollected != nil {
		mc.onCollected(msg.Message, mc)
	}

	if mc.filter(msg.Message) {
		mc.collectedMatches = append(mc.collectedMatches, msg.Message)
		if mc.onMatched != nil {
			mc.onMatched(msg.Message, mc)
		}
	}

	if mc.options.MaxMessages != 0 && len(mc.collectedMessages) >= mc.options.MaxMessages {
		mc.Close("maxMessagesReached")
	} else if mc.options.MaxMatches != 0 && len(mc.collectedMatches) >= mc.options.MaxMatches {
		mc.Close("maxMatchesReached")
	}
}
