package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/config"
	"github.com/zekroTJA/shinpuru/internal/middleware"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/services/webserver/v1/models"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekrotja/discordgo"
)

type UsersettingsController struct {
	session *discordgo.Session
	cfg     *config.Config
	db      database.Database
	pmw     *middleware.PermissionsMiddleware
}

func (c *UsersettingsController) Setup(container di.Container, router fiber.Router) {
	c.session = container.Get(static.DiDiscordSession).(*discordgo.Session)
	c.cfg = container.Get(static.DiConfig).(*config.Config)
	c.db = container.Get(static.DiDatabase).(database.Database)
	c.pmw = container.Get(static.DiPermissionMiddleware).(*middleware.PermissionsMiddleware)

	router.Get("/ota", c.getOTA)
	router.Post("/ota", c.postOTA)
}

func (c *UsersettingsController) getOTA(ctx *fiber.Ctx) error {
	uid := ctx.Locals("uid").(string)

	enabled, err := c.db.GetUserOTAEnabled(uid)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return err
	}

	return ctx.JSON(&models.UsersettingsOTA{Enabled: enabled})
}

func (c *UsersettingsController) postOTA(ctx *fiber.Ctx) error {
	uid := ctx.Locals("uid").(string)

	var err error

	data := new(models.UsersettingsOTA)
	if err = ctx.BodyParser(data); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err = c.db.SetUserOTAEnabled(uid, data.Enabled); err != nil {
		return err
	}

	return ctx.JSON(data)
}
