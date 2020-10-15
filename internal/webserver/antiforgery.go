package webserver

import (
	"errors"
	"strings"
	"time"

	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"github.com/zekroTJA/shinpuru/pkg/random"
	"github.com/zekroTJA/timedmap"
)

const (
	afLoginTokenLifetime   = 10 * time.Minute
	afSessionTokenLifetime = 2 * time.Hour
	afTokenCookieName      = "XSRF-TOKEN"
	afTokenHeaderName      = "x-xsrf-token"
	afUidIdent             = "uid"
)

// AntiForgery wraps HTTP handlers binding and
// recovering anti forgery tokens to the requests
// and from responses to protect against XSRF
// attacks.
type AntiForgery struct {
	tm *timedmap.TimedMap

	loginTokens   timedmap.Section
	sessionTokens timedmap.Section
}

// NewAntiForgery initializes a new instance
// of AntiForgery.
func NewAntiForgery() (af *AntiForgery) {
	af = new(AntiForgery)

	af.tm = timedmap.New(10 * time.Minute)
	af.loginTokens = af.tm.Section(0)
	af.sessionTokens = af.tm.Section(1)

	return
}

// Handler is the single-point HTTP handler to bind
// and recover anti-forgery tokens to requests and
// from responses.
func (af *AntiForgery) Handler(ctx *routing.Context) (err error) {
	var token string

	path := string(ctx.Path())
	method := string(ctx.Method())

	// The temporarily used generic af-token should be set on the
	// first GET request. Because an SPA is used, the first request
	// might not be '/' or '/index.html'. Because the SPA first
	// calls 'GET /main.js', the temporary token is set here.
	if method == "GET" && strings.HasPrefix(path, "/main") && strings.HasSuffix(path, ".js") {
		token, err = af.generateToken()
		if err != nil {
			ctx.Abort()
			return jsonError(ctx, err, fasthttp.StatusInternalServerError)
		}
		af.loginTokens.Set(token, struct{}{}, afLoginTokenLifetime)
		af.setCookie(ctx, afLoginTokenLifetime, token)
		return
	}

	userId, _ := ctx.Get("uid").(string)
	token = af.recoverToken(ctx)

	// If a user has authenticated successfully, the request method
	// is GET and no token was recovered from headers, a session
	// af-token is generated, stored and set as cookie to the response.
	if userId != "" && method == "GET" && token == "" {
		token, err = af.generateToken()
		if err != nil {
			ctx.Abort()
			return jsonError(ctx, err, fasthttp.StatusInternalServerError)
		}
		af.sessionTokens.Set(token, userId, afSessionTokenLifetime)
		af.setCookie(ctx, afSessionTokenLifetime, token)
		return
	}

	// Further, only POST, DELETE and PUT requests must be validated
	// using the anti-forgery token header.
	if method != "POST" && method != "DELETE" && method != "PUT" {
		return
	}

	// If no token was recovered, error with code 400.
	if token == "" {
		return af.errorReqeust(ctx)
	}

	// When user is not authenticated, validate the recovered token
	// against the list of temporary login tokens. Otherwise, validate
	// against the session tokens which are bound to the user IDs.
	if userId == "" {
		if !af.loginTokens.Contains(token) {
			return af.errorReqeust(ctx)
		}
	} else {
		recUserId, ok := af.sessionTokens.GetValue(token).(string)
		if !ok || recUserId != userId {
			return af.errorReqeust(ctx)
		}
	}

	return
}

// generateToken generates a cryptographically random
// base64 token with a length of 32 characters.
func (af *AntiForgery) generateToken() (string, error) {
	return random.GetRandBase64Str(32)
}

// setCookie sets the passed anti-forgery token as cookie
// with the specified lifetime to the response.
func (af *AntiForgery) setCookie(ctx *routing.Context, lifetime time.Duration, token string) {
	cookie := new(fasthttp.Cookie)
	cookie.SetKey(afTokenCookieName)
	cookie.SetExpire(time.Now().Add(lifetime))
	cookie.SetHTTPOnly(false)
	cookie.SetSameSite(fasthttp.CookieSameSiteStrictMode)
	cookie.SetValue(token)
	ctx.Response.Header.SetCookie(cookie)
}

// recoverToken tries to get the anti-forgery token form
// the requests and returns it. If no token was found, an
// empty string is returned.
func (af *AntiForgery) recoverToken(ctx *routing.Context) string {
	res := ctx.Request.Header.Peek(afTokenHeaderName)
	if res == nil {
		return ""
	}
	return string(res)
}

// errorRequest sends a 400 bad request response to the
// passed request context when the XSRF token validation
// failed.
func (af *AntiForgery) errorReqeust(ctx *routing.Context) error {
	ctx.Abort()
	return jsonError(ctx, errors.New("xsrf validation failed"), fasthttp.StatusBadRequest)
}
