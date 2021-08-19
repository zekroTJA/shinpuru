// Package random provides some general purpose
// cryptographically pseudo-random utilities.
package random

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
)

var ErrInvalidLen = errors.New("invalid length")

// GetRandBase64Str creates a cryptographically randomly
// generated set of bytes with the length of len which
// is returned as base64 encoded string.
func GetRandBase64Str(len int) (string, error) {
	if len <= 0 {
		return "", ErrInvalidLen
	}

	data, err := GetRandByteArray(len)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(data)[:len], nil
}

// GetRandByteArray creates a cryptographically randomly
// generated set of bytes with the length of len.
func GetRandByteArray(len int) (data []byte, err error) {
	if len <= 0 {
		return nil, ErrInvalidLen
	}

	data = make([]byte, len)
	_, err = rand.Read(data)

	return
}

// MustGetRandBase64Str executes GetRandBase64Str and
// panics if an error was returned.
func MustGetRandBase64Str(len int) string {
	v, err := GetRandBase64Str(len)
	if err != nil {
		panic(err)
	}

	return v
}

// MustGetRandByteArray executes GetRandByteArray and
// panics if an error was returned.
func MustGetRandByteArray(len int) []byte {
	v, err := GetRandByteArray(len)
	if err != nil {
		panic(err)
	}

	return v
}
