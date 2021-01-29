// Package timeutil provides some general purpose
// functionalities around the time package.
package timeutil

import "time"

// FromUnix returns a time.Time struct from
// the passed unix timestamp t.
func FromUnix(t int) time.Time {
	return time.Unix(int64(t/1000), 0)
}

// ToUnix returns the passed time.Time struct
// as unix milliseconds timestamp.
func ToUnix(t time.Time) int {
	return int(t.UnixNano() / 1_000_000)
}
