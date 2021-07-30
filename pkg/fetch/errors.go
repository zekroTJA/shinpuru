package fetch

import "errors"

var (
	// ErrNotFound is returned when an object could
	// not be found
	ErrNotFound = errors.New("not found")
)
