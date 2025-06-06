package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/models"
	"github.com/zekroTJA/shinpuru/internal/services/config"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/services/timeprovider"
	"github.com/zekroTJA/shinpuru/internal/util/embedded"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/random"
)

// DatabaseRefreshTokenHandler implements RefreshTokenHandler
// for a base64 encoded token stored in the database
type DatabaseRefreshTokenHandler struct {
	db      database.Database
	session *discordgo.Session
	tp      timeprovider.Provider
}

// NewDatabaseRefreshTokenHandler returns a new instance
// of DatabaseRefreshTokenHandler
func NewDatabaseRefreshTokenHandler(container di.Container) *DatabaseRefreshTokenHandler {
	return &DatabaseRefreshTokenHandler{
		db:      container.Get(static.DiDatabase).(database.Database),
		session: container.Get(static.DiDiscordSession).(*discordgo.Session),
		tp:      container.Get(static.DiTimeProvider).(timeprovider.Provider),
	}
}

func (rth *DatabaseRefreshTokenHandler) GetRefreshToken(ident string) (token string, err error) {
	token, err = random.GetRandBase64Str(64)
	if err != nil {
		return
	}

	err = rth.db.SetUserRefreshToken(ident, token, rth.tp.Now().Add(static.AuthSessionExpiration))
	return
}

func (rth *DatabaseRefreshTokenHandler) ValidateRefreshToken(token string) (ident string, err error) {
	ident, expires, err := rth.db.GetUserByRefreshToken(token)
	if err != nil {
		return
	}

	if rth.tp.Now().After(expires) {
		err = errors.New("expired")
	}

	u, _ := rth.session.User(ident)
	if u == nil {
		err = errors.New("invalid user")
		return
	}

	return
}

func (rth *DatabaseRefreshTokenHandler) RevokeToken(ident string) error {
	err := rth.db.RevokeUserRefreshToken(ident)
	if database.IsErrDatabaseNotFound(err) {
		err = nil
	}
	return err
}

var (
	jwtGenerationMethod = jwt.SigningMethodHS256
)

// JWTAccessTokenHandler implements AccessTokenHandler
// for a JWT based access token
type JWTAccessTokenHandler struct {
	sessionExpiration time.Duration
	sessionSecret     []byte
	tp                timeprovider.Provider
}

// NewJWTAccessTokenHandler returns a new instance
// of JWTAccessTokenHandler
func NewJWTAccessTokenHandler(container di.Container) (ath *JWTAccessTokenHandler, err error) {
	cfg := container.Get(static.DiConfig).(config.Provider).
		Config().WebServer.AccessToken
	if err != nil {
		return nil, err
	}
	ath = &JWTAccessTokenHandler{
		sessionExpiration: time.Duration(cfg.LifetimeSeconds) * time.Second,
		sessionSecret:     []byte(cfg.Secret),
		tp:                container.Get(static.DiTimeProvider).(timeprovider.Provider),
	}
	return
}

func (ath *JWTAccessTokenHandler) GetAccessToken(ident string) (token string, expires time.Time, err error) {
	now := ath.tp.Now()
	expires = now.Add(ath.sessionExpiration)

	claims := jwt.RegisteredClaims{}
	claims.Issuer = fmt.Sprintf("shinpuru v.%s", embedded.AppVersion)
	claims.Subject = ident
	claims.ExpiresAt = jwt.NewNumericDate(expires)
	claims.NotBefore = jwt.NewNumericDate(now)
	claims.IssuedAt = jwt.NewNumericDate(now)

	token, err = jwt.NewWithClaims(jwtGenerationMethod, claims).
		SignedString(ath.sessionSecret)
	return
}

func (ath *JWTAccessTokenHandler) ValidateAccessToken(token string) (ident string, err error) {
	jwtToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return ath.sessionSecret, nil
	})
	if jwtToken == nil || err != nil || !jwtToken.Valid || jwtToken.Claims.Valid() != nil {
		return
	}

	claimsMap, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		err = errors.New("invalid claims")
		return
	}

	ident, _ = claimsMap["sub"].(string)

	return
}

type apiTokenClaims struct {
	jwt.RegisteredClaims

	Salt string `json:"sp_salt,omitempty"`
}

func apiTokenClaimsFromMap(m jwt.MapClaims) apiTokenClaims {
	c := apiTokenClaims{
		RegisteredClaims: standardClaimsFromMap(m),
	}

	c.Salt, _ = m["sp_salt"].(string)

	return c
}

func standardClaimsFromMap(m jwt.MapClaims) jwt.RegisteredClaims {
	c := jwt.RegisteredClaims{}

	c.Issuer, _ = m["iss"].(string)
	c.Subject, _ = m["sub"].(string)
	c.ExpiresAt, _ = m["exp"].(*jwt.NumericDate)
	c.NotBefore, _ = m["nbf"].(*jwt.NumericDate)
	c.IssuedAt, _ = m["iat"].(*jwt.NumericDate)

	return c
}

// DatabaseAPITokenHandler implements APITokenHandler
// for a JWT token stored in the database
type DatabaseAPITokenHandler struct {
	db      database.Database
	session *discordgo.Session
	secret  []byte
	tp      timeprovider.Provider
}

// NewDatabaseAPITokenHandler returns a new instance
// of DatabaseAPITokenHandler
func NewDatabaseAPITokenHandler(container di.Container) (*DatabaseAPITokenHandler, error) {
	cfg := container.Get(static.DiConfig).(config.Provider)
	secret := []byte(cfg.Config().WebServer.APITokenKey)

	return &DatabaseAPITokenHandler{
		db:      container.Get(static.DiDatabase).(database.Database),
		session: container.Get(static.DiDiscordSession).(*discordgo.Session),
		tp:      container.Get(static.DiTimeProvider).(timeprovider.Provider),
		secret:  secret,
	}, nil
}

func (apith *DatabaseAPITokenHandler) GetAPIToken(ident string) (token string, expires time.Time, err error) {
	now := apith.tp.Now()
	expires = now.Add(static.ApiTokenExpiration)

	salt, err := random.GetRandBase64Str(16)
	if err != nil {
		return
	}

	claims := apiTokenClaims{}
	claims.Issuer = fmt.Sprintf("shinpuru v.%s", embedded.AppVersion)
	claims.Subject = ident
	claims.ExpiresAt = jwt.NewNumericDate(expires)
	claims.NotBefore = jwt.NewNumericDate(now)
	claims.IssuedAt = jwt.NewNumericDate(now)
	claims.Salt = salt

	token, err = jwt.NewWithClaims(jwtGenerationMethod, claims).
		SignedString(apith.secret)
	if err != nil {
		return
	}

	tokenEntry := models.APITokenEntry{
		Salt:    salt,
		Created: now,
		Expires: expires,
		UserID:  ident,
	}

	if err = apith.db.SetAPIToken(tokenEntry); err != nil {
		return
	}

	return
}

func (apith *DatabaseAPITokenHandler) ValidateAPIToken(token string) (ident string, err error) {
	jwtToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return apith.secret, nil
	})
	if jwtToken == nil && err != nil {
		return "", err
	}
	if !jwtToken.Valid || jwtToken.Claims.Valid() != nil {
		return "", nil
	}

	claimsMap, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return "", nil
	}

	claims := apiTokenClaimsFromMap(claimsMap)

	tokenEntry, err := apith.db.GetAPIToken(claims.Subject)
	if database.IsErrDatabaseNotFound(err) {
		return "", nil
	} else if err != nil {
		return "", err
	}

	if tokenEntry.Salt != claims.Salt {
		return "", err
	}

	tokenEntry.Hits++
	tokenEntry.LastAccess = apith.tp.Now()
	apith.db.SetAPIToken(tokenEntry)

	return claims.Subject, nil
}

func (apith *DatabaseAPITokenHandler) RevokeToken(ident string) error {
	return apith.db.DeleteAPIToken(ident)
}
