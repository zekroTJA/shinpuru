package webserver

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	jwt "github.com/dgrijalva/jwt-go"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"

	"github.com/zekroTJA/shinpuru/internal/core/database"
	"github.com/zekroTJA/shinpuru/internal/shared/models"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/pkg/lctimer"
	"github.com/zekroTJA/shinpuru/pkg/random"
)

var (
	sessionExpiration       = 7 * 24 * time.Hour
	sessionSecretExpiration = sessionExpiration * 2
	apiTokenExpiration      = 365 * 24 * time.Hour

	jwtGenerationMethod = jwt.SigningMethodHS256
)

const (
	sessionKeyLength = 128

	jwtIssuer = "shinpuru v.%v"
)

// Auth provides handlers and untilities to authorize a
// HTTP request by session cookie.
type Auth struct {
	db      database.Database
	session *discordgo.Session

	tokenSecret            []byte
	sessionSecret          []byte
	sessionSecretRefreshed time.Time
}

// NewAuth initializes a new Auth instance with the passed
// database provider and discordgo session.
func NewAuth(db database.Database, s *discordgo.Session,
	lct *lctimer.LifeCycleTimer, tokenSecret []byte) (auth *Auth, err error) {

	auth = &Auth{
		db:          db,
		session:     s,
		tokenSecret: tokenSecret,
	}

	err = auth.RefreshSessionSecret()
	auth.sessionSecretRefreshed = time.Now()

	// Refresh sessionSecret key everytime after sessionSecretExpiration
	// has expired from auth.sessionSecretRefreshed.
	lct.OnTick(func(now time.Time) {
		if now.Sub(auth.sessionSecretRefreshed) > sessionSecretExpiration {
			if err := auth.RefreshSessionSecret(); err != nil {
				util.Log.Errorf("failed refreshing auth session secret: %s", err.Error())
			}
			auth.sessionSecretRefreshed = now
		}
	})

	return
}

// RefreshSessionSecret randomly generates a new JWT key
// for session key generation and signing.
func (auth *Auth) RefreshSessionSecret() (err error) {
	auth.sessionSecret, err = random.GetRandByteArray(32)
	return
}

// LoginFailedHandler returns a 401 Unauthorized response.
func (auth *Auth) LoginFailedHandler(ctx *routing.Context, status int, msg string) error {
	return jsonResponse(ctx, nil, fasthttp.StatusUnauthorized)
}

// LoginSuccessHandler fetches the user by uid. If the user
// exists, the user ID is set to the context as value for "id".
func (auth *Auth) LoginSuccessHandler(ctx *routing.Context, uid string) error {
	user, _ := auth.session.User(uid)
	if user == nil {
		return jsonError(ctx, errUnauthorized, fasthttp.StatusUnauthorized)
	}

	ctx.Set("uid", uid)

	sessionKey, err := auth.createSessionKey(uid)
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	expires := time.Now().Add(sessionExpiration)

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
	cookie := "__session=; Expires=Thu, 01 Jan 1970 00:00:00 GMT; Path=/; HttpOnly"
	ctx.Response.Header.Set("Set-Cookie", cookie)

	return jsonResponse(ctx, nil, fasthttp.StatusOK)

}

// CreateAPIToken creates a JWT API token for the passed userID,
// sets it to the database and returns the token information.
func (auth *Auth) CreateAPIToken(userID string) (*APITokenResponse, error) {
	now := time.Now()
	expires := now.Add(apiTokenExpiration)

	salt, err := random.GetRandBase64Str(16)
	if err != nil {
		return nil, err
	}

	claims := APITokenClaims{}
	claims.Issuer = fmt.Sprintf(jwtIssuer, util.AppVersion)
	claims.Subject = userID
	claims.ExpiresAt = expires.Unix()
	claims.NotBefore = now.Unix()
	claims.IssuedAt = now.Unix()
	claims.Salt = salt

	token, err := jwt.NewWithClaims(jwtGenerationMethod, claims).
		SignedString(auth.tokenSecret)
	if err != nil {
		return nil, err
	}

	tokenEntry := &models.APITokenEntry{
		Salt:    salt,
		Created: now,
		Expires: expires,
		UserID:  userID,
	}

	if err = auth.db.SetAPIToken(tokenEntry); err != nil {
		return nil, err
	}

	return &APITokenResponse{
		Created: now,
		Expires: expires,
		Token:   token,
	}, nil
}

// createSessionKey creates a JWT string with the passed
// userID and sessionExpiration as expiration value.
func (auth *Auth) createSessionKey(userID string) (string, error) {
	now := time.Now()
	expires := now.Add(sessionExpiration)

	claims := APITokenClaims{}
	claims.Issuer = fmt.Sprintf(jwtIssuer, util.AppVersion)
	claims.Subject = userID
	claims.ExpiresAt = expires.Unix()
	claims.NotBefore = now.Unix()
	claims.IssuedAt = now.Unix()

	token, err := jwt.NewWithClaims(jwtGenerationMethod, claims).
		SignedString(auth.sessionSecret)
	if err != nil {
		return "", err
	}

	return token, nil
}

// checkSessionCookie checks the set cookie for session key of
// the request by validating the JWT signature against the
// specified signing key and obtains the user ID from the JWT.
//
// This function only returns an error when the check fails
// unexpectedly. When the key was invalid, an empty string and
// no error is returned.
func (auth *Auth) checkSessionCookie(ctx *routing.Context) (string, error) {
	key := ctx.Request.Header.Cookie("__session")
	if key == nil || len(key) == 0 {
		return "", nil
	}

	keyStr := string(key)

	token, err := jwt.Parse(keyStr, func(t *jwt.Token) (interface{}, error) {
		return auth.sessionSecret, nil
	})
	if token == nil || err != nil || !token.Valid || token.Claims.Valid() != nil {
		return "", nil
	}

	claimsMap, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", nil
	}

	claims := SessionTokenClaimsFromMap(claimsMap)

	return claims.Subject, nil
}

// checkAPIToken checks for a set Authorization header with a
// bearer token and tries to authenticate this token.
// If the token is valid, the user ID is returned.
//
// Occuring errors during authenticatrion are returned.
// Invalid authentication does not return any errors.
func (auth *Auth) checkAPIToken(ctx *routing.Context) (string, error) {
	key := ctx.Request.Header.Peek("Authorization")
	if key == nil || len(key) == 0 {
		return "", nil
	}

	keyStr := string(key)
	if !strings.HasPrefix(strings.ToLower(keyStr), "bearer ") {
		return "", nil
	}
	keyStr = keyStr[7:]

	token, err := jwt.Parse(keyStr, func(t *jwt.Token) (interface{}, error) {
		return auth.tokenSecret, nil
	})
	if token == nil && err != nil {
		return "", err
	}
	if !token.Valid || token.Claims.Valid() != nil {
		return "", nil
	}

	claimsMap, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", nil
	}

	claims := APITokenClaimsFromMap(claimsMap)

	tokenEntry, err := auth.db.GetAPIToken(claims.Subject)
	if database.IsErrDatabaseNotFound(err) {
		return "", nil
	} else if err != nil {
		return "", err
	}

	if tokenEntry.Salt != claims.Salt {
		return "", err
	}

	tokenEntry.Hits++
	tokenEntry.LastAccess = time.Now()
	auth.db.SetAPIToken(tokenEntry)

	return claims.Subject, nil
}

// checkAuth wraps checkSessionCookie as request handler.
//
// This handler also ignores requests to static files
// which are accessable even with no authorized session.
//
// If the authorization failed, a redirect to /login is sent.
func (auth *Auth) checkAuth(ctx *routing.Context) error {
	var usedAPIToken bool

	uid, err := auth.checkSessionCookie(ctx)
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	if uid == "" {
		uid, err = auth.checkAPIToken(ctx)
		if err != nil {
			return jsonError(ctx, err, fasthttp.StatusInternalServerError)
		}
		usedAPIToken = true
	}

	if uid != "" {
		ctx.Set("uid", uid)
		ctx.Set("usedApiToken", usedAPIToken)
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
