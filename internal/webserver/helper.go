package webserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/zekroTJA/shinpuru/internal/core"

	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"github.com/zekroTJA/shinpuru/internal/util"
)

var emptyResponseBody = []byte("{}")

var (
	headerXForwardedFor = []byte("X-Forwarded-For")
	headerUserAgent     = []byte("User-Agent")
)

var defStatusBoddies = map[int][]byte{
	http.StatusOK:           []byte("{\n  \"code\": 200,\n  \"message\": \"ok\"\n}"),
	http.StatusCreated:      []byte("{\n  \"code\": 201,\n  \"message\": \"created\"\n}"),
	http.StatusNotFound:     []byte("{\n  \"code\": 404,\n  \"message\": \"not found\"\n}"),
	http.StatusUnauthorized: []byte("{\n  \"code\": 401,\n  \"message\": \"unauthorized\"\n}"),
}

// jsonError writes the error message of err and the
// passed status to response context and aborts the
// execution of following registered handlers ONLY IF
// err != nil.
// This function always returns a nil error that the
// default error handler can be bypassed.
func jsonError(ctx *routing.Context, err error, status int) error {
	if err != nil {
		ctx.Response.Header.SetContentType("application/json")
		ctx.SetStatusCode(status)
		ctx.SetBodyString(fmt.Sprintf("{\n  \"code\": %d,\n  \"message\": \"%s\"\n}",
			status, err.Error()))
		ctx.Abort()
	}
	return nil
}

// jsonResponse tries to parse the passed interface v
// to JSON and writes it to the response context body
// as same as the passed status code.
// If the parsing fails, this will result in a jsonError
// output of the error with status 500.
// This function always returns a nil error.
func jsonResponse(ctx *routing.Context, v interface{}, status int) error {
	var err error
	data := emptyResponseBody

	if v == nil {
		if d, ok := defStatusBoddies[status]; ok {
			data = d
		}
	} else {
		if util.Release != "TRUE" {
			data, err = json.MarshalIndent(v, "", "  ")
		} else {
			data, err = json.Marshal(v)
		}
		if err != nil {
			return jsonError(ctx, err, fasthttp.StatusInternalServerError)
		}
	}

	ctx.Response.Header.SetContentType("application/json")
	ctx.SetStatusCode(status)
	_, err = ctx.Write(data)

	return jsonError(ctx, err, fasthttp.StatusInternalServerError)
}

// parseJSONBody tries to parse a requests JSON
// body to the passed object pointer. If the
// parsing fails, this will result in a jsonError
// output with status 400.
// This function always returns a nil error.
func parseJSONBody(ctx *routing.Context, v interface{}) error {
	data := ctx.PostBody()
	err := json.Unmarshal(data, v)
	if err != nil {
		jsonError(ctx, err, fasthttp.StatusBadRequest)
	}
	return err
}

func (ws *WebServer) addHeaders(ctx *routing.Context) error {
	ctx.Response.Header.SetServer("shinpuru v." + util.AppVersion)

	if util.Release != "TRUE" {
		ctx.Response.Header.Set("Access-Control-Allow-Origin", "http://localhost:8081")
		ctx.Response.Header.Set("Access-Control-Allow-Headers", "authorization, content-type, set-cookie, cookie, server")
		ctx.Response.Header.Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, POST, DELETE, OPTIONS")
		ctx.Response.Header.Set("Access-Control-Allow-Credentials", "true")
	}

	return nil
}

func getIPAddr(ctx *routing.Context) string {
	forwardedfor := ctx.Request.Header.PeekBytes(headerXForwardedFor)
	if forwardedfor != nil && len(forwardedfor) > 0 {
		return string(forwardedfor)
	}

	return ctx.RemoteIP().String()
}

func (ws *WebServer) handlerFiles(ctx *routing.Context) error {
	path := string(ctx.Path())

	if strings.HasPrefix(path, "/api") || strings.HasPrefix(path, "/_/") {
		ctx.Next()
		return nil
	}

	if strings.HasSuffix(path, ".js") ||
		strings.HasSuffix(path, ".css") ||
		strings.HasPrefix(path, "/assets") ||
		strings.HasPrefix(path, "/favicon.ico") {

		fileHandlerStatic.NewRequestHandler()(ctx.RequestCtx)
		ctx.Abort()
		return nil
	}

	ctx.SendFile("./web/dist/web/index.html")
	ctx.Abort()
	return nil
}

func (ws *WebServer) optionsHandler(ctx *routing.Context) error {
	if string(ctx.Method()) == "OPTIONS" {
		ctx.SetStatusCode(fasthttp.StatusOK)
		ctx.Abort()
	}
	return nil
}

func errInternalOrNotFound(ctx *routing.Context, err error) error {
	if core.IsErrDatabaseNotFound(err) {
		return jsonError(ctx, errNotFound, fasthttp.StatusNotFound)
	}
	return jsonError(ctx, err, fasthttp.StatusInternalServerError)
}

func errInternalIgnoreNotFound(ctx *routing.Context, err error) (bool, error) {
	if core.IsErrDatabaseNotFound(err) {
		return false, nil
	}
	return true, jsonError(ctx, err, fasthttp.StatusInternalServerError)
}
