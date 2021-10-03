// Package multierror impements handling
// multiple errors as one error object.
package multierror

import (
	"fmt"
	"strings"
)

// FormatFunc is the function used to stringify
// a MultiError.
type FormatFunc func(errors []error) string

var (
	defaultFormatFunc = func(errors []error) string {
		ln := len(errors)
		lines := make([]string, ln+1)
		lines[0] = fmt.Sprintf("MultiError stack [%d inner errors]:", ln)
		for i := 1; i < ln+1; i++ {
			lines[i] = fmt.Sprintf("  [%02d] %s", i-1, errors[i-1].Error())
		}
		return strings.Join(lines, "\n")
	}
)

// MultiError implements the error interface
// and can contain and merge multiple errors.
type MultiError struct {
	errors     []error
	formatFunc FormatFunc
}

// New creates a new instance of MultiError
// using the passed format function. If this argument
// is nil, the default format function will be used.
func New(formatFunc ...FormatFunc) (m *MultiError) {
	m = new(MultiError)

	if formatFunc != nil && len(formatFunc) > 0 && formatFunc[0] != nil {
		m.formatFunc = formatFunc[0]
	} else {
		m.formatFunc = defaultFormatFunc
	}

	return
}

func (m *MultiError) Error() string {
	if m.Len() == 0 {
		return ""
	}
	return m.formatFunc(m.errors)
}

// Errors returns the internal list of errors.
func (m *MultiError) Errors() []error {
	return m.errors
}

// ForEach iterates over the list of errors and
// executes f for each error in the list.
func (m *MultiError) ForEach(f func(err error, i int)) {
	for i, err := range m.errors {
		f(err, i)
	}
}

// Append adds an error object to the
// MultiError cotainer if the error
// is != nil
func (m *MultiError) Append(err ...error) {
	for _, e := range err {
		if e != nil {
			m.errors = append(m.errors, e)
		}
	}
}

// Len returns the amount of errors contained
// in the MultiError container
func (m *MultiError) Len() int {
	return len(m.errors)
}

// DEPRECATED
//
// Returns the MultiError object as
// error interface.
//
// This function is deprecated. Please simply
// use the MultiError object itself as error,
// because it implements the error interface.
func (m *MultiError) Concat() error {
	return m
}

// Nillify returns nil if the MultiError does
// not contain any error objects, otherwise
// the MultiError instance is returned.
func (m *MultiError) Nillify() error {
	if m.Len() > 0 {
		return m
	}
	return nil
}
