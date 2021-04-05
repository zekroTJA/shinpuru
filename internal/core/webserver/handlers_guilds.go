package webserver

import (
	"fmt"
	"strconv"

	"github.com/bwmarrin/discordgo"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"github.com/zekroTJA/shinpuru/internal/core/database"
	"github.com/zekroTJA/shinpuru/internal/shared/models"
	"github.com/zekroTJA/shinpuru/pkg/discordutil"
)

// ---------------------------------------------------------------------------
// - GET /api/guilds

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

// ---------------------------------------------------------------------------
// - GET /api/guilds/:guildid

func (ws *WebServer) handlerGuildsGetGuild(ctx *routing.Context) error {
	userID := ctx.Get("uid").(string)

	guildID := ctx.Param("guildid")

	memb, _ := ws.session.GuildMember(guildID, userID)
	if memb == nil {
		return jsonError(ctx, errNotFound, fasthttp.StatusNotFound)
	}

	guild, err := discordutil.GetGuild(ws.session, guildID)
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	gRes := GuildFromGuild(guild, memb, ws.db, ws.config.Discord.OwnerID)
	return jsonResponse(ctx, gRes, fasthttp.StatusOK)
}

// ---------------------------------------------------------------------------
// - GET /api/guilds/:guildid/scoreboard

func (ws *WebServer) handlerGetGuildScoreboard(ctx *routing.Context) error {
	guildID := ctx.Param("guildid")
	limitQ := ctx.QueryArgs().Peek("limit")

	limit := 25

	if len(limitQ) > 0 {
		limit, err := strconv.Atoi(string(limitQ))
		if err != nil {
			return jsonError(ctx, err, fasthttp.StatusBadRequest)
		}
		if limit < 0 || limit > 100 {
			return jsonError(ctx,
				fmt.Errorf("limit must be in range [0, 100]"), fasthttp.StatusBadRequest)
		}
	}

	karmaList, err := ws.db.GetKarmaGuild(guildID, limit)

	if err == database.ErrDatabaseNotFound {
		return jsonError(ctx, nil, fasthttp.StatusNotFound)
	} else if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	results := make([]*GuildKarmaEntry, len(karmaList))

	var i int
	for _, e := range karmaList {
		member, err := discordutil.GetMember(ws.session, guildID, e.UserID)
		if err != nil {
			continue
		}
		results[i] = &GuildKarmaEntry{
			Member: MemberFromMember(member),
			Value:  e.Value,
		}
		i++
	}

	return jsonResponse(ctx, &ListResponse{N: i, Data: results[:i]}, fasthttp.StatusOK)
}

// ---------------------------------------------------------------------------
// - GET /api/guilds/:guildid/antiraid/joinlog

func (ws *WebServer) handlerGetGuildAntiraidJoinlog(ctx *routing.Context) error {
	userID := ctx.Get("uid").(string)

	guildID := ctx.Param("guildid")

	if ok, _, err := ws.pmw.CheckPermissions(ws.session, guildID, userID, "sp.guild.config.antiraid"); err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	} else if !ok {
		return jsonError(ctx, errUnauthorized, fasthttp.StatusUnauthorized)
	}

	joinlog, err := ws.db.GetAntiraidJoinList(guildID)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	if joinlog == nil {
		joinlog = make([]*models.JoinLogEntry, 0)
	}

	return jsonResponse(ctx, &ListResponse{
		N:    len(joinlog),
		Data: joinlog,
	}, fasthttp.StatusOK)
}

// ---------------------------------------------------------------------------
// - DELETE /api/guilds/:guildid/antiraid/joinlog

func (ws *WebServer) handlerDeleteGuildAntiraidJoinlog(ctx *routing.Context) error {
	userID := ctx.Get("uid").(string)

	guildID := ctx.Param("guildid")

	if ok, _, err := ws.pmw.CheckPermissions(ws.session, guildID, userID, "sp.guild.config.antiraid"); err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	} else if !ok {
		return jsonError(ctx, errUnauthorized, fasthttp.StatusUnauthorized)
	}

	if err := ws.db.FlushAntiraidJoinList(guildID); err != nil && !database.IsErrDatabaseNotFound(err) {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	return jsonResponse(ctx, nil, fasthttp.StatusOK)
}

// ---------------------------------------------------------------------------
// - GET /api/guilds/:guildid/starboard

func (ws *WebServer) handlerGetGuildStarboard(ctx *routing.Context) error {
	guildID := ctx.Param("guildid")
	limitQ := ctx.QueryArgs().Peek("limit")
	offsetQ := ctx.QueryArgs().Peek("offset")
	sortQ := ctx.QueryArgs().Peek("sort")

	limit := 20
	offset := 0
	sort := models.StarboardSortByLatest

	if len(limitQ) > 0 {
		limit, err := strconv.Atoi(string(limitQ))
		if err != nil {
			return jsonError(ctx, err, fasthttp.StatusBadRequest)
		}
		if limit < 0 || limit > 100 {
			return jsonError(ctx,
				fmt.Errorf("limit must be in range [0, 100]"), fasthttp.StatusBadRequest)
		}
	}

	if len(offsetQ) > 0 {
		offset, err := strconv.Atoi(string(offsetQ))
		if err != nil {
			return jsonError(ctx, err, fasthttp.StatusBadRequest)
		}
		if offset < 0 {
			return jsonError(ctx,
				fmt.Errorf("offset must be larger or equal 0"), fasthttp.StatusBadRequest)
		}
	}

	if len(sortQ) > 0 {
		switch string(sortQ) {
		case "latest":
			sort = models.StarboardSortByLatest
		case "top":
			sort = models.StarboardSortByMostRated
		default:
			return jsonError(ctx,
				fmt.Errorf("invalid sort property"), fasthttp.StatusBadRequest)
		}
	}

	entries, err := ws.db.GetStarboardEntries(guildID, sort, limit, offset)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	results := make([]*StarboardEntryResponse, len(entries))

	var i int
	for _, e := range entries {
		if e.Deleted {
			continue
		}

		member, err := discordutil.GetMember(ws.session, guildID, e.AuthorID)
		if err != nil {
			continue
		}

		results[i] = &StarboardEntryResponse{
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

	return jsonResponse(ctx, &ListResponse{N: i, Data: results[:i]}, fasthttp.StatusOK)
}
