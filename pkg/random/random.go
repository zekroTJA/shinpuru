// Package random provides some general purpose
// cryptographically pseudo-random utilities.
package random

import (
	"crypto/rand"
	"encoding/base64"
)

// GetRandBase64Str creates a cryptographically randomly
// generated set of bytes with the length of len which
// is returned as base64 encoded string.
func GetRandBase64Str(len int) (string, error) {
	data, err := GetRandByteArray(len)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(data), nil
}

// GetRandByteArray creates a cryptographically randomly
// generated set of bytes with the length of len.
func GetRandByteArray(len int) (data []byte, err error) {
	data = make([]byte, len)
	_, err = rand.Read(data)

	return
}
