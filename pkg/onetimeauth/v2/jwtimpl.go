// Package onetimeout provides short duration valid
// JWT tokens which are only valid exactly once.
package onetimeauth

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/zekroTJA/shinpuru/pkg/random"
	"github.com/zekroTJA/timedmap"
)

var (
	ErrInvalidToken  = errors.New("invalid token")
	ErrInvalidClaims = errors.New("invalid claims")
)

// otaClaims extends jwt.StandardClaims by
// the token string.
type otaClaims struct {
	jwt.StandardClaims

	Token string `json:"tkn"`
}

// JwtOneTimeAuth implements OneTimeAuth
// using JWT tokens.
type JwtOneTimeAuth struct {
	signingKey []byte
	options    *JwtOptions

	tokens *timedmap.TimedMap
}

// NewJwt initializes a new JwtOneTimeAuth with a signing
// key generated on initialization.
//
// The passed options configure the OneTimeAuth
// instances. Non-set values in the options are
// defaulted with values of DefaultOptions.
func NewJwt(options *JwtOptions) (a *JwtOneTimeAuth, err error) {
	if options == nil {
		options = DefaultOptions
	} else {
		options.complete()
	}

	key, err := random.GetRandByteArray(options.SigningKeyLength)
	if err != nil {
		return
	}

	a = &JwtOneTimeAuth{
		signingKey: key,
		options:    options,
		tokens:     timedmap.New(10 * time.Minute),
	}

	return
}

func (a *JwtOneTimeAuth) GetKey(ident string) (token string, err error) {
	now := time.Now()

	claims := otaClaims{}
	claims.Issuer = a.options.Issuer
	claims.Subject = ident
	claims.ExpiresAt = now.Add(a.options.Lifetime).Unix()
	claims.NotBefore = now.Unix()
	claims.IssuedAt = now.Unix()
	if claims.Token, err = random.GetRandBase64Str(32); err != nil {
		return
	}

	token, err = jwt.NewWithClaims(a.options.SigningMethod, claims).
		SignedString(a.signingKey)

	a.tokens.Set(claims.Token, ident, a.options.Lifetime)

	return
}

func (a *JwtOneTimeAuth) ValidateKey(key string) (ident string, err error) {
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
		err = ErrInvalidClaims
		return
	}

	ident, okS := claims["sub"].(string)
	tkn, okT := claims["tkn"].(string)
	if !okS || !okT {
		err = ErrInvalidClaims
		return
	}

	if uid, ok := a.tokens.GetValue(tkn).(string); !ok || uid != ident {
		err = ErrInvalidToken
		return
	}

	a.tokens.Remove(tkn)

	return
}
