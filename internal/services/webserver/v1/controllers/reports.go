package controllers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/bwmarrin/snowflake"
	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/services/config"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/services/permissions"
	"github.com/zekroTJA/shinpuru/internal/services/report"
	"github.com/zekroTJA/shinpuru/internal/services/webserver/v1/models"
	"github.com/zekroTJA/shinpuru/internal/util/static"
)

type ReportsController struct {
	session *discordgo.Session
	cfg     config.Provider
	db      database.Database
	repSvc  *report.ReportService
	pmw     *permissions.Permissions
}

func (c *ReportsController) Setup(container di.Container, router fiber.Router) {
	c.session = container.Get(static.DiDiscordSession).(*discordgo.Session)
	c.cfg = container.Get(static.DiConfig).(config.Provider)
	c.db = container.Get(static.DiDatabase).(database.Database)
	c.repSvc = container.Get(static.DiReport).(*report.ReportService)
	c.pmw = container.Get(static.DiPermissions).(*permissions.Permissions)

	router.Get("/:id", c.getReport)
	router.Post("/:id/revoke", c.postRevoke)
}

// @Summary Get Report
// @Description Returns a single report object by its ID.
// @Tags Reports
// @Accept json
// @Produce json
// @Param id path string true "The report ID."
// @Success 200 {object} models.Report
// @Failure 400 {object} models.Error
// @Failure 401 {object} models.Error
// @Failure 404 {object} models.Error
// @Router /reports/{id} [get]
func (c *ReportsController) getReport(ctx *fiber.Ctx) (err error) {
	_id := ctx.Params("id")

	id, err := snowflake.ParseString(_id)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	rep, err := c.db.GetReport(id)
	if database.IsErrDatabaseNotFound(err) {
		return fiber.ErrNotFound
	}
	if err != nil {
		return err
	}

	return ctx.JSON(models.ReportFromReport(rep, c.cfg.Config().WebServer.PublicAddr))
}

// @Summary Revoke Report
// @Description Revokes a given report by ID.
// @Tags Reports
// @Accept json
// @Produce json
// @Param id path string true "The report ID."
// @Param payload body models.ReasonRequest true "The revoke reason payload."
// @Success 200 {object} models.Report
// @Failure 400 {object} models.Error
// @Failure 401 {object} models.Error
// @Failure 404 {object} models.Error
// @Router /reports/{id}/revoke [post]
func (c *ReportsController) postRevoke(ctx *fiber.Ctx) (err error) {
	uid := ctx.Locals("uid").(string)

	_id := ctx.Params("id")

	id, err := snowflake.ParseString(_id)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	rep, err := c.db.GetReport(id)
	if database.IsErrDatabaseNotFound(err) {
		return fiber.ErrNotFound
	}
	if err != nil {
		return err
	}

	ok, _, err := c.pmw.CheckPermissions(c.session, rep.GuildID, uid, "sp.guild.mod.report.revoke")
	if err != nil {
		return err
	}
	if !ok {
		return fiber.ErrForbidden
	}

	var reason models.ReasonRequest
	if err := ctx.BodyParser(&reason); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	_, err = c.repSvc.RevokeReport(
		rep,
		uid,
		reason.Reason,
		c.cfg.Config().WebServer.Addr,
	)

	if err != nil {
		return err
	}

	return ctx.JSON(models.Ok)
}
