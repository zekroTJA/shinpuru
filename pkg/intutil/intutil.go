// Package intutil provides some utility
// functionalities for integers.
package intutil

// FromBool returns ifTrue if cond is true
// else returns ifFalse.
func FromBool(cond bool, ifTrue, ifFalse int) int {
	if cond {
		return ifTrue
	}
	return ifFalse
}
