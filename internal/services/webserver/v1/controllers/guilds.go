package controllers

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/bwmarrin/snowflake"
	"github.com/gofiber/fiber/v2"
	"github.com/makeworld-the-better-one/go-isemoji"
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/config"
	"github.com/zekroTJA/shinpuru/internal/middleware"
	sharedmodels "github.com/zekroTJA/shinpuru/internal/models"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/services/kvcache"
	"github.com/zekroTJA/shinpuru/internal/services/storage"
	"github.com/zekroTJA/shinpuru/internal/services/webserver/v1/models"
	"github.com/zekroTJA/shinpuru/internal/services/webserver/wsutil"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/snowflakenodes"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/discordutil"
	"github.com/zekroTJA/shinpuru/pkg/fetch"
	"github.com/zekroTJA/shinpuru/pkg/permissions"
	"github.com/zekroTJA/shinpuru/pkg/stringutil"
	"github.com/zekrotja/dgrs"
)

type GuildsController struct {
	db      database.Database
	st      storage.Storage
	kvc     kvcache.Provider
	session *discordgo.Session
	cfg     *config.Config
	pmw     *middleware.PermissionsMiddleware
	state   *dgrs.State
}

func (c *GuildsController) Setup(container di.Container, router fiber.Router) {
	c.session = container.Get(static.DiDiscordSession).(*discordgo.Session)
	c.cfg = container.Get(static.DiConfig).(*config.Config)
	c.db = container.Get(static.DiDatabase).(database.Database)
	c.pmw = container.Get(static.DiPermissionMiddleware).(*middleware.PermissionsMiddleware)
	c.kvc = container.Get(static.DiKVCache).(kvcache.Provider)
	c.st = container.Get(static.DiObjectStorage).(storage.Storage)
	c.state = container.Get(static.DiState).(*dgrs.State)

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
	router.Get("/:guildid/settings", c.getGuildSettings)
	router.Post("/:guildid/settings", c.postGuildSettings)
	router.Get("/:guildid/settings/karma", c.pmw.HandleWs(c.session, "sp.guild.config.karma"), c.getGuildSettingsKarma)
	router.Post("/:guildid/settings/karma", c.pmw.HandleWs(c.session, "sp.guild.config.karma"), c.postGuildSettingsKarma)
	router.Get("/:guildid/settings/karma/blocklist", c.pmw.HandleWs(c.session, "sp.guild.config.karma"), c.getGuildSettingsKarmaBlocklist)
	router.Put("/:guildid/settings/karma/blocklist/:memberid", c.pmw.HandleWs(c.session, "sp.guild.config.karma"), c.putGuildSettingsKarmaBlocklist)
	router.Delete("/:guildid/settings/karma/blocklist/:memberid", c.pmw.HandleWs(c.session, "sp.guild.config.karma"), c.deleteGuildSettingsKarmaBlocklist)
	router.Get("/:guildid/settings/karma/rules", c.pmw.HandleWs(c.session, "sp.guild.config.karma"), c.getGuildSettingsKarmaRules)
	router.Post("/:guildid/settings/karma/rules", c.pmw.HandleWs(c.session, "sp.guild.config.karma"), c.createGuildSettingsKrameRule)
	router.Post("/:guildid/settings/karma/rules/:id", c.pmw.HandleWs(c.session, "sp.guild.config.karma"), c.updateGuildSettingsKrameRule)
	router.Delete("/:guildid/settings/karma/rules/:id", c.pmw.HandleWs(c.session, "sp.guild.config.karma"), c.deleteGuildSettingsKrameRule)
	router.Get("/:guildid/settings/antiraid", c.pmw.HandleWs(c.session, "sp.guild.config.antiraid"), c.getGuildSettingsAntiraid)
	router.Post("/:guildid/settings/antiraid", c.pmw.HandleWs(c.session, "sp.guild.config.antiraid"), c.postGuildSettingsAntiraid)
	router.Get("/:guildid/settings/logs", c.pmw.HandleWs(c.session, "sp.guild.config.logs"), c.getGuildSettingsLogs)
	router.Get("/:guildid/settings/logs/count", c.pmw.HandleWs(c.session, "sp.guild.config.logs"), c.getGuildSettingsLogsCount)
	router.Delete("/:guildid/settings/logs", c.pmw.HandleWs(c.session, "sp.guild.config.logs"), c.deleteGuildSettingsLogEntry)
	router.Delete("/:guildid/settings/logs/:id", c.pmw.HandleWs(c.session, "sp.guild.config.logs"), c.deleteGuildSettingsLogEntry)
	router.Get("/:guildid/settings/logs/state", c.pmw.HandleWs(c.session, "sp.guild.config.logs"), c.getGuildSettingsLogsState)
	router.Post("/:guildid/settings/logs/state", c.pmw.HandleWs(c.session, "sp.guild.config.logs"), c.postGuildSettingsLogsState)
	router.Post("/:guildid/settings/flushguilddata", c.pmw.HandleWs(c.session, "sp.guild.admin.flushdata"), c.postFlushGuildData)
}

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

	return ctx.JSON(&models.ListResponse{N: len(guildRs), Data: guildRs})
}

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

	gRes, err := models.GuildFromGuild(guild, memb, c.db, c.cfg.Discord.OwnerID)
	if err != nil {
		return err
	}

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

	return ctx.JSON(&models.ListResponse{N: i, Data: results[:i]})
}

func (c *GuildsController) getGuildAntiraidJoinlog(ctx *fiber.Ctx) error {
	guildID := ctx.Params("guildid")

	joinlog, err := c.db.GetAntiraidJoinList(guildID)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return err
	}

	if joinlog == nil {
		joinlog = make([]*sharedmodels.JoinLogEntry, 0)
	}

	return ctx.JSON(&models.ListResponse{N: len(joinlog), Data: joinlog})
}

func (c *GuildsController) deleteGuildAntiraidJoinlog(ctx *fiber.Ctx) error {
	guildID := ctx.Params("guildid")

	if err := c.db.FlushAntiraidJoinList(guildID); err != nil && !database.IsErrDatabaseNotFound(err) {
		return err
	}

	return ctx.JSON(models.Ok)
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

	var reps []*sharedmodels.Report

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

func (c *GuildsController) getGuildSettings(ctx *fiber.Ctx) error {
	guildID := ctx.Params("guildid")

	gs := new(models.GuildSettings)
	var err error

	if gs.Prefix, err = c.db.GetGuildPrefix(guildID); err != nil && !database.IsErrDatabaseNotFound(err) {
		return err
	}

	if gs.Perms, err = c.db.GetGuildPermissions(guildID); err != nil && !database.IsErrDatabaseNotFound(err) {
		return err
	}

	if gs.AutoRoles, err = c.db.GetGuildAutoRole(guildID); err != nil && !database.IsErrDatabaseNotFound(err) {
		return err
	}

	if gs.ModLogChannel, err = c.db.GetGuildModLog(guildID); err != nil && !database.IsErrDatabaseNotFound(err) {
		return err
	}

	if gs.VoiceLogChannel, err = c.db.GetGuildVoiceLog(guildID); err != nil && !database.IsErrDatabaseNotFound(err) {
		return err
	}

	if gs.JoinMessageChannel, gs.JoinMessageText, err = c.db.GetGuildJoinMsg(guildID); err != nil && !database.IsErrDatabaseNotFound(err) {
		return err
	}

	if gs.LeaveMessageChannel, gs.LeaveMessageText, err = c.db.GetGuildLeaveMsg(guildID); err != nil && !database.IsErrDatabaseNotFound(err) {
		return err
	}

	return ctx.JSON(gs)
}

func (c *GuildsController) postGuildSettings(ctx *fiber.Ctx) error {
	uid := ctx.Locals("uid").(string)

	guildID := ctx.Params("guildid")

	var err error

	gs := new(models.GuildSettings)
	if err = ctx.BodyParser(gs); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if gs.AutoRoles != nil {
		if ok, _, err := c.pmw.CheckPermissions(c.session, guildID, uid, "sp.guild.config.autorole"); err != nil {
			return wsutil.ErrInternalOrNotFound(err)
		} else if !ok {
			return fiber.ErrUnauthorized
		}

		if err = c.db.SetGuildAutoRole(guildID, gs.AutoRoles); err != nil {
			return wsutil.ErrInternalOrNotFound(err)
		}
	}

	if gs.ModLogChannel != "" {
		if ok, _, err := c.pmw.CheckPermissions(c.session, guildID, uid, "sp.guild.config.modlog"); err != nil {
			return wsutil.ErrInternalOrNotFound(err)
		} else if !ok {
			return fiber.ErrUnauthorized
		}

		if gs.ModLogChannel == "__RESET__" {
			gs.ModLogChannel = ""
		}

		if err = c.db.SetGuildModLog(guildID, gs.ModLogChannel); err != nil {
			return wsutil.ErrInternalOrNotFound(err)
		}
	}

	if gs.Prefix != "" {
		if ok, _, err := c.pmw.CheckPermissions(c.session, guildID, uid, "sp.guild.config.prefix"); err != nil {
			return wsutil.ErrInternalOrNotFound(err)
		} else if !ok {
			return fiber.ErrUnauthorized
		}

		if gs.Prefix == "__RESET__" {
			gs.Prefix = ""
		}

		if err = c.db.SetGuildPrefix(guildID, gs.Prefix); err != nil {
			return wsutil.ErrInternalOrNotFound(err)
		}
	}

	if gs.VoiceLogChannel != "" {
		if ok, _, err := c.pmw.CheckPermissions(c.session, guildID, uid, "sp.guild.config.voicelog"); err != nil {
			return wsutil.ErrInternalOrNotFound(err)
		} else if !ok {
			return fiber.ErrUnauthorized
		}

		if gs.VoiceLogChannel == "__RESET__" {
			gs.VoiceLogChannel = ""
		}

		if err = c.db.SetGuildVoiceLog(guildID, gs.VoiceLogChannel); err != nil {
			return wsutil.ErrInternalOrNotFound(err)
		}
	}

	if gs.JoinMessageChannel != "" && gs.JoinMessageText != "" {
		if ok, _, err := c.pmw.CheckPermissions(c.session, guildID, uid, "sp.guild.config.joinmsg"); err != nil {
			return wsutil.ErrInternalOrNotFound(err)
		} else if !ok {
			return fiber.ErrUnauthorized
		}

		if gs.JoinMessageChannel == "__RESET__" && gs.JoinMessageText == "__RESET__" {
			gs.JoinMessageChannel = ""
			gs.JoinMessageText = ""
		}

		if err = c.db.SetGuildJoinMsg(guildID, gs.JoinMessageChannel, gs.JoinMessageText); err != nil {
			return wsutil.ErrInternalOrNotFound(err)
		}
	}

	if gs.LeaveMessageChannel != "" && gs.LeaveMessageText != "" {
		if ok, _, err := c.pmw.CheckPermissions(c.session, guildID, uid, "sp.guild.config.leavemsg"); err != nil {
			return wsutil.ErrInternalOrNotFound(err)
		} else if !ok {
			return fiber.ErrUnauthorized
		}

		if gs.LeaveMessageChannel == "__RESET__" && gs.LeaveMessageText == "__RESET__" {
			gs.LeaveMessageChannel = ""
			gs.LeaveMessageText = ""
		}

		if err = c.db.SetGuildLeaveMsg(guildID, gs.LeaveMessageChannel, gs.LeaveMessageText); err != nil {
			return wsutil.ErrInternalOrNotFound(err)
		}
	}

	return ctx.JSON(models.Ok)
}

func (c *GuildsController) getGuildPermissions(ctx *fiber.Ctx) error {
	uid := ctx.Locals("uid").(string)

	guildID := ctx.Params("guildid")

	if memb, _ := c.session.GuildMember(guildID, uid); memb == nil {
		return fiber.ErrNotFound
	}

	var perms map[string]permissions.PermissionArray
	var err error

	if perms, err = c.db.GetGuildPermissions(guildID); err != nil && !database.IsErrDatabaseNotFound(err) {
		return err
	}

	return ctx.JSON(perms)
}

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
			rperms = make(permissions.PermissionArray, 0)
		}

		rperms, changed := rperms.Update(update.Perm, false)

		if changed {
			if err = c.db.SetGuildRolePermission(guildID, roleID, rperms); err != nil {
				return err
			}
		}
	}

	return ctx.JSON(models.Ok)
}

func (c *GuildsController) postGuildToggleInviteblock(ctx *fiber.Ctx) error {
	guildID := ctx.Params("guildid")

	var data struct {
		Enabled bool `json:"enabled"`
	}

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

func (c *GuildsController) getGuildSettingsKarma(ctx *fiber.Ctx) error {
	guildID := ctx.Params("guildid")

	settings := new(models.KarmaSettings)

	var err error

	if settings.State, err = c.db.GetKarmaState(guildID); err != nil && !database.IsErrDatabaseNotFound(err) {
		return err
	}

	if settings.Tokens, err = c.db.GetKarmaTokens(guildID); err != nil && !database.IsErrDatabaseNotFound(err) {
		return err
	}

	emotesInc, emotesDec, err := c.db.GetKarmaEmotes(guildID)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return err
	}
	settings.EmotesIncrease = strings.Split(emotesInc, "")
	settings.EmotesDecrease = strings.Split(emotesDec, "")

	if settings.Penalty, err = c.db.GetKarmaPenalty(guildID); err != nil && !database.IsErrDatabaseNotFound(err) {
		return err
	}

	return ctx.JSON(settings)
}

func (c *GuildsController) postGuildSettingsKarma(ctx *fiber.Ctx) error {
	guildID := ctx.Params("guildid")

	settings := new(models.KarmaSettings)
	var err error

	if err = ctx.BodyParser(settings); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err = c.db.SetKarmaState(guildID, settings.State); err != nil {
		return err
	}

	if !checkEmojis(settings.EmotesIncrease) || !checkEmojis(settings.EmotesDecrease) {
		return fiber.NewError(fiber.StatusBadRequest, "invalid emoji")
	}

	emotesInc := strings.Join(settings.EmotesIncrease, "")
	emotesDec := strings.Join(settings.EmotesDecrease, "")
	if err = c.db.SetKarmaEmotes(guildID, emotesInc, emotesDec); err != nil {
		return err
	}

	if err = c.db.SetKarmaTokens(guildID, settings.Tokens); err != nil {
		return err
	}

	fmt.Printf("%+v\n", settings)
	if err = c.db.SetKarmaPenalty(guildID, settings.Penalty); err != nil {
		return err
	}

	return ctx.JSON(models.Ok)
}

func (c *GuildsController) getGuildSettingsKarmaBlocklist(ctx *fiber.Ctx) error {
	guildID := ctx.Params("guildid")

	idList, err := c.db.GetKarmaBlockList(guildID)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return err
	}

	memberList := make([]*models.Member, len(idList))
	var m *discordgo.Member
	var i int
	for _, id := range idList {
		if m, err = c.state.Member(guildID, id); err != nil {
			continue
		}
		memberList[i] = models.MemberFromMember(m)
		i++
	}

	memberList = memberList[:i]

	return ctx.JSON(&models.ListResponse{N: len(memberList), Data: memberList})
}

func (c *GuildsController) putGuildSettingsKarmaBlocklist(ctx *fiber.Ctx) error {
	guildID := ctx.Params("guildid")
	memberID := ctx.Params("memberid")

	memb, err := fetch.FetchMember(c.session, guildID, memberID)
	if err == fetch.ErrNotFound {
		return fiber.ErrNotFound
	}
	if err != nil {
		return err
	}

	ok, err := c.db.IsKarmaBlockListed(guildID, memb.User.ID)
	if err != nil {
		return err
	}
	if ok {
		return fiber.NewError(fiber.StatusBadRequest, "member is already blocklisted")
	}

	if err = c.db.AddKarmaBlockList(guildID, memb.User.ID); err != nil {
		return err
	}

	return ctx.JSON(models.Ok)
}

func (c *GuildsController) deleteGuildSettingsKarmaBlocklist(ctx *fiber.Ctx) error {
	guildID := ctx.Params("guildid")
	memberID := ctx.Params("memberid")

	ok, err := c.db.IsKarmaBlockListed(guildID, memberID)
	if err != nil {
		return err
	}
	if !ok {
		return fiber.NewError(fiber.StatusBadRequest, "member is not blocklisted")
	}

	if err = c.db.RemoveKarmaBlockList(guildID, memberID); err != nil {
		return err
	}

	return ctx.JSON(models.Ok)
}

func (c *GuildsController) getGuildSettingsAntiraid(ctx *fiber.Ctx) error {
	guildID := ctx.Params("guildid")

	settings := new(models.AntiraidSettings)

	var err error
	if settings.State, err = c.db.GetAntiraidState(guildID); err != nil && !database.IsErrDatabaseNotFound(err) {
		return err
	}

	if settings.RegenerationPeriod, err = c.db.GetAntiraidRegeneration(guildID); err != nil && !database.IsErrDatabaseNotFound(err) {
		return err
	}

	if settings.Burst, err = c.db.GetAntiraidBurst(guildID); err != nil && !database.IsErrDatabaseNotFound(err) {
		return err
	}

	return ctx.JSON(settings)
}

func (c *GuildsController) postGuildSettingsAntiraid(ctx *fiber.Ctx) error {
	guildID := ctx.Params("guildid")

	settings := new(models.AntiraidSettings)
	if err := ctx.BodyParser(settings); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if settings.RegenerationPeriod < 1 {
		return fiber.NewError(fiber.StatusBadRequest, "regeneration period must be larger than 0")
	}
	if settings.Burst < 1 {
		return fiber.NewError(fiber.StatusBadRequest, "burst must be larger than 0")
	}

	var err error

	if err = c.db.SetAntiraidState(guildID, settings.State); err != nil {
		return err
	}

	if err = c.db.SetAntiraidRegeneration(guildID, settings.RegenerationPeriod); err != nil {
		return err
	}

	if err = c.db.SetAntiraidBurst(guildID, settings.Burst); err != nil {
		return err
	}

	return ctx.JSON(models.Ok)
}

func (c *GuildsController) getGuildUnbanrequests(ctx *fiber.Ctx) error {
	guildID := ctx.Params("guildid")

	requests, err := c.db.GetGuildUnbanRequests(guildID)
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

func (c *GuildsController) getGuildUnbanrequest(ctx *fiber.Ctx) error {
	guildID := ctx.Params("guildid")
	id := ctx.Params("id")

	request, err := c.db.GetUnbanRequest(id)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return err
	}
	if request == nil || request.GuildID != guildID {
		return fiber.ErrNotFound
	}

	return ctx.JSON(request.Hydrate())
}

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
	if request == nil || request.GuildID != guildID {
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
	request.Processed = time.Now()
	request.ProcessedMessage = rUpdate.ProcessedMessage

	if err = c.db.UpdateUnbanRequest(request); err != nil {
		return err
	}

	if request.Status == sharedmodels.UnbanRequestStateAccepted {
		if err = c.session.GuildBanDelete(request.GuildID, request.UserID); err != nil {
			return err
		}
	}

	return ctx.JSON(request.Hydrate())
}

func (c *GuildsController) getGuildSettingsKarmaRules(ctx *fiber.Ctx) error {
	guildID := ctx.Params("guildid")

	rules, err := c.db.GetKarmaRules(guildID)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return err
	}

	return ctx.JSON(models.ListResponse{N: len(rules), Data: rules})
}

func (c *GuildsController) createGuildSettingsKrameRule(ctx *fiber.Ctx) error {
	guildID := ctx.Params("guildid")

	rule := new(sharedmodels.KarmaRule)
	if err := ctx.BodyParser(rule); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := rule.Validate(); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	rule.GuildID = guildID
	rule.ID = snowflakenodes.NodeKarmaRules.Generate()

	if rule.Action == sharedmodels.KarmaActionToggleRole {
		role, err := fetch.FetchRole(c.session, guildID, rule.Argument)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		rule.Argument = role.ID
	}

	sum := rule.CalculateChecksum()
	ok, err := c.db.CheckKarmaRule(guildID, sum)
	if err != nil {
		return err
	}
	if ok {
		return fiber.NewError(fiber.StatusBadRequest, "same rule already exists")
	}

	if err := c.db.AddOrUpdateKarmaRule(rule); err != nil {
		return err
	}

	return ctx.JSON(rule)
}

func (c *GuildsController) updateGuildSettingsKrameRule(ctx *fiber.Ctx) (err error) {
	guildID := ctx.Params("guildid")
	id := ctx.Params("id")

	rule := new(sharedmodels.KarmaRule)
	if err := ctx.BodyParser(rule); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := rule.Validate(); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	rule.GuildID = guildID
	rule.ID, err = snowflake.ParseString(id)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if rule.Action == sharedmodels.KarmaActionToggleRole {
		role, err := fetch.FetchRole(c.session, guildID, rule.Argument)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		rule.Argument = role.ID
	}

	sum := rule.CalculateChecksum()
	ok, err := c.db.CheckKarmaRule(guildID, sum)
	if err != nil {
		return err
	}
	if ok {
		return fiber.NewError(fiber.StatusBadRequest, "same rule already exists")
	}

	if err := c.db.AddOrUpdateKarmaRule(rule); err != nil {
		return err
	}

	return ctx.JSON(rule)
}

func (c *GuildsController) deleteGuildSettingsKrameRule(ctx *fiber.Ctx) error {
	guildID := ctx.Params("guildid")
	id := ctx.Params("id")

	sfId, err := snowflake.ParseString(id)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := c.db.RemoveKarmaRule(guildID, sfId); err != nil {
		return err
	}

	return ctx.JSON(models.Ok)
}

func (c *GuildsController) getGuildSettingsLogs(ctx *fiber.Ctx) error {
	guildID := ctx.Params("guildid")

	limit, err := wsutil.GetQueryInt(ctx, "limit", 50, 1, 1000)
	if err != nil {
		return err
	}
	offset, err := wsutil.GetQueryInt(ctx, "offset", 0, 0, 0)
	if err != nil {
		return err
	}
	severity, err := wsutil.GetQueryInt(ctx, "severity",
		int(sharedmodels.GLAll), int(sharedmodels.GLAll), int(sharedmodels.GLFatal))
	if err != nil {
		return err
	}

	res, err := c.db.GetGuildLogEntries(guildID, offset, limit, sharedmodels.GuildLogSeverity(severity))
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return err
	}

	return ctx.JSON(&models.ListResponse{N: len(res), Data: res})
}

func (c *GuildsController) getGuildSettingsLogsCount(ctx *fiber.Ctx) error {
	guildID := ctx.Params("guildid")

	severity, err := wsutil.GetQueryInt(ctx, "severity",
		int(sharedmodels.GLAll), int(sharedmodels.GLAll), int(sharedmodels.GLFatal))
	if err != nil {
		return err
	}

	res, err := c.db.GetGuildLogEntriesCount(guildID, sharedmodels.GuildLogSeverity(severity))
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return err
	}

	return ctx.JSON(&models.Count{Count: res})
}

func (c *GuildsController) getGuildSettingsLogsState(ctx *fiber.Ctx) error {
	guildID := ctx.Params("guildid")

	disabled, err := c.db.GetGuildLogDisable(guildID)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return err
	}

	return ctx.JSON(&models.State{
		State: !disabled,
	})
}

func (c *GuildsController) postGuildSettingsLogsState(ctx *fiber.Ctx) error {
	guildID := ctx.Params("guildid")

	state := new(models.State)
	if err := ctx.BodyParser(state); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	err := c.db.SetGuildLogDisable(guildID, !state.State)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return err
	}

	return ctx.JSON(state)
}

func (c *GuildsController) deleteGuildSettingsLogEntry(ctx *fiber.Ctx) (err error) {
	guildID := ctx.Params("guildid")
	id := ctx.Params("id")

	if id != "" {
		var ids snowflake.ID
		ids, err = snowflake.ParseString(id)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		err = c.db.DeleteLogEntry(guildID, ids)
	} else {
		err = c.db.DeleteLogEntries(guildID)
	}

	if database.IsErrDatabaseNotFound(err) {
		return fiber.ErrNotFound
	}
	if err != nil {
		return
	}

	return ctx.JSON(models.Ok)
}

func (c *GuildsController) postFlushGuildData(ctx *fiber.Ctx) (err error) {
	guildID := ctx.Params("guildid")

	timeoutKey := "GUILDFLUSH:" + guildID
	if reset, ok := c.kvc.Get(timeoutKey).(bool); reset && ok {
		return fiber.NewError(fiber.StatusTooManyRequests, "this action can only be performed every 24 hours")
	}

	guild, err := c.state.Guild(guildID)
	if err != nil {
		return
	}

	payload := struct {
		Validation string `json:"validation"`
		LeaveAfter bool   `json:"leave_after"`
	}{}

	if err = ctx.BodyParser(&payload); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if payload.Validation != guild.Name {
		return fiber.NewError(fiber.StatusBadRequest, "invalid validation")
	}

	if err = util.FlushAllGuildData(c.db, c.st, guildID); err != nil {
		return
	}

	if payload.LeaveAfter {
		if err = c.session.GuildLeave(guildID); err != nil {
			return
		}
	}

	c.kvc.Set(timeoutKey, true, 24*time.Hour)

	return ctx.JSON(models.Ok)
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
