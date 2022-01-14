// Package onetimeout provides short duration valid
// JWT tokens which are only valid exactly once.
package onetimeauth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zekroTJA/shinpuru/pkg/random"
	"github.com/zekroTJA/timedmap"
)

// otaClaims extends jwt.StandardClaims by
// the token string.
type otaClaims struct {
	jwt.StandardClaims

	Token string `json:"tkn"`
}

// OneTimeAuth provides functionalities to generate
// and validate a one time authentication key based
// on a passed user ID.
type OneTimeAuth struct {
	signingKey []byte
	options    *Options

	tokens *timedmap.TimedMap
}

// New initializes a new OneTimeAuth with a signing
// key generated on initialization.
//
// The passed options configure the OneTimeAuth
// instances. Non-set values in the options are
// defaulted with values of DefaultOptions.
func New(options *Options) (a *OneTimeAuth, err error) {
	if options == nil {
		options = DefaultOptions
	} else {
		options.complete()
	}

	key, err := random.GetRandByteArray(options.SigningKeyLength)
	if err != nil {
		return
	}

	a = &OneTimeAuth{
		signingKey: key,
		options:    options,
		tokens:     timedmap.New(10 * time.Minute),
	}

	return
}

// GetKey generates and registers a new OTA key
// based on the passed userID.
func (a *OneTimeAuth) GetKey(userID string) (token string, err error) {
	now := time.Now()

	claims := otaClaims{}
	claims.Issuer = a.options.Issuer
	claims.Subject = userID
	claims.ExpiresAt = now.Add(a.options.Lifetime).Unix()
	claims.NotBefore = now.Unix()
	claims.IssuedAt = now.Unix()
	if claims.Token, err = random.GetRandBase64Str(32); err != nil {
		return
	}

	token, err = jwt.NewWithClaims(a.options.SigningMethod, claims).
		SignedString(a.signingKey)

	a.tokens.Set(claims.Token, userID, a.options.Lifetime)

	return
}

// ValidateKey tries to validate a given key. If
// the validation fails, an error is returned with
// details why the validation has failed.
//
// If the token is valid, the recovered userID and
// a nil error is returned.
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
