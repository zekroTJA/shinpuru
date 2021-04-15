package auth

import "time"

type RefreshTokenHandler interface {
	GetRefreshToken(ident string) (token string, err error)
	ValidateRefreshToken(token string) (ident string, err error)
	RevokeToken(ident string) error
}

type AccessTokenHandler interface {
	GetAccessToken(ident string) (token string, expires time.Time, err error)
	ValidateAccessToken(token string) (ident string, err error)
}
