package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/dgrijalva/jwt-go"
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/core/database"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/random"
)

type DatabaseRefreshTokenHandler struct {
	db      database.Database
	session *discordgo.Session
}

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

type JWTAccessTokenHandler struct {
	sessionExpiration time.Duration
	sessionSecret     []byte
}

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
