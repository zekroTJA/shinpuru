// Package boolutil provides simple utility
// functions around booleans.
package boolutil

// AsInt returns 1 if v is true and
// 0 if v is false.
func AsInt(v bool) int {
	if v {
		return 1
	}
	return 0
}

// FromInt returns true if v > 0,
// otherwise returns false.
func FromInt(v int) bool {
	return v > 0
}
