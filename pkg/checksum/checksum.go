package checksum

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"hash"
)

func Sum(v interface{}, hasher hash.Hash) (hash string, err error) {
	data, err := json.Marshal(v)
	if err != nil {
		return
	}

	hash = hex.EncodeToString(hasher.Sum(data))

	return
}

func SumSha1(v interface{}) (string, error) {
	return Sum(v, sha1.New())
}

func SumSha256(v interface{}) (string, error) {
	return Sum(v, sha256.New())
}

func SumMd5(v interface{}) (string, error) {
	return Sum(v, md5.New())
}

func Must(hash string, err error) string {
	if err != nil {
		panic(err)
	}
	return hash
}

// func Sum(v interface{})
