package listeners

import (
	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/core/config"
	"github.com/zekroTJA/shinpuru/internal/core/database"
	"github.com/zekroTJA/shinpuru/internal/core/twitchnotify"
	"github.com/zekroTJA/shinpuru/internal/util"
)

type ListenerTwitchNotify struct {
	config    *config.Config
	db        database.Database
	session   *discordgo.Session
	notMsgIDs map[string][]*discordgo.Message
}

func NewListenerTwitchNotify(session *discordgo.Session, config *config.Config, db database.Database) *ListenerTwitchNotify {
	return &ListenerTwitchNotify{
		config:    config,
		db:        db,
		session:   session,
		notMsgIDs: make(map[string][]*discordgo.Message),
	}
}

func (l *ListenerTwitchNotify) HandlerWentOnline(d *twitchnotify.NotifyData, u *twitchnotify.NotifyUser) {
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

func (l *ListenerTwitchNotify) HandlerWentOffline(d *twitchnotify.NotifyData, u *twitchnotify.NotifyUser) {
	if l.session == nil {
		return
	}

	if msgs, ok := l.notMsgIDs[d.ID]; ok {
		for _, msg := range msgs {
			l.session.ChannelMessageDelete(msg.ChannelID, msg.ID)
		}
	}
}
