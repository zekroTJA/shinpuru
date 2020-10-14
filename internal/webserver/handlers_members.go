package webserver

import (
	"github.com/bwmarrin/discordgo"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"github.com/zekroTJA/shinpuru/internal/core/database"
	"github.com/zekroTJA/shinpuru/internal/util/report"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/discordutil"
)

// ---------------------------------------------------------------------------
// - GET /api/guilds/:guildid/members

func (ws *WebServer) handlerGetGuildMembers(ctx *routing.Context) error {
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

// ---------------------------------------------------------------------------
// - GET /api/guilds/:guildid/:memberid

func (ws *WebServer) handlerGuildsGetMember(ctx *routing.Context) error {
	userID := ctx.Get("uid").(string)

	guildID := ctx.Param("guildid")
	memberID := ctx.Param("memberid")

	var memb *discordgo.Member

	if memb, _ = ws.session.GuildMember(guildID, userID); memb == nil {
		return jsonError(ctx, errNotFound, fasthttp.StatusNotFound)
	}

	guild, err := discordutil.GetGuild(ws.session, guildID)
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	memb, _ = ws.session.GuildMember(guildID, memberID)
	if memb == nil {
		return jsonError(ctx, errNotFound, fasthttp.StatusNotFound)
	}

	memb.GuildID = guildID

	mm := MemberFromMember(memb)

	switch {
	case discordutil.IsAdmin(guild, memb):
		mm.Dominance = 1
	case guild.OwnerID == memberID:
		mm.Dominance = 2
	case ws.config.Discord.OwnerID == memb.User.ID:
		mm.Dominance = 3
	}

	mm.Karma, err = ws.db.GetKarma(memberID, guildID)
	if !database.IsErrDatabaseNotFound(err) && err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	mm.KarmaTotal, err = ws.db.GetKarmaSum(memberID)
	if !database.IsErrDatabaseNotFound(err) && err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	return jsonResponse(ctx, mm, fasthttp.StatusOK)
}

// ---------------------------------------------------------------------------
// - GET /api/guilds/:guildid/:memberid/permissions

func (ws *WebServer) handlerGetMemberPermissions(ctx *routing.Context) error {
	userID := ctx.Get("uid").(string)

	guildID := ctx.Param("guildid")
	memberID := ctx.Param("memberid")

	if memb, _ := ws.session.GuildMember(guildID, userID); memb == nil {
		return jsonError(ctx, errNotFound, fasthttp.StatusNotFound)
	}

	perm, _, err := ws.pmw.GetPermissions(ws.session, guildID, memberID)
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusBadRequest)
	}

	return jsonResponse(ctx, &PermissionsResponse{
		Permissions: perm,
	}, fasthttp.StatusOK)
}

// ---------------------------------------------------------------------------
// - GET /api/guilds/:guildid/reports
// - GET /api/guilds/:guildid/:memberid/reports

func (ws *WebServer) handlerGetReports(ctx *routing.Context) error {
	userID := ctx.Get("uid").(string)

	guildID := ctx.Param("guildid")
	memberID := ctx.Param("memberid")

	offset := ctx.QueryArgs().GetUintOrZero("offset")
	limit := ctx.QueryArgs().GetUintOrZero("limit")

	if memb, _ := ws.session.GuildMember(guildID, userID); memb == nil {
		return jsonError(ctx, errNotFound, fasthttp.StatusNotFound)
	}

	var reps []*report.Report
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
			resReps[i] = ReportFromReport(r, ws.config.WebServer.PublicAddr)
		}
	}

	return jsonResponse(ctx, &ListResponse{
		N:    len(resReps),
		Data: resReps,
	}, fasthttp.StatusOK)
}

// ---------------------------------------------------------------------------
// - GET /api/guilds/:guildid/reports/count
// - GET /api/guilds/:guildid/:memberid/reports/count

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

// ---------------------------------------------------------------------------
// - GET /api/guilds/:guildid/:memberid/permissions/allowed

func (ws *WebServer) handlerGetMemberPermissionsAllowed(ctx *routing.Context) error {
	// userID := ctx.Get("uid").(string)

	guildID := ctx.Param("guildid")
	memberID := ctx.Param("memberid")

	perms, _, err := ws.pmw.GetPermissions(ws.session, guildID, memberID)
	if database.IsErrDatabaseNotFound(err) {
		return jsonError(ctx, errNotFound, fasthttp.StatusNotFound)
	}
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	cmds := ws.cmdhandler.GetCommandInstances()

	allowed := make([]string, len(cmds) + len(static.AdditionalPermissions))
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

	return jsonResponse(ctx, &ListResponse{
		N:    i,
		Data: allowed[:i],
	}, fasthttp.StatusOK)
}
