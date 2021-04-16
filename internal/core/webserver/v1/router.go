package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/core/webserver/auth"
	"github.com/zekroTJA/shinpuru/internal/core/webserver/v1/controllers"
	"github.com/zekroTJA/shinpuru/internal/util/static"
)

type Router struct {
	container di.Container
}

func (r *Router) SetContainer(container di.Container) {
	r.container = container
}

func (r *Router) Route(router fiber.Router) {
	authMw := r.container.Get(static.DiAuthMiddleware).(auth.Middleware)

	new(controllers.EtcController).Setup(r.container, router)
	new(controllers.UtilController).Setup(r.container, router.Group("/util"))
	new(controllers.AuthController).Setup(r.container, router.Group("/auth"))
	new(controllers.OTAController).Setup(r.container, router.Group("/ota"))

	// --- REQUIRES ACCESS TOKEN AUTH ---

	new(controllers.TokenController).Setup(r.container, router.Group("/token", authMw.Handle))
	new(controllers.ReportsController).Setup(r.container, router.Group("/reports", authMw.Handle))
	new(controllers.GuildsController).Setup(r.container, router.Group("/guilds", authMw.Handle))
	new(controllers.UnbanrequestsController).Setup(r.container, router.Group("/unbanrequests", authMw.Handle))
	new(controllers.MemberReportingController).Setup(r.container, router.Group("/guilds/:guildid/:memberid", authMw.Handle))
	new(controllers.GuildBackupsController).Setup(r.container, router.Group("/guilds/:guildid/backups", authMw.Handle))
	new(controllers.GuildMembersController).Setup(r.container, router.Group("/guilds/:guildid", authMw.Handle))
}
