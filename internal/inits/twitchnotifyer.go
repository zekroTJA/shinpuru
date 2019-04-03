package inits

import (
	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/core"
	"github.com/zekroTJA/shinpuru/internal/listeners"
	"github.com/zekroTJA/shinpuru/internal/util"
)

func InitTwitchNotifyer(session *discordgo.Session, config *core.Config, db core.Database) *core.TwitchNotifyWorker {
	if config.Etc == nil || config.Etc.TwitchAppID == "" {
		return nil
	}

	listener := listeners.NewListenerTwitchNotify(session, config, db)
	tnw := core.NewTwitchNotifyWorker(config.Etc.TwitchAppID,
		listener.HandlerWentOnline, listener.HandlerWentOffline)

	notifies, err := db.GetAllTwitchNotifies("")
	if err == nil {
		for _, notify := range notifies {
			if u, err := tnw.GetUser(notify.TwitchUserID, core.TwitchNotifyIdentID); err == nil {
				tnw.AddUser(u)
			}
		}
	} else {
		util.Log.Error("failed getting Twitch notify entreis: ", err)
	}

	return tnw
}
