// Package checksum provides functions to generate a hash sum
// from any given object.
package checksum

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"hash"
)

// Sum returns a hash string from the given object v using
// the given hash function hasher.
func Sum(v interface{}, hasher hash.Hash) (hash string, err error) {
	data, err := json.Marshal(v)
	if err != nil {
		return
	}

	hash = hex.EncodeToString(hasher.Sum(data))

	return
}

// SumSha1 is shorthand for Sum using the sha1 hash function.
func SumSha1(v interface{}) (string, error) {
	return Sum(v, sha1.New())
}

// SumSha256 is shorthand for Sum using the sha256 hash function.
func SumSha256(v interface{}) (string, error) {
	return Sum(v, sha256.New())
}

// SumMd5 is shorthand for Sum using the md5 hash function.
func SumMd5(v interface{}) (string, error) {
	return Sum(v, md5.New())
}

// Must can be used to execute a hash function and panics
// if the hash generation returned an error.
//
// Example:
//   myObject := &MyObject{}
//   hash := checksum.Must(checksum.SumSha1(myObject))
func Must(hash string, err error) string {
	if err != nil {
		panic(err)
	}
	return hash
}
