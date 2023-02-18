package controllers

import (
	"strings"

	_ "crypto/sha512"

	"github.com/bwmarrin/discordgo"
	"github.com/bwmarrin/snowflake"
	"github.com/gofiber/fiber/v2"
	"github.com/makeworld-the-better-one/go-isemoji"
	"github.com/sarulabs/di/v2"
	sharedmodels "github.com/zekroTJA/shinpuru/internal/models"
	"github.com/zekroTJA/shinpuru/internal/services/codeexec"
	"github.com/zekroTJA/shinpuru/internal/services/config"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/services/kvcache"
	permservice "github.com/zekroTJA/shinpuru/internal/services/permissions"
	"github.com/zekroTJA/shinpuru/internal/services/storage"
	"github.com/zekroTJA/shinpuru/internal/services/timeprovider"
	"github.com/zekroTJA/shinpuru/internal/services/verification"
	"github.com/zekroTJA/shinpuru/internal/services/webserver/v1/models"
	"github.com/zekroTJA/shinpuru/internal/services/webserver/wsutil"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/discordutil"
	"github.com/zekroTJA/shinpuru/pkg/permissions"
	"github.com/zekroTJA/shinpuru/pkg/stringutil"
	"github.com/zekrotja/dgrs"
	"github.com/zekrotja/sop"
)

type GuildsController struct {
	db      database.Database
	st      storage.Storage
	kvc     kvcache.Provider
	session *discordgo.Session
	cfg     config.Provider
	pmw     *permservice.Permissions
	state   *dgrs.State
	vs      verification.Provider
	cef     codeexec.Factory
	tp      timeprovider.Provider
}

func (c *GuildsController) Setup(container di.Container, router fiber.Router) {
	c.session = container.Get(static.DiDiscordSession).(*discordgo.Session)
	c.cfg = container.Get(static.DiConfig).(config.Provider)
	c.db = container.Get(static.DiDatabase).(database.Database)
	c.pmw = container.Get(static.DiPermissions).(*permservice.Permissions)
	c.kvc = container.Get(static.DiKVCache).(kvcache.Provider)
	c.st = container.Get(static.DiObjectStorage).(storage.Storage)
	c.state = container.Get(static.DiState).(*dgrs.State)
	c.vs = container.Get(static.DiVerification).(verification.Provider)
	c.cef = container.Get(static.DiCodeExecFactory).(codeexec.Factory)
	c.tp = container.Get(static.DiTimeProvider).(timeprovider.Provider)

	router.Get("", c.getGuilds)
	router.Get("/:guildid", c.getGuild)
	router.Get("/:guildid/scoreboard", c.getGuildScoreboard)
	router.Get("/:guildid/starboard", c.getGuildStarboard)
	router.Get("/:guildid/antiraid/joinlog", c.pmw.HandleWs(c.session, "sp.guild.config.antiraid"), c.getGuildAntiraidJoinlog)
	router.Delete("/:guildid/antiraid/joinlog", c.pmw.HandleWs(c.session, "sp.guild.config.antiraid"), c.deleteGuildAntiraidJoinlog)
	router.Get("/:guildid/reports", c.getReports)
	router.Get("/:guildid/reports/count", c.getReportsCount)
	router.Get("/:guildid/permissions", c.getGuildPermissions)
	router.Post("/:guildid/permissions", c.pmw.HandleWs(c.session, "sp.guild.config.perms"), c.postGuildPermissions)
	router.Post("/:guildid/inviteblock", c.pmw.HandleWs(c.session, "sp.guild.mod.inviteblock"), c.postGuildToggleInviteblock)
	router.Get("/:guildid/unbanrequests", c.pmw.HandleWs(c.session, "sp.guild.mod.unbanrequests"), c.getGuildUnbanrequests)
	router.Get("/:guildid/unbanrequests/count", c.pmw.HandleWs(c.session, "sp.guild.mod.unbanrequests"), c.getGuildUnbanrequestsCount)
	router.Get("/:guildid/unbanrequests/:id", c.pmw.HandleWs(c.session, "sp.guild.mod.unbanrequests"), c.getGuildUnbanrequest)
	router.Post("/:guildid/unbanrequests/:id", c.pmw.HandleWs(c.session, "sp.guild.mod.unbanrequests"), c.postGuildUnbanrequest)
}

// @Summary List Guilds
// @Description Returns a list of guilds the authenticated user has in common with shinpuru.
// @Tags Guilds
// @Accept json
// @Produce json
// @Success 200 {array} models.GuildReduced "Wrapped in models.ListResponse"
// @Failure 401 {object} models.Error
// @Router /guilds [get]
func (c *GuildsController) getGuilds(ctx *fiber.Ctx) (err error) {
	uid := ctx.Locals("uid").(string)

	guilds, err := c.state.Guilds()
	if err != nil {
		return err
	}

	userGuilds, err := c.state.UserGuilds(uid)
	if err != nil {
		return
	}

	guildRs := make([]*models.GuildReduced, len(userGuilds))
	i := 0
	for _, guild := range guilds {
		if stringutil.ContainsAny(guild.ID, userGuilds) {
			guildRs[i] = models.GuildReducedFromGuild(guild)
			i++
		}
	}
	guildRs = guildRs[:i]

	return ctx.JSON(models.NewListResponse(guildRs))
}

// @Summary Get Guild
// @Description Returns a single guild object by it's ID.
// @Tags Guilds
// @Accept json
// @Produce json
// @Param id path string true "The ID of the guild."
// @Success 200 {object} models.Guild
// @Failure 401 {object} models.Error
// @Failure 404 {object} models.Error
// @Router /guilds/{id} [get]
func (c *GuildsController) getGuild(ctx *fiber.Ctx) error {
	uid := ctx.Locals("uid").(string)

	guildID := ctx.Params("guildid")

	memb, _ := c.state.Member(guildID, uid)
	if memb == nil {
		return fiber.ErrNotFound
	}

	guild, err := c.state.Guild(guildID, true)
	if err != nil {
		return err
	}

	gRes, err := models.GuildFromGuild(guild, memb, c.db, c.cfg.Config().Discord.OwnerID)
	if err != nil {
		return err
	}

	return ctx.JSON(gRes)
}

// @Summary Get Guild Scoreboard
// @Description Returns a list of scoreboard entries for the given guild.
// @Tags Guilds
// @Accept json
// @Produce json
// @Param id path string true "The ID of the guild."
// @Param limit query int false "Limit the amount of result values" default(25) minimum(1) maximum(100)
// @Success 200 {array} models.GuildKarmaEntry "Wrapped in models.ListResponse"
// @Failure 401 {object} models.Error
// @Failure 404 {object} models.Error
// @Router /guilds/{id}/scoreboard [get]
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
		member, err := c.state.Member(guildID, e.UserID)
		if err != nil {
			continue
		}
		results[i] = &models.GuildKarmaEntry{
			Member: models.MemberFromMember(member),
			Value:  e.Value,
		}
		i++
	}

	return ctx.JSON(models.NewListResponse(results[:i]))
}

// @Summary Get Antiraid Joinlog
// @Description Returns a list of joined members during an antiraid trigger.
// @Tags Guilds
// @Accept json
// @Produce json
// @Param id path string true "The ID of the guild."
// @Success 200 {array} sharedmodels.JoinLogEntry "Wrapped in models.ListResponse"
// @Failure 401 {object} models.Error
// @Failure 404 {object} models.Error
// @Router /guilds/{id}/antiraid/joinlog [get]
func (c *GuildsController) getGuildAntiraidJoinlog(ctx *fiber.Ctx) error {
	guildID := ctx.Params("guildid")

	joinlog, err := c.db.GetAntiraidJoinList(guildID)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return err
	}

	if joinlog == nil {
		joinlog = make([]sharedmodels.JoinLogEntry, 0)
	}

	return ctx.JSON(models.NewListResponse(joinlog))
}

// @Summary Reset Antiraid Joinlog
// @Description Deletes all entries of the antiraid joinlog.
// @Tags Guilds
// @Accept json
// @Produce json
// @Param id path string true "The ID of the guild."
// @Success 200 {object} models.Status
// @Failure 401 {object} models.Error
// @Failure 404 {object} models.Error
// @Router /guilds/{id}/antiraid/joinlog [delete]
func (c *GuildsController) deleteGuildAntiraidJoinlog(ctx *fiber.Ctx) error {
	guildID := ctx.Params("guildid")

	if err := c.db.FlushAntiraidJoinList(guildID); err != nil && !database.IsErrDatabaseNotFound(err) {
		return err
	}

	return ctx.JSON(models.Ok)
}

// @Summary Get Guild Starboard
// @Description Returns a list of starboard entries for the given guild.
// @Tags Guilds
// @Accept json
// @Produce json
// @Param id path string true "The ID of the guild."
// @Success 200 {array} models.StarboardEntryResponse "Wrapped in models.ListResponse"
// @Failure 401 {object} models.Error
// @Failure 404 {object} models.Error
// @Router /guilds/{id}/starboard [get]
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

		member, err := c.state.Member(guildID, e.AuthorID)
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

	return ctx.JSON(models.NewListResponse(results[:i]))
}

// @Summary Get Guild Modlog
// @Description Returns a list of guild modlog entries for the given guild.
// @Tags Guilds
// @Accept json
// @Produce json
// @Param id path string true "The ID of the guild."
// @Param offset query int false "The offset of returned entries" default(0)
// @Param limit query int false "The amount of returned entries (0 = all)" default(0)
// @Success 200 {array} models.Report "Wrapped in models.ListResponse"
// @Failure 401 {object} models.Error
// @Failure 404 {object} models.Error
// @Router /guilds/{id}/reports [get]
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

	var reps []sharedmodels.Report

	reps, err = c.db.GetReportsGuild(guildID, offset, limit)
	if err != nil {
		return err
	}

	resReps := make([]models.Report, 0)
	if reps != nil {
		resReps = make([]models.Report, len(reps))
		for i, r := range reps {
			resReps[i] = models.ReportFromReport(r, c.cfg.Config().WebServer.PublicAddr)
			user, err := c.state.User(r.VictimID)
			if err == nil {
				resReps[i].Victim = models.FlatUserFromUser(user)
			}
			user, err = c.state.User(r.ExecutorID)
			if err == nil {
				resReps[i].Executor = models.FlatUserFromUser(user)
			}
		}
	}

	return ctx.JSON(models.NewListResponse(resReps))
}

// @Summary Get Guild Modlog Count
// @Description Returns the total count of entries in the guild mod log.
// @Tags Guilds
// @Accept json
// @Produce json
// @Param id path string true "The ID of the guild."
// @Success 200 {object} models.Count
// @Failure 401 {object} models.Error
// @Failure 404 {object} models.Error
// @Router /guilds/{id}/reports/count [get]
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

// @Summary Get Guild Permission Settings
// @Description Returns the specified guild permission settings.
// @Tags Guilds
// @Accept json
// @Produce json
// @Param id path string true "The ID of the guild."
// @Success 200 {object} models.PermissionsMap
// @Failure 401 {object} models.Error
// @Failure 404 {object} models.Error
// @Router /guilds/{id}/permissions [get]
func (c *GuildsController) getGuildPermissions(ctx *fiber.Ctx) error {
	uid := ctx.Locals("uid").(string)

	guildID := ctx.Params("guildid")

	if memb, _ := c.session.GuildMember(guildID, uid); memb == nil {
		return fiber.ErrNotFound
	}

	var perms models.PermissionsMap
	var err error

	if perms, err = c.db.GetGuildPermissions(guildID); err != nil && !database.IsErrDatabaseNotFound(err) {
		return err
	}

	return ctx.JSON(perms)
}

// @Summary Apply Guild Permission Rule
// @Description Apply a new guild permission rule for a specified role.
// @Tags Guilds
// @Accept json
// @Produce json
// @Param id path string true "The ID of the guild."
// @Param payload body models.PermissionsUpdate true "The permission rule payload."
// @Success 200 {object} models.Status
// @Failure 400 {object} models.Error
// @Failure 401 {object} models.Error
// @Failure 404 {object} models.Error
// @Router /guilds/{id}/permissions [post]
func (c *GuildsController) postGuildPermissions(ctx *fiber.Ctx) error {
	guildID := ctx.Params("guildid")

	update := new(models.PermissionsUpdate)
	if err := ctx.BodyParser(update); err != nil {
		return fiber.ErrBadRequest
	}

	sperm := update.Perm[1:]
	if !strings.HasPrefix(sperm, "sp.guild") && !strings.HasPrefix(sperm, "sp.etc") && !strings.HasPrefix(sperm, "sp.chat") {
		return fiber.NewError(fiber.StatusBadRequest, "you can only give permissions over the domains 'sp.guild', 'sp.etc' and 'sp.chat'")
	}

	perms, err := c.db.GetGuildPermissions(guildID)
	if err != nil {
		if database.IsErrDatabaseNotFound(err) {
			return fiber.ErrNotFound
		}
		return err
	}

	for _, roleID := range update.RoleIDs {
		rperms, ok := perms[roleID]
		if !ok {
			rperms = permissions.PermissionArray{}
		}

		rperms, changed := rperms.Update(update.Perm, update.Override)

		if len(rperms) == 0 {
			delete(perms, roleID)
		} else {
			perms[roleID] = rperms
		}

		if changed {
			if err = c.db.SetGuildRolePermission(guildID, roleID, rperms); err != nil {
				return err
			}
		}
	}

	return ctx.JSON(perms)
}

// @Summary Toggle Guild Inviteblock Enable
// @Description Toggle enabled state of the guild invite block system.
// @Tags Guilds
// @Accept json
// @Produce json
// @Param id path string true "The ID of the guild."
// @Param payload body models.EnableStatus true "The enable status payload."
// @Success 200 {object} models.Status
// @Failure 400 {object} models.Error
// @Failure 401 {object} models.Error
// @Failure 404 {object} models.Error
// @Router /guilds/{id}/inviteblock [post]
func (c *GuildsController) postGuildToggleInviteblock(ctx *fiber.Ctx) error {
	guildID := ctx.Params("guildid")

	var data models.EnableStatus
	if err := ctx.BodyParser(&data); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	val := ""
	if data.Enabled {
		val = "1"
	}

	if err := c.db.SetGuildInviteBlock(guildID, val); err != nil {
		return err
	}

	return ctx.JSON(models.Ok)
}

// @Summary Get Guild Unbanrequests
// @Description Returns the list of the guild unban requests.
// @Tags Guilds
// @Accept json
// @Produce json
// @Param id path string true "The ID of the guild."
// @Success 200 {array} models.RichUnbanRequest "Wrapped in models.ListReponse"
// @Failure 401 {object} models.Error
// @Failure 404 {object} models.Error
// @Router /guilds/{id}/unbanrequests [get]
func (c *GuildsController) getGuildUnbanrequests(ctx *fiber.Ctx) error {
	guildID := ctx.Params("guildid")

	requests, err := c.db.GetGuildUnbanRequests(guildID)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return err
	}
	if requests == nil {
		requests = make([]sharedmodels.UnbanRequest, 0)
	}

	res := sop.Map[sharedmodels.UnbanRequest](sop.Slice(requests),
		func(r sharedmodels.UnbanRequest, i int) *models.RichUnbanRequest {
			r.Hydrate()
			rub := &models.RichUnbanRequest{
				UnbanRequest: r,
			}
			if creator, _ := c.state.User(rub.UserID); creator != nil {
				rub.Creator = models.FlatUserFromUser(creator)
			}
			if proc, _ := c.state.User(rub.ProcessedBy); proc != nil {
				rub.Processor = models.FlatUserFromUser(proc)
			}
			return rub
		})

	return ctx.JSON(models.NewListResponse(res.Unwrap()))
}

// @Summary Get Guild Unbanrequests Count
// @Description Returns the total or filtered count of guild unban requests.
// @Tags Guilds
// @Accept json
// @Produce json
// @Param id path string true "The ID of the guild."
// @Param state query sharedmodels.UnbanRequestState false "Filter count by given state."
// @Success 200 {object} models.Count
// @Failure 400 {object} models.Error
// @Failure 401 {object} models.Error
// @Failure 404 {object} models.Error
// @Router /guilds/{id}/unbanrequests/count [get]
func (c *GuildsController) getGuildUnbanrequestsCount(ctx *fiber.Ctx) error {
	guildID := ctx.Params("guildid")

	stateFilter, err := wsutil.GetQueryInt(ctx, "state", -1, 0, 0)
	if err != nil {
		return err
	}

	requests, err := c.db.GetGuildUnbanRequests(guildID)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return err
	}
	if requests == nil {
		requests = make([]sharedmodels.UnbanRequest, 0)
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

// @Summary Get Single Guild Unbanrequest
// @Description Returns a single guild unban request by ID.
// @Tags Guilds
// @Accept json
// @Produce json
// @Param id path string true "The ID of the guild."
// @Param requestid path string true "The ID of the unbanrequest."
// @Success 200 {object} models.RichUnbanRequest
// @Failure 401 {object} models.Error
// @Failure 404 {object} models.Error
// @Router /guilds/{id}/unbanrequests/{requestid} [get]
func (c *GuildsController) getGuildUnbanrequest(ctx *fiber.Ctx) error {
	guildID := ctx.Params("guildid")
	id := ctx.Params("id")

	request, err := c.db.GetUnbanRequest(id)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return err
	}
	if request.GuildID != guildID {
		return fiber.ErrNotFound
	}

	request.Hydrate()
	rub := &models.RichUnbanRequest{
		UnbanRequest: request,
	}
	if creator, _ := c.state.User(rub.UserID); creator != nil {
		rub.Processor = models.FlatUserFromUser(creator)
	}
	if proc, _ := c.state.User(rub.ProcessedBy); proc != nil {
		rub.Processor = models.FlatUserFromUser(proc)
	}

	return ctx.JSON(rub)
}

// @Summary Process Guild Unbanrequest
// @Description Process a guild unban request.
// @Tags Guilds
// @Accept json
// @Produce json
// @Param id path string true "The ID of the guild."
// @Param requestid path string true "The ID of the unbanrequest."
// @Success 200 {object} models.RichUnbanRequest
// @Failure 401 {object} models.Error
// @Failure 404 {object} models.Error
// @Router /guilds/{id}/unbanrequests/{requestid} [post]
func (c *GuildsController) postGuildUnbanrequest(ctx *fiber.Ctx) error {
	uid := ctx.Locals("uid").(string)

	guildID := ctx.Params("guildid")
	id := ctx.Params("id")

	rUpdate := new(sharedmodels.UnbanRequest)
	if err := ctx.BodyParser(rUpdate); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	request, err := c.db.GetUnbanRequest(id)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return err
	}
	if request.GuildID != guildID {
		return fiber.ErrNotFound
	}

	if rUpdate.ProcessedMessage == "" {
		return fiber.NewError(fiber.StatusBadRequest, "process reason message must be provided")
	}

	if request.ID, err = snowflake.ParseString(id); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	request.ProcessedBy = uid
	request.Status = rUpdate.Status
	request.Processed = c.tp.Now()
	request.ProcessedMessage = rUpdate.ProcessedMessage

	if err = c.db.UpdateUnbanRequest(request); err != nil {
		return err
	}

	if request.Status == sharedmodels.UnbanRequestStateAccepted {
		if err = c.session.GuildBanDelete(request.GuildID, request.UserID); err != nil {
			return err
		}
	}

	request.Hydrate()
	rub := &models.RichUnbanRequest{
		UnbanRequest: request,
	}
	if creator, _ := c.state.User(rub.UserID); creator != nil {
		rub.Processor = models.FlatUserFromUser(creator)
	}
	if proc, _ := c.state.User(rub.ProcessedBy); proc != nil {
		rub.Processor = models.FlatUserFromUser(proc)
	}

	return ctx.JSON(rub)
}

// ---------------------------------------------------------------------------
// - HELPERS

func checkEmojis(emojis []string) bool {
	for _, e := range emojis {
		if !isemoji.IsEmojiNonStrict(e) {
			return false
		}
	}
	return true
}
