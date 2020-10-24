package inits

import (
	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/core/config"
	"github.com/zekroTJA/shinpuru/internal/core/database"
	"github.com/zekroTJA/shinpuru/internal/core/listeners"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/pkg/twitchnotify"
)

func InitTwitchNotifyer(session *discordgo.Session, config *config.Config, db database.Database) (*twitchnotify.NotifyWorker, *listeners.ListenerTwitchNotify) {
	if config.TwitchApp == nil {
		return nil, nil
	}

	listener := listeners.NewListenerTwitchNotify(session, config, db)
	tnw, err := twitchnotify.New(twitchnotify.Credentials{
		ClientID:     config.TwitchApp.ClientID,
		ClientSecret: config.TwitchApp.ClientSecret,
	}, listener.HandlerWentOnline, listener.HandlerWentOffline)

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

	return tnw, listener
}
