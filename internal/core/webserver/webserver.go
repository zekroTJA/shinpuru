package webserver

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/core/config"
	v1 "github.com/zekroTJA/shinpuru/internal/core/webserver/v1"
	"github.com/zekroTJA/shinpuru/internal/core/webserver/v1/controllers"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/static"
)

type WebServer struct {
	app       *fiber.App
	cfg       *config.Config
	container di.Container
}

type errorModel struct {
	Error   string `json:"error"`
	Code    int    `json:"code"`
	Context string `json:"context,omitempty"`
}

func New(container di.Container) (ws *WebServer, err error) {
	ws = new(WebServer)

	ws.container = container
	ws.cfg = container.Get(static.DiConfig).(*config.Config)

	ws.app = fiber.New(fiber.Config{
		ErrorHandler:          ws.errorHandler,
		ServerHeader:          fmt.Sprintf("shinpuru v%s", util.AppVersion),
		DisableStartupMessage: util.IsRelease(),
		ProxyHeader:           "X-Forwarded-For",
	})

	if !util.IsRelease() {
		ws.app.Use(cors.New(cors.Config{
			AllowOrigins:     ws.cfg.WebServer.DebugPublicAddr,
			AllowHeaders:     "authorization, content-type, set-cookie, cookie, server",
			AllowMethods:     "GET, POST, PUT, PATCH, POST, DELETE, OPTIONS",
			AllowCredentials: true,
		}))
		ws.app.Use(logger.New())
	}

	if rl := ws.cfg.WebServer.RateLimit; rl != nil && rl.Enabled {
		ws.app.Use(limiter.New(limiter.Config{
			Max:        rl.Max,
			Expiration: time.Duration(rl.DurationSeconds) * time.Second,
			LimitReached: func(c *fiber.Ctx) error {
				return fiber.ErrTooManyRequests
			},
		}))
	}

	new(controllers.ImagestoreController).Setup(ws.container, ws.app.Group("/imagestore"))
	ws.registerRouter(new(v1.Router), "/api/v1", "/api")

	ws.app.Use(filesystem.New(filesystem.Config{
		Root:         http.Dir("web/dist/web"),
		Browse:       true,
		Index:        "index.html",
		MaxAge:       3600,
		NotFoundFile: "index.html",
	}))

	return
}

func (ws *WebServer) ListenAndServeBlocking() error {
	tls := ws.cfg.WebServer.TLS

	if tls != nil && tls.Enabled {
		if tls.Cert == "" || tls.Key == "" {
			return errors.New("cert file and key file must be specified")
		}
		return ws.app.ListenTLS(ws.cfg.WebServer.Addr, tls.Cert, tls.Key)
	}

	return ws.app.Listen(ws.cfg.WebServer.Addr)
}

func (ws *WebServer) registerRouter(router Router, route ...string) {
	router.SetContainer(ws.container)
	for _, r := range route {
		router.Route(ws.app.Group(r))
	}
}

func (ws *WebServer) errorHandler(ctx *fiber.Ctx, err error) error {
	if fErr, ok := err.(*fiber.Error); ok {
		ctx.Status(fErr.Code)
		return ctx.JSON(&errorModel{
			Error: fErr.Message,
			Code:  fErr.Code,
		})
	}
	return fiber.DefaultErrorHandler(ctx, err)
}
