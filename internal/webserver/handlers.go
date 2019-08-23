package webserver

import (
	"fmt"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/zekroTJA/shinpuru/internal/shared"

	"github.com/zekroTJA/shinpuru/internal/core"

	"github.com/bwmarrin/discordgo"
	"github.com/bwmarrin/snowflake"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"github.com/zekroTJA/shinpuru/internal/util"
)

func (ws *WebServer) handlerGetMe(ctx *routing.Context) error {
	userID := ctx.Get("uid").(string)

	user, err := ws.session.User(userID)
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	created, _ := util.GetDiscordSnowflakeCreationTime(user.ID)

	res := &User{
		User:      user,
		AvatarURL: user.AvatarURL(""),
		CreatedAt: created,
		BotOwner:  userID == ws.config.Discord.OwnerID,
	}

	return jsonResponse(ctx, res, fasthttp.StatusOK)
}

func (ws *WebServer) handlerGuildsGet(ctx *routing.Context) error {
	userID := ctx.Get("uid").(string)

	guilds := make([]*GuildReduced, len(ws.session.State.Guilds))
	i := 0
	for _, g := range ws.session.State.Guilds {
		if g.MemberCount < 10000 {
			for _, m := range g.Members {
				if m.User.ID == userID {
					guilds[i] = GuildReducedFromGuild(g)
					i++
					break
				}
			}
		} else {
			if gm, _ := ws.session.GuildMember(g.ID, userID); gm != nil {
				guilds[i] = GuildReducedFromGuild(g)
				i++
			}
		}
	}
	guilds = guilds[:i]

	return jsonResponse(ctx, &ListResponse{
		N:    i,
		Data: guilds,
	}, fasthttp.StatusOK)
}

func (ws *WebServer) handlerGuildsGetGuild(ctx *routing.Context) error {
	userID := ctx.Get("uid").(string)

	guildID := ctx.Param("guildid")

	memb, _ := ws.session.GuildMember(guildID, userID)
	if memb == nil {
		return jsonError(ctx, errNotFound, fasthttp.StatusNotFound)
	}

	guild, err := ws.session.Guild(guildID)
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	return jsonResponse(ctx, GuildFromGuild(guild, memb, ws.cmdhandler), fasthttp.StatusOK)
}

func (ws *WebServer) handlerGuildGetMembers(ctx *routing.Context) error {
	userID := ctx.Get("uid").(string)

	guildID := ctx.Param("guildid")

	memb, _ := ws.session.GuildMember(guildID, userID)
	if memb == nil {
		return jsonError(ctx, errNotFound, fasthttp.StatusNotFound)
	}

	after := ""
	limit := 0

	after = string(ctx.QueryArgs().Peek("after"))
	limit, _ = ctx.QueryArgs().GetUint("limit")

	if limit < 1 {
		limit = 100
	}

	members, err := ws.session.GuildMembers(guildID, after, limit)
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	memblen := len(members)
	fhmembers := make([]*Member, memblen)

	for i, m := range members {
		fhmembers[i] = MemberFromMember(m)
	}

	return jsonResponse(ctx, &ListResponse{
		N:    memblen,
		Data: fhmembers,
	}, fasthttp.StatusOK)
}

func (ws *WebServer) handlerGuildsGetMember(ctx *routing.Context) error {
	userID := ctx.Get("uid").(string)

	guildID := ctx.Param("guildid")
	memberID := ctx.Param("memberid")

	var memb *discordgo.Member

	if memb, _ = ws.session.GuildMember(guildID, userID); memb == nil {
		return jsonError(ctx, errNotFound, fasthttp.StatusNotFound)
	}

	guild, err := ws.session.Guild(guildID)
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	memb, _ = ws.session.GuildMember(guildID, memberID)
	if memb == nil {
		return jsonError(ctx, errNotFound, fasthttp.StatusNotFound)
	}

	mm := MemberFromMember(memb)

	switch {
	case util.IsAdmin(guild, memb):
		mm.Dominance = 1
	case guild.OwnerID == memberID:
		mm.Dominance = 2
	case ws.cmdhandler.IsBotOwner(memberID):
		mm.Dominance = 3
	}

	return jsonResponse(ctx, mm, fasthttp.StatusOK)
}

func (ws *WebServer) handlerGetPermissions(ctx *routing.Context) error {
	userID := ctx.Get("uid").(string)

	guildID := ctx.Param("guildid")
	memberID := ctx.Param("memberid")

	if memb, _ := ws.session.GuildMember(guildID, userID); memb == nil {
		return jsonError(ctx, errNotFound, fasthttp.StatusNotFound)
	}

	perm, err := ws.cmdhandler.GetPermissions(ws.session, guildID, memberID)
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusBadRequest)
	}

	return jsonResponse(ctx, &PermissionsResponse{
		Permissions: perm,
	}, fasthttp.StatusOK)
}

func (ws *WebServer) handlerGetReports(ctx *routing.Context) error {
	userID := ctx.Get("uid").(string)

	guildID := ctx.Param("guildid")
	memberID := ctx.Param("memberid")

	sortBy := string(ctx.QueryArgs().Peek("sortBy"))
	offset := ctx.QueryArgs().GetUintOrZero("offset")
	limit := ctx.QueryArgs().GetUintOrZero("limit")

	if memb, _ := ws.session.GuildMember(guildID, userID); memb == nil {
		return jsonError(ctx, errNotFound, fasthttp.StatusNotFound)
	}

	var reps []*util.Report
	var err error

	if memberID != "" {
		reps, err = ws.db.GetReportsFiltered(guildID, memberID, -1)
	} else {
		reps, err = ws.db.GetReportsGuild(guildID, offset, limit)
	}

	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	resReps := make([]*Report, 0)
	if reps != nil {
		resReps = make([]*Report, len(reps))
		for i, r := range reps {
			resReps[i] = ReportFromReport(r)
		}
	}

	if sortBy == "created" {
		sort.SliceStable(resReps, func(i, j int) bool {
			return resReps[i].Created.After(resReps[j].Created)
		})
	}

	return jsonResponse(ctx, &ListResponse{
		N:    len(resReps),
		Data: resReps,
	}, fasthttp.StatusOK)
}

func (ws *WebServer) handlerGetReportsCount(ctx *routing.Context) error {
	userID := ctx.Get("uid").(string)

	guildID := ctx.Param("guildid")
	memberID := ctx.Param("memberid")

	if memb, _ := ws.session.GuildMember(guildID, userID); memb == nil {
		return jsonError(ctx, errNotFound, fasthttp.StatusNotFound)
	}

	var count int
	var err error

	if memberID != "" {
		count, err = ws.db.GetReportsFilteredCount(guildID, memberID, -1)
	} else {
		count, err = ws.db.GetReportsGuildCount(guildID)
	}

	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	return jsonResponse(ctx, &Count{Count: count}, fasthttp.StatusOK)
}

func (ws *WebServer) handlerGetReport(ctx *routing.Context) error {
	// userID := ctx.Get("uid").(string)

	_id := ctx.Param("id")

	id, err := snowflake.ParseString(_id)
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusBadRequest)
	}

	rep, err := ws.db.GetReport(id)
	if err == core.ErrDatabaseNotFound {
		return jsonError(ctx, errNotFound, fasthttp.StatusNotFound)
	}
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	return jsonResponse(ctx, ReportFromReport(rep), fasthttp.StatusOK)
}

func (ws *WebServer) handlerGetPermissionsAllowed(ctx *routing.Context) error {
	// userID := ctx.Get("uid").(string)

	guildID := ctx.Param("guildid")
	memberID := ctx.Param("memberid")

	perms, err := ws.cmdhandler.GetPermissions(ws.session, guildID, memberID)
	if core.IsErrDatabaseNotFound(err) {
		return jsonError(ctx, errNotFound, fasthttp.StatusNotFound)
	}
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	cmds := ws.cmdhandler.GetCmdInstances()

	allowed := make([]string, len(cmds))
	i := 0
	for _, cmd := range cmds {
		if core.PermissionCheck(cmd.GetDomainName(), perms) {
			allowed[i] = cmd.GetDomainName()
			i++
		}
	}

	return jsonResponse(ctx, &ListResponse{
		N:    i,
		Data: allowed[:i],
	}, fasthttp.StatusOK)
}

func (ws *WebServer) handlerGetGuildSettings(ctx *routing.Context) error {
	gs := new(GuildSettings)

	guildID := ctx.Param("guildid")

	var err error

	if gs.Prefix, err = ws.db.GetGuildPrefix(guildID); err != nil && !core.IsErrDatabaseNotFound(err) {
		return errInternalOrNotFound(ctx, err)
	}

	if gs.Perms, err = ws.db.GetGuildPermissions(guildID); err != nil && !core.IsErrDatabaseNotFound(err) {
		return errInternalOrNotFound(ctx, err)
	}

	if gs.AutoRole, err = ws.db.GetGuildAutoRole(guildID); err != nil && !core.IsErrDatabaseNotFound(err) {
		return errInternalOrNotFound(ctx, err)
	}

	if gs.ModLogChannel, err = ws.db.GetGuildModLog(guildID); err != nil && !core.IsErrDatabaseNotFound(err) {
		return errInternalOrNotFound(ctx, err)
	}

	if gs.VoiceLogChannel, err = ws.db.GetGuildVoiceLog(guildID); err != nil && !core.IsErrDatabaseNotFound(err) {
		return errInternalOrNotFound(ctx, err)
	}

	if gs.JoinMessageText, gs.JoinMessageChannel, err = ws.db.GetGuildJoinMsg(guildID); err != nil && !core.IsErrDatabaseNotFound(err) {
		return errInternalOrNotFound(ctx, err)
	}

	if gs.LeaveMessageText, gs.LeaveMessageChannel, err = ws.db.GetGuildLeaveMsg(guildID); err != nil && !core.IsErrDatabaseNotFound(err) {
		return errInternalOrNotFound(ctx, err)
	}

	return jsonResponse(ctx, gs, fasthttp.StatusOK)
}

func (ws *WebServer) handlerPostGuildSettings(ctx *routing.Context) error {
	userID := ctx.Get("uid").(string)

	guildID := ctx.Param("guildid")

	var err error

	gs := new(GuildSettings)
	if err = parseJSONBody(ctx, gs); err != nil {
		return jsonError(ctx, err, fasthttp.StatusBadRequest)
	}

	if gs.AutoRole != "" {
		if ok, err := ws.cmdhandler.CheckPermissions(ws.session, guildID, userID, "sp.guild.config.autorole"); err != nil {
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
		if ok, err := ws.cmdhandler.CheckPermissions(ws.session, guildID, userID, "sp.guild.config.modlog"); err != nil {
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
		if ok, err := ws.cmdhandler.CheckPermissions(ws.session, guildID, userID, "sp.guild.config.prefix"); err != nil {
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
		if ok, err := ws.cmdhandler.CheckPermissions(ws.session, guildID, userID, "sp.guild.config.voicelog"); err != nil {
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
		if ok, err := ws.cmdhandler.CheckPermissions(ws.session, guildID, userID, "sp.guild.config.joinmsg"); err != nil {
			return errInternalOrNotFound(ctx, err)
		} else if !ok {
			return jsonError(ctx, errUnauthorized, fasthttp.StatusUnauthorized)
		}

		if gs.JoinMessageChannel == "__RESET__" && gs.JoinMessageText == "__RESET__" {
			gs.JoinMessageChannel = ""
			gs.JoinMessageText = ""
		}

		if err = ws.db.SetGuildJoinMsg(guildID, gs.JoinMessageText, gs.JoinMessageChannel); err != nil {
			return errInternalOrNotFound(ctx, err)
		}
	}

	if gs.LeaveMessageChannel != "" && gs.LeaveMessageText != "" {
		if ok, err := ws.cmdhandler.CheckPermissions(ws.session, guildID, userID, "sp.guild.config.leavemsg"); err != nil {
			return errInternalOrNotFound(ctx, err)
		} else if !ok {
			return jsonError(ctx, errUnauthorized, fasthttp.StatusUnauthorized)
		}

		if gs.LeaveMessageChannel == "__RESET__" && gs.LeaveMessageText == "__RESET__" {
			gs.LeaveMessageChannel = ""
			gs.LeaveMessageText = ""
		}

		if err = ws.db.SetGuildLeaveMsg(guildID, gs.LeaveMessageText, gs.LeaveMessageChannel); err != nil {
			return errInternalOrNotFound(ctx, err)
		}
	}

	return jsonResponse(ctx, nil, fasthttp.StatusOK)
}

func (ws *WebServer) handlerGetGuildPermissions(ctx *routing.Context) error {
	userID := ctx.Get("uid").(string)

	guildID := ctx.Param("guildid")

	if memb, _ := ws.session.GuildMember(guildID, userID); memb == nil {
		return jsonError(ctx, errNotFound, fasthttp.StatusNotFound)
	}

	var perms map[string]core.PermissionArray
	var err error

	if perms, err = ws.db.GetGuildPermissions(guildID); err != nil && !core.IsErrDatabaseNotFound(err) {
		return errInternalOrNotFound(ctx, err)
	}

	return jsonResponse(ctx, perms, fasthttp.StatusOK)
}

func (ws *WebServer) handlerPostGuildPermissions(ctx *routing.Context) error {
	userID := ctx.Get("uid").(string)

	guildID := ctx.Param("guildid")

	if ok, err := ws.cmdhandler.CheckPermissions(ws.session, guildID, userID, "sp.guild.config.perms"); err != nil {
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
		if core.IsErrDatabaseNotFound(err) {
			return jsonError(ctx, errNotFound, fasthttp.StatusNotFound)
		}
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	for _, roleID := range update.RoleIDs {
		rperms, ok := perms[roleID]
		if !ok {
			rperms = make(core.PermissionArray, 0)
		}

		rperms = rperms.Update(update.Perm)

		if err = ws.db.SetGuildRolePermission(guildID, roleID, rperms); err != nil {
			return jsonError(ctx, err, fasthttp.StatusInternalServerError)
		}
	}

	return jsonResponse(ctx, nil, fasthttp.StatusOK)
}

func (ws *WebServer) handlerPostGuildMemberReport(ctx *routing.Context) error {
	userID := ctx.Get("uid").(string)

	guildID := ctx.Param("guildid")

	memberID := ctx.Param("memberid")

	if ok, err := ws.cmdhandler.CheckPermissions(ws.session, guildID, userID, "sp.guild.mod.report"); err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	} else if !ok {
		return jsonError(ctx, errUnauthorized, fasthttp.StatusUnauthorized)
	}

	repReq := new(ReportRequest)
	if err := parseJSONBody(ctx, repReq); err != nil {
		return jsonError(ctx, err, fasthttp.StatusBadRequest)
	}

	if memberID == userID {
		return jsonError(ctx, fmt.Errorf("you can not report yourself"), fasthttp.StatusBadRequest)
	}

	if ok, err := repReq.Validate(ctx); !ok {
		return err
	}

	rep, err := shared.PushReport(
		ws.session,
		ws.db,
		guildID,
		userID,
		memberID,
		repReq.Reason,
		repReq.Attachment,
		repReq.Type)

	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	return jsonResponse(ctx, ReportFromReport(rep), fasthttp.StatusCreated)
}

func (ws *WebServer) handlerPostGuildMemberKick(ctx *routing.Context) error {
	userID := ctx.Get("uid").(string)

	guildID := ctx.Param("guildid")

	memberID := ctx.Param("memberid")

	if ok, err := ws.cmdhandler.CheckPermissions(ws.session, guildID, userID, "sp.guild.mod.kick"); err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	} else if !ok {
		return jsonError(ctx, errUnauthorized, fasthttp.StatusUnauthorized)
	}

	req := new(ReasonRequest)
	if err := parseJSONBody(ctx, req); err != nil {
		return jsonError(ctx, err, fasthttp.StatusBadRequest)
	}

	if memberID == userID {
		return jsonError(ctx, fmt.Errorf("you can not kick yourself"), fasthttp.StatusBadRequest)
	}

	guild, err := ws.session.Guild(guildID)
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	executor, err := ws.session.GuildMember(guildID, userID)
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	victim, err := ws.session.GuildMember(guildID, memberID)
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	if util.RolePosDiff(victim, executor, guild) >= 0 {
		return jsonError(ctx, fmt.Errorf("you can not kick members with higher or same permissions than/as yours"), fasthttp.StatusBadRequest)
	}

	if ok, err := req.Validate(ctx); !ok {
		return err
	}

	rep, err := shared.PushKick(
		ws.session,
		ws.db,
		guildID,
		userID,
		memberID,
		req.Reason,
		req.Attachment)

	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	return jsonResponse(ctx, ReportFromReport(rep), fasthttp.StatusCreated)
}

func (ws *WebServer) handlerPostGuildMemberBan(ctx *routing.Context) error {
	userID := ctx.Get("uid").(string)

	guildID := ctx.Param("guildid")

	memberID := ctx.Param("memberid")

	if ok, err := ws.cmdhandler.CheckPermissions(ws.session, guildID, userID, "sp.guild.mod.ban"); err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	} else if !ok {
		return jsonError(ctx, errUnauthorized, fasthttp.StatusUnauthorized)
	}

	req := new(ReasonRequest)
	if err := parseJSONBody(ctx, req); err != nil {
		return jsonError(ctx, err, fasthttp.StatusBadRequest)
	}

	if memberID == userID {
		return jsonError(ctx, fmt.Errorf("you can not ban yourself"), fasthttp.StatusBadRequest)
	}

	guild, err := ws.session.Guild(guildID)
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	executor, err := ws.session.GuildMember(guildID, userID)
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	victim, err := ws.session.GuildMember(guildID, memberID)
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	if util.RolePosDiff(victim, executor, guild) >= 0 {
		return jsonError(ctx, fmt.Errorf("you can not ban members with higher or same permissions than/as yours"), fasthttp.StatusBadRequest)
	}

	if ok, err := req.Validate(ctx); !ok {
		return err
	}

	rep, err := shared.PushBan(
		ws.session,
		ws.db,
		guildID,
		userID,
		memberID,
		req.Reason,
		req.Attachment)

	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	return jsonResponse(ctx, ReportFromReport(rep), fasthttp.StatusCreated)
}

func (ws *WebServer) handlerGetPresence(ctx *routing.Context) error {
	presenceRaw, err := ws.db.GetSetting(util.SettingPresence)
	if err != nil {
		if core.IsErrDatabaseNotFound(err) {
			return jsonResponse(ctx, &util.Presence{
				Game:   util.StdMotd,
				Status: "online",
			}, fasthttp.StatusOK)
		}
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	presence, err := util.UnmarshalPresence(presenceRaw)
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	return jsonResponse(ctx, presence, fasthttp.StatusOK)
}

func (ws *WebServer) handlerPostPresence(ctx *routing.Context) error {
	userID := ctx.Get("uid").(string)

	if ok, err := ws.cmdhandler.CheckPermissions(ws.session, "", userID, "sp.game"); err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	} else if !ok {
		return jsonError(ctx, errUnauthorized, fasthttp.StatusUnauthorized)
	}

	presence := new(util.Presence)
	if err := parseJSONBody(ctx, presence); err != nil {
		return jsonError(ctx, err, fasthttp.StatusBadRequest)
	}

	if strings.Contains(presence.Game, util.PresenceSeperator) {
		return jsonError(ctx,
			fmt.Errorf("'%s' is used as seperator for the presence settings save and can not be used in the actual game message",
				util.PresenceSeperator), fasthttp.StatusBadRequest)
	}

	if err := presence.Validate(); err != nil {
		return jsonError(ctx, err, fasthttp.StatusBadRequest)
	}

	if err := ws.db.SetSetting(util.SettingPresence, presence.Marshal()); err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	if err := ws.session.UpdateStatusComplex(presence.ToUpdateStatusData()); err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	return jsonResponse(ctx, presence, fasthttp.StatusOK)
}

func (ws *WebServer) handlerGetInviteSettings(ctx *routing.Context) error {
	var guildID, message, inviteCode string
	var err error

	if guildID, err = ws.db.GetSetting(util.SettingWIInviteGuildID); err != nil {
		if isErr, err := errInternalIgnoreNotFound(ctx, err); isErr {
			return err
		}
	}

	if guildID == "" {
		return jsonResponse(ctx, &InviteSettingsResponse{
			Guild:     nil,
			InviteURL: "",
			Message:   "",
		}, fasthttp.StatusOK)
	}

	if message, err = ws.db.GetSetting(util.SettingWIInviteText); err != nil {
		if isErr, err := errInternalIgnoreNotFound(ctx, err); isErr {
			return err
		}
	}

	if inviteCode, err = ws.db.GetSetting(util.SettingWIInviteCode); err != nil {
		if isErr, err := errInternalIgnoreNotFound(ctx, err); isErr {
			return err
		}
	}

	guild, err := ws.session.Guild(guildID)
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	invites, err := ws.session.GuildInvites(guildID)
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	if inviteCode != "" {
		for _, inv := range invites {
			if inv.Inviter != nil && inv.Inviter.ID == ws.session.State.User.ID && !inv.Revoked {
				inviteCode = inv.Code
				break
			}
		}
	}

	if inviteCode == "" {
		var channel *discordgo.Channel
		for _, c := range guild.Channels {
			if c.Type == discordgo.ChannelTypeGuildText {
				channel = c
			}
		}
		if channel == nil {
			return jsonError(ctx, fmt.Errorf("could not find any channel to create invite for"), fasthttp.StatusConflict)
		}

		invite, err := ws.session.ChannelInviteCreate(channel.ID, discordgo.Invite{
			Temporary: false,
		})
		if err != nil {
			return jsonError(ctx, err, fasthttp.StatusInternalServerError)
		}

		inviteCode = invite.Code
		if err = ws.db.SetSetting(util.SettingWIInviteCode, inviteCode); err != nil {
			return jsonError(ctx, err, fasthttp.StatusInternalServerError)
		}
	}

	res := &InviteSettingsResponse{
		Guild:     GuildFromGuild(guild, nil, nil),
		Message:   message,
		InviteURL: fmt.Sprintf("https://discord.gg/%s", inviteCode),
	}

	return jsonResponse(ctx, res, fasthttp.StatusOK)
}

func (ws *WebServer) handlerPostInviteSettings(ctx *routing.Context) error {
	userID := ctx.Get("uid").(string)

	if ok, err := ws.cmdhandler.CheckPermissions(ws.session, "", userID, "sp.noguildinvite"); err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	} else if !ok {
		return jsonError(ctx, errUnauthorized, fasthttp.StatusUnauthorized)
	}

	req := new(InviteSettingsRequest)
	if err := parseJSONBody(ctx, req); err != nil {
		return jsonError(ctx, err, fasthttp.StatusBadRequest)
	}

	var err error

	if req.GuildID != "" {

		guild, err := ws.session.Guild(req.GuildID)
		if err != nil {
			return jsonError(ctx, err, fasthttp.StatusBadRequest)
		}

		if req.InviteCode != "" {
			invites, err := ws.session.GuildInvites(req.GuildID)
			if err != nil {
				return jsonError(ctx, err, fasthttp.StatusInternalServerError)
			}

			var valid bool
			for _, inv := range invites {
				if inv.Code == req.InviteCode && !inv.Revoked {
					valid = true
					break
				}
			}

			if !valid {
				return jsonError(ctx, fmt.Errorf("invalid invite code"), fasthttp.StatusBadRequest)
			}
		} else {
			var channel *discordgo.Channel
			for _, c := range guild.Channels {
				if c.Type == discordgo.ChannelTypeGuildText {
					channel = c
				}
			}
			if channel == nil {
				return jsonError(ctx, fmt.Errorf("could not find any channel to create invite for"), fasthttp.StatusConflict)
			}

			invite, err := ws.session.ChannelInviteCreate(channel.ID, discordgo.Invite{
				Temporary: false,
			})
			if err != nil {
				return jsonError(ctx, err, fasthttp.StatusInternalServerError)
			}

			req.InviteCode = invite.Code
		}
	}

	if err = ws.db.SetSetting(util.SettingWIInviteCode, req.InviteCode); err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	if err = ws.db.SetSetting(util.SettingWIInviteGuildID, req.GuildID); err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	if err = ws.db.SetSetting(util.SettingWIInviteText, req.Messsage); err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	return jsonResponse(ctx, nil, fasthttp.StatusOK)
}

func (ws *WebServer) handlerGetSystemInfo(ctx *routing.Context) error {

	buildTS, _ := strconv.Atoi(util.AppDate)
	buildDate := time.Unix(int64(buildTS), 0)

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	uptime := int64(time.Since(util.StatsStartupTime).Seconds())

	info := &SystemInfo{
		Version:    util.AppVersion,
		CommitHash: util.AppCommit,
		BuildDate:  buildDate,
		GoVersion:  runtime.Version(),

		Uptime:    uptime,
		UptimeStr: fmt.Sprintf("%d", uptime),

		OS:          runtime.GOOS,
		Arch:        runtime.GOARCH,
		CPUs:        runtime.NumCPU(),
		GoRoutines:  runtime.NumGoroutine(),
		StackUse:    memStats.StackInuse,
		StackUseStr: fmt.Sprintf("%d", memStats.StackInuse),
		HeapUse:     memStats.HeapInuse,
		HeapUseStr:  fmt.Sprintf("%d", memStats.HeapInuse),

		BotUserID: ws.session.State.User.ID,
		BotInvite: fmt.Sprintf("https://discordapp.com/api/oauth2/authorize?client_id=%s&scope=bot&permissions=%d",
			ws.session.State.User.ID, util.InvitePermission),

		Guilds: len(ws.session.State.Guilds),
	}

	return jsonResponse(ctx, info, fasthttp.StatusOK)
}
