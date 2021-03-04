package webserver

import (
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"github.com/zekroTJA/shinpuru/internal/core/database"
)

// ---------------------------------------------------------------------------
// - GET /api/usersettings/ota

func (ws *WebServer) handlerGetUsersettingsOta(ctx *routing.Context) error {
	userID := ctx.Get("uid").(string)

	enabled, err := ws.db.GetUserOTAEnabled(userID)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	return jsonResponse(ctx, &UsersettingsOTA{enabled}, fasthttp.StatusOK)
}

// ---------------------------------------------------------------------------
// - POST /api/usersettings/ota

func (ws *WebServer) handlerPostUsersettingsOta(ctx *routing.Context) (err error) {
	userID := ctx.Get("uid").(string)

	data := new(UsersettingsOTA)
	if err = parseJSONBody(ctx, data); err != nil {
		return jsonResponse(ctx, err, fasthttp.StatusBadRequest)
	}

	if err = ws.db.SetUserOTAEnabled(userID, data.Enabled); err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	return jsonResponse(ctx, data, fasthttp.StatusOK)
}
