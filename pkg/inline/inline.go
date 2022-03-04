// Package inline provides general inline
// operation functions like inline if or
// null coalescence.
package inline

// II (inline-if) takes a comparable value
// v. If the value v equals the default value of
// TIn, s is returned. Otherwise, p is returned.
//
// This is comparable to a syntax like
//   res = v ? p : s
func II[TOut any](v bool, p, s TOut) TOut {
	if v {
		return p
	}
	return s
}

// NC (nil coalescence) takes a comparable
// value v. If v equals the default value
// of T, s is returned. Otherwise, v is
// returned.
//
// This is comparable to a syntax like
//   res = v ?? s
func NC[T comparable](v, s T) T {
	var def T
	return II(v == def, s, v)
}
