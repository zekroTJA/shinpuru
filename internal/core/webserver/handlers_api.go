package webserver

import (
	"errors"
	"fmt"
	"io"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/gabriel-vasile/mimetype"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"

	"github.com/zekroTJA/shinpuru/internal/core/database"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/imgstore"
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
		BotInvite: fmt.Sprintf("https://discord.com/api/oauth2/authorize?client_id=%s&scope=bot&permissions=%d",
			ws.session.State.User.ID, static.InvitePermission),

		Guilds: len(ws.session.State.Guilds),
	}

	return jsonResponse(ctx, info, fasthttp.StatusOK)
}

// ---------------------------------------------------------------------------
// - GET /imagestore/:id

func (ws *WebServer) handlerGetImage(ctx *routing.Context) error {
	path := ctx.Param("id")

	pathSplit := strings.Split(path, ".")
	imageIDstr := pathSplit[0]

	imageID, err := snowflake.ParseString(imageIDstr)
	if err != nil {
		return jsonError(ctx, fmt.Errorf("invalid snowflake ID"), fasthttp.StatusBadRequest)
	}

	reader, size, err := ws.st.GetObject(static.StorageBucketImages, imageID.String())
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusBadRequest)
	}

	defer reader.Close()

	img := new(imgstore.Image)

	img.Size = int(size)
	img.Data = make([]byte, img.Size)
	_, err = reader.Read(img.Data)
	if err != nil && err != io.EOF {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	img.MimeType = mimetype.Detect(img.Data).String()

	etag := etag.Generate(img.Data, false)

	ctx.Response.Header.SetContentType(img.MimeType)
	// 30 days browser caching
	ctx.Response.Header.Set("Cache-Control", "public, max-age=2592000, immutable")
	ctx.Response.Header.Set("ETag", etag)
	ctx.SetBody(img.Data)

	return nil
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
