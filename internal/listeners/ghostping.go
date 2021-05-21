package listeners

import (
	"regexp"
	"strings"
	"time"

	"github.com/sarulabs/di/v2"
	"github.com/sirupsen/logrus"
	"github.com/zekroTJA/shinpuru/internal/middleware"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/services/guildlog"
	"github.com/zekroTJA/shinpuru/internal/util/static"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/timedmap"
)

const (
	gpDelay = 6 * time.Hour
	gpTick  = 5 * time.Minute
)

type ListenerGhostPing struct {
	db              database.Database
	gl              guildlog.Logger
	msgCache        *timedmap.TimedMap
	recentlyDeleted map[string]struct{}
	gpim            *middleware.GhostPingIgnoreMiddleware
}

func NewListenerGhostPing(container di.Container) *ListenerGhostPing {
	return &ListenerGhostPing{
		db:       container.Get(static.DiDatabase).(database.Database),
		gl:       container.Get(static.DiGuildLog).(guildlog.Logger).Section("ghostping"),
		gpim:     container.Get(static.DiGhostpingIgnoreMiddleware).(*middleware.GhostPingIgnoreMiddleware),
		msgCache: timedmap.New(gpTick),
	}
}

func (l *ListenerGhostPing) HandlerMessageCreate(s *discordgo.Session, e *discordgo.MessageCreate) {
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
}

func (l *ListenerGhostPing) HandlerMessageDelete(s *discordgo.Session, e *discordgo.MessageDelete) {
	rx := regexp.MustCompile(`(@here)|(@everyone)`)

	l.msgCache.Set(e.ID, e.Message, gpDelay)

	if l.gpim.ContainsAndRemove(e.ID) {
		return
	}

	v := l.msgCache.GetValue(e.ID)
	if v == nil {
		return
	}

	deletedMsg, ok := v.(*discordgo.Message)
	if !ok {
		return
	}

	if deletedMsg.Author == nil || deletedMsg.Author.ID == s.State.User.ID {
		return
	}

	gpMsg, err := l.db.GetGuildGhostpingMsg(deletedMsg.GuildID)
	if err != nil {
		if !database.IsErrDatabaseNotFound(err) {
			logrus.WithError(err).WithField("gid", deletedMsg.GuildID).Error("GHOSTPING :: failed getting ghost ping msg")
			l.gl.Errorf(deletedMsg.GuildID, "Failed getting ghost ping message: %s", err.Error())
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

	_, err = s.ChannelMessageSend(deletedMsg.ChannelID, gpMsg)
	if err != nil {
		l.gl.Errorf(deletedMsg.GuildID, "Failed sending ghost ping message: %s", err.Error())
	}

	l.msgCache.Remove(e.ID)
}
