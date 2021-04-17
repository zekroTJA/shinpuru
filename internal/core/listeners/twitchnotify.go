package listeners

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/core/config"
	"github.com/zekroTJA/shinpuru/internal/core/database"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/twitchnotify"
)

type ListenerTwitchNotify struct {
	config    *config.Config
	db        database.Database
	session   *discordgo.Session
	notMsgIDs map[string][]*discordgo.Message
}

func NewListenerTwitchNotify(container di.Container) *ListenerTwitchNotify {
	return &ListenerTwitchNotify{
		config:    container.Get(static.DiConfig).(*config.Config),
		db:        container.Get(static.DiDatabase).(database.Database),
		session:   container.Get(static.DiDiscordSession).(*discordgo.Session),
		notMsgIDs: make(map[string][]*discordgo.Message),
	}
}

func (l *ListenerTwitchNotify) TearDown() {
	if l == nil {
		return
	}

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
		util.Log.Error("Faield getting Twitch notify entries from database: ", err)
		return
	}

	msgs := make([]*discordgo.Message, 0)
	for _, not := range nots {
		emb := twitchnotify.GetEmbed(d, u)
		msg, err := l.session.ChannelMessageSendEmbed(not.ChannelID, emb)
		if err != nil {
			if err = l.db.DeleteTwitchNotify(u.ID, not.GuildID); err != nil {
				util.Log.Error("Failed removing Twitch notify entry from database: ", err)
			}
			return
		}
		msgs = append(msgs, msg)
	}
	l.notMsgIDs[d.ID] = msgs

}

func (l *ListenerTwitchNotify) HandlerWentOffline(d *twitchnotify.Stream, u *twitchnotify.User) {
	if l.session == nil {
		return
	}

	if msgs, ok := l.notMsgIDs[d.ID]; ok {
		for _, msg := range msgs {
			l.session.ChannelMessageDelete(msg.ChannelID, msg.ID)
		}
	}
}
