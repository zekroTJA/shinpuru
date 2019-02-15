package listeners

import (
	"strings"
	"time"

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
	deleteHandlerAdded bool
	msgCache           *timedmap.TimedMap
}

func NewListenerGhostPing(db core.Database) *ListenerGhostPing {
	return &ListenerGhostPing{
		db:       db,
		msgCache: timedmap.New(gpTick),
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

	l.msgCache.Set(e.ID, e.Message, gpDelay)

	if !l.deleteHandlerAdded {
		s.AddHandler(func(_ *discordgo.Session, eDel *discordgo.MessageDelete) {
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

			gpMsg, err := l.db.GetGuildGhostpingMsg(e.GuildID)
			if err != nil {
				if !core.IsErrDatabaseNotFound(err) {
					util.Log.Errorf("failed getting ghost ping msg for guild %s: %s\n", e.GuildID, err.Error())
				}
				return
			}

			uPinged := e.Mentions[0]

			if uPinged.Bot {
				return
			}

			gpMsg = strings.Replace(gpMsg, "{pinger}", deletedMsg.Author.Mention(), -1)
			gpMsg = strings.Replace(gpMsg, "{pinged}", uPinged.Mention(), -1)
			gpMsg = strings.Replace(gpMsg, "{msg}", deletedMsg.Content, -1)

			s.ChannelMessageSend(e.ChannelID, gpMsg)

			l.msgCache.Remove(eDel.ID)
		})

		l.deleteHandlerAdded = true
	}
}
