package webserver

import (
	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

func (ws *WebServer) handlerGetMe(ctx *routing.Context) error {
	userID := ctx.Get("uid").(string)

	user, err := ws.session.User(userID)
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	res := &User{
		User:      user,
		AvatarURL: user.AvatarURL(""),
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

	guildID := ctx.Param("id")

	if memb, _ := ws.session.GuildMember(guildID, userID); memb == nil {
		return jsonError(ctx, errNotFound, fasthttp.StatusNotFound)
	}

	guild, err := ws.session.Guild(guildID)
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	return jsonResponse(ctx, GuildFromGuild(guild), fasthttp.StatusOK)
}
