package random

import (
	"crypto/rand"
	"encoding/base64"
)

// GetRandBase64Str creates a cryptographically randomly
// generated set of bytes with the length of len which
// is returned as base64 encoded string.
func GetRandBase64Str(len int) (string, error) {
	str := make([]byte, len)

	if _, err := rand.Read(str); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(str), nil
}
