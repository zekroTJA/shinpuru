package inits

import (
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/listeners"
	"github.com/zekroTJA/shinpuru/internal/services/config"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/twitchnotify"
	"github.com/zekrotja/rogu/log"
)

func InitTwitchNotifyWorker(container di.Container) *twitchnotify.NotifyWorker {

	listener := container.Get(static.DiTwitchNotifyListener).(*listeners.ListenerTwitchNotify)
	cfg := container.Get(static.DiConfig).(config.Provider)
	db := container.Get(static.DiDatabase).(database.Database)

	if cfg.Config().TwitchApp.ClientID == "" || cfg.Config().TwitchApp.ClientSecret == "" {
		return nil
	}

	log := log.Tagged("TwitchNotify")
	log.Info().Msg("Initializing twitch notifications ...")

	tnw, err := twitchnotify.New(
		twitchnotify.Credentials{
			ClientID:     cfg.Config().TwitchApp.ClientID,
			ClientSecret: cfg.Config().TwitchApp.ClientSecret,
		},
		listener.HandlerWentOnline,
		listener.HandlerWentOffline,
		twitchnotify.Config{
			TimerDelay: 0,
		},
	)

	if err != nil {
		log.Fatal().Err(err).Msg("Twitch app credentials are invalid")
	}

	notifies, err := db.GetAllTwitchNotifies("")
	if err == nil {
		for _, notify := range notifies {
			if u, err := tnw.GetUser(notify.TwitchUserID, twitchnotify.IdentID); err == nil {
				tnw.AddUser(u)
			}
		}
	} else {
		log.Fatal().Err(err).Msg("Failed getting Twitch notify entreis")
	}

	return tnw
}
