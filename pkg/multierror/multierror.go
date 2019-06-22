// Package multierror impements handling
// multiple errors as one error object.
//   Authors: Ringo Hoffmann.
package multierror

import (
	"errors"
	"fmt"
	"strings"
)

// FormatFunc is the function used to format an
// error element of a MulitError
type FormatFunc func(i int, e error) string

// MultiError contains multiple error objects
// and the frmatting function for concating into
// one handable error object
type MultiError struct {
	errors     []error
	formatFunc FormatFunc
}

// New creates a new instance of MultiError
// using the passed format function. If this argument
// is nil, the default format function will be used.
func New(formatFunc FormatFunc) *MultiError {
	if formatFunc == nil {
		formatFunc = func(i int, e error) string {
			return fmt.Sprintf("MultiError %02x: %s\n", i, e)
		}
	}

	return &MultiError{
		formatFunc: formatFunc,
	}
}

// Append adds an error object to the
// MultiError cotainer if the error
// is != nil
func (m *MultiError) Append(err error) {
	if err != nil {
		m.errors = append(m.errors, err)
	}
}

// Len returns the ammount of errors contained
// in the MultiError container
func (m *MultiError) Len() int {
	return len(m.errors)
}

// Concat creates one handable error object
// from all errors in the MultiError container
// using the formatting function.
func (m *MultiError) Concat() error {
	if m.Len() == 0 {
		return nil
	}

	if m.Len() == 1 {
		return m.errors[0]
	}

	strErrArr := make([]string, m.Len())
	for i, e := range m.errors {
		strErrArr[i] = m.formatFunc(i, e)
	}

	return errors.New(strings.Join(strErrArr, ""))
}
