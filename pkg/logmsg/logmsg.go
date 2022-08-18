package logmsg

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/pkg/voidbuffer/v2"
)

type LogMessage struct {
	*discordgo.Message
	session *discordgo.Session
	content *voidbuffer.VoidBuffer[string]
	emb     *discordgo.MessageEmbed
	cMsgs   <-chan string
	cErr    <-chan error
	cClose  chan string
}

func New(
	s *discordgo.Session,
	channelId string,
	emb *discordgo.MessageEmbed,
	cMsgs <-chan string,
	cErr <-chan error,
	initMsg string,
) (lm *LogMessage, err error) {
	return NewWithSender(s, func(e *discordgo.MessageEmbed) (*discordgo.Message, error) {
		return s.ChannelMessageSendEmbed(channelId, e)
	}, emb, cMsgs, cErr, initMsg)
}

func NewWithSender(
	s *discordgo.Session,
	sender func(emb *discordgo.MessageEmbed) (*discordgo.Message, error),
	emb *discordgo.MessageEmbed,
	cMsgs <-chan string,
	cErr <-chan error,
	initMsg string,
) (lm *LogMessage, err error) {
	lm = &LogMessage{
		session: s,
		content: voidbuffer.New[string](20),
		emb:     emb,
		cMsgs:   cMsgs,
		cErr:    cErr,
		cClose:  make(chan string, 1),
	}

	lm.content.Push("ℹ️ " + initMsg)

	lm.updateEmbed()
	lm.Message, err = sender(emb)
	if err != nil {
		return
	}
	go lm.watcher()
	return
}

func (lm *LogMessage) Close(msg string) {
	lm.cClose <- "✔️ " + msg
}

func (lm *LogMessage) watcher() {
	for {
		select {
		case <-lm.cClose:
			return
		case v := <-lm.cMsgs:
			lm.updateMessage(fmt.Sprintf("ℹ️ %s", v))
		case err := <-lm.cErr:
			if err == nil {
				continue
			}
			lm.updateMessage(fmt.Sprintf("⚠️ %s", err.Error()))
		}
	}
}

func (lm *LogMessage) updateEmbed() {
	var desc strings.Builder
	for _, e := range lm.content.Snapshot() {
		desc.WriteString(e)
		desc.WriteRune('\n')
	}
	lm.emb.Description = desc.String()
}

func (lm *LogMessage) updateMessage(v string) {
	lm.content.Push(v)
	lm.updateEmbed()
	lm.Message, _ = lm.session.ChannelMessageEditEmbed(
		lm.Message.ChannelID, lm.Message.ID, lm.emb)
}
