// Package timeutil provides some general purpose
// functionalities around the time package.
package timeutil

import (
	"errors"
	"regexp"
	"strconv"
	"time"

	"github.com/zekroTJA/shinpuru/pkg/regexputil"
)

var (
	ErrInvalidDurationFormat = errors.New("invalid duration format")

	rxDuration = regexp.MustCompile(`(?is)^\s*(?:(?P<w>-?\d+)w)?\s*(?:(?P<d>-?\d+)d)?\s*(?:(?P<h>-?\d+)h)?\s*(?:(?P<m>-?\d+)m)?\s*(?:(?P<s>-?\d+)s)?\s*(?:(?P<ms>-?\d+)ms)?\s*(?:(?P<us>-?\d+)[uµ]s)?\s*(?:(?P<ns>-?\d+)ns)?\s*$`)
)

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

// DateOnly returns the given DateTime
// with time set to 00:00:00.
func DateOnly(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}

// ParseDuration tries to parse a duration from the passed
// string s. The format is composed of an integer number
// in combination with a time unit suffix. Following
// suffixes are supported:
//
//   - w       (weeks)
//   - d       (days)
//   - h       (hours)
//   - m       (minutes)
//   - s       (seconds)
//   - s       (seconds)
//   - ms      (milliseconds)
//   - us / µs (microseconds)
//   - ns      (nanoseconds)
//
// Also spaces and tab spaces are supported between the
// elements of the duration string. A valid example would
// be following duration string:
//
//  "3w1d 4h12m3s40ms"
//
// Also, substractions inside the duration strings are
// possible. The following example results in a furation
// of 23 hours.
//
//   "1d -1h"
func ParseDuration(s string) (time.Duration, error) {
	matches := regexputil.FindNamedSubmatchMap(rxDuration, s)
	if len(matches) == 0 {
		return 0, errors.New("invalid duration format")
	}

	var d time.Duration

	if wStr, ok := matches["w"]; ok {
		v, _ := strconv.Atoi(wStr)
		d += time.Duration(v) * 7 * 24 * time.Hour
	}
	if dStr, ok := matches["d"]; ok {
		v, _ := strconv.Atoi(dStr)
		d += time.Duration(v) * 24 * time.Hour
	}
	if hStr, ok := matches["h"]; ok {
		v, _ := strconv.Atoi(hStr)
		d += time.Duration(v) * time.Hour
	}
	if mStr, ok := matches["m"]; ok {
		v, _ := strconv.Atoi(mStr)
		d += time.Duration(v) * time.Minute
	}
	if sStr, ok := matches["s"]; ok {
		v, _ := strconv.Atoi(sStr)
		d += time.Duration(v) * time.Second
	}
	if msStr, ok := matches["ms"]; ok {
		v, _ := strconv.Atoi(msStr)
		d += time.Duration(v) * time.Millisecond
	}
	if usStr, ok := matches["us"]; ok {
		v, _ := strconv.Atoi(usStr)
		d += time.Duration(v) * time.Microsecond
	}
	if nsStr, ok := matches["ns"]; ok {
		v, _ := strconv.Atoi(nsStr)
		d += time.Duration(v) * time.Nanosecond
	}

	return d, nil
}
