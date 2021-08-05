package onetimeauth

import "errors"

var (
	// ErrInvalidToken is returned when an invalid
	// token was specified.
	ErrInvalidToken = errors.New("invalid token")
	// ErrInvalidClaims is returned when a token with
	// insufficient claims was specified.
	ErrInvalidClaims = errors.New("invalid claims")
	// ErrInvalidScopes is returned when a token with
	// insufficient scopes was specified.
	ErrInvalidScopes = errors.New("invalid scopes")
)
