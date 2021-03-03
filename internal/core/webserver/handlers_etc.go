package webserver

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/bwmarrin/snowflake"
	"github.com/gabriel-vasile/mimetype"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"github.com/zekroTJA/shinpuru/internal/util/imgstore"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/etag"
)

var errInvalidOTAToken = errors.New("invalid one time auth token")

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

	return ws.auth.LoginSuccessHandler(ctx, userID)
}
