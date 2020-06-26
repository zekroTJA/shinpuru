// Package etag implements generation functionalities
// for the ETag specification of RFC7273 2.3.
// https://tools.ietf.org/html/rfc7232#section-2.3.1
package etag

import (
	"crypto/sha1"
	"fmt"
)

// Generate an ETag by body byte array.
// weak specifies if the generated ETag should
// be tagged as "weak".
func Generate(body []byte, weak bool) string {
	hash := sha1.Sum(body)

	weakTag := ""
	if weak {
		weakTag = "W/"
	}

	tag := fmt.Sprintf("%s\"%x\"", weakTag, hash)

	return tag
}
