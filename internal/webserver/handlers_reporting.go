package webserver

import (
	"bytes"
	"fmt"

	"github.com/bwmarrin/snowflake"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"github.com/zekroTJA/shinpuru/internal/core/database"
	"github.com/zekroTJA/shinpuru/internal/shared"
	"github.com/zekroTJA/shinpuru/internal/util/imgstore"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/discordutil"
	"github.com/zekroTJA/shinpuru/pkg/roleutil"
)

// ---------------------------------------------------------------------------
// - POST /api/guilds/:guildid/:memberid/reports

func (ws *WebServer) handlerPostGuildMemberReport(ctx *routing.Context) error {
	userID := ctx.Get("uid").(string)

	guildID := ctx.Param("guildid")

	memberID := ctx.Param("memberid")

	if ok, _, err := ws.pmw.CheckPermissions(ws.session, guildID, userID, "sp.guild.mod.report"); err != nil {
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

	if repReq.Attachment != "" {
		img, err := imgstore.DownloadFromURL(repReq.Attachment)
		if err != nil {
			return jsonError(ctx, err, fasthttp.StatusBadRequest)
		}
		err = ws.st.PutObject(static.StorageBucketImages, img.ID.String(),
			bytes.NewReader(img.Data), int64(img.Size), img.MimeType)
		if err != nil {
			return jsonError(ctx, err, fasthttp.StatusInternalServerError)
		}
		repReq.Attachment = img.ID.String()
	}

	rep, err := shared.PushReport(
		ws.session,
		ws.db,
		ws.config.WebServer.PublicAddr,
		guildID,
		userID,
		memberID,
		repReq.Reason,
		repReq.Attachment,
		repReq.Type)

	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	return jsonResponse(ctx, ReportFromReport(rep, ws.config.WebServer.PublicAddr), fasthttp.StatusCreated)
}

// ---------------------------------------------------------------------------
// - POST /api/guilds/:guildid/:memberid/kick

func (ws *WebServer) handlerPostGuildMemberKick(ctx *routing.Context) error {
	userID := ctx.Get("uid").(string)

	guildID := ctx.Param("guildid")

	memberID := ctx.Param("memberid")

	if ok, _, err := ws.pmw.CheckPermissions(ws.session, guildID, userID, "sp.guild.mod.kick"); err != nil {
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

	guild, err := discordutil.GetGuild(ws.session, guildID)
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

	if roleutil.PositionDiff(victim, executor, guild) >= 0 {
		return jsonError(ctx, fmt.Errorf("you can not kick members with higher or same permissions than/as yours"), fasthttp.StatusBadRequest)
	}

	if ok, err := req.Validate(ctx); !ok {
		return err
	}

	if req.Attachment != "" {
		img, err := imgstore.DownloadFromURL(req.Attachment)
		if err != nil {
			return jsonError(ctx, err, fasthttp.StatusBadRequest)
		}
		err = ws.st.PutObject(static.StorageBucketImages, img.ID.String(),
			bytes.NewReader(img.Data), int64(img.Size), img.MimeType)
		if err != nil {
			return jsonError(ctx, err, fasthttp.StatusInternalServerError)
		}
		req.Attachment = img.ID.String()
	}

	rep, err := shared.PushKick(
		ws.session,
		ws.db,
		ws.config.WebServer.PublicAddr,
		guildID,
		userID,
		memberID,
		req.Reason,
		req.Attachment)

	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	return jsonResponse(ctx, ReportFromReport(rep, ws.config.WebServer.PublicAddr), fasthttp.StatusCreated)
}

// ---------------------------------------------------------------------------
// - POST /api/guilds/:guildid/:memberid/ban

func (ws *WebServer) handlerPostGuildMemberBan(ctx *routing.Context) error {
	userID := ctx.Get("uid").(string)

	guildID := ctx.Param("guildid")

	memberID := ctx.Param("memberid")

	if ok, _, err := ws.pmw.CheckPermissions(ws.session, guildID, userID, "sp.guild.mod.ban"); err != nil {
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

	guild, err := discordutil.GetGuild(ws.session, guildID)
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

	if roleutil.PositionDiff(victim, executor, guild) >= 0 {
		return jsonError(ctx, fmt.Errorf("you can not ban members with higher or same permissions than/as yours"), fasthttp.StatusBadRequest)
	}

	if ok, err := req.Validate(ctx); !ok {
		return err
	}

	if req.Attachment != "" {
		img, err := imgstore.DownloadFromURL(req.Attachment)
		if err != nil {
			return jsonError(ctx, err, fasthttp.StatusBadRequest)
		}
		err = ws.st.PutObject(static.StorageBucketImages, img.ID.String(),
			bytes.NewReader(img.Data), int64(img.Size), img.MimeType)
		if err != nil {
			return jsonError(ctx, err, fasthttp.StatusInternalServerError)
		}
		req.Attachment = img.ID.String()
	}

	rep, err := shared.PushBan(
		ws.session,
		ws.db,
		ws.config.WebServer.PublicAddr,
		guildID,
		userID,
		memberID,
		req.Reason,
		req.Attachment)

	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	return jsonResponse(ctx, ReportFromReport(rep, ws.config.WebServer.PublicAddr), fasthttp.StatusCreated)
}

// ---------------------------------------------------------------------------
// - GET /api/reports/:id

func (ws *WebServer) handlerGetReport(ctx *routing.Context) error {
	// userID := ctx.Get("uid").(string)

	_id := ctx.Param("id")

	id, err := snowflake.ParseString(_id)
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusBadRequest)
	}

	rep, err := ws.db.GetReport(id)
	if database.IsErrDatabaseNotFound(err) {
		return jsonError(ctx, errNotFound, fasthttp.StatusNotFound)
	}
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	return jsonResponse(ctx, ReportFromReport(rep, ws.config.WebServer.PublicAddr), fasthttp.StatusOK)
}

// ---------------------------------------------------------------------------
// - GET /api/reports/:id/revoke

func (ws *WebServer) handlerPostReportRevoke(ctx *routing.Context) error {
	userID := ctx.Get("uid").(string)

	_id := ctx.Param("id")

	id, err := snowflake.ParseString(_id)
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusBadRequest)
	}

	rep, err := ws.db.GetReport(id)
	if database.IsErrDatabaseNotFound(err) {
		return jsonError(ctx, errNotFound, fasthttp.StatusNotFound)
	}
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	if ok, _, err := ws.pmw.CheckPermissions(ws.session, rep.GuildID, userID, "sp.guild.mod.report"); err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	} else if !ok {
		return jsonError(ctx, errUnauthorized, fasthttp.StatusUnauthorized)
	}

	var reason struct {
		Reason string `json:"reason"`
	}

	if err := parseJSONBody(ctx, &reason); err != nil {
		return jsonError(ctx, err, fasthttp.StatusBadRequest)
	}

	_, err = shared.RevokeReport(
		rep,
		userID,
		reason.Reason,
		ws.config.WebServer.Addr,
		ws.db,
		ws.session)

	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	return jsonResponse(ctx, nil, fasthttp.StatusOK)
}
