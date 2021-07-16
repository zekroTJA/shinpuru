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

// NowAddPtr adds t to now and returns the resulting
// time as *time.Time. If d is <= 0, nil is returned.
func NowAddPtr(d time.Duration) *time.Time {
	if d <= 0 {
		return nil
	}
	t := time.Now().Add(d)
	return &t
}
