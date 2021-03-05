package webserver

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/bwmarrin/snowflake"
	"github.com/gabriel-vasile/mimetype"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"github.com/zekroTJA/shinpuru/internal/core/database"
	"github.com/zekroTJA/shinpuru/internal/util/imgstore"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/etag"
)

var (
	errInvalidOTAToken = errors.New("invalid one time auth token")
	errOTADisabled     = errors.New("ota disabled by user")
)

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
// - GET /ota

func (ws *WebServer) handlerGetOta(ctx *routing.Context) error {
	token := string(ctx.QueryArgs().Peek("token"))

	if token == "" {
		return jsonError(ctx, errInvalidOTAToken, fasthttp.StatusUnauthorized)
	}

	userID, err := ws.ota.ValidateKey(token)
	if err != nil {
		return jsonError(ctx, errInvalidOTAToken, fasthttp.StatusUnauthorized)
	}

	enabled, err := ws.db.GetUserOTAEnabled(userID)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	if !enabled {
		return jsonError(ctx, errOTADisabled, fasthttp.StatusUnauthorized)
	}

	if ch, err := ws.session.UserChannelCreate(userID); err == nil {
		ipaddr := getIPAddr(ctx)
		useragent := string(ctx.Request.Header.UserAgent())
		emb := &discordgo.MessageEmbed{
			Color: static.ColorEmbedOrange,
			Description: fmt.Sprintf("Someone logged in to the web interface as you.\n"+
				"\n**Details:**\nIP Address: ||`%s`||\nUser Agent: `%s`\n\n"+
				"If this was not you, consider disabling OTA [**here**](%s/usersettings).",
				ipaddr, useragent, ws.config.WebServer.PublicAddr),
			Timestamp: time.Now().Format(time.RFC3339),
		}
		ws.session.ChannelMessageSendEmbed(ch.ID, emb)
	}

	return ws.auth.LoginSuccessHandler(ctx, userID)
}
