package controllers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/bwmarrin/snowflake"
	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/config"
	"github.com/zekroTJA/shinpuru/internal/middleware"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/services/report"
	"github.com/zekroTJA/shinpuru/internal/services/webserver/v1/models"
	"github.com/zekroTJA/shinpuru/internal/util/static"
)

type ReportsController struct {
	session *discordgo.Session
	cfg     *config.Config
	db      database.Database
	repSvc  *report.ReportService
}

func (c *ReportsController) Setup(container di.Container, router fiber.Router) {
	c.session = container.Get(static.DiDiscordSession).(*discordgo.Session)
	c.cfg = container.Get(static.DiConfig).(*config.Config)
	c.db = container.Get(static.DiDatabase).(database.Database)
	c.repSvc = container.Get(static.DiReport).(*report.ReportService)

	pmw := container.Get(static.DiPermissionMiddleware).(*middleware.PermissionsMiddleware)

	router.Get("/:id", c.getReport)
	router.Post("/:id/revoke", pmw.HandleWs(c.session, "sp.guild.mod.report"), c.postRevoke)
}

// @Summary Get Report
// @Description Returns a single report object by its ID.
// @Accept json
// @Produce json
// @Param id path string true "The report ID."
// @Success 200 {object} models.Report
// @Failure 400 {object} models.Error
// @Failure 401 {object} models.Error
// @Failure 404 {object} models.Error
// @Router /reports/:id [get]
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

	return ctx.JSON(models.ReportFromReport(rep, c.cfg.WebServer.PublicAddr))
}

// @Summary Revoke Report
// @Description Revokes a given report by ID.
// @Accept json
// @Produce json
// @Param id path string true "The report ID."
// @Param payload body models.ReasonRequest true "The revoke reason payload."
// @Success 200 {object} models.Report
// @Failure 400 {object} models.Error
// @Failure 401 {object} models.Error
// @Failure 404 {object} models.Error
// @Router /reports/:id/revoke [post]
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

	var reason models.ReasonRequest
	if err := ctx.BodyParser(&reason); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	_, err = c.repSvc.RevokeReport(
		rep,
		uid,
		reason.Reason,
		c.cfg.WebServer.Addr,
		c.db,
		c.session)

	if err != nil {
		return err
	}

	return ctx.JSON(models.Ok)
}
