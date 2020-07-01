package util

import (
	"regexp"
)

var (
	rxNumber = regexp.MustCompile(`^\d+$`)
)

// IsNumber returns true if the passed string is
// a valid number.
func IsNumber(str string) bool {
	return rxNumber.MatchString(str)
}

// EnsureNotEmpty returns def if str is empty.
func EnsureNotEmpty(str, def string) string {
	if str == "" {
		return def
	}
	return str
}

// BoolAsString returns ifTrue if cond is true
// else returns ifFalse.
func BoolAsString(cond bool, ifTrue, ifFalse string) string {
	if cond {
		return ifTrue
	}
	return ifFalse
}

// IndexOfStringArray returns the index of the
// passed str in the passed arr. If str is not
// in arr, -1 is returned.
func IndexOfStrArray(str string, arr []string) int {
	for i, v := range arr {
		if v == str {
			return i
		}
	}
	return -1
}

// StringArrayContains returns true if str is
// in arr.
func StringArrayContains(str string, arr []string) bool {
	return IndexOfStrArray(str, arr) > -1
}
