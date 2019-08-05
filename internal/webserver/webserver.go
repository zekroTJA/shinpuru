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

	config *core.ConfigWS
}

func NewWebServer(db core.Database, s *discordgo.Session, cmd *commands.CmdHandler, config *core.ConfigWS, clientID, clientSecret string) (ws *WebServer) {
	ws = new(WebServer)

	if !strings.HasPrefix(config.PublicAddr, "http") {
		protocol := "http"
		if config.TLS != nil && config.TLS.Enabled {
			protocol += "s"
		}
		config.PublicAddr = fmt.Sprintf("%s://%s", protocol, config.PublicAddr)
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
		config.PublicAddr+endpointAuthCB,
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

	ws.router.Use(ws.addHeaders, ws.auth.checkAuth, ws.handlerFiles)

	ws.router.Get(endpointLogInWithDC, ws.dcoauth.HandlerInit)
	ws.router.Get(endpointAuthCB, ws.dcoauth.HandlerCallback)

	api := ws.router.Group("/api")
	api.
		Get("/me", ws.handlerGetMe)

	api.
		Post("/logout", ws.auth.LogOutHandler)

	guilds := api.Group("/guilds")
	guilds.
		Get("", ws.handlerGuildsGet)
	guilds.
		Get("/<id:[0-9]+>", ws.handlerGuildsGetGuild)
	guilds.
		Get("/<guildid:[0-9]+>/reports", ws.handlerGetReports)

	member := guilds.Group("/<guildid:[0-9]+>/<memberid:[0-9]+>")
	member.
		Get("", ws.handlerGuildsGetMember)
	member.
		Get("/reports", ws.handlerGetReports)
	member.
		Get("/permissions", ws.handlerGetPermissions)
	member.
		Get("/permissions/allowed", ws.handlerGetPermissionsAllowed)

	reports := api.Group("/reports")
	reports.
		Get("/<id:[0-9]+>", ws.handlerGetReport)
}

func (ws *WebServer) ListenAndServeBlocking() error {
	tls := ws.config.TLS

	if tls.Enabled {
		if tls.Cert == "" || tls.Key == "" {
			return errors.New("cert file and key file must be specified")
		}
		return ws.server.ListenAndServeTLS(ws.config.Addr, tls.Cert, tls.Key)
	}

	return ws.server.ListenAndServe(ws.config.Addr)
}
