package webserver

import (
	"errors"
	"time"

	"github.com/bwmarrin/snowflake"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"github.com/zekroTJA/shinpuru/internal/core/database"
	"github.com/zekroTJA/shinpuru/internal/util/report"
	"github.com/zekroTJA/shinpuru/internal/util/snowflakenodes"
	"github.com/zekroTJA/shinpuru/pkg/discordutil"
)

// ---------------------------------------------------------------------------
// - GET /api/guilds/:guildid/unbanrequests

func (ws *WebServer) handlerGetGuildUnbanrequests(ctx *routing.Context) error {
	userID := ctx.Get("uid").(string)

	guildID := ctx.Param("guildid")

	if ok, _, err := ws.pmw.CheckPermissions(ws.session, guildID, userID, "sp.guild.mod.unbanrequests"); err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	} else if !ok {
		return jsonError(ctx, errUnauthorized, fasthttp.StatusUnauthorized)
	}

	requests, err := ws.db.GetGuildUnbanRequests(guildID)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}
	if requests == nil {
		requests = make([]*report.UnbanRequest, 0)
	}

	for _, r := range requests {
		r.Hydrate()
	}

	return jsonResponse(ctx, &ListResponse{len(requests), requests}, fasthttp.StatusOK)
}

// ---------------------------------------------------------------------------
// - GET /api/guilds/:guildid/unbanrequests/count

func (ws *WebServer) handlerGetGuildUnbanrequestsCount(ctx *routing.Context) error {
	userID := ctx.Get("uid").(string)

	guildID := ctx.Param("guildid")

	stateFilter, err := ctx.QueryArgs().GetUint("state")
	if err != nil {
		stateFilter = -1
	}

	if ok, _, err := ws.pmw.CheckPermissions(ws.session, guildID, userID, "sp.guild.mod.unbanrequests"); err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	} else if !ok {
		return jsonError(ctx, errUnauthorized, fasthttp.StatusUnauthorized)
	}

	requests, err := ws.db.GetGuildUnbanRequests(guildID)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}
	if requests == nil {
		requests = make([]*report.UnbanRequest, 0)
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

	return jsonResponse(ctx, &Count{count}, fasthttp.StatusOK)
}

// ---------------------------------------------------------------------------
// - GET /api/guilds/:guildid/:memberid/unbanrequests

func (ws *WebServer) handlerGetGuildMemberUnbanrequests(ctx *routing.Context) error {
	userID := ctx.Get("uid").(string)

	guildID := ctx.Param("guildid")
	memberID := ctx.Param("memberid")

	if ok, _, err := ws.pmw.CheckPermissions(ws.session, guildID, userID, "sp.guild.mod.unbanrequests"); err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	} else if !ok {
		return jsonError(ctx, errUnauthorized, fasthttp.StatusUnauthorized)
	}

	requests, err := ws.db.GetGuildUserUnbanRequests(guildID, memberID)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}
	if requests == nil {
		requests = make([]*report.UnbanRequest, 0)
	}

	for _, r := range requests {
		r.Hydrate()
	}

	return jsonResponse(ctx, &ListResponse{len(requests), requests}, fasthttp.StatusOK)
}

// ---------------------------------------------------------------------------
// - GET /api/guilds/:guildid/:memberid/unbanrequests/count

func (ws *WebServer) handlerGetGuildMemberUnbanrequestsCount(ctx *routing.Context) error {
	userID := ctx.Get("uid").(string)

	guildID := ctx.Param("guildid")
	memberID := ctx.Param("memberid")

	stateFilter, err := ctx.QueryArgs().GetUint("state")
	if err != nil {
		stateFilter = -1
	}

	if ok, _, err := ws.pmw.CheckPermissions(ws.session, guildID, userID, "sp.guild.mod.unbanrequests"); err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	} else if !ok {
		return jsonError(ctx, errUnauthorized, fasthttp.StatusUnauthorized)
	}

	requests, err := ws.db.GetGuildUserUnbanRequests(guildID, memberID)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}
	if requests == nil {
		requests = make([]*report.UnbanRequest, 0)
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

	return jsonResponse(ctx, &Count{count}, fasthttp.StatusOK)
}

// ---------------------------------------------------------------------------
// - GET /api/guilds/:guildid/unbanrequests/:id

func (ws *WebServer) handlerGetGuildUnbanrequest(ctx *routing.Context) error {
	userID := ctx.Get("uid").(string)

	guildID := ctx.Param("guildid")
	id := ctx.Param("id")

	if ok, _, err := ws.pmw.CheckPermissions(ws.session, guildID, userID, "sp.guild.mod.unbanrequests"); err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	} else if !ok {
		return jsonError(ctx, errUnauthorized, fasthttp.StatusUnauthorized)
	}

	request, err := ws.db.GetUnbanRequest(id)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}
	if request == nil || request.GuildID != guildID {
		return jsonResponse(ctx, nil, fasthttp.StatusNotFound)
	}

	return jsonResponse(ctx, request.Hydrate(), fasthttp.StatusOK)
}

// ---------------------------------------------------------------------------
// - POST /api/guilds/:guildid/unbanrequests/:id

func (ws *WebServer) handlerPostGuildUnbanrequest(ctx *routing.Context) error {
	userID := ctx.Get("uid").(string)

	guildID := ctx.Param("guildid")
	id := ctx.Param("id")

	if ok, _, err := ws.pmw.CheckPermissions(ws.session, guildID, userID, "sp.guild.mod.unbanrequests"); err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	} else if !ok {
		return jsonError(ctx, errUnauthorized, fasthttp.StatusUnauthorized)
	}

	rUpdate := new(report.UnbanRequest)
	if err := parseJSONBody(ctx, rUpdate); err != nil {
		return jsonError(ctx, err, fasthttp.StatusBadRequest)
	}

	request, err := ws.db.GetUnbanRequest(id)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}
	if request == nil || request.GuildID != guildID {
		return jsonResponse(ctx, nil, fasthttp.StatusNotFound)
	}

	if rUpdate.ProcessedMessage == "" {
		return jsonError(ctx, errors.New("process reason message must be provided"), fasthttp.StatusBadRequest)
	}

	if request.ID, err = snowflake.ParseString(id); err != nil {
		return jsonError(ctx, err, fasthttp.StatusBadRequest)
	}
	request.ProcessedBy = userID
	request.Status = rUpdate.Status
	request.Processed = time.Now()
	request.ProcessedMessage = rUpdate.ProcessedMessage

	if err = ws.db.UpdateUnbanRequest(request); err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	if request.Status == report.UnbanRequestStateAccepted {
		if err = ws.session.GuildBanDelete(request.GuildID, request.UserID); err != nil {
			return jsonError(ctx, err, fasthttp.StatusInternalServerError)
		}
	}

	return jsonResponse(ctx, request.Hydrate(), fasthttp.StatusOK)
}

// ---------------------------------------------------------------------------
// - GET /api/unbanrequests/bannedguilds

func (ws *WebServer) handlerGetUnbanrequestBannedguilds(ctx *routing.Context) error {
	userID := ctx.Get("uid").(string)

	guildsArr, err := ws.getUserBannedGuilds(userID)
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	return jsonResponse(ctx, &ListResponse{len(guildsArr), guildsArr}, fasthttp.StatusOK)
}

// ---------------------------------------------------------------------------
// - GET /api/unbanrequests

func (ws *WebServer) handlerGetUnbanrequest(ctx *routing.Context) error {
	userID := ctx.Get("uid").(string)

	requests, err := ws.db.GetGuildUserUnbanRequests(userID, "")
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}
	if requests == nil {
		requests = make([]*report.UnbanRequest, 0)
	}

	for _, r := range requests {
		r.Hydrate()
	}

	return jsonResponse(ctx, &ListResponse{len(requests), requests}, fasthttp.StatusCreated)
}

// ---------------------------------------------------------------------------
// - POST /api/unbanrequests

func (ws *WebServer) handlerPostUnbanrequest(ctx *routing.Context) error {
	userID := ctx.Get("uid").(string)

	user, err := ws.session.User(userID)
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	req := new(report.UnbanRequest)
	if err := parseJSONBody(ctx, req); err != nil {
		return jsonError(ctx, err, fasthttp.StatusBadRequest)
	}
	if err := req.Validate(); err != nil {
		return jsonError(ctx, err, fasthttp.StatusBadRequest)
	}

	rep, err := ws.db.GetReportsFiltered(req.GuildID, userID, 1)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	if rep == nil || len(rep) == 0 {
		return jsonError(ctx, errors.New("you have no filed ban reports on this guild"), fasthttp.StatusBadRequest)
	}

	requests, err := ws.db.GetGuildUserUnbanRequests(userID, req.GuildID)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	if requests != nil {
		for _, r := range requests {
			if r.Status == report.UnbanRequestStatePending {
				return jsonError(ctx, errors.New("there is still one open unban request to be proceed"), fasthttp.StatusBadRequest)
			}
		}
	}

	finalReq := &report.UnbanRequest{
		ID:      snowflakenodes.NodeUnbanRequests.Generate(),
		UserID:  userID,
		GuildID: req.GuildID,
		UserTag: user.String(),
		Message: req.Message,
		Status:  report.UnbanRequestStatePending,
	}

	if err := ws.db.AddUnbanRequest(finalReq); err != nil && !database.IsErrDatabaseNotFound(err) {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	return jsonResponse(ctx, finalReq.Hydrate(), fasthttp.StatusCreated)
}

// --- HELPERS ------------

func (ws *WebServer) getUserBannedGuilds(userID string) ([]*GuildReduced, error) {
	reps, err := ws.db.GetReportsFiltered("", userID, 1)
	if err != nil {
		if database.IsErrDatabaseNotFound(err) {
			return []*GuildReduced{}, nil
		}
		return nil, err
	}

	guilds := make(map[string]*GuildReduced)
	for _, r := range reps {
		if _, ok := guilds[r.GuildID]; ok {
			continue
		}
		guild, err := discordutil.GetGuild(ws.session, r.GuildID)
		if err != nil {
			return nil, err
		}
		guilds[r.GuildID] = GuildReducedFromGuild(guild)
	}

	guildsArr := make([]*GuildReduced, len(guilds))
	i := 0
	for _, g := range guilds {
		guildsArr[i] = g
		i++
	}

	return guildsArr, nil
}
