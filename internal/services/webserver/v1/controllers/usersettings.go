package controllers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/services/permissions"
	"github.com/zekroTJA/shinpuru/internal/services/webserver/v1/models"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekrotja/dgrs"
)

type UsersettingsController struct {
	session *discordgo.Session
	db      database.Database
	state   *dgrs.State
	pmw     *permissions.Permissions
}

func (c *UsersettingsController) Setup(container di.Container, router fiber.Router) {
	c.session = container.Get(static.DiDiscordSession).(*discordgo.Session)
	c.db = container.Get(static.DiDatabase).(database.Database)
	c.state = container.Get(static.DiState).(*dgrs.State)
	c.pmw = container.Get(static.DiPermissions).(*permissions.Permissions)

	router.Get("/ota", c.getOTA)
	router.Post("/ota", c.postOTA)
	router.Post("/flush", c.postFlush)
}

// @Summary Get OTA Usersettings State
// @Description Returns the current state of the OTA user setting.
// @Tags User Settings
// @Accept json
// @Produce json
// @Success 200 {object} models.UsersettingsOTA
// @Failure 401 {object} models.Error
// @Failure 404 {object} models.Error
// @Router /usersettings/ota [get]
func (c *UsersettingsController) getOTA(ctx *fiber.Ctx) error {
	uid := ctx.Locals("uid").(string)

	enabled, err := c.db.GetUserOTAEnabled(uid)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return err
	}

	return ctx.JSON(&models.UsersettingsOTA{Enabled: enabled})
}

// @Summary Update OTA Usersettings State
// @Description Update the OTA user settings state.
// @Tags User Settings
// @Accept json
// @Produce json
// @Param payload body models.UsersettingsOTA true "The OTA settings payload."
// @Success 200 {object} models.UsersettingsOTA
// @Failure 400 {object} models.Error
// @Failure 401 {object} models.Error
// @Failure 404 {object} models.Error
// @Router /usersettings/ota [post]
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

// @Summary FLush all user data
// @Description Flush all user data.
// @Tags User Settings
// @Accept json
// @Produce json
// @Success 200 {object} models.UsersettingsOTA
// @Failure 400 {object} models.Error
// @Failure 401 {object} models.Error
// @Router /usersettings/flush [post]
func (c *UsersettingsController) postFlush(ctx *fiber.Ctx) error {
	uid := ctx.Locals("uid").(string)

	res, err := util.FlushAllUserData(c.db, c.state, uid)
	if err != nil {
		return err
	}

	return ctx.JSON(res)
}
