package listeners

import (
	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/core"
	"github.com/zekroTJA/shinpuru/internal/util"
)

type ListenerTwitchNotify struct {
	config  *core.Config
	db      core.Database
	session *discordgo.Session
}

func NewListenerTwitchNotify(session *discordgo.Session, config *core.Config, db core.Database) *ListenerTwitchNotify {
	return &ListenerTwitchNotify{
		config:  config,
		db:      db,
		session: session,
	}
}

func (l *ListenerTwitchNotify) Handler(d *core.TwitchNotifyData, u *core.TwitchNotifyUser) {
	if l.session == nil {
		return
	}

	nots, err := l.db.GetAllTwitchNotifies(u.ID)
	if err != nil {
		util.Log.Error("Faield getting Twitch notify entries from database: ", err)
		return
	}

	for _, not := range nots {
		emb := core.TwitchNotifyGetEmbed(d, u)
		_, err := l.session.ChannelMessageSendEmbed(not.ChannelID, emb)
		if err != nil {
			if err = l.db.DeleteTwitchNotify(u.ID, not.GuildID); err != nil {
				util.Log.Error("Failed removing Twitch notify entry from database: ", err)
			}
		}
	}
}
