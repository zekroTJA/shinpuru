// Package validators provides some (more or less)
// general purpose validator functions for user
// inputs.
package validators

import (
	"errors"
	"regexp"
	"strconv"
)

// Lengthable is a constraint for types passable
// to the len() function.
type Lengthable interface {
	~string | ~[]any
}

// IntOrString is a constraint for a type which
// is either an integer or string.
type IntOrString interface {
	~int | ~string
}

// Combine can be used to check multiple validators
// in one function call.
//
// Example:
//   Combine(Length[string](10, 100), IsSimpleUrl())
func Combine[T any](validators ...func(T) error) func(T) error {
	return func(v T) (err error) {
		for _, validator := range validators {
			if err = validator(v); err != nil {
				break
			}
		}
		return
	}
}

// Length returns a function which validates the
// given input to be at least min and at maximum
// max long.
//
// When max is 0, no maximum length is checked.
func Length[T Lengthable](min, max int) func(T) error {
	return func(v T) error {
		if len(v) < min {
			return errors.New("value too short")
		}
		if max > 0 && len(v) > max {
			return errors.New("value too long")
		}
		return nil
	}
}

// IsInteger returns a function which validates
// that the passed string is a valid integer.
func IsInteger(allowEmpty ...bool) func(string) error {
	return func(s string) error {
		if opt(allowEmpty, false) && s == "" {
			return nil
		}
		_, err := strconv.Atoi(s)
		if err != nil {
			err = errors.New("value is not a number")
		}
		return err
	}
}

// InRange returns a function which validates
// a string or integer input to be in numeral
// range of [min, max].
//
// When max is 0, no maximum value is checked.
func InRange[T IntOrString](min, max int) func(T) error {
	return func(v T) (err error) {
		var vi interface{} = v
		switch vs := vi.(type) {
		case string:
			i, err := strconv.Atoi(vs)
			if err != nil {
				return err
			}
			return InRange[int](min, max)(i)
		case int:
			if vs < min {
				return errors.New("value is too low")
			}
			if max > 0 && vs > max {
				return errors.New("value is too large")
			}
			return nil
		}
		return errors.New("invalid type")
	}
}

// MatchesRegex returns a function validating that
// the passed string value matches the given regex.
func MatchesRegex(rx string) func(string) error {
	crx := regexp.MustCompile(rx)
	return func(s string) error {
		if !crx.MatchString(s) {
			return errors.New("value does not match regex " + rx)
		}
		return nil
	}
}

// IsDomain returns a function validating that
// the passed string value matches a domain regex.
func IsDomain() func(string) error {
	return MatchesRegex(`^(?:[\w\-]+\.)+[a-z]{2,6}$`)
}

// IsEmailAddress returns a function validating that
// the passed string value matches an e-mail regex.
func IsEmailAddress() func(string) error {
	return MatchesRegex(`^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+$`)
}

// IsSimpleUrl returns a function validating that
// the passed string value matches a simple URL regex.
func IsSimpleUrl() func(string) error {
	return MatchesRegex(`^https?:\/\/(?:www\.)?(?:[\w_-]+\.)+[\w_-]+.*$`)
}

func opt[T any](v []T, def T) T {
	if len(v) == 0 {
		return def
	}
	return v[0]
}
