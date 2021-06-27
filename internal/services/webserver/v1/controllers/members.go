package controllers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/config"
	"github.com/zekroTJA/shinpuru/internal/middleware"
	sharedmodels "github.com/zekroTJA/shinpuru/internal/models"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/services/webserver/v1/models"
	"github.com/zekroTJA/shinpuru/internal/services/webserver/wsutil"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/discordutil"
	"github.com/zekroTJA/shireikan"
)

type GuildMembersController struct {
	session    *discordgo.Session
	cfg        *config.Config
	db         database.Database
	pmw        *middleware.PermissionsMiddleware
	cmdHandler shireikan.Handler
}

func (c *GuildMembersController) Setup(container di.Container, router fiber.Router) {
	c.session = container.Get(static.DiDiscordSession).(*discordgo.Session)
	c.cfg = container.Get(static.DiConfig).(*config.Config)
	c.db = container.Get(static.DiDatabase).(database.Database)
	c.pmw = container.Get(static.DiPermissionMiddleware).(*middleware.PermissionsMiddleware)
	c.cmdHandler = container.Get(static.DiCommandHandler).(shireikan.Handler)

	router.Get("/members", c.getMembers)
	router.Get("/:memberid", c.getMember)
	router.Get("/:memberid/permissions", c.getMemberPermissions)
	router.Get("/:memberid/permissions/allowed", c.getMemberPermissionsAllowed)
	router.Get("/:memberid/reports", c.getReports)
	router.Get("/:memberid/reports/count", c.getReportsCount)
	router.Get("/:memberid/unbanrequests", c.pmw.HandleWs(c.session, "sp.guild.mod.unbanrequests"), c.getMemberUnbanrequests)
	router.Get("/:memberid/unbanrequests/count", c.pmw.HandleWs(c.session, "sp.guild.mod.unbanrequests"), c.getMemberUnbanrequestsCount)
}

func (c *GuildMembersController) getMembers(ctx *fiber.Ctx) (err error) {
	uid := ctx.Locals("uid").(string)

	guildID := ctx.Params("guildid")

	memb, _ := c.session.GuildMember(guildID, uid)
	if memb == nil {
		return fiber.ErrNotFound
	}

	after := ""
	limit := 0

	after = ctx.Query("after")
	limit, err = wsutil.GetQueryInt(ctx, "limit", 100, 0, 2000)
	if err != nil {
		return err
	}

	members, err := c.session.GuildMembers(guildID, after, limit)
	if err != nil {
		return err
	}

	memblen := len(members)
	fhmembers := make([]*models.Member, memblen)

	for i, m := range members {
		fhmembers[i] = models.MemberFromMember(m)
	}

	return ctx.JSON(&models.ListResponse{N: len(fhmembers), Data: fhmembers})
}

func (c *GuildMembersController) getMember(ctx *fiber.Ctx) (err error) {
	uid := ctx.Locals("uid").(string)

	guildID := ctx.Params("guildid")
	memberID := ctx.Params("memberid")

	var memb *discordgo.Member

	if memb, _ = c.session.GuildMember(guildID, uid); memb == nil {
		return fiber.ErrNotFound
	}

	guild, err := discordutil.GetGuild(c.session, guildID)
	if err != nil {
		return err
	}

	memb, _ = c.session.GuildMember(guildID, memberID)
	if memb == nil {
		return fiber.ErrNotFound
	}

	memb.GuildID = guildID

	mm := models.MemberFromMember(memb)

	switch {
	case discordutil.IsAdmin(guild, memb):
		mm.Dominance = 1
	case guild.OwnerID == memberID:
		mm.Dominance = 2
	case c.cfg.Discord.OwnerID == memb.User.ID:
		mm.Dominance = 3
	}

	mm.Karma, err = c.db.GetKarma(memberID, guildID)
	if !database.IsErrDatabaseNotFound(err) && err != nil {
		return err
	}

	mm.KarmaTotal, err = c.db.GetKarmaSum(memberID)
	if !database.IsErrDatabaseNotFound(err) && err != nil {
		return err
	}

	if muteRoleID, err := c.db.GetGuildMuteRole(guildID); err == nil {
		for _, roleID := range memb.Roles {
			if roleID == muteRoleID {
				mm.ChatMuted = true
				break
			}
		}
	}

	return ctx.JSON(mm)
}

func (c *GuildMembersController) getMemberPermissions(ctx *fiber.Ctx) (err error) {
	uid := ctx.Locals("uid").(string)

	guildID := ctx.Params("guildid")
	memberID := ctx.Params("memberid")

	if memb, _ := c.session.GuildMember(guildID, uid); memb == nil {
		return fiber.ErrNotFound
	}

	perm, _, err := c.pmw.GetPermissions(c.session, guildID, memberID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return ctx.JSON(&models.PermissionsResponse{
		Permissions: perm,
	})
}

func (c *GuildMembersController) getMemberPermissionsAllowed(ctx *fiber.Ctx) (err error) {
	guildID := ctx.Params("guildid")
	memberID := ctx.Params("memberid")

	perms, _, err := c.pmw.GetPermissions(c.session, guildID, memberID)
	if database.IsErrDatabaseNotFound(err) {
		return fiber.ErrNotFound
	}
	if err != nil {
		return err
	}

	cmds := c.cmdHandler.GetCommandInstances()

	allowed := make([]string, len(cmds)+len(static.AdditionalPermissions))
	i := 0
	for _, cmd := range cmds {
		if perms.Check(cmd.GetDomainName()) {
			allowed[i] = cmd.GetDomainName()
			i++
		}
	}

	for _, p := range static.AdditionalPermissions {
		if perms.Check(p) {
			allowed[i] = p
			i++
		}
	}

	return ctx.JSON(&models.ListResponse{N: i, Data: allowed[:i]})
}

func (c *GuildMembersController) getReports(ctx *fiber.Ctx) (err error) {
	uid := ctx.Locals("uid").(string)

	guildID := ctx.Params("guildid")
	memberID := ctx.Params("memberid")

	if memb, _ := c.session.GuildMember(guildID, uid); memb == nil {
		return fiber.ErrNotFound
	}

	reps, err := c.db.GetReportsFiltered(guildID, memberID, -1)
	if err != nil {
		return err
	}

	resReps := make([]*models.Report, 0)
	if reps != nil {
		resReps = make([]*models.Report, len(reps))
		for i, r := range reps {
			resReps[i] = models.ReportFromReport(r, c.cfg.WebServer.PublicAddr)
		}
	}

	return ctx.JSON(&models.ListResponse{N: len(resReps), Data: resReps})
}

func (c *GuildMembersController) getReportsCount(ctx *fiber.Ctx) (err error) {
	uid := ctx.Locals("uid").(string)

	guildID := ctx.Params("guildid")
	memberID := ctx.Params("memberid")

	if memb, _ := c.session.GuildMember(guildID, uid); memb == nil {
		return fiber.ErrNotFound
	}

	count, err := c.db.GetReportsFilteredCount(guildID, memberID, -1)
	if err != nil {
		return err
	}

	return ctx.JSON(&models.Count{Count: count})
}

func (c *GuildMembersController) getMemberUnbanrequests(ctx *fiber.Ctx) (err error) {
	guildID := ctx.Params("guildid")
	memberID := ctx.Params("memberid")

	requests, err := c.db.GetGuildUserUnbanRequests(guildID, memberID)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return err
	}
	if requests == nil {
		requests = make([]*sharedmodels.UnbanRequest, 0)
	}

	for _, r := range requests {
		r.Hydrate()
	}

	return ctx.JSON(&models.ListResponse{N: len(requests), Data: requests})
}

func (c *GuildMembersController) getMemberUnbanrequestsCount(ctx *fiber.Ctx) (err error) {
	guildID := ctx.Params("guildid")
	memberID := ctx.Params("memberid")

	stateFilter, err := wsutil.GetQueryInt(ctx, "state", -1, 0, 0)
	if err != nil {
		return err
	}

	requests, err := c.db.GetGuildUserUnbanRequests(guildID, memberID)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return err
	}
	if requests == nil {
		requests = make([]*sharedmodels.UnbanRequest, 0)
	}

	count := len(requests)
	if stateFilter > -1 {
		count = 0
		for _, r := range requests {
			if int(r.Status) == stateFilter {
				count++
			}
		}
	}

	return ctx.JSON(&models.Count{Count: count})
}
