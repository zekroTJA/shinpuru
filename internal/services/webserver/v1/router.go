package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/services/webserver/auth"
	"github.com/zekroTJA/shinpuru/internal/services/webserver/v1/controllers"
	"github.com/zekroTJA/shinpuru/internal/util/embedded"
	"github.com/zekroTJA/shinpuru/internal/util/static"
)

type Router struct {
	container di.Container
}

func (r *Router) SetContainer(container di.Container) {
	r.container = container
}

// Route registrar.
//
// @title shinpuru main API
// @version 1.0
// @description The shinpuru main REST API.
//
// @Tag.Name Etc
// @tag.Description General root API functionalities.
//
// @Tag.Name Utilities
// @tag.Description General utility functionalities.
//
// @Tag.Name Authorization
// @tag.Description Authorization endpoints.
//
// @Tag.Name OTA
// @tag.Description One Time Auth token endpoints.
//
// @Tag.Name Public
// @tag.Description Public API endpoints.
//
// @Tag.Name Search
// @tag.Description Search endpoints.
//
// @Tag.Name Tokens
// @tag.Description API token endpoints.
//
// @Tag.Name Global Settings
// @tag.Description Global bot settings endpoints.
//
// @Tag.Name Reports
// @tag.Description General reports endpoints.
//
// @Tag.Name Guilds
// @tag.Description Guild specific endpoints.
//
// @Tag.Name Guild Settings
// @Tag.Description Guild specific settings endpoints.
//
// @Tag.Name Guild Backups
// @tag.Description Guild backup endpoints.
//
// @Tag.Name Unban Requests
// @tag.Description Unban requests endpoints.
//
// @Tag.Name User Settings
// @tag.Description User specific settings endpoints.
//
// @Tag.Name Member Reporting
// @tag.Description Member reporting endpoints.
//
// @Tag.Name Members
// @tag.Description Members specific endpoints.
//
// @Tag.Name Channels
// @tag.Description Channels specific endpoints.
//
// @Tag.Name Verification
// @tag.Description User verification endpoints.
//
// @BasePath /api/v1
func (r *Router) Route(router fiber.Router) {
	authMw := r.container.Get(static.DiAuthMiddleware).(auth.Middleware)

	if !embedded.IsRelease() {
		new(controllers.DebugController).Setup(r.container, router.Group("debug"))
	}

	new(controllers.EtcController).Setup(r.container, router)
	new(controllers.UtilController).Setup(r.container, router.Group("/util"))
	new(controllers.AuthController).Setup(r.container, router.Group("/auth"))
	new(controllers.OTAController).Setup(r.container, router.Group("/ota"))
	new(controllers.PublicController).Setup(r.container, router.Group("/public"))

	router.Get("/stack", func(ctx *fiber.Ctx) error { return ctx.JSON(ctx.App().Stack()) })

	// --- REQUIRES ACCESS TOKEN AUTH ---

	router.Use(authMw.Handle)

	new(controllers.SearchController).Setup(r.container, router.Group("/search"))
	new(controllers.TokenController).Setup(r.container, router.Group("/token"))
	new(controllers.GlobalSettingsController).Setup(r.container, router.Group("/settings"))
	new(controllers.ReportsController).Setup(r.container, router.Group("/reports"))
	new(controllers.GuildsController).Setup(r.container, router.Group("/guilds"))
	new(controllers.MemberReportingController).Setup(r.container, router.Group("/guilds/:guildid/:memberid"))
	new(controllers.GuildBackupsController).Setup(r.container, router.Group("/guilds/:guildid/backups"))
	new(controllers.GuildsSettingsController).Setup(r.container, router.Group("/guilds/:guildid/settings"))
	new(controllers.GuildMembersController).Setup(r.container, router.Group("/guilds/:guildid"))
	new(controllers.ChannelController).Setup(r.container, router.Group("/channels/:guildid"))
	new(controllers.UsersController).Setup(r.container, router.Group("/users"))
	new(controllers.UsersettingsController).Setup(r.container, router.Group("/usersettings"))
	new(controllers.UnbanrequestsController).Setup(r.container, router.Group("/unbanrequests"))
	new(controllers.VerificationController).Setup(r.container, router.Group("/verification"))
}
