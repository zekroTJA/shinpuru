package controllers

import (
	"crypto"
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/bwmarrin/snowflake"
	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
	sharedmodels "github.com/zekroTJA/shinpuru/internal/models"
	"github.com/zekroTJA/shinpuru/internal/services/codeexec"
	"github.com/zekroTJA/shinpuru/internal/services/config"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/services/kvcache"
	permservice "github.com/zekroTJA/shinpuru/internal/services/permissions"
	"github.com/zekroTJA/shinpuru/internal/services/storage"
	"github.com/zekroTJA/shinpuru/internal/services/verification"
	"github.com/zekroTJA/shinpuru/internal/services/webserver/v1/models"
	"github.com/zekroTJA/shinpuru/internal/services/webserver/wsutil"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/snowflakenodes"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/fetch"
	"github.com/zekroTJA/shinpuru/pkg/hashutil"
	"github.com/zekroTJA/shinpuru/pkg/jdoodle"
	"github.com/zekroTJA/shinpuru/pkg/stringutil"
	"github.com/zekrotja/dgrs"
)

type GuildsSettingsController struct {
	db      database.Database
	st      storage.Storage
	kvc     kvcache.Provider
	session *discordgo.Session
	cfg     config.Provider
	pmw     *permservice.Permissions
	state   *dgrs.State
	vs      verification.Provider
	cef     codeexec.Factory
}

func (c *GuildsSettingsController) Setup(container di.Container, router fiber.Router) {
	c.session = container.Get(static.DiDiscordSession).(*discordgo.Session)
	c.cfg = container.Get(static.DiConfig).(config.Provider)
	c.db = container.Get(static.DiDatabase).(database.Database)
	c.pmw = container.Get(static.DiPermissions).(*permservice.Permissions)
	c.kvc = container.Get(static.DiKVCache).(kvcache.Provider)
	c.st = container.Get(static.DiObjectStorage).(storage.Storage)
	c.state = container.Get(static.DiState).(*dgrs.State)
	c.vs = container.Get(static.DiVerification).(verification.Provider)
	c.cef = container.Get(static.DiCodeExecFactory).(codeexec.Factory)

	router.Get("", c.getGuildSettings)
	router.Post("", c.postGuildSettings)
	router.Get("/karma", c.pmw.HandleWs(c.session, "sp.guild.config.karma"), c.getGuildSettingsKarma)
	router.Post("/karma", c.pmw.HandleWs(c.session, "sp.guild.config.karma"), c.postGuildSettingsKarma)
	router.Get("/karma/blocklist", c.pmw.HandleWs(c.session, "sp.guild.config.karma"), c.getGuildSettingsKarmaBlocklist)
	router.Put("/karma/blocklist/:memberid", c.pmw.HandleWs(c.session, "sp.guild.config.karma"), c.putGuildSettingsKarmaBlocklist)
	router.Delete("/karma/blocklist/:memberid", c.pmw.HandleWs(c.session, "sp.guild.config.karma"), c.deleteGuildSettingsKarmaBlocklist)
	router.Get("/karma/rules", c.pmw.HandleWs(c.session, "sp.guild.config.karma"), c.getGuildSettingsKarmaRules)
	router.Post("/karma/rules", c.pmw.HandleWs(c.session, "sp.guild.config.karma"), c.createGuildSettingsKrameRule)
	router.Post("/karma/rules/:id", c.pmw.HandleWs(c.session, "sp.guild.config.karma"), c.updateGuildSettingsKrameRule)
	router.Delete("/karma/rules/:id", c.pmw.HandleWs(c.session, "sp.guild.config.karma"), c.deleteGuildSettingsKrameRule)
	router.Get("/antiraid", c.pmw.HandleWs(c.session, "sp.guild.config.antiraid"), c.getGuildSettingsAntiraid)
	router.Post("/antiraid", c.pmw.HandleWs(c.session, "sp.guild.config.antiraid"), c.postGuildSettingsAntiraid)
	router.Post("/antiraid/action", c.pmw.HandleWs(c.session, "sp.guild.config.antiraid"), c.postGuildSettingsAntiraidAction)
	router.Get("/logs", c.pmw.HandleWs(c.session, "sp.guild.config.logs"), c.getGuildSettingsLogs)
	router.Get("/logs/count", c.pmw.HandleWs(c.session, "sp.guild.config.logs"), c.getGuildSettingsLogsCount)
	router.Delete("/logs", c.pmw.HandleWs(c.session, "sp.guild.config.logs"), c.deleteGuildSettingsLogEntries)
	router.Delete("/logs/:id", c.pmw.HandleWs(c.session, "sp.guild.config.logs"), c.deleteGuildSettingsLogEntries)
	router.Get("/logs/state", c.pmw.HandleWs(c.session, "sp.guild.config.logs"), c.getGuildSettingsLogsState)
	router.Post("/logs/state", c.pmw.HandleWs(c.session, "sp.guild.config.logs"), c.postGuildSettingsLogsState)
	router.Post("/flushguilddata", c.pmw.HandleWs(c.session, "sp.guild.admin.flushdata"), c.postFlushGuildData)
	router.Get("/api", c.pmw.HandleWs(c.session, "sp.guild.config.api"), c.getGuildSettingsAPI)
	router.Post("/api", c.pmw.HandleWs(c.session, "sp.guild.config.api"), c.postGuildSettingsAPI)
	router.Get("/verification", c.pmw.HandleWs(c.session, "sp.guild.config.verification"), c.getGuildSettingsVerification)
	router.Post("/verification", c.pmw.HandleWs(c.session, "sp.guild.config.verification"), c.postGuildSettingsVerification)
	router.Get("/codeexec", c.pmw.HandleWs(c.session, "sp.guild.config.exec"), c.getGuildSettingsCodeExec)
	router.Post("/codeexec", c.pmw.HandleWs(c.session, "sp.guild.config.exec"), c.postGuildSettingsCodeExec)
}

// @Summary Get Guild Settings
// @Description Returns the specified general guild settings.
// @Tags Guild Settings
// @Accept json
// @Produce json
// @Param id path string true "The ID of the guild."
// @Success 200 {object} models.GuildSettings
// @Failure 401 {object} models.Error
// @Failure 404 {object} models.Error
// @Router /guilds/{id}/settings [get]
func (c *GuildsSettingsController) getGuildSettings(ctx *fiber.Ctx) error {
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

	if gs.ModNotChannel, err = c.db.GetGuildModNot(guildID); err != nil && !database.IsErrDatabaseNotFound(err) {
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

// @Summary Get Guild Settings
// @Description Returns the specified general guild settings.
// @Tags Guild Settings
// @Accept json
// @Produce json
// @Param id path string true "The ID of the guild."
// @Param payload body models.GuildSettings true "Modified guild settings payload."
// @Success 200 {object} models.Status
// @Failure 400 {object} models.Error
// @Failure 401 {object} models.Error
// @Failure 404 {object} models.Error
// @Router /guilds/{id}/settings [post]
func (c *GuildsSettingsController) postGuildSettings(ctx *fiber.Ctx) error {
	uid := ctx.Locals("uid").(string)

	guildID := ctx.Params("guildid")

	var err error

	gs := new(models.GuildSettings)
	if err = ctx.BodyParser(gs); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// TODO: Change `fiber.ErrUnauthorized` to `fiber.ErrForbidden` 👇

	if gs.AutoRoles != nil {
		if ok, _, err := c.pmw.CheckPermissions(c.session, guildID, uid, "sp.guild.config.autorole"); err != nil {
			return wsutil.ErrInternalOrNotFound(err)
		} else if !ok {
			return fiber.ErrUnauthorized
		}

		if stringutil.ContainsAny("@everyone", gs.AutoRoles) {
			return fiber.NewError(fiber.StatusBadRequest,
				"@everyone can not be set as autorole")
		}

		guildRoles, err := c.state.Roles(guildID, true)
		if err != nil {
			return err
		}
		guildRoleIDs := make([]string, len(guildRoles))
		for i, role := range guildRoles {
			guildRoleIDs[i] = role.ID
		}

		if nc := stringutil.NotContained(gs.AutoRoles, guildRoleIDs); len(nc) > 0 {
			return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf(
				"Following RoleIDs are not existent on this guild: [%s]", strings.Join(nc, ", ")))
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

	if gs.ModNotChannel != "" {
		if ok, _, err := c.pmw.CheckPermissions(c.session, guildID, uid, "sp.guild.config.modnot"); err != nil {
			return wsutil.ErrInternalOrNotFound(err)
		} else if !ok {
			return fiber.ErrForbidden
		}

		if gs.ModNotChannel == "__RESET__" {
			gs.ModNotChannel = ""
		}

		if err = c.db.SetGuildModNot(guildID, gs.ModNotChannel); err != nil {
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
		if ok, _, err := c.pmw.CheckPermissions(c.session, guildID, uid, "sp.guild.config.announcements"); err != nil {
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
		if ok, _, err := c.pmw.CheckPermissions(c.session, guildID, uid, "sp.guild.config.announcements"); err != nil {
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

// @Summary Get Guild Karma Settings
// @Description Returns the specified guild karma settings.
// @Tags Guild Settings
// @Accept json
// @Produce json
// @Param id path string true "The ID of the guild."
// @Success 200 {object} models.KarmaSettings
// @Failure 401 {object} models.Error
// @Failure 404 {object} models.Error
// @Router /guilds/{id}/settings/karma [get]
func (c *GuildsSettingsController) getGuildSettingsKarma(ctx *fiber.Ctx) error {
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

// @Summary Update Guild Karma Settings
// @Description Update the guild karma settings specification.
// @Tags Guild Settings
// @Accept json
// @Produce json
// @Param id path string true "The ID of the guild."
// @Param payload body models.KarmaSettings true "The guild karma settings payload."
// @Success 200 {object} models.Status
// @Failure 400 {object} models.Error
// @Failure 401 {object} models.Error
// @Failure 404 {object} models.Error
// @Router /guilds/{id}/settings/karma [post]
func (c *GuildsSettingsController) postGuildSettingsKarma(ctx *fiber.Ctx) error {
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

	if err = c.db.SetKarmaPenalty(guildID, settings.Penalty); err != nil {
		return err
	}

	return ctx.JSON(models.Ok)
}

// @Summary Get Guild Karma Blocklist
// @Description Returns the specified guild karma blocklist entries.
// @Tags Guild Settings
// @Accept json
// @Produce json
// @Param id path string true "The ID of the guild."
// @Success 200 {array} models.Member "Wrapped in models.ListResponse"
// @Failure 401 {object} models.Error
// @Failure 404 {object} models.Error
// @Router /guilds/{id}/settings/karma/blocklist [get]
func (c *GuildsSettingsController) getGuildSettingsKarmaBlocklist(ctx *fiber.Ctx) error {
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

	return ctx.JSON(models.NewListResponse(memberList))
}

// @Summary Add Guild Karma Blocklist Entry
// @Description Add a guild karma blocklist entry.
// @Tags Guild Settings
// @Accept json
// @Produce json
// @Param id path string true "The ID of the guild."
// @Param memberid path string true "The ID of the guild."
// @Success 200 {object} models.Member
// @Failure 400 {object} models.Error
// @Failure 401 {object} models.Error
// @Failure 404 {object} models.Error
// @Router /guilds/{id}/settings/karma/blocklist/{memberid} [put]
func (c *GuildsSettingsController) putGuildSettingsKarmaBlocklist(ctx *fiber.Ctx) error {
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

	return ctx.JSON(memb)
}

// @Summary Remove Guild Karma Blocklist Entry
// @Description Remove a guild karma blocklist entry.
// @Tags Guild Settings
// @Accept json
// @Produce json
// @Param id path string true "The ID of the guild."
// @Param memberid path string true "The ID of the guild."
// @Success 200 {object} models.Status
// @Failure 400 {object} models.Error
// @Failure 401 {object} models.Error
// @Failure 404 {object} models.Error
// @Router /guilds/{id}/settings/karma/blocklist/{memberid} [delete]
func (c *GuildsSettingsController) deleteGuildSettingsKarmaBlocklist(ctx *fiber.Ctx) error {
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

// @Summary Get Guild Antiraid Settings
// @Description Returns the specified guild antiraid settings.
// @Tags Guild Settings
// @Accept json
// @Produce json
// @Param id path string true "The ID of the guild."
// @Success 200 {object} models.AntiraidSettings
// @Failure 401 {object} models.Error
// @Failure 404 {object} models.Error
// @Router /guilds/{id}/settings/antiraid [get]
func (c *GuildsSettingsController) getGuildSettingsAntiraid(ctx *fiber.Ctx) error {
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

	if settings.Verification, err = c.db.GetAntiraidVerification(guildID); err != nil && !database.IsErrDatabaseNotFound(err) {
		return err
	}

	return ctx.JSON(settings)
}

// @Summary Update Guild Antiraid Settings
// @Description Update the guild antiraid settings specification.
// @Tags Guild Settings
// @Accept json
// @Produce json
// @Param id path string true "The ID of the guild."
// @Param payload body models.AntiraidSettings true "The guild antiraid settings payload."
// @Success 200 {object} models.Status
// @Failure 400 {object} models.Error
// @Failure 401 {object} models.Error
// @Failure 404 {object} models.Error
// @Router /guilds/{id}/settings/antiraid [post]
func (c *GuildsSettingsController) postGuildSettingsAntiraid(ctx *fiber.Ctx) error {
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

	if err = c.db.SetAntiraidVerification(guildID, settings.Verification); err != nil {
		return err
	}

	return ctx.JSON(models.Ok)
}

// @Summary Guild Antiraid Bulk Action
// @Description Execute a specific action on antiraid listed users
// @Tags Guild Settings
// @Accept json
// @Produce json
// @Param id path string true "The ID of the guild."
// @Param payload body models.AntiraidAction true "The antiraid action payload."
// @Success 200 {object} models.Status
// @Failure 400 {object} models.Error
// @Failure 401 {object} models.Error
// @Failure 404 {object} models.Error
// @Router /guilds/{id}/settings/antiraid/action [post]
func (c *GuildsSettingsController) postGuildSettingsAntiraidAction(ctx *fiber.Ctx) (err error) {
	guildID := ctx.Params("guildid")

	var action models.AntiraidAction
	if err = ctx.BodyParser(&action); err != nil {
		return
	}

	var actF func(id string) error
	switch action.Type {
	case models.AntiraidActionTypeKick:
		actF = func(id string) error {
			return c.session.GuildMemberDelete(guildID, id)
		}
	case models.AntiraidActionTypeBan:
		actF = func(id string) error {
			return c.session.GuildBanCreateWithReason(guildID, id, "antiraid purge", 7)
		}
	default:
		return fiber.NewError(fiber.StatusBadRequest, "invalid action type")
	}

	if len(action.IDs) == 0 {
		return
	}

	joinList, err := c.db.GetAntiraidJoinList(guildID)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return fiber.NewError(fiber.StatusBadRequest, "ID list must contain entries")
	}

	var contained int
	for _, e := range joinList {
	inner:
		for _, id := range action.IDs {
			if e.UserID == id {
				contained++
				break inner
			}
		}
	}
	if contained != len(action.IDs) {
		return fiber.NewError(fiber.StatusBadRequest, "ID list contains entry not contained in antiraid joinlist")
	}

	for _, id := range action.IDs {
		if err = actF(id); err != nil {
			return
		}
		if err = c.db.RemoveAntiraidJoinList(guildID, id); err != nil {
			return
		}
	}

	return ctx.JSON(models.Ok)
}

// @Summary Get Guild Settings Karma Rules
// @Description Returns a list of specified guild karma rules.
// @Tags Guild Settings
// @Accept json
// @Produce json
// @Param id path string true "The ID of the guild."
// @Success 200 {array} sharedmodels.KarmaRule "Wrapped in models.ListResponse"
// @Failure 401 {object} models.Error
// @Failure 404 {object} models.Error
// @Router /guilds/{id}/settings/karma/rules [get]
func (c *GuildsSettingsController) getGuildSettingsKarmaRules(ctx *fiber.Ctx) error {
	guildID := ctx.Params("guildid")

	rules, err := c.db.GetKarmaRules(guildID)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return err
	}

	return ctx.JSON(models.NewListResponse(rules))
}

// @Summary Create Guild Settings Karma
// @Description Create a guild karma rule.
// @Tags Guild Settings
// @Accept json
// @Produce json
// @Param id path string true "The ID of the guild."
// @Param payload body sharedmodels.KarmaRule true "The karma rule payload."
// @Success 200 {object} sharedmodels.KarmaRule
// @Failure 400 {object} models.Error
// @Failure 401 {object} models.Error
// @Failure 404 {object} models.Error
// @Router /guilds/{id}/settings/karma/rules [post]
func (c *GuildsSettingsController) createGuildSettingsKrameRule(ctx *fiber.Ctx) error {
	guildID := ctx.Params("guildid")

	var rule sharedmodels.KarmaRule
	if err := ctx.BodyParser(&rule); err != nil {
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

// @Summary Update Guild Settings Karma
// @Description Update a karma rule by ID.
// @Tags Guild Settings
// @Accept json
// @Produce json
// @Param id path string true "The ID of the guild."
// @Param ruleid path string true "The ID of the rule."
// @Param payload body sharedmodels.KarmaRule true "The karma rule update payload."
// @Success 200 {object} sharedmodels.KarmaRule
// @Failure 400 {object} models.Error
// @Failure 401 {object} models.Error
// @Failure 404 {object} models.Error
// @Router /guilds/{id}/settings/karma/rules/{ruleid} [post]
func (c *GuildsSettingsController) updateGuildSettingsKrameRule(ctx *fiber.Ctx) (err error) {
	guildID := ctx.Params("guildid")
	id := ctx.Params("id")

	var rule sharedmodels.KarmaRule
	if err := ctx.BodyParser(&rule); err != nil {
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

// @Summary Remove Guild Settings Karma
// @Description Remove a guild karma rule by ID.
// @Tags Guild Settings
// @Accept json
// @Produce json
// @Param id path string true "The ID of the guild."
// @Param ruleid path string true "The ID of the rule."
// @Success 200 {object} models.State
// @Failure 400 {object} models.Error
// @Failure 401 {object} models.Error
// @Failure 404 {object} models.Error
// @Router /guilds/{id}/settings/karma/rules/{ruleid} [delete]
func (c *GuildsSettingsController) deleteGuildSettingsKrameRule(ctx *fiber.Ctx) error {
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

// @Summary Get Guild Log
// @Description Returns a list of entries of the guild log.
// @Tags Guild Settings
// @Accept json
// @Produce json
// @Param id path string true "The ID of the guild."
// @Param limit query int false "The amount of values returned." default(50) minimum(1) maximum(1000)
// @Param offset query int false "The amount of values to be skipped." default(0)
// @Param severity query sharedmodels.GuildLogSeverity false "Filter by log severity." default(sharedmodels.GLAll)
// @Success 200 {array} sharedmodels.GuildLogEntry "Wrapped in models.ListResponse"
// @Failure 400 {object} models.Error
// @Failure 401 {object} models.Error
// @Failure 404 {object} models.Error
// @Router /guilds/{id}/settings/logs [get]
func (c *GuildsSettingsController) getGuildSettingsLogs(ctx *fiber.Ctx) error {
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
	order := ctx.Query("order", "desc")
	ascending := order == "asc"

	res, err := c.db.GetGuildLogEntries(
		guildID, offset, limit, sharedmodels.GuildLogSeverity(severity), ascending)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return err
	}

	return ctx.JSON(models.NewListResponse(res))
}

// @Summary Get Guild Log Count
// @Description Returns the total or filtered count of guild log entries.
// @Tags Guild Settings
// @Accept json
// @Produce json
// @Param id path string true "The ID of the guild."
// @Param severity query sharedmodels.GuildLogSeverity false "Filter by log severity." default(sharedmodels.GLAll)
// @Success 200 {object} models.Count
// @Failure 400 {object} models.Error
// @Failure 401 {object} models.Error
// @Failure 404 {object} models.Error
// @Router /guilds/{id}/settings/logs [get]
func (c *GuildsSettingsController) getGuildSettingsLogsCount(ctx *fiber.Ctx) error {
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

// @Summary Get Guild Settings Log State
// @Description Returns the enabled state of the guild log setting.
// @Tags Guild Settings
// @Accept json
// @Produce json
// @Param id path string true "The ID of the guild."
// @Success 200 {object} models.State
// @Failure 401 {object} models.Error
// @Failure 404 {object} models.Error
// @Router /guilds/{id}/settings/logs/state [get]
func (c *GuildsSettingsController) getGuildSettingsLogsState(ctx *fiber.Ctx) error {
	guildID := ctx.Params("guildid")

	disabled, err := c.db.GetGuildLogDisable(guildID)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return err
	}

	return ctx.JSON(&models.State{
		State: !disabled,
	})
}

// @Summary Update Guild Settings Log State
// @Description Update the enabled state of the log state guild setting.
// @Tags Guild Settings
// @Accept json
// @Produce json
// @Param id path string true "The ID of the guild."
// @Param payload body models.State true "The state payload."
// @Success 200 {object} models.State
// @Failure 400 {object} models.Error
// @Failure 401 {object} models.Error
// @Failure 404 {object} models.Error
// @Router /guilds/{id}/settings/logs/state [post]
func (c *GuildsSettingsController) postGuildSettingsLogsState(ctx *fiber.Ctx) error {
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

// @Summary Delete Guild Log Entries
// @Description Delete all guild log entries.
// @Tags Guild Settings
// @Accept json
// @Produce json
// @Param id path string true "The ID of the guild."
// @Success 200 {object} models.State
// @Failure 400 {object} models.Error
// @Failure 401 {object} models.Error
// @Failure 404 {object} models.Error
// @Router /guilds/{id}/settings/logs [delete]
//
// This is a dummy method for API doc generation.
func (*GuildsSettingsController) _(*fiber.Ctx) error {
	return nil
}

// @Summary Delete Guild Log Entries
// @Description Delete a single log entry.
// @Tags Guild Settings
// @Accept json
// @Produce json
// @Param id path string true "The ID of the guild."
// @Param entryid path string true "The ID of the entry to be deleted."
// @Success 200 {object} models.State
// @Failure 400 {object} models.Error
// @Failure 401 {object} models.Error
// @Failure 404 {object} models.Error
// @Router /guilds/{id}/settings/logs/{entryid} [delete]
func (c *GuildsSettingsController) deleteGuildSettingsLogEntries(ctx *fiber.Ctx) (err error) {
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

// @Summary Flush Guild Data
// @Description Flushes all guild data from the database.
// @Tags Guild Settings
// @Accept json
// @Produce json
// @Param id path string true "The ID of the guild."
// @Param payload body models.FlushGuildRequest true "The guild flush payload."
// @Success 200 {object} models.State
// @Failure 400 {object} models.Error
// @Failure 401 {object} models.Error
// @Failure 404 {object} models.Error
// @Router /guilds/{id}/settings/flushguilddata [post]
func (c *GuildsSettingsController) postFlushGuildData(ctx *fiber.Ctx) (err error) {
	guildID := ctx.Params("guildid")

	timeoutKey := "GUILDFLUSH:" + guildID
	if reset, ok := c.kvc.Get(timeoutKey).(bool); reset && ok {
		return fiber.NewError(fiber.StatusTooManyRequests, "this action can only be performed every 24 hours")
	}

	guild, err := c.state.Guild(guildID)
	if err != nil {
		return
	}

	var payload models.FlushGuildRequest
	if err = ctx.BodyParser(&payload); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if payload.Validation != guild.Name {
		return fiber.NewError(fiber.StatusBadRequest, "invalid validation")
	}

	if err = util.FlushAllGuildData(c.session, c.db, c.st, c.state, guildID); err != nil {
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

// @Summary Get Guild Settings API State
// @Description Returns the settings state of the Guild API.
// @Tags Guild Settings
// @Accept json
// @Produce json
// @Param id path string true "The ID of the guild."
// @Success 200 {object} sharedmodels.GuildAPISettings
// @Failure 401 {object} models.Error
// @Failure 404 {object} models.Error
// @Router /guilds/{id}/settings/api [get]
func (c *GuildsSettingsController) getGuildSettingsAPI(ctx *fiber.Ctx) error {
	guildID := ctx.Params("guildid")

	state, err := c.db.GetGuildAPI(guildID)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return err
	}

	state.Hydrate()
	state.TokenHash = ""
	return ctx.JSON(state)
}

// @Summary Set Guild Settings API State
// @Description Set the settings state of the Guild API.
// @Tags Guild Settings
// @Accept json
// @Produce json
// @Param id path string true "The ID of the guild."
// @Param payload body models.GuildAPISettingsRequest true "The guild API settings payload."
// @Success 200 {object} sharedmodels.GuildAPISettings
// @Failure 401 {object} models.Error
// @Failure 404 {object} models.Error
// @Router /guilds/{id}/settings/api [post]
func (c *GuildsSettingsController) postGuildSettingsAPI(ctx *fiber.Ctx) (err error) {
	guildID := ctx.Params("guildid")

	state, err := c.db.GetGuildAPI(guildID)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return err
	}

	newState := new(models.GuildAPISettingsRequest)
	if err = ctx.BodyParser(newState); err != nil {
		return
	}

	newState.TokenHash = state.TokenHash

	if newState.ResetToken {
		newState.TokenHash = ""
	} else if newState.NewToken != "" {
		hasher := hashutil.Hasher{HashFunc: crypto.SHA512, SaltSize: 128}
		newState.TokenHash, err = hasher.Hash(newState.NewToken)
	}

	if err = c.db.SetGuildAPI(guildID, newState.GuildAPISettings); err != nil {
		return
	}

	newState.Hydrate()
	newState.TokenHash = ""
	return ctx.JSON(newState.GuildAPISettings)
}

// @Summary Get Guild Settings Verification State
// @Description Returns the settings state of the Guild Verification.
// @Tags Guild Settings
// @Accept json
// @Produce json
// @Param id path string true "The ID of the guild."
// @Success 200 {object} models.EnableStatus
// @Failure 401 {object} models.Error
// @Failure 404 {object} models.Error
// @Router /guilds/{id}/settings/verification [get]
func (c *GuildsSettingsController) getGuildSettingsVerification(ctx *fiber.Ctx) error {
	guildID := ctx.Params("guildid")

	state, err := c.vs.GetEnabled(guildID)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return err
	}

	res := models.EnableStatus{
		Enabled: state,
	}

	return ctx.JSON(res)
}

// @Summary Set Guild Settings Verification State
// @Description Set the settings state of the Guild Verification.
// @Tags Guild Settings
// @Accept json
// @Produce json
// @Param id path string true "The ID of the guild."
// @Param payload body models.EnableStatus true "The guild API settings payload."
// @Success 200 {object} models.EnableStatus
// @Failure 401 {object} models.Error
// @Failure 404 {object} models.Error
// @Router /guilds/{id}/settings/verification [post]
func (c *GuildsSettingsController) postGuildSettingsVerification(ctx *fiber.Ctx) (err error) {
	guildID := ctx.Params("guildid")

	var state models.EnableStatus
	if err = ctx.BodyParser(&state); err != nil {
		return
	}

	err = c.vs.SetEnabled(guildID, state.Enabled)
	if err != nil {
		return
	}

	return ctx.JSON(state)
}

// @Summary Get Guild Settings Code Exec State
// @Description Returns the settings state of the Guild Code Exec.
// @Tags Guild Settings
// @Accept json
// @Produce json
// @Param id path string true "The ID of the guild."
// @Success 200 {object} models.EnableStatus
// @Failure 401 {object} models.Error
// @Failure 404 {object} models.Error
// @Router /guilds/{id}/settings/codeexec [get]
func (c *GuildsSettingsController) getGuildSettingsCodeExec(ctx *fiber.Ctx) error {
	guildID := ctx.Params("guildid")

	var (
		res models.CodeExecSettings
		err error
	)

	res.Enabled, err = c.db.GetGuildCodeExecEnabled(guildID)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return err
	}

	res.Type = c.cef.Name()

	if res.Type == "jdoodle" {
		creds, err := c.db.GetGuildJdoodleKey(guildID)
		if err != nil && !database.IsErrDatabaseNotFound(err) {
			return err
		}
		credsSplit := strings.Split(creds, "#")
		if len(credsSplit) == 2 {
			res.JdoodleClientId = credsSplit[0]
			res.JdoodleClientSecret = credsSplit[1]
		}
	}

	return ctx.JSON(res)
}

// @Summary Set Guild Settings Code Exec State
// @Description Set the settings state of the Guild Code Exec.
// @Tags Guild Settings
// @Accept json
// @Produce json
// @Param id path string true "The ID of the guild."
// @Param payload body models.EnableStatus true "The guild API settings payload."
// @Success 200 {object} models.EnableStatus
// @Failure 401 {object} models.Error
// @Failure 404 {object} models.Error
// @Router /guilds/{id}/settings/codeexec [post]
func (c *GuildsSettingsController) postGuildSettingsCodeExec(ctx *fiber.Ctx) (err error) {
	guildID := ctx.Params("guildid")

	var state models.CodeExecSettings
	if err = ctx.BodyParser(&state); err != nil {
		return
	}

	err = c.db.SetGuildCodeExecEnabled(guildID, state.Enabled)
	if err != nil {
		return
	}

	if c.cef.Name() == "jdoodle" {
		var creds string
		if state.JdoodleClientId == "" && state.JdoodleClientSecret == "" {
		} else if state.JdoodleClientId != "" && state.JdoodleClientSecret != "" {
			_, err = jdoodle.NewWrapper(state.JdoodleClientId, state.JdoodleClientSecret).CreditsSpent()
			if err != nil {
				return fiber.NewError(fiber.StatusBadRequest, "The JDoodle credentials are invalid.")
			}
			creds = fmt.Sprintf("%s#%s", state.JdoodleClientId, state.JdoodleClientSecret)
		} else {
			return fiber.NewError(fiber.StatusBadRequest, "Either both credential values must be empty or both must be defined!")
		}

		err = c.db.SetGuildJdoodleKey(guildID, creds)
		if err != nil {
			return
		}
	}

	return ctx.JSON(state)
}
