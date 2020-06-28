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
	if config.TwitchApp == nil {
		return nil
	}

	listener := listeners.NewListenerTwitchNotify(session, config, db)
	tnw, err := twitchnotify.New(config.TwitchApp,
		listener.HandlerWentOnline, listener.HandlerWentOffline)

	if err != nil {
		util.Log.Fatalf("twitch app credentials are invalid: %s", err)
	}

	notifies, err := db.GetAllTwitchNotifies("")
	if err == nil {
		for _, notify := range notifies {
			if u, err := tnw.GetUser(notify.TwitchUserID, twitchnotify.IdentID); err == nil {
				tnw.AddUser(u)
			}
		}
	} else {
		util.Log.Error("failed getting Twitch notify entreis: ", err)
	}

	return tnw
}
