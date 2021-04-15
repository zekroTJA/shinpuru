package webserver

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/core/config"
	v1 "github.com/zekroTJA/shinpuru/internal/core/webserver/v1"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/static"
)

type WebServer struct {
	app       *fiber.App
	cfg       *config.Config
	container di.Container
}

func New(container di.Container) (ws *WebServer, err error) {
	ws = new(WebServer)

	ws.container = container
	ws.cfg = container.Get(static.DiConfig).(*config.Config)

	ws.app = fiber.New(fiber.Config{
		ErrorHandler:          ws.errorHandler,
		DisableStartupMessage: util.IsRelease(),
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
	if fErr, ok := err.(*fiber.Error); (ok && fErr.Code >= 500) || !ok {
		// custom error handling here
	}
	return fiber.DefaultErrorHandler(ctx, err)
}
