package onetimeauth

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/pkg/random"
)

const (
	keyLength = 128

	jwtIssuer = "shinpuru v.%v"
)

var jwtGenerationMethod = jwt.SigningMethodHS256

type OneTimeAuth struct {
	signingKey []byte
	duration   time.Duration
}

func New(duration time.Duration) (a *OneTimeAuth, err error) {
	key, err := random.GetRandByteArray(keyLength)
	if err != nil {
		return
	}

	a = &OneTimeAuth{key, duration}

	return
}

func (a *OneTimeAuth) GetKey(userID string) (token string, err error) {
	now := time.Now()

	claims := jwt.StandardClaims{}
	claims.Issuer = fmt.Sprintf(jwtIssuer, util.AppVersion)
	claims.Subject = userID
	claims.ExpiresAt = now.Add(a.duration).Unix()
	claims.NotBefore = now.Unix()
	claims.IssuedAt = now.Unix()

	token, err = jwt.NewWithClaims(jwtGenerationMethod, claims).
		SignedString(a.signingKey)

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

	userID, ok = claims["sub"].(string)
	if !ok {
		err = errors.New("invalid claims")
		return
	}

	return
}
