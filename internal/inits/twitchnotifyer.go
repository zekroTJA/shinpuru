package inits

import (
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/core/config"
	"github.com/zekroTJA/shinpuru/internal/core/database"
	"github.com/zekroTJA/shinpuru/internal/core/listeners"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/twitchnotify"
)

func InitTwitchNotifyWorker(container di.Container) *twitchnotify.NotifyWorker {

	listener := container.Get(static.DiTwitchNotifyListener).(*listeners.ListenerTwitchNotify)
	cfg := container.Get(static.DiConfig).(*config.Config)
	db := container.Get(static.DiDatabase).(database.Database)

	if cfg.TwitchApp == nil {
		return nil
	}

	tnw, err := twitchnotify.New(twitchnotify.Credentials{
		ClientID:     cfg.TwitchApp.ClientID,
		ClientSecret: cfg.TwitchApp.ClientSecret,
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

	return tnw
}
