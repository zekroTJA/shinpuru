package inits

import (
	"fmt"

	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/core/config"
	"github.com/zekroTJA/shinpuru/internal/core/webserver"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/mimefix"
)

func InitWebServer(container di.Container) (ws *webserver.WebServer) {

	cfg := container.Get(static.DiConfig).(*config.Config)

	if cfg.WebServer != nil && cfg.WebServer.Enabled {
		curr, ok := mimefix.Check()
		if !ok {
			util.Log.Infof("Mime check of .js returned invalid mime value '%s', trying to fix this now...", curr)
			if err := mimefix.Fix(); err != nil {
				util.Log.Errorf("Fixing .js mime value failed (maybe run as admin to fix this): %s", err.Error())
				util.Log.Warning("Mime value of .js was not fixed. This may lead to erroneous behaviour of the web server")
			} else {
				util.Log.Info("Successfully fixed .js mime value")
			}
		}

		ws, err := webserver.New(container)
		if err != nil {
			util.Log.Fatalf(fmt.Sprintf("Failed initializing web server: %s", err.Error()))
		}

		go func() {
			if err = ws.ListenAndServeBlocking(); err != nil {
				util.Log.Fatalf("Failed starting up web server: %s", err.Error())
			}
		}()
		util.Log.Info(fmt.Sprintf("Web server running on address %s (%s)...", cfg.WebServer.Addr, cfg.WebServer.PublicAddr))
	}
	return
}
