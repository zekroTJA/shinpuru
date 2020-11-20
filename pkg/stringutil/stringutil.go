// Package stringutil provides generl string
// utility functions.
package stringutil

import (
	"regexp"
	"strings"
)

var (
	rxNumber = regexp.MustCompile(`^-?\d+$`)
)

// IsInteger returns true if the passed string is
// a valid number.
func IsInteger(str string) bool {
	return rxNumber.MatchString(str)
}

// EnsureNotEmpty returns def if str is empty.
func EnsureNotEmpty(str, def string) string {
	if str == "" {
		return def
	}
	return str
}

// FromBool returns ifTrue if cond is true
// else returns ifFalse.
func FromBool(cond bool, ifTrue, ifFalse string) string {
	if cond {
		return ifTrue
	}
	return ifFalse
}

// IndexOfStringArray returns the index of the
// passed str in the passed arr. If str is not
// in arr, -1 is returned.
func IndexOf(str string, arr []string) int {
	for i, v := range arr {
		if v == str {
			return i
		}
	}
	return -1
}

// ContainsAny returns true if str is
// in arr.
func ContainsAny(str string, arr []string) bool {
	return IndexOf(str, arr) > -1
}

// HasPrefixAny returns true if the given str has
// any of the given prefixes.
func HasPrefixAny(str string, prefixes ...string) bool {
	for _, p := range prefixes {
		if strings.HasPrefix(str, p) {
			return true
		}
	}

	return false
}

// HasSuffixAny returns true if the given str has
// any of the given suffixes.
func HasSuffixAny(str string, suffixes ...string) bool {
	for _, s := range suffixes {
		if strings.HasSuffix(str, s) {
			return true
		}
	}

	return false
}
