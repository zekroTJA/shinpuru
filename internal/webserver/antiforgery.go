package webserver

import (
	"errors"
	"time"

	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"github.com/zekroTJA/shinpuru/pkg/random"
	"github.com/zekroTJA/timedmap"
)

const (
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

	sessionTokens timedmap.Section
}

// NewAntiForgery initializes a new instance
// of AntiForgery.
func NewAntiForgery() (af *AntiForgery) {
	af = new(AntiForgery)

	af.tm = timedmap.New(10 * time.Minute)
	af.sessionTokens = af.tm.Section(0)

	return
}

// SessionSetHandler sets session af-tokens on each
// GET request to the endpoint(s) bound to this handler.
func (af *AntiForgery) SessionSetHandler(ctx *routing.Context) (err error) {
	if af.isApiToken(ctx) {
		return
	}

	var token string

	method := string(ctx.Method())
	userId, _ := ctx.Get("uid").(string)

	// If a user has authenticated successfully, the request method
	// is GET and no token was recovered from headers, a session
	// af-token is generated, stored and set as cookie to the response.
	if userId != "" && method == "GET" {
		token, err = af.generateToken()
		if err != nil {
			ctx.Abort()
			return jsonError(ctx, err, fasthttp.StatusInternalServerError)
		}
		af.sessionTokens.Set(userId, token, afSessionTokenLifetime)
		af.setCookie(ctx, afSessionTokenLifetime, token)
		return
	}

	return
}

// Handler checks for the af-token header if the
// request is a POST, PUT or DELETE request and errors
// with code 400 when no session af-token was recovered
// or the recovered token did not match the saved token
// for the authenticated user.
func (af *AntiForgery) Handler(ctx *routing.Context) (err error) {
	if af.isApiToken(ctx) {
		return
	}

	method := string(ctx.Method())

	// Further, only POST, DELETE and PUT requests must be validated
	// using the anti-forgery token header.
	if method != "POST" && method != "DELETE" && method != "PUT" {
		return
	}

	userId, _ := ctx.Get("uid").(string)
	token := af.recoverToken(ctx)

	// If no token was recovered, error with code 400.
	if token == "" {
		return af.errorReqeust(ctx)
	}

	// When user is not authenticated, validate the recovered token
	// against the list of temporary login tokens. Otherwise, validate
	// against the session tokens which are bound to the user IDs.
	if userId != "" {
		recToken, ok := af.sessionTokens.GetValue(userId).(string)
		if !ok || recToken != token {
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
	cookie.SetPath("/")
	cookie.SetExpire(time.Now().Add(lifetime * 2))
	cookie.SetHTTPOnly(false)
	cookie.SetSameSite(fasthttp.CookieSameSiteStrictMode)
	cookie.SetValue(token)
	ctx.Response.Header.SetCookie(cookie)
}

// isApiToken is true if the request wwas authenticated
// using an API bearer token.
func (af *AntiForgery) isApiToken(ctx *routing.Context) bool {
	usedApiToken, _ := ctx.Get("usedApiToken").(bool)
	return usedApiToken
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
