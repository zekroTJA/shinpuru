package onetimeauth

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/pkg/random"
	"github.com/zekroTJA/timedmap"
)

const (
	keyLength   = 128
	tokenLength = 32

	jwtIssuer = "shinpuru v.%v"
)

var jwtGenerationMethod = jwt.SigningMethodHS256

type otaClaims struct {
	jwt.StandardClaims

	Token string `json:"tkn"`
}

type OneTimeAuth struct {
	signingKey []byte
	duration   time.Duration

	tokens *timedmap.TimedMap
}

func New(duration time.Duration) (a *OneTimeAuth, err error) {
	key, err := random.GetRandByteArray(keyLength)
	if err != nil {
		return
	}

	a = &OneTimeAuth{
		signingKey: key,
		duration:   duration,
		tokens:     timedmap.New(10 * time.Minute),
	}

	return
}

func (a *OneTimeAuth) GetKey(userID string) (token string, err error) {
	now := time.Now()

	claims := otaClaims{}
	claims.Issuer = fmt.Sprintf(jwtIssuer, util.AppVersion)
	claims.Subject = userID
	claims.ExpiresAt = now.Add(a.duration).Unix()
	claims.NotBefore = now.Unix()
	claims.IssuedAt = now.Unix()
	if claims.Token, err = random.GetRandBase64Str(32); err != nil {
		return
	}

	token, err = jwt.NewWithClaims(jwtGenerationMethod, claims).
		SignedString(a.signingKey)

	a.tokens.Set(claims.Token, userID, a.duration)

	return
}

func (a *OneTimeAuth) ValidateKey(key string) (userID string, err error) {
	token, err := jwt.Parse(key, func(t *jwt.Token) (interface{}, error) {
		return a.signingKey, nil
	})
	if err != nil {
		return
	}
	if err = token.Claims.Valid(); err != nil {
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		err = errors.New("invalid claims")
		return
	}

	userID, okS := claims["sub"].(string)
	tkn, okT := claims["tkn"].(string)
	if !okS || !okT {
		err = errors.New("invalid claims")
		return
	}

	if uid, ok := a.tokens.GetValue(tkn).(string); !ok || uid != userID {
		err = errors.New("invalid token")
		return
	}

	a.tokens.Remove(tkn)

	return
}
