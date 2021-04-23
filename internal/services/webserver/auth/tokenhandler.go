package auth

import "time"

// RefreshTokenHandler provides functionalities
// to manage refresh tokens.
type RefreshTokenHandler interface {

	// GetRefreshToken takes an indent and genereates
	// a unique and secure token which can be used to
	// recover the passed ident from it.
	GetRefreshToken(ident string) (token string, err error)

	// ValidateRefreshToken takes a token and validates
	// it. When the token is invalid, an error is returned.
	// Otherwise, the ident is recovered from the token
	// and returned.
	ValidateRefreshToken(token string) (ident string, err error)

	// RevokeToken marks the token linked to the passed
	// ident as invalid so it can not be validated
	// anymore.
	RevokeToken(ident string) error
}

// AccessTokenHandler provides functionalities
// to manage access tokens.
type AccessTokenHandler interface {

	// GetAccessToken takes an indent and genereates
	// a unique and secure token which can be used to
	// recover the passed ident from it.
	//
	// Also, an expiration time is returned after which
	// the token will become invalid.
	GetAccessToken(ident string) (token string, expires time.Time, err error)

	// ValidateAccessToken takes a token and validates
	// it. When the token is invalid, an error is returned.
	// Otherwise, the ident is recovered from the token
	// and returned.
	ValidateAccessToken(token string) (ident string, err error)
}

// APITokenHandler provides functionalities
// to manage API tokens.
type APITokenHandler interface {

	// GetAPIToken takes an indent and genereates
	// a unique and secure token which can be used to
	// recover the passed ident from it.
	//
	// Also, an expiration time is returned after which
	// the token will become invalid.
	GetAPIToken(ident string) (token string, expires time.Time, err error)

	// ValidateAPIToken takes a token and validates
	// it. When the token is invalid, an error is returned.
	// Otherwise, the ident is recovered from the token
	// and returned.
	ValidateAPIToken(token string) (ident string, err error)

	// RevokeToken marks the token linked to the passed
	// ident as invalid so it can not be validated
	// anymore.
	RevokeToken(ident string) error
}
