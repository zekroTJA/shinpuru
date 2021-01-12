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
