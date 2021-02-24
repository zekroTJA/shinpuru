package inits

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/core/config"
	"github.com/zekroTJA/shinpuru/internal/core/database"
	"github.com/zekroTJA/shinpuru/internal/core/middleware"
	"github.com/zekroTJA/shinpuru/internal/core/storage"
	"github.com/zekroTJA/shinpuru/internal/core/webserver"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/pkg/lctimer"
	"github.com/zekroTJA/shinpuru/pkg/mimefix"
	"github.com/zekroTJA/shireikan"
)

func InitWebServer(s *discordgo.Session, db database.Database, st storage.Storage,
	cmdHandler shireikan.Handler, lct *lctimer.LifeCycleTimer, cfg *config.Config, pmw *middleware.PermissionsMiddleware) (ws *webserver.WebServer) {

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

		ws, err := webserver.New(db, st, s, cmdHandler, lct, cfg, pmw)
		if err != nil {
			util.Log.Fatalf(fmt.Sprintf("Failed initializing web server: %s", err.Error()))
		}

		go ws.ListenAndServeBlocking()
		util.Log.Info(fmt.Sprintf("Web server running on address %s (%s)...", cfg.WebServer.Addr, cfg.WebServer.PublicAddr))
	}
	return
}
