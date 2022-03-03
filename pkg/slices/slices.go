// Package slices adds generic utility
// functionalities for slices.
package slices

// IndexOf returns the first index of occurence
// of the given element v in the slice s.
//
// If no occurences were found, -1 is returned.
func IndexOf[T comparable](s []T, v T) (i int) {
	var c T
	for i, c = range s {
		if v == c {
			return
		}
	}
	return -1
}

// Contains returns true if the given element
// v occurs in the slice s.
func Contains[T comparable](s []T, v T) bool {
	return IndexOf(s, v) != -1
}

// Splice takes a slice and returns two new arrays
// where the first slice contains the elements of s
// with the elements from i to i+n cut out and the
// second slice contains the cut out elements.
func Splice[T any](s []T, i, n int) (ns, rest []T) {
	if i < 0 {
		i = 0
	}
	if i+n >= len(s) {
		n = len(s) - i
	}
	cs := make([]T, len(s))
	copy(cs, s)
	ns = append(cs[0:i], cs[i+n:]...)
	rest = s[i : i+n]
	return
}
