package webserver

import (
	"errors"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/zekroTJA/shinpuru/internal/commands"
	"github.com/zekroTJA/shinpuru/internal/core"
	"github.com/zekroTJA/shinpuru/pkg/discordoauth"

	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

// Error Objects
var (
	errNotFound         = errors.New("not found")
	errInvalidArguments = errors.New("invalid arguments")
	errNoAccess         = errors.New("access denied")
	errUnauthorized     = errors.New("unauthorized")
)

const (
	endpointLogInWithDC = "/_/loginwithdiscord"
	endpointAuthCB      = "/_/authorizationcallback"
)

// Static File Handlers
var (
	fileHandlerStatic = fasthttp.FS{
		Root:       "./web/dist/web",
		IndexNames: []string{"index.html"},
		Compress:   true,
		// PathRewrite: func(ctx *fasthttp.RequestCtx) []byte {
		// 	return ctx.Path()[7:]
		// },
	}
)

type WebServer struct {
	server *fasthttp.Server
	router *routing.Router

	db         core.Database
	rlm        *RateLimitManager
	auth       *Auth
	dcoauth    *discordoauth.DiscordOAuth
	session    *discordgo.Session
	cmdhandler *commands.CmdHandler

	config *core.Config
}

func NewWebServer(db core.Database, s *discordgo.Session, cmd *commands.CmdHandler, config *core.Config, clientID, clientSecret string) (ws *WebServer) {
	ws = new(WebServer)

	if !strings.HasPrefix(config.WebServer.PublicAddr, "http") {
		protocol := "http"
		if config.WebServer.TLS != nil && config.WebServer.TLS.Enabled {
			protocol += "s"
		}
		config.WebServer.PublicAddr = fmt.Sprintf("%s://%s", protocol, config.WebServer.PublicAddr)
	}

	ws.config = config
	ws.db = db
	ws.session = s
	ws.cmdhandler = cmd
	ws.rlm = NewRateLimitManager()
	ws.router = routing.New()
	ws.server = &fasthttp.Server{
		Handler: ws.router.HandleRequest,
	}

	ws.auth = NewAuth(db, s)

	ws.dcoauth = discordoauth.NewDiscordOAuth(
		clientID,
		clientSecret,
		config.WebServer.PublicAddr+endpointAuthCB,
		ws.auth.LoginFailedHandler,
		ws.auth.LoginSuccessHandler,
	)

	ws.registerHandlers()

	return
}

func (ws *WebServer) registerHandlers() {
	// rlGlobal := ws.rlm.GetHandler(500*time.Millisecond, 50)
	// rlUsersCreate := ws.rlm.GetHandler(15*time.Second, 1)
	// rlPageCreate := ws.rlm.GetHandler(5*time.Second, 5)

	ws.router.Use(ws.addHeaders, ws.optionsHandler, ws.auth.checkAuth, ws.handlerFiles)

	ws.router.Get(endpointLogInWithDC, ws.dcoauth.HandlerInit)
	ws.router.Get(endpointAuthCB, ws.dcoauth.HandlerCallback)

	api := ws.router.Group("/api")
	api.
		Get("/me", ws.handlerGetMe)
	api.
		Post("/logout", ws.auth.LogOutHandler)
	api.
		Get("/sysinfo", ws.handlerGetSystemInfo)

	settings := api.Group("/settings")
	settings.
		Get("/presence", ws.handlerGetPresence).
		Post(ws.handlerPostPresence)
	settings.
		Get("/noguildinvite", ws.handlerGetInviteSettings).
		Post(ws.handlerPostInviteSettings)

	guilds := api.Group("/guilds")
	guilds.
		Get("", ws.handlerGuildsGet)

	guild := guilds.Group("/<guildid:[0-9]+>")
	guild.
		Get("", ws.handlerGuildsGetGuild)
	guild.
		Get("/settings", ws.handlerGetGuildSettings).
		Post(ws.handlerPostGuildSettings)
	guild.
		Get("/permissions", ws.handlerGetGuildPermissions).
		Post(ws.handlerPostGuildPermissions)
	guild.
		Get("/members", ws.handlerGuildGetMembers)

	guildReports := guild.Group("/reports")
	guildReports.
		Get("", ws.handlerGetReports)
	guildReports.
		Get("/count", ws.handlerGetReportsCount)

	member := guilds.Group("/<guildid:[0-9]+>/<memberid:[0-9]+>")
	member.
		Get("", ws.handlerGuildsGetMember)
	member.
		Get("/permissions", ws.handlerGetPermissions)
	member.
		Get("/permissions/allowed", ws.handlerGetPermissionsAllowed)
	member.
		Post("/kick", ws.handlerPostGuildMemberKick)
	member.
		Post("/ban", ws.handlerPostGuildMemberBan)

	memberReports := member.Group("/reports")
	memberReports.
		Get("", ws.handlerGetReports).
		Post(ws.handlerPostGuildMemberReport)
	memberReports.
		Get("/count", ws.handlerGetReportsCount)

	reports := api.Group("/reports")
	reports.
		Get("/<id:[0-9]+>", ws.handlerGetReport)
}

func (ws *WebServer) ListenAndServeBlocking() error {
	tls := ws.config.WebServer.TLS

	if tls.Enabled {
		if tls.Cert == "" || tls.Key == "" {
			return errors.New("cert file and key file must be specified")
		}
		return ws.server.ListenAndServeTLS(ws.config.WebServer.Addr, tls.Cert, tls.Key)
	}

	return ws.server.ListenAndServe(ws.config.WebServer.Addr)
}
