package controllers

import (
	"bytes"

	"github.com/bwmarrin/discordgo"
	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/config"
	"github.com/zekroTJA/shinpuru/internal/middleware"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/services/storage"
	"github.com/zekroTJA/shinpuru/internal/services/webserver/v1/models"
	"github.com/zekroTJA/shinpuru/internal/services/webserver/wsutil"
	"github.com/zekroTJA/shinpuru/internal/util/imgstore"
	"github.com/zekroTJA/shinpuru/internal/util/report"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/discordutil"
	"github.com/zekroTJA/shinpuru/pkg/roleutil"
)

type MemberReportingController struct {
	session *discordgo.Session
	cfg     *config.Config
	db      database.Database
	st      storage.Storage
}

func (c *MemberReportingController) Setup(container di.Container, router fiber.Router) {
	c.session = container.Get(static.DiDiscordSession).(*discordgo.Session)
	c.cfg = container.Get(static.DiConfig).(*config.Config)
	c.db = container.Get(static.DiDatabase).(database.Database)
	c.st = container.Get(static.DiObjectStorage).(storage.Storage)

	pmw := container.Get(static.DiPermissionMiddleware).(*middleware.PermissionsMiddleware)

	router.Post("/reports", pmw.HandleWs(c.session, "sp.guild.mod.report"), c.postReport)
	router.Post("/kick", pmw.HandleWs(c.session, "sp.guild.mod.kick"), c.postKick)
	router.Post("/ban", pmw.HandleWs(c.session, "sp.guild.mod.ban"), c.postBan)
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

	rep, err := report.PushReport(
		c.session,
		c.db,
		c.cfg.WebServer.PublicAddr,
		guildID,
		uid,
		memberID,
		repReq.Reason,
		repReq.Attachment,
		repReq.Type)

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

	guild, err := discordutil.GetGuild(c.session, guildID)
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

	rep, err := report.PushKick(
		c.session,
		c.db,
		c.cfg.WebServer.PublicAddr,
		guildID,
		uid,
		memberID,
		req.Reason,
		req.Attachment)

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

	guild, err := discordutil.GetGuild(c.session, guildID)
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

	rep, err := report.PushBan(
		c.session,
		c.db,
		c.cfg.WebServer.PublicAddr,
		guildID,
		uid,
		memberID,
		req.Reason,
		req.Attachment)

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

	guild, err := discordutil.GetGuild(c.session, guildID)
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

	rep, err := report.PushMute(
		c.session,
		c.db,
		c.cfg.WebServer.PublicAddr,
		guildID,
		uid,
		memberID,
		req.Reason,
		req.Attachment,
		muteRoleID)

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

	guild, err := discordutil.GetGuild(c.session, guildID)
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

	_, err = report.RevokeMute(
		c.session,
		c.db,
		c.cfg.WebServer.PublicAddr,
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
