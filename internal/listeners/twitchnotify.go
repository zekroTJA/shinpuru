package listeners

import (
	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/core"
	"github.com/zekroTJA/shinpuru/internal/util"
)

type ListenerTwitchNotify struct {
	config    *core.Config
	db        core.Database
	session   *discordgo.Session
	notMsgIDs map[string][]*discordgo.Message
}

func NewListenerTwitchNotify(session *discordgo.Session, config *core.Config, db core.Database) *ListenerTwitchNotify {
	return &ListenerTwitchNotify{
		config:    config,
		db:        db,
		session:   session,
		notMsgIDs: make(map[string][]*discordgo.Message),
	}
}

func (l *ListenerTwitchNotify) HandlerWentOnline(d *core.TwitchNotifyData, u *core.TwitchNotifyUser) {
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
		emb := core.TwitchNotifyGetEmbed(d, u)
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

func (l *ListenerTwitchNotify) HandlerWentOffline(d *core.TwitchNotifyData, u *core.TwitchNotifyUser) {
	if l.session == nil {
		return
	}

	if msgs, ok := l.notMsgIDs[d.ID]; ok {
		for _, msg := range msgs {
			l.session.ChannelMessageDelete(msg.ChannelID, msg.ID)
		}
	}
}
