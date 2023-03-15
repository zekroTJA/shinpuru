package inits

import (
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/services/config"
	"github.com/zekroTJA/shinpuru/internal/services/webserver"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekrotja/rogu/log"
)

func InitWebServer(container di.Container) (ws *webserver.WebServer) {

	cfg := container.Get(static.DiConfig).(config.Provider)

	if cfg.Config().WebServer.Enabled {
		log := log.Tagged("WebServer")
		log.Info().Msg("Initializing web server ...")

		ws, err := webserver.New(container)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed initializing web server")
		}

		go func() {
			if err = ws.ListenAndServeBlocking(); err != nil {
				log.Fatal().Err(err).Msg("Failed starting up web server")
			}
		}()
		log.Info().Fields(
			"bindAddr", cfg.Config().WebServer.Addr,
			"publicAddr", cfg.Config().WebServer.PublicAddr,
		).Msg("Web server running")
	}
	return
}
