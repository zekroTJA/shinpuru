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

type AntiForgery struct {
	tm *timedmap.TimedMap

	loginTokens   timedmap.Section
	sessionTokens timedmap.Section
}

func NewAntiForgery() (af *AntiForgery) {
	af = new(AntiForgery)

	af.tm = timedmap.New(10 * time.Minute)
	af.loginTokens = af.tm.Section(0)
	af.sessionTokens = af.tm.Section(1)

	return
}

func (af *AntiForgery) Handler(ctx *routing.Context) (err error) {
	var token string

	path := string(ctx.Path())
	method := string(ctx.Method())

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

	if method != "POST" && method != "DELETE" && method != "PUT" {
		return
	}

	if token == "" {
		return af.errorReqeust(ctx)
	}

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

func (af *AntiForgery) generateToken() (string, error) {
	return random.GetRandBase64Str(32)
}

func (af *AntiForgery) setCookie(ctx *routing.Context, lifetime time.Duration, token string) {
	cookie := new(fasthttp.Cookie)
	cookie.SetKey(afTokenCookieName)
	cookie.SetExpire(time.Now().Add(lifetime))
	cookie.SetHTTPOnly(false)
	cookie.SetSameSite(fasthttp.CookieSameSiteStrictMode)
	cookie.SetValue(token)
	ctx.Response.Header.SetCookie(cookie)
}

func (af *AntiForgery) recoverToken(ctx *routing.Context) string {
	res := ctx.Request.Header.Peek(afTokenHeaderName)
	if res == nil {
		return ""
	}
	return string(res)
}

func (af *AntiForgery) errorReqeust(ctx *routing.Context) error {
	ctx.Abort()
	return jsonError(ctx, errors.New("xsrf validation failed"), fasthttp.StatusBadRequest)
}
