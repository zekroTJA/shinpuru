// Package inline provides general inline
// operation functions like inline if or
// null coalescence.
package inline

// II (inline-if) takes a comparable value
// v. If the value v equals the default value of
// TIn, s is returned. Otherwise, p is returned.
//
// This is comparable to a syntax like
//   res = v == default(TIn) ? p : s
func II[TIn comparable, TOut any](v TIn, p, s TOut) TOut {
	var def TIn
	if v == def {
		return s
	}
	return p
}

// NC (nil coalescence) takes a comparable
// value v. If v equals the default value
// of T, s is returned. Otherwise, v is
// returned.
//
// This is comparable to a syntax like
//   res = v ?? s
func NC[T comparable](v, s T) T {
	return II(v, v, s)
}
