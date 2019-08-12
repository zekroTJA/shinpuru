package inits

import (
	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/commands"
	"github.com/zekroTJA/shinpuru/internal/core"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/webserver"
)

func InitWebServer(s *discordgo.Session, db core.Database, cmdHandler *commands.CmdHandler, cfg *core.Config) (ws *webserver.WebServer) {
	if cfg.WebServer != nil && cfg.WebServer.Enabled {
		ws = webserver.NewWebServer(db, s, cmdHandler, cfg, cfg.Discord.ClientID, cfg.Discord.ClientSecret)
		go ws.ListenAndServeBlocking()
		util.Log.Info("Web server running on address %s (%s)...", cfg.WebServer.Addr, cfg.WebServer.PublicAddr)
	}
	return
}
