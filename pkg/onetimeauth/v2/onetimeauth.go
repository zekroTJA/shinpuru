// Package onetimeout provides short duration valid
// tokens which are only valid exactly once.
package onetimeauth

import "time"

// OneTimeAuth provides functionalities to generate
// and validate a one time authentication key based
// on a passed ident.
type OneTimeAuth interface {

	// GetKey generates and registers a new OTA key
	// based on the passed ident.
	//
	// You can also specify scopes for which the token
	// can be used.
	GetKey(ident string, scopes ...string) (token string, expires time.Time, err error)

	// ValidateKey tries to validate a given key. If
	// the validation fails, an error is returned with
	// details why the validation has failed.
	//
	// You can also pass scopes which the token must be
	// explicitly be issued for. If one of the required
	// scopes does not match with the token's scopes, the
	// validation fails with ErrInvalidScopes.
	//
	// If the token is valid, the recovered ident and
	// a nil error is returned.
	ValidateKey(key string, scopes ...string) (ident string, err error)
}
