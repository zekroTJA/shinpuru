package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/dgrijalva/jwt-go"
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/config"
	"github.com/zekroTJA/shinpuru/internal/models"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/random"
)

// DatabaseRefreshTokenHandler implements RefreshTokenHandler
// for a base64 encoded token stored in the database
type DatabaseRefreshTokenHandler struct {
	db      database.Database
	session *discordgo.Session
}

// NewDatabaseRefreshTokenHandler returns a new instance
// of DatabaseRefreshTokenHandler
func NewDatabaseRefreshTokenHandler(container di.Container) *DatabaseRefreshTokenHandler {
	return &DatabaseRefreshTokenHandler{
		db:      container.Get(static.DiDatabase).(database.Database),
		session: container.Get(static.DiDiscordSession).(*discordgo.Session),
	}
}

func (rth *DatabaseRefreshTokenHandler) GetRefreshToken(ident string) (token string, err error) {
	token, err = random.GetRandBase64Str(64)
	if err != nil {
		return
	}

	err = rth.db.SetUserRefreshToken(ident, token, time.Now().Add(static.AuthSessionExpiration))
	return
}

func (rth *DatabaseRefreshTokenHandler) ValidateRefreshToken(token string) (ident string, err error) {
	ident, expires, err := rth.db.GetUserByRefreshToken(token)
	if err != nil {
		return
	}

	if time.Now().After(expires) {
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
}

// NewJWTAccessTokenHandler returns a new instance
// of JWTAccessTokenHandler
func NewJWTAccessTokenHandler(container di.Container) (ath *JWTAccessTokenHandler, err error) {
	secret, err := random.GetRandByteArray(32)
	if err != nil {
		return nil, err
	}
	ath = &JWTAccessTokenHandler{
		sessionExpiration: 10 * time.Minute, // TODO: Get from config
		sessionSecret:     secret,           // TODO: Get from config
	}
	return
}

func (ath *JWTAccessTokenHandler) GetAccessToken(ident string) (token string, expires time.Time, err error) {
	now := time.Now()
	expires = now.Add(ath.sessionExpiration)

	claims := jwt.StandardClaims{}
	claims.Issuer = fmt.Sprintf("shinpuru v.%s", util.AppVersion)
	claims.Subject = ident
	claims.ExpiresAt = expires.Unix()
	claims.NotBefore = now.Unix()
	claims.IssuedAt = now.Unix()

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
	jwt.StandardClaims

	Salt string `json:"sp_salt,omitempty"`
}

func apiTokenClaimsFromMap(m jwt.MapClaims) apiTokenClaims {
	c := apiTokenClaims{
		StandardClaims: standardClaimsFromMap(m),
	}

	c.Salt, _ = m["sp_salt"].(string)

	return c
}

func standardClaimsFromMap(m jwt.MapClaims) jwt.StandardClaims {
	c := jwt.StandardClaims{}

	c.Issuer, _ = m["iss"].(string)
	c.Subject, _ = m["sub"].(string)
	c.ExpiresAt, _ = m["exp"].(int64)
	c.NotBefore, _ = m["nbf"].(int64)
	c.IssuedAt, _ = m["iat"].(int64)

	return c
}

// DatabaseAPITokenHandler implements APITokenHandler
// for a JWT token stored in the database
type DatabaseAPITokenHandler struct {
	db      database.Database
	session *discordgo.Session
	secret  []byte
}

// NewDatabaseAPITokenHandler returns a new instance
// of DatabaseAPITokenHandler
func NewDatabaseAPITokenHandler(container di.Container) (*DatabaseAPITokenHandler, error) {
	cfg := container.Get(static.DiConfig).(*config.Config)
	secret := []byte(cfg.WebServer.APITokenKey)

	return &DatabaseAPITokenHandler{
		db:      container.Get(static.DiDatabase).(database.Database),
		session: container.Get(static.DiDiscordSession).(*discordgo.Session),
		secret:  secret,
	}, nil
}

func (apith *DatabaseAPITokenHandler) GetAPIToken(ident string) (token string, expires time.Time, err error) {
	now := time.Now()
	expires = now.Add(static.ApiTokenExpiration)

	salt, err := random.GetRandBase64Str(16)
	if err != nil {
		return
	}

	claims := apiTokenClaims{}
	claims.Issuer = fmt.Sprintf("shinpuru v.%s", util.AppVersion)
	claims.Subject = ident
	claims.ExpiresAt = expires.Unix()
	claims.NotBefore = now.Unix()
	claims.IssuedAt = now.Unix()
	claims.Salt = salt

	token, err = jwt.NewWithClaims(jwtGenerationMethod, claims).
		SignedString(apith.secret)
	if err != nil {
		return
	}

	tokenEntry := &models.APITokenEntry{
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
	tokenEntry.LastAccess = time.Now()
	apith.db.SetAPIToken(tokenEntry)

	return claims.Subject, nil
}

func (apith *DatabaseAPITokenHandler) RevokeToken(ident string) error {
	return apith.db.DeleteAPIToken(ident)
}
