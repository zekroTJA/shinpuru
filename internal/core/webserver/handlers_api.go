package webserver

import (
	"errors"
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"time"

	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"

	"github.com/zekroTJA/shinpuru/internal/core/database"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/colors"
	"github.com/zekroTJA/shinpuru/pkg/discordutil"
	"github.com/zekroTJA/shinpuru/pkg/etag"
)

// ---------------------------------------------------------------------------
// - GET /api/me

func (ws *WebServer) handlerGetMe(ctx *routing.Context) error {
	userID := ctx.Get("uid").(string)

	user, err := ws.session.User(userID)
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	created, _ := discordutil.GetDiscordSnowflakeCreationTime(user.ID)

	res := &User{
		User:      user,
		AvatarURL: user.AvatarURL(""),
		CreatedAt: created,
		BotOwner:  userID == ws.config.Discord.OwnerID,
	}

	return jsonResponse(ctx, res, fasthttp.StatusOK)
}

// ---------------------------------------------------------------------------
// - GET /api/util/landingpageinfo

func (ws *WebServer) handlerGetLandingPageInfo(ctx *routing.Context) error {
	res := new(LandingPageResponse)

	publicInvites := true
	localInvite := true

	if ws.config.WebServer.LandingPage != nil {
		publicInvites = ws.config.WebServer.LandingPage.ShowPublicInvites
		localInvite = ws.config.WebServer.LandingPage.ShowLocalInvite
	}

	if publicInvites {
		res.PublicCanaryInvite = static.PublicCanaryInvite
		res.PublicMainInvite = static.PublicMainInvite
	}

	if localInvite {
		res.LocalInvite = util.GetInviteLink(ws.session)
	}

	return jsonResponse(ctx, res, fasthttp.StatusOK)
}

// ---------------------------------------------------------------------------
// - GET /api/sysinfo

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
		BotInvite: util.GetInviteLink(ws.session),

		Guilds: len(ws.session.State.Guilds),
	}

	return jsonResponse(ctx, info, fasthttp.StatusOK)
}

// ---------------------------------------------------------------------------
// - GET /api/util/color/:hexcode

func (ws *WebServer) handlerGetColor(ctx *routing.Context) error {
	hexcode := ctx.Param("hexcode")
	size := strings.ToLower(
		string(ctx.QueryArgs().Peek("size")))

	var xSize, ySize int
	var err error

	if size == "" {
		xSize, ySize = 24, 24
	} else if strings.Contains(size, "x") {
		split := strings.Split(size, "x")
		if len(split) != 2 {
			return jsonError(ctx, errors.New("invalid size parameter; must provide two size dimensions"), fasthttp.StatusBadRequest)
		}
		if xSize, err = strconv.Atoi(split[0]); err != nil {
			return jsonError(ctx, err, fasthttp.StatusBadRequest)
		}
		if ySize, err = strconv.Atoi(split[1]); err != nil {
			return jsonError(ctx, err, fasthttp.StatusBadRequest)
		}
	} else {
		if xSize, err = strconv.Atoi(size); err != nil {
			return jsonError(ctx, err, fasthttp.StatusBadRequest)
		}
		ySize = xSize
	}

	if xSize < 1 || ySize < 1 || xSize > 5000 || ySize > 5000 {
		return jsonError(ctx, errors.New("invalid size parameter; value must be in range [1..5000]"), fasthttp.StatusBadRequest)
	}

	clr, err := colors.FromHex(hexcode)
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusBadRequest)
	}

	buff, err := colors.CreateImage(clr, xSize, ySize)
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusBadRequest)
	}

	data := buff.Bytes()

	etag := etag.Generate(data, false)

	ctx.Response.Header.SetContentType("image/png")
	// 365 days browser caching
	ctx.Response.Header.Set("Cache-Control", "public, max-age=31536000, immutable")
	ctx.Response.Header.Set("ETag", etag)
	ctx.SetBody(data)

	return nil
}

// ---------------------------------------------------------------------------
// - GET /api/token

func (ws *WebServer) handlerGetToken(ctx *routing.Context) error {
	userID := ctx.Get("uid").(string)

	token, err := ws.db.GetAPIToken(userID)
	if database.IsErrDatabaseNotFound(err) {
		return jsonError(ctx, fmt.Errorf("no token found"), fasthttp.StatusNotFound)
	} else if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	tokenResp := &APITokenResponse{
		Created:    token.Created,
		Expires:    token.Expires,
		Hits:       token.Hits,
		LastAccess: token.LastAccess,
	}

	return jsonResponse(ctx, tokenResp, fasthttp.StatusOK)
}

// ---------------------------------------------------------------------------
// - POST /api/token

func (ws *WebServer) handlerPostToken(ctx *routing.Context) error {
	userID := ctx.Get("uid").(string)

	token, err := ws.auth.CreateAPIToken(userID)
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	return jsonResponse(ctx, token, fasthttp.StatusOK)
}

// ---------------------------------------------------------------------------
// - DELETE /api/token

func (ws *WebServer) handlerDeleteToken(ctx *routing.Context) error {
	userID := ctx.Get("uid").(string)

	err := ws.auth.db.DeleteAPIToken(userID)
	if database.IsErrDatabaseNotFound(err) {
		return jsonError(ctx, fmt.Errorf("no token found"), fasthttp.StatusNotFound)
	} else if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	return jsonResponse(ctx, nil, fasthttp.StatusOK)
}

// ---------------------------------------------------------------------------
// - GET /api/commands

func (ws *WebServer) handlerGetCommands(ctx *routing.Context) error {
	cmdInstances := ws.cmdhandler.GetCommandInstances()
	cmdInfos := make([]*CommandInfo, len(cmdInstances))

	for i, c := range cmdInstances {
		cmdInfo := GetCommandInfoFromCommand(c)
		cmdInfos[i] = cmdInfo
	}

	list := ListResponse{N: len(cmdInfos), Data: cmdInfos}
	return jsonResponse(ctx, list, fasthttp.StatusOK)
}
