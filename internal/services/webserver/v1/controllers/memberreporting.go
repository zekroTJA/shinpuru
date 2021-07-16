package controllers

import (
	"bytes"

	"github.com/bwmarrin/discordgo"
	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/config"
	"github.com/zekroTJA/shinpuru/internal/middleware"
	sharedmodels "github.com/zekroTJA/shinpuru/internal/models"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/services/report"
	"github.com/zekroTJA/shinpuru/internal/services/storage"
	"github.com/zekroTJA/shinpuru/internal/services/webserver/v1/models"
	"github.com/zekroTJA/shinpuru/internal/services/webserver/wsutil"
	"github.com/zekroTJA/shinpuru/internal/util/imgstore"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/roleutil"
	"github.com/zekrotja/dgrs"
)

type MemberReportingController struct {
	session *discordgo.Session
	cfg     *config.Config
	db      database.Database
	st      storage.Storage
	repSvc  *report.ReportService
	state   *dgrs.State
}

func (c *MemberReportingController) Setup(container di.Container, router fiber.Router) {
	c.session = container.Get(static.DiDiscordSession).(*discordgo.Session)
	c.cfg = container.Get(static.DiConfig).(*config.Config)
	c.db = container.Get(static.DiDatabase).(database.Database)
	c.st = container.Get(static.DiObjectStorage).(storage.Storage)
	c.repSvc = container.Get(static.DiReport).(*report.ReportService)
	c.state = container.Get(static.DiState).(*dgrs.State)

	pmw := container.Get(static.DiPermissionMiddleware).(*middleware.PermissionsMiddleware)

	router.Post("/reports", pmw.HandleWs(c.session, "sp.guild.mod.report"), c.postReport)
	router.Post("/kick", pmw.HandleWs(c.session, "sp.guild.mod.kick"), c.postKick)
	router.Post("/ban", pmw.HandleWs(c.session, "sp.guild.mod.ban"), c.postBan)
	router.Post("/mute", pmw.HandleWs(c.session, "sp.guild.mod.mute"), c.postMute)
	router.Post("/unmute", pmw.HandleWs(c.session, "sp.guild.mod.mute"), c.postUnmute)
}

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
		AttachmehtURL: repReq.Attachment,
		Type:          repReq.Type,
	})

	if err != nil {
		return err
	}

	return ctx.JSON(models.ReportFromReport(rep, c.cfg.WebServer.PublicAddr))
}

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

	guild, err := c.state.Guild(guildID)
	if err != nil {
		return err
	}

	executor, err := c.session.GuildMember(guildID, uid)
	if err != nil {
		return err
	}

	victim, err := c.session.GuildMember(guildID, memberID)
	if err != nil {
		return err
	}

	if roleutil.PositionDiff(victim, executor, guild) >= 0 {
		return fiber.NewError(fiber.StatusBadRequest, "you can not kick members with higher or same permissions than/as yours")
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
		AttachmehtURL: req.Attachment,
	})

	if err != nil {
		return err
	}

	return ctx.JSON(models.ReportFromReport(rep, c.cfg.WebServer.PublicAddr))
}

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

	guild, err := c.state.Guild(guildID)
	if err != nil {
		return err
	}

	executor, err := c.session.GuildMember(guildID, uid)
	if err != nil {
		return err
	}

	var victim *discordgo.Member
	if !anonymous {
		victim, err = c.session.GuildMember(guildID, memberID)
		if err != nil {
			return err
		}
	}

	if !anonymous && roleutil.PositionDiff(victim, executor, guild) >= 0 {
		return fiber.NewError(fiber.StatusBadRequest, "you can not ban members with higher or same permissions than/as yours")
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
		AttachmehtURL: req.Attachment,
	})

	if err != nil {
		return err
	}

	return ctx.JSON(models.ReportFromReport(rep, c.cfg.WebServer.PublicAddr))
}

func (c *MemberReportingController) postMute(ctx *fiber.Ctx) (err error) {
	uid := ctx.Locals("uid").(string)

	guildID := ctx.Params("guildid")
	memberID := ctx.Params("memberid")

	req := new(models.ReasonRequest)
	if err := ctx.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	muteRoleID, err := c.db.GetGuildMuteRole(guildID)
	if database.IsErrDatabaseNotFound(err) {
		return fiber.NewError(fiber.StatusBadRequest, "mute role is not set up on this guild")
	} else if err != nil {
		return err
	}

	if memberID == uid {
		return fiber.NewError(fiber.StatusBadRequest, "you can not mute yourself")
	}

	guild, err := c.state.Guild(guildID)
	if err != nil {
		return err
	}

	executor, err := c.session.GuildMember(guildID, uid)
	if err != nil {
		return err
	}

	victim, err := c.session.GuildMember(guildID, memberID)
	if err != nil {
		return err
	}

	if roleutil.PositionDiff(victim, executor, guild) >= 0 {
		return fiber.NewError(fiber.StatusBadRequest, "you can not mute members with higher or same permissions than/as yours")
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

	rep, err := c.repSvc.PushMute(&sharedmodels.Report{
		GuildID:       guildID,
		ExecutorID:    uid,
		VictimID:      memberID,
		Msg:           req.Reason,
		AttachmehtURL: req.Attachment,
	}, muteRoleID)

	if err != nil {
		return err
	}

	return ctx.JSON(models.ReportFromReport(rep, c.cfg.WebServer.PublicAddr))
}

func (c *MemberReportingController) postUnmute(ctx *fiber.Ctx) (err error) {
	uid := ctx.Locals("uid").(string)

	guildID := ctx.Params("guildid")
	memberID := ctx.Params("memberid")

	req := new(models.ReasonRequest)
	if err := ctx.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	muteRoleID, err := c.db.GetGuildMuteRole(guildID)
	if database.IsErrDatabaseNotFound(err) {
		return fiber.NewError(fiber.StatusBadRequest, "mute role is not set up on this guild")
	} else if err != nil {
		return err
	}

	if memberID == uid {
		return fiber.NewError(fiber.StatusBadRequest, "you can not mute yourself")
	}

	guild, err := c.state.Guild(guildID)
	if err != nil {
		return err
	}

	executor, err := c.session.GuildMember(guildID, uid)
	if err != nil {
		return err
	}

	victim, err := c.session.GuildMember(guildID, memberID)
	if err != nil {
		return err
	}

	if roleutil.PositionDiff(victim, executor, guild) >= 0 {
		return fiber.NewError(fiber.StatusBadRequest, "you can not mute members with higher or same permissions than/as yours")
	}

	_, err = c.repSvc.RevokeMute(
		guildID,
		uid,
		memberID,
		req.Reason,
		muteRoleID)

	if err != nil {
		return err
	}

	return ctx.JSON(models.Ok)
}
