package webserver

import (
	"errors"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/zekroTJA/shinpuru/internal/commands"
	"github.com/zekroTJA/shinpuru/internal/core/config"
	"github.com/zekroTJA/shinpuru/internal/core/database"
	"github.com/zekroTJA/shinpuru/internal/core/storage"
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
	}
)

// WebServer exposes HTTP REST API endpoints to
// access shinpurus functionalities via a web app.
type WebServer struct {
	server *fasthttp.Server
	router *routing.Router

	db         database.Database
	st         storage.Storage
	rlm        *RateLimitManager
	auth       *Auth
	dcoauth    *discordoauth.DiscordOAuth
	session    *discordgo.Session
	cmdhandler *commands.CmdHandler

	config *config.Config
}

// New creates a new instance of WebServer consuming the passed
// database provider, storage provider, discordgo session, command
// handler and configuration.
func New(db database.Database, st storage.Storage, s *discordgo.Session,
	cmd *commands.CmdHandler, config *config.Config) (ws *WebServer) {

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
	ws.st = st
	ws.session = s
	ws.cmdhandler = cmd
	ws.rlm = NewRateLimitManager()
	ws.router = routing.New()
	ws.server = &fasthttp.Server{
		Handler: ws.router.HandleRequest,
	}

	ws.auth = NewAuth(db, s)

	ws.dcoauth = discordoauth.NewDiscordOAuth(
		config.Discord.ClientID,
		config.Discord.ClientSecret,
		config.WebServer.PublicAddr+endpointAuthCB,
		ws.auth.LoginFailedHandler,
		ws.auth.LoginSuccessHandler,
	)

	ws.registerHandlers()

	return
}

// ListenAndServeBlocking starts the listening and serving
// loop of the web server which blocks the current goroutine.
//
// If an error is returned, the startup failed with the
// specified error.
func (ws *WebServer) ListenAndServeBlocking() error {
	tls := ws.config.WebServer.TLS

	if tls != nil && tls.Enabled {
		if tls.Cert == "" || tls.Key == "" {
			return errors.New("cert file and key file must be specified")
		}
		return ws.server.ListenAndServeTLS(ws.config.WebServer.Addr, tls.Cert, tls.Key)
	}

	return ws.server.ListenAndServe(ws.config.WebServer.Addr)
}

// registerHandlers registers all request handler for the
// request URL specified match tree.
func (ws *WebServer) registerHandlers() {
	// --------------------------------
	// AVAILABLE WITHOUT AUTH

	ws.router.Use(ws.addHeaders, ws.optionsHandler, ws.handlerFiles)

	imagestore := ws.router.Group("/imagestore")
	imagestore.
		Get("/<id>", ws.handlerGetImage)

	ws.router.Get(endpointLogInWithDC, ws.dcoauth.HandlerInit)
	ws.router.Get(endpointAuthCB, ws.dcoauth.HandlerCallback)

	// --------------------------------
	// ONLY AVAILABLE AFTER AUTH

	ws.router.Use(ws.auth.checkAuth)

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
