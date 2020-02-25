package inits

import (
	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/core/config"
	"github.com/zekroTJA/shinpuru/internal/core/database"
	"github.com/zekroTJA/shinpuru/internal/core/twitchnotify"
	"github.com/zekroTJA/shinpuru/internal/listeners"
	"github.com/zekroTJA/shinpuru/internal/util"
)

func InitTwitchNotifyer(session *discordgo.Session, config *config.Config, db database.Database) *twitchnotify.NotifyWorker {
	if config.Etc == nil || config.Etc.TwitchAppID == "" {
		return nil
	}

	listener := listeners.NewListenerTwitchNotify(session, config, db)
	tnw := twitchnotify.NewNotifyWorker(config.Etc.TwitchAppID,
		listener.HandlerWentOnline, listener.HandlerWentOffline)

	notifies, err := db.GetAllTwitchNotifies("")
	if err == nil {
		for _, notify := range notifies {
			if u, err := tnw.GetUser(notify.TwitchUserID, twitchnotify.TwitchNotifyIdentID); err == nil {
				tnw.AddUser(u)
			}
		}
	} else {
		util.Log.Error("failed getting Twitch notify entreis: ", err)
	}

	return tnw
}
