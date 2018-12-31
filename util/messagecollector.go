package util

import (
	"errors"
	"time"

	"github.com/bwmarrin/discordgo"
)

type MessageCollectorOptions struct {
	MaxMessages        int
	MaxMatches         int
	Timeout            time.Duration
	DeleteMatchesAfter bool
}

type MessageCollector struct {
	CollectedMessages []*discordgo.Message
	CollectedMatches  []*discordgo.Message
	Closed            bool
	channelID         string
	session           *discordgo.Session
	options           *MessageCollectorOptions
	filter            func(*discordgo.Message) bool
	eventUnsub        func()
	onCollected       func(*discordgo.Message, *MessageCollector)
	onMatched         func(*discordgo.Message, *MessageCollector)
	onClosed          func(string, *MessageCollector)
}

func NewMessageCollector(s *discordgo.Session, channelID string, filter func(*discordgo.Message) bool, options *MessageCollectorOptions) (*MessageCollector, error) {
	if s == nil {
		return nil, errors.New("session is not defined")
	}
	if options == nil {
		options = new(MessageCollectorOptions)
	}
	mc := &MessageCollector{
		session:           s,
		channelID:         channelID,
		filter:            filter,
		options:           options,
		CollectedMessages: make([]*discordgo.Message, 0),
	}
	mc.eventUnsub = mc.session.AddHandler(func(s *discordgo.Session, msg *discordgo.MessageCreate) {
		if msg.ChannelID != mc.channelID {
			return
		}
		mc.CollectedMessages = append(mc.CollectedMessages, msg.Message)
		if mc.onCollected != nil {
			mc.onCollected(msg.Message, mc)
		}
		if mc.filter(msg.Message) {
			mc.CollectedMatches = append(mc.CollectedMatches, msg.Message)
			if mc.onMatched != nil {
				mc.onMatched(msg.Message, mc)
			}
		}
		if mc.options.MaxMessages != 0 && len(mc.CollectedMessages) >= mc.options.MaxMessages {
			mc.Close("maxMessagesReached")
		} else if mc.options.MaxMatches != 0 && len(mc.CollectedMatches) >= mc.options.MaxMatches {
			mc.Close("maxMatchesReached")
		}
	})
	if mc.options.Timeout != 0*time.Second {
		time.AfterFunc(mc.options.Timeout, func() {
			// TODO: test if throws nil pointer exception
			mc.Close("timeout")
		})
	}
	return mc, nil
}

func (mc *MessageCollector) Close(reason string) {
	if mc.Closed {
		return
	}
	if mc.eventUnsub != nil {
		mc.eventUnsub()
	}
	if mc.onClosed != nil {
		mc.onClosed(reason, mc)
	}
	mc.Closed = true
	if mc.options.DeleteMatchesAfter {
		for _, msg := range mc.CollectedMatches {
			mc.session.ChannelMessageDelete(mc.channelID, msg.ID)
		}
	}
}

func (mc *MessageCollector) OnColelcted(handler func(*discordgo.Message, *MessageCollector)) {
	mc.onCollected = handler
}

func (mc *MessageCollector) OnMatched(handler func(*discordgo.Message, *MessageCollector)) {
	mc.onMatched = handler
}

func (mc *MessageCollector) OnClosed(handler func(string, *MessageCollector)) {
	mc.onClosed = handler
}
