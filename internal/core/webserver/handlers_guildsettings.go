package webserver

import (
	"errors"
	"fmt"
	"strings"

	"github.com/makeworld-the-better-one/go-isemoji"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"github.com/zekroTJA/shinpuru/internal/core/database"
	"github.com/zekroTJA/shinpuru/pkg/permissions"
)

// ---------------------------------------------------------------------------
// - GET /api/guilds/:guildid/settings

func (ws *WebServer) handlerGetGuildSettings(ctx *routing.Context) error {
	gs := new(GuildSettings)

	guildID := ctx.Param("guildid")

	var err error

	if gs.Prefix, err = ws.db.GetGuildPrefix(guildID); err != nil && !database.IsErrDatabaseNotFound(err) {
		return errInternalOrNotFound(ctx, err)
	}

	if gs.Perms, err = ws.db.GetGuildPermissions(guildID); err != nil && !database.IsErrDatabaseNotFound(err) {
		return errInternalOrNotFound(ctx, err)
	}

	if gs.AutoRole, err = ws.db.GetGuildAutoRole(guildID); err != nil && !database.IsErrDatabaseNotFound(err) {
		return errInternalOrNotFound(ctx, err)
	}

	if gs.ModLogChannel, err = ws.db.GetGuildModLog(guildID); err != nil && !database.IsErrDatabaseNotFound(err) {
		return errInternalOrNotFound(ctx, err)
	}

	if gs.VoiceLogChannel, err = ws.db.GetGuildVoiceLog(guildID); err != nil && !database.IsErrDatabaseNotFound(err) {
		return errInternalOrNotFound(ctx, err)
	}

	if gs.JoinMessageChannel, gs.JoinMessageText, err = ws.db.GetGuildJoinMsg(guildID); err != nil && !database.IsErrDatabaseNotFound(err) {
		return errInternalOrNotFound(ctx, err)
	}

	if gs.LeaveMessageChannel, gs.LeaveMessageText, err = ws.db.GetGuildLeaveMsg(guildID); err != nil && !database.IsErrDatabaseNotFound(err) {
		return errInternalOrNotFound(ctx, err)
	}

	return jsonResponse(ctx, gs, fasthttp.StatusOK)
}

// ---------------------------------------------------------------------------
// - POST /api/guilds/:guildid/settings

func (ws *WebServer) handlerPostGuildSettings(ctx *routing.Context) error {
	userID := ctx.Get("uid").(string)

	guildID := ctx.Param("guildid")

	var err error

	gs := new(GuildSettings)
	if err = parseJSONBody(ctx, gs); err != nil {
		return jsonError(ctx, err, fasthttp.StatusBadRequest)
	}

	if gs.AutoRole != "" {
		if ok, _, err := ws.pmw.CheckPermissions(ws.session, guildID, userID, "sp.guild.config.autorole"); err != nil {
			return errInternalOrNotFound(ctx, err)
		} else if !ok {
			return jsonError(ctx, errUnauthorized, fasthttp.StatusUnauthorized)
		}

		if gs.AutoRole == "__RESET__" {
			gs.AutoRole = ""
		}

		if err = ws.db.SetGuildAutoRole(guildID, gs.AutoRole); err != nil {
			return errInternalOrNotFound(ctx, err)
		}
	}

	if gs.ModLogChannel != "" {
		if ok, _, err := ws.pmw.CheckPermissions(ws.session, guildID, userID, "sp.guild.config.modlog"); err != nil {
			return errInternalOrNotFound(ctx, err)
		} else if !ok {
			return jsonError(ctx, errUnauthorized, fasthttp.StatusUnauthorized)
		}

		if gs.ModLogChannel == "__RESET__" {
			gs.ModLogChannel = ""
		}

		if err = ws.db.SetGuildModLog(guildID, gs.ModLogChannel); err != nil {
			return errInternalOrNotFound(ctx, err)
		}
	}

	if gs.Prefix != "" {
		if ok, _, err := ws.pmw.CheckPermissions(ws.session, guildID, userID, "sp.guild.config.prefix"); err != nil {
			return errInternalOrNotFound(ctx, err)
		} else if !ok {
			return jsonError(ctx, errUnauthorized, fasthttp.StatusUnauthorized)
		}

		if gs.Prefix == "__RESET__" {
			gs.Prefix = ""
		}

		if err = ws.db.SetGuildPrefix(guildID, gs.Prefix); err != nil {
			return errInternalOrNotFound(ctx, err)
		}
	}

	if gs.VoiceLogChannel != "" {
		if ok, _, err := ws.pmw.CheckPermissions(ws.session, guildID, userID, "sp.guild.config.voicelog"); err != nil {
			return errInternalOrNotFound(ctx, err)
		} else if !ok {
			return jsonError(ctx, errUnauthorized, fasthttp.StatusUnauthorized)
		}

		if gs.VoiceLogChannel == "__RESET__" {
			gs.VoiceLogChannel = ""
		}

		if err = ws.db.SetGuildVoiceLog(guildID, gs.VoiceLogChannel); err != nil {
			return errInternalOrNotFound(ctx, err)
		}
	}

	if gs.JoinMessageChannel != "" && gs.JoinMessageText != "" {
		if ok, _, err := ws.pmw.CheckPermissions(ws.session, guildID, userID, "sp.guild.config.joinmsg"); err != nil {
			return errInternalOrNotFound(ctx, err)
		} else if !ok {
			return jsonError(ctx, errUnauthorized, fasthttp.StatusUnauthorized)
		}

		if gs.JoinMessageChannel == "__RESET__" && gs.JoinMessageText == "__RESET__" {
			gs.JoinMessageChannel = ""
			gs.JoinMessageText = ""
		}

		if err = ws.db.SetGuildJoinMsg(guildID, gs.JoinMessageChannel, gs.JoinMessageText); err != nil {
			return errInternalOrNotFound(ctx, err)
		}
	}

	if gs.LeaveMessageChannel != "" && gs.LeaveMessageText != "" {
		if ok, _, err := ws.pmw.CheckPermissions(ws.session, guildID, userID, "sp.guild.config.leavemsg"); err != nil {
			return errInternalOrNotFound(ctx, err)
		} else if !ok {
			return jsonError(ctx, errUnauthorized, fasthttp.StatusUnauthorized)
		}

		if gs.LeaveMessageChannel == "__RESET__" && gs.LeaveMessageText == "__RESET__" {
			gs.LeaveMessageChannel = ""
			gs.LeaveMessageText = ""
		}

		if err = ws.db.SetGuildLeaveMsg(guildID, gs.LeaveMessageChannel, gs.LeaveMessageText); err != nil {
			return errInternalOrNotFound(ctx, err)
		}
	}

	return jsonResponse(ctx, nil, fasthttp.StatusOK)
}

// ---------------------------------------------------------------------------
// - GET /api/guilds/:guildid/permissions

func (ws *WebServer) handlerGetGuildPermissions(ctx *routing.Context) error {
	userID := ctx.Get("uid").(string)

	guildID := ctx.Param("guildid")

	if memb, _ := ws.session.GuildMember(guildID, userID); memb == nil {
		return jsonError(ctx, errNotFound, fasthttp.StatusNotFound)
	}

	var perms map[string]permissions.PermissionArray
	var err error

	if perms, err = ws.db.GetGuildPermissions(guildID); err != nil && !database.IsErrDatabaseNotFound(err) {
		return errInternalOrNotFound(ctx, err)
	}

	return jsonResponse(ctx, perms, fasthttp.StatusOK)
}

// ---------------------------------------------------------------------------
// - POST /api/guilds/:guildid/permissions

func (ws *WebServer) handlerPostGuildPermissions(ctx *routing.Context) error {
	userID := ctx.Get("uid").(string)

	guildID := ctx.Param("guildid")

	if ok, _, err := ws.pmw.CheckPermissions(ws.session, guildID, userID, "sp.guild.config.perms"); err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	} else if !ok {
		return jsonError(ctx, errUnauthorized, fasthttp.StatusUnauthorized)
	}

	update := new(PermissionsUpdate)
	if err := parseJSONBody(ctx, update); err != nil {
		return jsonError(ctx, errInvalidArguments, fasthttp.StatusBadRequest)
	}

	sperm := update.Perm[1:]
	if !strings.HasPrefix(sperm, "sp.guild") && !strings.HasPrefix(sperm, "sp.etc") && !strings.HasPrefix(sperm, "sp.chat") {
		return jsonError(ctx, fmt.Errorf("you can only give permissions over the domains 'sp.guild', 'sp.etc' and 'sp.chat'"), fasthttp.StatusBadRequest)
	}

	perms, err := ws.db.GetGuildPermissions(guildID)
	if err != nil {
		if database.IsErrDatabaseNotFound(err) {
			return jsonError(ctx, errNotFound, fasthttp.StatusNotFound)
		}
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	for _, roleID := range update.RoleIDs {
		rperms, ok := perms[roleID]
		if !ok {
			rperms = make(permissions.PermissionArray, 0)
		}

		rperms, changed := rperms.Update(update.Perm, false)

		if changed {
			if err = ws.db.SetGuildRolePermission(guildID, roleID, rperms); err != nil {
				return jsonError(ctx, err, fasthttp.StatusInternalServerError)
			}
		}
	}

	return jsonResponse(ctx, nil, fasthttp.StatusOK)
}

// ---------------------------------------------------------------------------
// - POST /api/guilds/:guildid/inviteblock

func (ws *WebServer) handlerPostGuildInviteBlock(ctx *routing.Context) error {
	userID := ctx.Get("uid").(string)

	guildID := ctx.Param("guildid")

	if ok, _, err := ws.pmw.CheckPermissions(ws.session, guildID, userID, "sp.guild.mod.inviteblock"); err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	} else if !ok {
		return jsonError(ctx, errUnauthorized, fasthttp.StatusUnauthorized)
	}

	var data struct {
		Enabled bool `json:"enabled"`
	}

	if err := parseJSONBody(ctx, &data); err != nil {
		return jsonError(ctx, err, fasthttp.StatusBadRequest)
	}

	val := ""
	if data.Enabled {
		val = "1"
	}

	if err := ws.db.SetGuildInviteBlock(guildID, val); err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	return jsonResponse(ctx, nil, fasthttp.StatusOK)
}

// ---------------------------------------------------------------------------
// - GET /api/guilds/:guildid/settings/karma

func (ws *WebServer) handlerGetGuildSettingsKarma(ctx *routing.Context) (err error) {
	userID := ctx.Get("uid").(string)

	guildID := ctx.Param("guildid")

	if ok, _, err := ws.pmw.CheckPermissions(ws.session, guildID, userID, "sp.guild.config.karma"); err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	} else if !ok {
		return jsonError(ctx, errUnauthorized, fasthttp.StatusUnauthorized)
	}

	settings := new(KarmaSettings)

	if settings.State, err = ws.db.GetKarmaState(guildID); err != nil && !database.IsErrDatabaseNotFound(err) {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	if settings.Tokens, err = ws.db.GetKarmaTokens(guildID); err != nil && !database.IsErrDatabaseNotFound(err) {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	emotesInc, emotesDec, err := ws.db.GetKarmaEmotes(guildID)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}
	settings.EmotesIncrease = strings.Split(emotesInc, "")
	settings.EmotesDecrease = strings.Split(emotesDec, "")

	return jsonResponse(ctx, settings, fasthttp.StatusOK)
}

// ---------------------------------------------------------------------------
// - POST /api/guilds/:guildid/settings/karma

func (ws *WebServer) handlerPostGuildSettingsKarma(ctx *routing.Context) (err error) {
	userID := ctx.Get("uid").(string)

	guildID := ctx.Param("guildid")

	if ok, _, err := ws.pmw.CheckPermissions(ws.session, guildID, userID, "sp.guild.config.karma"); err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	} else if !ok {
		return jsonError(ctx, errUnauthorized, fasthttp.StatusUnauthorized)
	}

	settings := new(KarmaSettings)
	if err := parseJSONBody(ctx, settings); err != nil {
		return jsonError(ctx, err, fasthttp.StatusBadRequest)
	}

	if err = ws.db.SetKarmaState(guildID, settings.State); err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	if !checkEmojis(settings.EmotesIncrease) || !checkEmojis(settings.EmotesDecrease) {
		return jsonError(ctx, errors.New("invalid emoji"), fasthttp.StatusBadRequest)
	}

	emotesInc := strings.Join(settings.EmotesIncrease, "")
	emotesDec := strings.Join(settings.EmotesDecrease, "")
	if err = ws.db.SetKarmaEmotes(guildID, emotesInc, emotesDec); err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	if err = ws.db.SetKarmaTokens(guildID, settings.Tokens); err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	return jsonResponse(ctx, nil, fasthttp.StatusOK)
}

// ---------------------------------------------------------------------------
// - GET /api/guilds/:guildid/settings/antiraid

func (ws *WebServer) handlerGetGuildSettingsAntiraid(ctx *routing.Context) (err error) {
	userID := ctx.Get("uid").(string)

	guildID := ctx.Param("guildid")

	if ok, _, err := ws.pmw.CheckPermissions(ws.session, guildID, userID, "sp.guild.config.antiraid"); err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	} else if !ok {
		return jsonError(ctx, errUnauthorized, fasthttp.StatusUnauthorized)
	}

	settings := new(AntiraidSettings)

	if settings.State, err = ws.db.GetAntiraidState(guildID); err != nil && !database.IsErrDatabaseNotFound(err) {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	if settings.RegenerationPeriod, err = ws.db.GetAntiraidRegeneration(guildID); err != nil && !database.IsErrDatabaseNotFound(err) {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	if settings.Burst, err = ws.db.GetAntiraidBurst(guildID); err != nil && !database.IsErrDatabaseNotFound(err) {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	return jsonResponse(ctx, settings, fasthttp.StatusOK)
}

// ---------------------------------------------------------------------------
// - POST /api/guilds/:guildid/settings/antiraid

func (ws *WebServer) handlerPostGuildSettingsAntiraid(ctx *routing.Context) (err error) {
	userID := ctx.Get("uid").(string)

	guildID := ctx.Param("guildid")

	if ok, _, err := ws.pmw.CheckPermissions(ws.session, guildID, userID, "sp.guild.config.antiraid"); err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	} else if !ok {
		return jsonError(ctx, errUnauthorized, fasthttp.StatusUnauthorized)
	}

	settings := new(AntiraidSettings)
	if err := parseJSONBody(ctx, settings); err != nil {
		return jsonError(ctx, err, fasthttp.StatusBadRequest)
	}

	if settings.RegenerationPeriod < 1 {
		return jsonError(ctx, errors.New("regeneration period must be larger than 0"), fasthttp.StatusBadRequest)
	}
	if settings.Burst < 1 {
		return jsonError(ctx, errors.New("burst must be larger than 0"), fasthttp.StatusBadRequest)
	}

	if err = ws.db.SetAntiraidState(guildID, settings.State); err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	if err = ws.db.SetAntiraidRegeneration(guildID, settings.RegenerationPeriod); err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	if err = ws.db.SetAntiraidBurst(guildID, settings.Burst); err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	return jsonResponse(ctx, nil, fasthttp.StatusOK)
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
