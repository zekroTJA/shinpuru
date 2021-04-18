// Package onetimeout provides short duration valid
// tokens which are only valid exactly once.
package onetimeauth

// OneTimeAuth provides functionalities to generate
// and validate a one time authentication key based
// on a passed ident.
type OneTimeAuth interface {

	// GetKey generates and registers a new OTA key
	// based on the passed ident.
	GetKey(ident string) (token string, err error)

	// ValidateKey tries to validate a given key. If
	// the validation fails, an error is returned with
	// details why the validation has failed.
	//
	// If the token is valid, the recovered ident and
	// a nil error is returned.
	ValidateKey(key string) (ident string, err error)
}
