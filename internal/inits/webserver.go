package inits

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/core/config"
	"github.com/zekroTJA/shinpuru/internal/core/database"
	"github.com/zekroTJA/shinpuru/internal/core/middleware"
	"github.com/zekroTJA/shinpuru/internal/core/storage"
	"github.com/zekroTJA/shinpuru/internal/core/webserver"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/lctimer"
	"github.com/zekroTJA/shinpuru/pkg/mimefix"
	"github.com/zekroTJA/shinpuru/pkg/onetimeauth"
	"github.com/zekroTJA/shireikan"
)

func InitWebServer(container di.Container) (ws *webserver.WebServer) {

	session := container.Get(static.DiDiscordSession).(*discordgo.Session)
	cfg := container.Get(static.DiConfig).(*config.Config)
	db := container.Get(static.DiDatabase).(database.Database)
	storage := container.Get(static.DiObjectStorage).(storage.Storage)
	lct := container.Get(static.DiLifecycleTimer).(*lctimer.LifeCycleTimer)
	pmw := container.Get(static.DiPermissionMiddleware).(*middleware.PermissionsMiddleware)
	ota := container.Get(static.DiOneTimeAuth).(*onetimeauth.OneTimeAuth)
	cmdHandler := container.Get(static.DiCommandHandler).(shireikan.Handler)

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

		ws, err := webserver.New(db, storage, session, cmdHandler, lct, cfg, pmw, ota)
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
