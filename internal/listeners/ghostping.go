package listeners

import (
	"regexp"
	"strings"
	"time"

	"github.com/zekroTJA/shinpuru/internal/commands"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/core"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/timedmap"
)

const (
	gpDelay = 6 * time.Hour
	gpTick  = 5 * time.Minute
)

type ListenerGhostPing struct {
	db                 core.Database
	cmdHandler         *commands.CmdHandler
	deleteHandlerAdded bool
	msgCache           *timedmap.TimedMap
}

func NewListenerGhostPing(db core.Database, cmdHandler *commands.CmdHandler) *ListenerGhostPing {
	return &ListenerGhostPing{
		db:         db,
		cmdHandler: cmdHandler,
		msgCache:   timedmap.New(gpTick),
	}
}

func (l *ListenerGhostPing) Handler(s *discordgo.Session, e *discordgo.MessageCreate) {
	if e.Author.Bot || len(e.Mentions) == 0 {
		return
	}

	userMentions := make([]*discordgo.User, 0)
	for _, ment := range e.Mentions {
		if !ment.Bot {
			userMentions = append(userMentions, ment)
		}
	}
	e.Mentions = userMentions

	if len(e.Mentions) == 0 {
		return
	}
	rx := regexp.MustCompile(`(@here)|(@everyone)`)

	l.msgCache.Set(e.ID, e.Message, gpDelay)

	if !l.deleteHandlerAdded {
		s.AddHandler(func(_ *discordgo.Session, eDel *discordgo.MessageDelete) {
			if l.cmdHandler.GetNotifiedCommandMsgs().Contains(eDel.ID) {
				l.cmdHandler.GetNotifiedCommandMsgs().Remove(eDel.ID)
				return
			}

			v := l.msgCache.GetValue(eDel.ID)
			if v == nil {
				return
			}

			deletedMsg, ok := v.(*discordgo.Message)
			if !ok {
				return
			}

			if deletedMsg.Author.ID == s.State.User.ID {
				return
			}

			gpMsg, err := l.db.GetGuildGhostpingMsg(deletedMsg.GuildID)
			if err != nil {
				if !core.IsErrDatabaseNotFound(err) {
					util.Log.Errorf("failed getting ghost ping msg for guild %s: %s\n", deletedMsg.GuildID, err.Error())
				}
				return
			}

			uPinged := deletedMsg.Mentions[0]

			if uPinged.ID == deletedMsg.Author.ID {
				return
			}

			deletedMsg.Content = rx.ReplaceAllStringFunc(deletedMsg.Content, func(s string) string {
				return "[@]" + s[1:]
			})

			if uPinged.Bot {
				return
			}

			gpMsg = strings.Replace(gpMsg, "{@pinger}", deletedMsg.Author.Mention(), -1)
			gpMsg = strings.Replace(gpMsg, "{@pinged}", uPinged.Mention(), -1)
			gpMsg = strings.Replace(gpMsg, "{pinger}", deletedMsg.Author.String(), -1)
			gpMsg = strings.Replace(gpMsg, "{pinged}", uPinged.String(), -1)
			gpMsg = strings.Replace(gpMsg, "{msg}", deletedMsg.Content, -1)

			s.ChannelMessageSend(deletedMsg.ChannelID, gpMsg)

			l.msgCache.Remove(eDel.ID)
		})

		l.deleteHandlerAdded = true
	}
}
