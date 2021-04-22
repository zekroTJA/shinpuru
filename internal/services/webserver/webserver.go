package webserver

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/config"
	v1 "github.com/zekroTJA/shinpuru/internal/services/webserver/v1"
	"github.com/zekroTJA/shinpuru/internal/services/webserver/v1/controllers"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/limiter"
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

	rlc := ws.cfg.WebServer.RateLimit
	rlh := limiter.New(limiter.Config{
		Next: func(ctx *fiber.Ctx) bool {
			return rlc == nil || !rlc.Enabled
		},
		Burst:           rlc.Burst,
		Duration:        time.Duration(rlc.LimitSeconds) * time.Second,
		CleanupInterval: 10 * time.Minute,
		KeyGenerator: func(ctx *fiber.Ctx) string {
			return ctx.IP()
		},
		OnLimitReached: func(ctx *fiber.Ctx) error {
			return fiber.ErrTooManyRequests
		},
	})

	new(controllers.ImagestoreController).Setup(ws.container, ws.app.Group("/imagestore"))
	ws.registerRouter(new(v1.Router), []string{"/api/v1", "/api"}, rlh)

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

func (ws *WebServer) registerRouter(router Router, routes []string, middlewares ...fiber.Handler) {
	router.SetContainer(ws.container)
	for _, r := range routes {
		router.Route(ws.app.Group(r, middlewares...))
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
