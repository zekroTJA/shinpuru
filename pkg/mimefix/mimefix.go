// Package mimefix provides functionalities to bypass
// this issue with fasthttp on windows hosts*:
// https://github.com/golang/go/issues/32350
package mimefix

import (
	"mime"
	"strings"
)

const (
	expectedJSMimeType = "application/javascript"
)

type iMimeFixer interface {
	Fix(expectedMime string) error
}

// Check gets the detected mime type for .js
// files and returns it. Also it is checked
// against the expected mime type so that true
// is returned if the value matches the expected
// mime type.
func Check() (curr string, ok bool) {
	curr = mime.TypeByExtension(".js")
	ok = strings.HasPrefix(curr, expectedJSMimeType)
	return
}

// Fix loads the os-dependent fixer and executes
// the Fix function. An error is returned when the
// fix routine failed.
func Fix() error {
	return new(mimeFixer).Fix(expectedJSMimeType)
}

// CheckFix first executes Check(). If the current
// value is alredy set ot the expected value, no
// further action is performed. Otherwise, fix is
// executed.
//
// Returns the current mime string value, true if
// Check() returned true and an error when Fix()
// was executed and has failed.
func CheckFix() (curr string, ok bool, err error) {
	if curr, ok = Check(); ok {
		return
	}

	err = Fix()
	return
}
