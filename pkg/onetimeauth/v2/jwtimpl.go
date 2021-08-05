// Package onetimeout provides short duration valid
// JWT tokens which are only valid exactly once.
package onetimeauth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/zekroTJA/shinpuru/pkg/random"
	"github.com/zekroTJA/timedmap"
)

// otaClaims extends jwt.StandardClaims by
// the token string.
type otaClaims struct {
	jwt.StandardClaims

	Token  string   `json:"tkn"`
	Scopes []string `json:"scp"`
}

// JwtOneTimeAuth implements OneTimeAuth
// using JWT tokens.
type JwtOneTimeAuth struct {
	signingKey []byte
	options    *JwtOptions

	tokens *timedmap.TimedMap
}

var _ OneTimeAuth = (*JwtOneTimeAuth)(nil)

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

func (a *JwtOneTimeAuth) GetKey(ident string, scopes ...string) (token string, expires time.Time, err error) {
	now := time.Now()
	expires = now.Add(a.options.Lifetime)

	claims := otaClaims{}
	claims.Issuer = a.options.Issuer
	claims.Subject = ident
	claims.ExpiresAt = expires.Unix()
	claims.NotBefore = now.Unix()
	claims.IssuedAt = now.Unix()
	claims.Scopes = scopes
	if claims.Token, err = random.GetRandBase64Str(32); err != nil {
		return
	}

	token, err = jwt.NewWithClaims(a.options.SigningMethod, claims).
		SignedString(a.signingKey)

	a.tokens.Set(claims.Token, ident, a.options.Lifetime)

	return
}

func (a *JwtOneTimeAuth) ValidateKey(key string, scopes ...string) (ident string, err error) {
	defer func() {
		if err != nil {
			ident = ""
		}
	}()

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
	scpi, _ := claims["scp"].([]interface{})
	if !okS || !okT {
		err = ErrInvalidClaims
		return
	}

	if uid, ok := a.tokens.GetValue(tkn).(string); !ok || uid != ident {
		err = ErrInvalidToken
		return
	}

	var scp []string
	if scpi != nil {
		scp = make([]string, len(scpi))
		for i := range scpi {
			if scp[i], ok = scpi[i].(string); !ok {
				err = ErrInvalidScopes
				return
			}
		}
	}

	if !validateScopes(scopes, scp) {
		err = ErrInvalidScopes
		return
	}

	a.tokens.Remove(tkn)

	return
}

func validateScopes(must []string, obtained []string) bool {
	for _, m := range must {
		if !contains(m, obtained) {
			return false
		}
	}
	return true
}

func contains(v string, arr []string) bool {
	for _, a := range arr {
		if a == v {
			return true
		}
	}
	return false
}
