package listeners

import (
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/sarulabs/di/v2"
	"github.com/sirupsen/logrus"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/services/guildlog"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/twitchnotify"
)

type ListenerTwitchNotify struct {
	db      database.Database
	gl      guildlog.Logger
	session *discordgo.Session

	mx        *sync.RWMutex
	notMsgIDs map[string][]*discordgo.Message
}

func NewListenerTwitchNotify(container di.Container) *ListenerTwitchNotify {
	return &ListenerTwitchNotify{
		db:        container.Get(static.DiDatabase).(database.Database),
		gl:        container.Get(static.DiGuildLog).(guildlog.Logger).Section("twitchnotify"),
		session:   container.Get(static.DiDiscordSession).(*discordgo.Session),
		mx:        &sync.RWMutex{},
		notMsgIDs: make(map[string][]*discordgo.Message),
	}
}

func (l *ListenerTwitchNotify) TearDown() {
	if l == nil {
		return
	}

	l.mx.RLock()
	defer l.mx.RUnlock()
	for _, msgs := range l.notMsgIDs {
		for _, msg := range msgs {
			l.session.ChannelMessageDelete(msg.ChannelID, msg.ID)
		}
	}
}

func (l *ListenerTwitchNotify) HandlerWentOnline(d *twitchnotify.Stream, u *twitchnotify.User) {
	if l.session == nil {
		return
	}

	nots, err := l.db.GetAllTwitchNotifies(u.ID)
	if err != nil {
		logrus.WithError(err).Fatal("Faield getting Twitch notify entries from database")
		return
	}

	msgs := make([]*discordgo.Message, 0)
	for _, not := range nots {
		emb := twitchnotify.GetEmbed(d, u)
		msg, err := l.session.ChannelMessageSendEmbed(not.ChannelID, emb)
		if err != nil {
			if err = l.db.DeleteTwitchNotify(u.ID, not.GuildID); err != nil {
				logrus.WithError(err).Fatal("Failed removing Twitch notify entry from database")
				l.gl.Errorf(not.GuildID, "Failed removing twitch notify entry from database (%s): %s", u.ID, err.Error())
			}
			return
		}
		msgs = append(msgs, msg)
	}

	l.mx.Lock()
	defer l.mx.Unlock()
	l.notMsgIDs[d.ID] = msgs
}

func (l *ListenerTwitchNotify) HandlerWentOffline(d *twitchnotify.Stream, u *twitchnotify.User) {
	if l.session == nil {
		return
	}

	l.mx.RLock()
	defer l.mx.RUnlock()
	if msgs, ok := l.notMsgIDs[d.ID]; ok {
		for _, msg := range msgs {
			l.session.ChannelMessageDelete(msg.ChannelID, msg.ID)
		}
	}
}
