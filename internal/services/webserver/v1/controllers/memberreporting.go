package controllers

import (
	"bytes"

	"github.com/bwmarrin/discordgo"
	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
	sharedmodels "github.com/zekroTJA/shinpuru/internal/models"
	"github.com/zekroTJA/shinpuru/internal/services/config"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/services/permissions"
	"github.com/zekroTJA/shinpuru/internal/services/report"
	"github.com/zekroTJA/shinpuru/internal/services/storage"
	"github.com/zekroTJA/shinpuru/internal/services/webserver/v1/models"
	"github.com/zekroTJA/shinpuru/internal/services/webserver/wsutil"
	"github.com/zekroTJA/shinpuru/internal/util/imgstore"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekrotja/dgrs"
)

type MemberReportingController struct {
	session *discordgo.Session
	cfg     config.Provider
	db      database.Database
	st      storage.Storage
	repSvc  *report.ReportService
	state   *dgrs.State
}

func (c *MemberReportingController) Setup(container di.Container, router fiber.Router) {
	c.session = container.Get(static.DiDiscordSession).(*discordgo.Session)
	c.cfg = container.Get(static.DiConfig).(config.Provider)
	c.db = container.Get(static.DiDatabase).(database.Database)
	c.st = container.Get(static.DiObjectStorage).(storage.Storage)
	c.repSvc = container.Get(static.DiReport).(*report.ReportService)
	c.state = container.Get(static.DiState).(*dgrs.State)

	pmw := container.Get(static.DiPermissions).(*permissions.Permissions)

	router.Post("/reports", pmw.HandleWs(c.session, "sp.guild.mod.report"), c.postReport)
	router.Post("/kick", pmw.HandleWs(c.session, "sp.guild.mod.kick"), c.postKick)
	router.Post("/ban", pmw.HandleWs(c.session, "sp.guild.mod.ban"), c.postBan)
	router.Post("/mute", pmw.HandleWs(c.session, "sp.guild.mod.mute"), c.postMute)
	router.Post("/unmute", pmw.HandleWs(c.session, "sp.guild.mod.mute"), c.postUnmute)
}

// @Summary Create A Member Report
// @Description Creates a member report.
// @Tags Member Reporting
// @Accept json
// @Produce json
// @Param id path string true "The ID of the guild."
// @Param memberid path string true "The ID of the victim member."
// @Param payload body models.ReportRequest true "The report payload."
// @Success 200 {object} models.Report
// @Failure 400 {object} models.Error
// @Failure 401 {object} models.Error
// @Failure 404 {object} models.Error
// @Router /guilds/{id}/{memberid}/reports [post]
func (c *MemberReportingController) postReport(ctx *fiber.Ctx) (err error) {
	uid := ctx.Locals("uid").(string)

	guildID := ctx.Params("guildid")
	memberID := ctx.Params("memberid")

	repReq := new(models.ReportRequest)
	if err := ctx.BodyParser(repReq); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if memberID == uid {
		return fiber.NewError(fiber.StatusBadRequest, "you can not report yourself")
	}

	if ok, err := repReq.Validate(false); !ok {
		return err
	}

	if repReq.Attachment != "" {
		img, err := imgstore.DownloadFromURL(repReq.Attachment)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		err = c.st.PutObject(static.StorageBucketImages, img.ID.String(),
			bytes.NewReader(img.Data), int64(img.Size), img.MimeType)
		if err != nil {
			return err
		}
		repReq.Attachment = img.ID.String()
	}

	rep, err := c.repSvc.PushReport(&sharedmodels.Report{
		GuildID:       guildID,
		ExecutorID:    uid,
		VictimID:      memberID,
		Msg:           repReq.Reason,
		AttachmentURL: repReq.Attachment,
		Type:          repReq.Type,
	})

	if err != nil {
		return err
	}

	return ctx.JSON(models.ReportFromReport(rep, c.cfg.Config().WebServer.PublicAddr))
}

// @Summary Create A Member Kick Report
// @Description Creates a member kick report.
// @Tags Member Reporting
// @Accept json
// @Produce json
// @Param id path string true "The ID of the guild."
// @Param memberid path string true "The ID of the victim member."
// @Param payload body models.ReasonRequest true "The report payload."
// @Success 200 {object} models.Report
// @Failure 400 {object} models.Error
// @Failure 401 {object} models.Error
// @Failure 404 {object} models.Error
// @Router /guilds/{id}/{memberid}/kick [post]
func (c *MemberReportingController) postKick(ctx *fiber.Ctx) (err error) {
	uid := ctx.Locals("uid").(string)

	guildID := ctx.Params("guildid")
	memberID := ctx.Params("memberid")

	req := new(models.ReasonRequest)
	if err := ctx.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if memberID == uid {
		return fiber.NewError(fiber.StatusBadRequest, "you can not kick yourself")
	}

	if ok, err := req.Validate(false); !ok {
		return err
	}

	if req.Attachment != "" {
		img, err := imgstore.DownloadFromURL(req.Attachment)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		err = c.st.PutObject(static.StorageBucketImages, img.ID.String(),
			bytes.NewReader(img.Data), int64(img.Size), img.MimeType)
		if err != nil {
			return err
		}
		req.Attachment = img.ID.String()
	}

	rep, err := c.repSvc.PushKick(&sharedmodels.Report{
		GuildID:       guildID,
		ExecutorID:    uid,
		VictimID:      memberID,
		Msg:           req.Reason,
		AttachmentURL: req.Attachment,
	})

	if err == report.ErrRoleDiff {
		return fiber.NewError(fiber.StatusBadRequest, "you can not kick members with higher or same permissions than/as yours")
	}

	if err != nil {
		return err
	}

	return ctx.JSON(models.ReportFromReport(rep, c.cfg.Config().WebServer.PublicAddr))
}

// @Summary Create A Member Ban Report
// @Description Creates a member ban report.
// @Tags Member Reporting
// @Accept json
// @Produce json
// @Param id path string true "The ID of the guild."
// @Param memberid path string true "The ID of the victim member."
// @Param payload body models.ReasonRequest true "The report payload."
// @Success 200 {object} models.Report
// @Failure 400 {object} models.Error
// @Failure 401 {object} models.Error
// @Failure 404 {object} models.Error
// @Router /guilds/{id}/{memberid}/ban [post]
func (c *MemberReportingController) postBan(ctx *fiber.Ctx) (err error) {
	uid := ctx.Locals("uid").(string)

	guildID := ctx.Params("guildid")
	memberID := ctx.Params("memberid")

	anonymous, err := wsutil.GetQueryBool(ctx, "anonymous", false)
	if err != nil {
		return
	}

	req := new(models.ReasonRequest)
	if err := ctx.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if memberID == uid {
		return fiber.NewError(fiber.StatusBadRequest, "you can not ban yourself")
	}

	if ok, err := req.Validate(false); !ok {
		return err
	}

	if req.Attachment != "" {
		img, err := imgstore.DownloadFromURL(req.Attachment)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		err = c.st.PutObject(static.StorageBucketImages, img.ID.String(),
			bytes.NewReader(img.Data), int64(img.Size), img.MimeType)
		if err != nil {
			return err
		}
		req.Attachment = img.ID.String()
	}

	rep, err := c.repSvc.PushBan(&sharedmodels.Report{
		GuildID:       guildID,
		ExecutorID:    uid,
		VictimID:      memberID,
		Msg:           req.Reason,
		AttachmentURL: req.Attachment,
		Timeout:       req.Timeout,
		Anonymous:     anonymous,
	})

	if err == report.ErrRoleDiff {
		return fiber.NewError(fiber.StatusBadRequest, "you can not ban members with higher or same permissions than/as yours")
	}

	if err != nil {
		return err
	}

	return ctx.JSON(models.ReportFromReport(rep, c.cfg.Config().WebServer.PublicAddr))
}

// @Summary Create A Member Mute Report
// @Description Creates a member mute report.
// @Tags Member Reporting
// @Accept json
// @Produce json
// @Param id path string true "The ID of the guild."
// @Param memberid path string true "The ID of the victim member."
// @Param payload body models.ReasonRequest true "The report payload."
// @Success 200 {object} models.Report
// @Failure 400 {object} models.Error
// @Failure 401 {object} models.Error
// @Failure 404 {object} models.Error
// @Router /guilds/{id}/{memberid}/mute [post]
func (c *MemberReportingController) postMute(ctx *fiber.Ctx) (err error) {
	uid := ctx.Locals("uid").(string)

	guildID := ctx.Params("guildid")
	memberID := ctx.Params("memberid")

	req := new(models.ReasonRequest)
	if err := ctx.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if memberID == uid {
		return fiber.NewError(fiber.StatusBadRequest, "you can not mute yourself")
	}

	if ok, err := req.Validate(true); !ok {
		return err
	}

	if req.Attachment != "" {
		img, err := imgstore.DownloadFromURL(req.Attachment)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		err = c.st.PutObject(static.StorageBucketImages, img.ID.String(),
			bytes.NewReader(img.Data), int64(img.Size), img.MimeType)
		if err != nil {
			return err
		}
		req.Attachment = img.ID.String()
	}

	if req.Timeout == nil {
		return fiber.NewError(fiber.StatusBadRequest, "You must pass a valid mute timeout duration")
	}

	rep, err := c.repSvc.PushMute(&sharedmodels.Report{
		GuildID:       guildID,
		ExecutorID:    uid,
		VictimID:      memberID,
		Msg:           req.Reason,
		AttachmentURL: req.Attachment,
		Timeout:       req.Timeout,
	})

	if err == report.ErrRoleDiff {
		return fiber.NewError(fiber.StatusBadRequest, "you can not mute members with higher or same permissions than/as yours")
	}

	if err != nil {
		return err
	}

	return ctx.JSON(models.ReportFromReport(rep, c.cfg.Config().WebServer.PublicAddr))
}

// @Summary Unmute A Member
// @Description Unmute a muted member.
// @Tags Member Reporting
// @Accept json
// @Produce json
// @Param id path string true "The ID of the guild."
// @Param memberid path string true "The ID of the victim member."
// @Param payload body models.ReasonRequest true "The unmute payload."
// @Success 200 {object} models.Status
// @Failure 400 {object} models.Error
// @Failure 401 {object} models.Error
// @Failure 404 {object} models.Error
// @Router /guilds/{id}/{memberid}/mute [post]
func (c *MemberReportingController) postUnmute(ctx *fiber.Ctx) (err error) {
	uid := ctx.Locals("uid").(string)

	guildID := ctx.Params("guildid")
	memberID := ctx.Params("memberid")

	req := new(models.ReasonRequest)
	if err := ctx.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if memberID == uid {
		return fiber.NewError(fiber.StatusBadRequest, "you can not mute yourself")
	}

	_, err = c.repSvc.RevokeMute(
		guildID,
		uid,
		memberID,
		req.Reason)

	if err == report.ErrRoleDiff {
		return fiber.NewError(fiber.StatusBadRequest, "you can not unmute members with higher or same permissions than/as yours")
	}

	if err != nil {
		return err
	}

	return ctx.JSON(models.Ok)
}
