package webserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"github.com/zekroTJA/shinpuru/internal/core/database"
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

// addHeaders adds the server header for shinpuru backend
// and 'X-Content-Type-Options' to 'nosniff'.
//
// If util.Release is not "TRUE", CORS headers are added to
// allow access from the angular dev webserver.
func (ws *WebServer) addHeaders(ctx *routing.Context) error {
	ctx.Response.Header.SetServer("shinpuru v." + util.AppVersion)
	ctx.Response.Header.Set("X-Content-Type-Options", "nosniff")

	if util.Release != "TRUE" {
		ctx.Response.Header.Set("Access-Control-Allow-Origin", ws.config.WebServer.DebugPublicAddr)
		ctx.Response.Header.Set("Access-Control-Allow-Headers", "authorization, content-type, set-cookie, cookie, server")
		ctx.Response.Header.Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, POST, DELETE, OPTIONS")
		ctx.Response.Header.Set("Access-Control-Allow-Credentials", "true")
	}

	return nil
}

// getIPAddr returns the IP address of the request origin.
// This method firstly consumes the "x-Forwarded-For" header,
// which is set by most reverse proxies. If this header value
// is empty, the actual origin address is returned.
func getIPAddr(ctx *routing.Context) string {
	forwardedfor := ctx.Request.Header.PeekBytes(headerXForwardedFor)
	if forwardedfor != nil && len(forwardedfor) > 0 {
		return string(forwardedfor)
	}

	return ctx.RemoteIP().String()
}

// handlerFiles is the request handler for SPA file routing.
func (ws *WebServer) handlerFiles(ctx *routing.Context) error {
	path := string(ctx.Path())

	if strings.HasPrefix(path, "/api/") || strings.HasPrefix(path, "/imagestore/") || strings.HasPrefix(path, "/_/") {
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

// optionsHandler handles OPTIONS requestst sent by
// brwoser as CORS preflight requests.
func (ws *WebServer) optionsHandler(ctx *routing.Context) error {
	if string(ctx.Method()) == "OPTIONS" {
		ctx.SetStatusCode(fasthttp.StatusOK)
		ctx.Abort()
	}
	return nil
}

// errInternalOrNotFound responds with a 404 Not Found
// if the passed err equals DatabaseNotFound.
// Otherwise, a 500 Internal Server Error is returned with
// the error informatio nas body data.
func errInternalOrNotFound(ctx *routing.Context, err error) error {
	if database.IsErrDatabaseNotFound(err) {
		return jsonError(ctx, errNotFound, fasthttp.StatusNotFound)
	}
	return jsonError(ctx, err, fasthttp.StatusInternalServerError)
}

// errInternalIgnoreNotFound only returns a 500 Internal
// Server Error if the passed err does not equal
// ErrDatabaseNotFound. Otherwise, no action is taken.
func errInternalIgnoreNotFound(ctx *routing.Context, err error) (bool, error) {
	if database.IsErrDatabaseNotFound(err) {
		return false, nil
	}
	return true, jsonError(ctx, err, fasthttp.StatusInternalServerError)
}
