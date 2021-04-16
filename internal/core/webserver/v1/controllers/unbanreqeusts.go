package controllers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/core/database"
	"github.com/zekroTJA/shinpuru/internal/core/middleware"
	"github.com/zekroTJA/shinpuru/internal/core/webserver/v1/models"
	"github.com/zekroTJA/shinpuru/internal/util/report"
	"github.com/zekroTJA/shinpuru/internal/util/snowflakenodes"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/discordutil"
)

type UnbanrequestsController struct {
	session *discordgo.Session
	db      database.Database
	pmw     *middleware.PermissionsMiddleware
}

func (c *UnbanrequestsController) Setup(container di.Container, router fiber.Router) {
	c.session = container.Get(static.DiDiscordSession).(*discordgo.Session)
	c.db = container.Get(static.DiDatabase).(database.Database)
	c.pmw = container.Get(static.DiPermissionMiddleware).(*middleware.PermissionsMiddleware)

	router.Get("", c.getUnbanrequests)
	router.Post("", c.postUnbanrequests)
	router.Get("/bannedguilds", c.getBannedGuilds)
}

func (c *UnbanrequestsController) getUnbanrequests(ctx *fiber.Ctx) error {
	uid := ctx.Locals("uid").(string)

	requests, err := c.db.GetGuildUserUnbanRequests(uid, "")
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return err
	}
	if requests == nil {
		requests = make([]*report.UnbanRequest, 0)
	}

	for _, r := range requests {
		r.Hydrate()
	}

	return ctx.JSON(&models.ListResponse{N: len(requests), Data: requests})
}

func (c *UnbanrequestsController) postUnbanrequests(ctx *fiber.Ctx) error {
	uid := ctx.Locals("uid").(string)

	user, err := c.session.User(uid)
	if err != nil {
		return err
	}

	req := new(report.UnbanRequest)
	if err := ctx.BodyParser(req); err != nil {
		return err
	}
	if err := req.Validate(); err != nil {
		return err
	}

	rep, err := c.db.GetReportsFiltered(req.GuildID, uid, 1)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return err
	}

	if rep == nil || len(rep) == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "you have no filed ban reports on this guild")
	}

	requests, err := c.db.GetGuildUserUnbanRequests(uid, req.GuildID)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return err
	}

	if requests != nil {
		for _, r := range requests {
			if r.Status == report.UnbanRequestStatePending {
				return fiber.NewError(fiber.StatusBadRequest, "there is still one open unban request to be proceed")
			}
		}
	}

	finalReq := &report.UnbanRequest{
		ID:      snowflakenodes.NodeUnbanRequests.Generate(),
		UserID:  uid,
		GuildID: req.GuildID,
		UserTag: user.String(),
		Message: req.Message,
		Status:  report.UnbanRequestStatePending,
	}

	if err := c.db.AddUnbanRequest(finalReq); err != nil && !database.IsErrDatabaseNotFound(err) {
		return err
	}

	return ctx.JSON(finalReq.Hydrate())
}

func (c *UnbanrequestsController) getBannedGuilds(ctx *fiber.Ctx) error {
	uid := ctx.Locals("uid").(string)

	guildsArr, err := c.getUserBannedGuilds(uid)
	if err != nil {
		return err
	}

	return ctx.JSON(&models.ListResponse{N: len(guildsArr), Data: guildsArr})
}

// --- HELPERS ------------

func (c *UnbanrequestsController) getUserBannedGuilds(userID string) ([]*models.GuildReduced, error) {
	reps, err := c.db.GetReportsFiltered("", userID, 1)
	if err != nil {
		if database.IsErrDatabaseNotFound(err) {
			return []*models.GuildReduced{}, nil
		}
		return nil, err
	}

	guilds := make(map[string]*models.GuildReduced)
	for _, r := range reps {
		if _, ok := guilds[r.GuildID]; ok {
			continue
		}
		guild, err := discordutil.GetGuild(c.session, r.GuildID)
		if err != nil {
			return nil, err
		}
		guilds[r.GuildID] = models.GuildReducedFromGuild(guild)
	}

	guildsArr := make([]*models.GuildReduced, len(guilds))
	i := 0
	for _, g := range guilds {
		guildsArr[i] = g
		i++
	}

	return guildsArr, nil
}
