package controllers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/core/config"
	"github.com/zekroTJA/shinpuru/internal/core/database"
	"github.com/zekroTJA/shinpuru/internal/core/middleware"
	"github.com/zekroTJA/shinpuru/internal/core/webserver/v1/models"
	"github.com/zekroTJA/shinpuru/internal/core/webserver/wsutil"
	sharedmodels "github.com/zekroTJA/shinpuru/internal/shared/models"
	"github.com/zekroTJA/shinpuru/internal/util/report"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/discordutil"
)

type GuildsController struct {
	session *discordgo.Session
	cfg     *config.Config
	db      database.Database
}

func (c *GuildsController) Setup(container di.Container, router fiber.Router) {
	c.session = container.Get(static.DiDiscordSession).(*discordgo.Session)
	c.cfg = container.Get(static.DiConfig).(*config.Config)
	c.db = container.Get(static.DiDatabase).(database.Database)

	pmw := container.Get(static.DiPermissionMiddleware).(*middleware.PermissionsMiddleware)

	router.Get("", c.getGuilds)
	router.Get("/:guildid", c.getGuild)
	router.Get("/:guildid/scoreboard", c.getGuildScoreboard)
	router.Get("/:guildid/starboard", c.getGuildStarboard)
	router.Delete("/:guildid/antiraid/joinlog", pmw.HandleWs(c.session, "sp.guild.config.antiraid"), c.deleteGuildAntiraidJoinlog)
	router.Get("/:guildid/reports", c.getReports)
	router.Get("/:guildid/reports/count", c.getReportsCount)
}

func (c *GuildsController) getGuilds(ctx *fiber.Ctx) (err error) {
	uid := ctx.Locals("uid").(string)

	guilds := make([]*models.GuildReduced, len(c.session.State.Guilds))
	i := 0
	for _, g := range c.session.State.Guilds {
		if g.MemberCount < 10000 {
			for _, m := range g.Members {
				if m.User.ID == uid {
					guilds[i] = models.GuildReducedFromGuild(g)
					i++
					break
				}
			}
		} else {
			if gm, _ := c.session.GuildMember(g.ID, uid); gm != nil {
				guilds[i] = models.GuildReducedFromGuild(g)
				i++
			}
		}
	}
	guilds = guilds[:i]

	return ctx.JSON(&models.ListResponse{N: len(guilds), Data: guilds})
}

func (c *GuildsController) getGuild(ctx *fiber.Ctx) error {
	uid := ctx.Locals("uid").(string)

	guildID := ctx.Params("guildid")

	memb, _ := c.session.GuildMember(guildID, uid)
	if memb == nil {
		return fiber.ErrNotFound
	}

	guild, err := discordutil.GetGuild(c.session, guildID)
	if err != nil {
		return err
	}

	gRes := models.GuildFromGuild(guild, memb, c.db, c.cfg.Discord.OwnerID)

	return ctx.JSON(gRes)
}

func (c *GuildsController) getGuildScoreboard(ctx *fiber.Ctx) error {
	guildID := ctx.Params("guildid")
	limit, err := wsutil.GetQueryInt(ctx, "limit", 25, 1, 100)
	if err != nil {
		return err
	}

	karmaList, err := c.db.GetKarmaGuild(guildID, limit)

	if err == database.ErrDatabaseNotFound {
		return fiber.ErrNotFound
	} else if err != nil {
		return err
	}

	results := make([]*models.GuildKarmaEntry, len(karmaList))

	var i int
	for _, e := range karmaList {
		member, err := discordutil.GetMember(c.session, guildID, e.UserID)
		if err != nil {
			continue
		}
		results[i] = &models.GuildKarmaEntry{
			Member: models.MemberFromMember(member),
			Value:  e.Value,
		}
		i++
	}

	return ctx.JSON(&models.ListResponse{N: len(results), Data: results})
}

func (c *GuildsController) deleteGuildAntiraidJoinlog(ctx *fiber.Ctx) error {
	guildID := ctx.Params("guildid")

	if err := c.db.FlushAntiraidJoinList(guildID); err != nil && !database.IsErrDatabaseNotFound(err) {
		return err
	}

	return ctx.SendStatus(fiber.StatusOK)
}

func (c *GuildsController) getGuildStarboard(ctx *fiber.Ctx) error {
	guildID := ctx.Params("guildid")
	limit, err := wsutil.GetQueryInt(ctx, "limit", 20, 1, 100)
	if err != nil {
		return err
	}
	offset, err := wsutil.GetQueryInt(ctx, "offset", 0, 0, 0)
	if err != nil {
		return err
	}
	sortQ := ctx.Query("sort")

	var sort sharedmodels.StarboardSortBy
	switch string(sortQ) {
	case "latest":
		sort = sharedmodels.StarboardSortByLatest
	case "top":
		sort = sharedmodels.StarboardSortByMostRated
	default:
		return fiber.NewError(fiber.StatusBadRequest, "invalid sort property")
	}

	entries, err := c.db.GetStarboardEntries(guildID, sort, limit, offset)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return err
	}

	results := make([]*models.StarboardEntryResponse, len(entries))

	var i int
	for _, e := range entries {
		if e.Deleted {
			continue
		}

		member, err := discordutil.GetMember(c.session, guildID, e.AuthorID)
		if err != nil {
			continue
		}

		results[i] = &models.StarboardEntryResponse{
			StarboardEntry: e,
			AuthorUsername: member.User.String(),
			AvatarURL:      member.User.AvatarURL(""),
			MessageURL: discordutil.GetMessageLink(&discordgo.Message{
				ChannelID: e.ChannelID,
				ID:        e.MessageID,
			}, guildID),
		}

		i++
	}

	return ctx.JSON(&models.ListResponse{N: i, Data: results[:i]})
}

func (c *GuildsController) getReports(ctx *fiber.Ctx) error {
	uid := ctx.Locals("uid").(string)

	guildID := ctx.Params("guildid")

	offset, err := wsutil.GetQueryInt(ctx, "offset", 0, 0, 0)
	if err != nil {
		return err
	}

	limit, err := wsutil.GetQueryInt(ctx, "limit", 0, 0, 0)
	if err != nil {
		return err
	}

	if memb, _ := c.session.GuildMember(guildID, uid); memb == nil {
		return fiber.ErrNotFound
	}

	var reps []*report.Report

	reps, err = c.db.GetReportsGuild(guildID, offset, limit)
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

func (c *GuildsController) getReportsCount(ctx *fiber.Ctx) error {
	uid := ctx.Locals("uid").(string)

	guildID := ctx.Params("guildid")

	if memb, _ := c.session.GuildMember(guildID, uid); memb == nil {
		return fiber.ErrNotFound
	}

	count, err := c.db.GetReportsGuildCount(guildID)
	if err != nil {
		return err
	}

	return ctx.JSON(&models.Count{Count: count})
}
