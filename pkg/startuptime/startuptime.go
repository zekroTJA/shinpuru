// Package startuptime provides simple
// functionalities to measure the startup
// time of an application.
package startuptime

import "time"

var startTs = time.Now()

// Took returns the duration since
// application startup.
func Took() time.Duration {
	return time.Since(startTs)
}
