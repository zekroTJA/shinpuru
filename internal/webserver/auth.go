package webserver

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"github.com/zekroTJA/shinpuru/internal/core"
	"github.com/zekroTJA/shinpuru/pkg/random"
)

var (
	sessionExpiration = 7 * 24 * time.Hour
)

const (
	sessionKeyLength = 128
)

type Auth struct {
	db      core.Database
	session *discordgo.Session
}

func NewAuth(db core.Database, s *discordgo.Session) *Auth {
	return &Auth{
		db:      db,
		session: s,
	}
}

func (auth *Auth) createSessionKey() (string, error) {
	return random.GetRandBase64Str(sessionKeyLength)
}

func (auth *Auth) checkSessionCookie(ctx *routing.Context) (string, error) {
	key := ctx.Request.Header.Cookie("__session")
	if key == nil || len(key) == 0 {
		return "", nil
	}

	skey := string(key)
	uid, err := auth.db.GetSession(skey)
	if err == core.ErrDatabaseNotFound {
		return "", nil
	}
	if err != nil {
		return "", err
	}

	return uid, nil
}

func (auth *Auth) checkAuth(ctx *routing.Context) error {
	uid, err := auth.checkSessionCookie(ctx)
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	if uid != "" {
		ctx.Set("uid", uid)
		return nil
	}

	path := string(ctx.Path())

	if strings.HasPrefix(path, "/api") {
		return jsonError(ctx, errUnauthorized, fasthttp.StatusUnauthorized)
	}

	if strings.HasSuffix(path, ".js") ||
		strings.HasSuffix(path, ".css") ||
		strings.HasPrefix(path, "/assets") ||
		strings.HasPrefix(path, "/favicon.ico") ||
		path == "/login" || path == endpointLogInWithDC || path == endpointAuthCB {
		return nil
	}

	ctx.Redirect("/login", fasthttp.StatusTemporaryRedirect)
	ctx.Abort()

	return nil
}

func (auth *Auth) LoginFailedHandler(ctx *routing.Context, status int, msg string) error {
	return jsonResponse(ctx, nil, fasthttp.StatusUnauthorized)
}

func (auth *Auth) LoginSuccessHandler(ctx *routing.Context, uid string) error {
	if u, _ := auth.session.User(uid); u == nil {
		return jsonError(ctx, errUnauthorized, fasthttp.StatusUnauthorized)
	}

	ctx.Set("uid", uid)

	sessionKey, err := auth.createSessionKey()
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	expires := time.Now().Add(sessionExpiration)

	if err = auth.db.SetSession(sessionKey, uid, expires); err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	cookie := fmt.Sprintf("__session=%s; Expires=%s; Path=/; HttpOnly",
		sessionKey, expires.Format(time.RFC1123))

	ctx.Response.Header.Set("Set-Cookie", cookie)

	ctx.Redirect("/", fasthttp.StatusTemporaryRedirect)
	ctx.Abort()
	return nil
}

func (auth *Auth) LogOutHandler(ctx *routing.Context) error {
	userID := ctx.Get("uid").(string)

	auth.db.DeleteSession(userID)

	cookie := "__session=; Expires=Thu, 01 Jan 1970 00:00:00 GMT; Path=/; HttpOnly"
	ctx.Response.Header.Set("Set-Cookie", cookie)

	return jsonResponse(ctx, nil, fasthttp.StatusOK)

}
