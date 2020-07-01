package webserver

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"

	"github.com/zekroTJA/shinpuru/internal/core/database"
	"github.com/zekroTJA/shinpuru/pkg/random"
)

var (
	sessionExpiration = 7 * 24 * time.Hour
)

const (
	sessionKeyLength = 128
)

// Auth provides handlers and untilities to authorize a
// HTTP request by session cookie.
type Auth struct {
	db      database.Database
	session *discordgo.Session
}

// NewAuth initializes a new Auth instance with the passed
// database provider and discordgo session.
func NewAuth(db database.Database, s *discordgo.Session) *Auth {
	return &Auth{
		db:      db,
		session: s,
	}
}

// LoginFailedHandler returns a 401 Unauthorized response.
func (auth *Auth) LoginFailedHandler(ctx *routing.Context, status int, msg string) error {
	return jsonResponse(ctx, nil, fasthttp.StatusUnauthorized)
}

// LoginSuccessHandler fetches the user by uid. If the user
// exists, the user ID is set to the context as value for "id".
// Also, a randomly generated session key is set as session key
// cookie and to the database to validate a session.
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

// LogOutHandler removes the session key of the authenticated
// user from the database and sends an unset cookie.
func (auth *Auth) LogOutHandler(ctx *routing.Context) error {
	userID := ctx.Get("uid").(string)

	auth.db.DeleteSession(userID)

	cookie := "__session=; Expires=Thu, 01 Jan 1970 00:00:00 GMT; Path=/; HttpOnly"
	ctx.Response.Header.Set("Set-Cookie", cookie)

	return jsonResponse(ctx, nil, fasthttp.StatusOK)

}

// createSessionKey randomly generates a base64 string with the
// length of sessionKeyLength.
func (auth *Auth) createSessionKey() (string, error) {
	return random.GetRandBase64Str(sessionKeyLength)
}

// checkSessionCookie checks the set cookie for session key of
// the request against the database. If the value exists and could
// be matched, the session is authorized.
func (auth *Auth) checkSessionCookie(ctx *routing.Context) (string, error) {
	key := ctx.Request.Header.Cookie("__session")
	if key == nil || len(key) == 0 {
		return "", nil
	}

	skey := string(key)
	uid, err := auth.db.GetSession(skey)
	if database.IsErrDatabaseNotFound(err) {
		return "", nil
	}
	if err != nil {
		return "", err
	}

	return uid, nil
}

// checkAuth wraps checkSessionCookie as request handler.
//
// This handler also ignores requests to static files
// which are accessable even with no authorized session.
//
// If the authorization failed, a redirect to /login is sent.
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

	ctx.Redirect("/login", fasthttp.StatusTemporaryRedirect)
	ctx.Abort()

	return nil
}
