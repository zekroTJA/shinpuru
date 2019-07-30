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
