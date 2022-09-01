package inits

import (
	"github.com/sarulabs/di/v2"
	"github.com/sirupsen/logrus"
	"github.com/zekroTJA/shinpuru/internal/services/config"
	"github.com/zekroTJA/shinpuru/internal/services/webserver"
	"github.com/zekroTJA/shinpuru/internal/util/static"
)

func InitWebServer(container di.Container) (ws *webserver.WebServer) {

	cfg := container.Get(static.DiConfig).(config.Provider)

	if cfg.Config().WebServer.Enabled {
		ws, err := webserver.New(container)
		if err != nil {
			logrus.WithError(err).Fatal("Failed initializing web server")
		}

		go func() {
			if err = ws.ListenAndServeBlocking(); err != nil {
				logrus.WithError(err).Fatal("Failed starting up web server")
			}
		}()
		logrus.WithFields(logrus.Fields{
			"bindAddr":   cfg.Config().WebServer.Addr,
			"publicAddr": cfg.Config().WebServer.PublicAddr,
		}).Info("Web server running")
	}
	return
}
