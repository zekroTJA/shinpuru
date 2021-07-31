package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/services/webserver/auth"
	"github.com/zekroTJA/shinpuru/internal/services/webserver/v1/controllers"
	"github.com/zekroTJA/shinpuru/internal/util/static"
)

type Router struct {
	container di.Container
}

func (r *Router) SetContainer(container di.Container) {
	r.container = container
}

// @title shinpuru main API
// @version 1.0
// @description The shinpuru main REST API.
// @BasePath /api/v1
func (r *Router) Route(router fiber.Router) {
	authMw := r.container.Get(static.DiAuthMiddleware).(auth.Middleware)

	new(controllers.EtcController).Setup(r.container, router)
	new(controllers.UtilController).Setup(r.container, router.Group("/util"))
	new(controllers.AuthController).Setup(r.container, router.Group("/auth"))
	new(controllers.OTAController).Setup(r.container, router.Group("/ota"))

	router.Get("/stack", func(ctx *fiber.Ctx) error { return ctx.JSON(ctx.App().Stack()) })

	// --- REQUIRES ACCESS TOKEN AUTH ---

	router.Use(authMw.Handle)

	new(controllers.TokenController).Setup(r.container, router.Group("/token"))
	new(controllers.GlobalSettingsController).Setup(r.container, router.Group("/settings"))
	new(controllers.ReportsController).Setup(r.container, router.Group("/reports"))
	new(controllers.GuildsController).Setup(r.container, router.Group("/guilds"))
	new(controllers.GuildBackupsController).Setup(r.container, router.Group("/guilds/:guildid/backups"))
	new(controllers.UnbanrequestsController).Setup(r.container, router.Group("/unbanrequests"))
	new(controllers.UsersettingsController).Setup(r.container, router.Group("/usersettings"))
	new(controllers.MemberReportingController).Setup(r.container, router.Group("/guilds/:guildid/:memberid"))
	new(controllers.GuildMembersController).Setup(r.container, router.Group("/guilds/:guildid"))
}
