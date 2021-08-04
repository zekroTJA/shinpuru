package hashutil

import (
	"bytes"
	"crypto"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/zekroTJA/shinpuru/pkg/random"
)

var (
	ErrInvalidFormat           = errors.New("invalid hash format")
	ErrInvalidSaltSize         = errors.New("invalid salt size: must be larger than 0")
	ErrUnsupportedHashFunction = errors.New("unsupported hash function")
)

type Hasher struct {
	HashFunc     crypto.Hash
	SaltSize     int
	PepperGetter func() ([]byte, error)
}

// Format: <hashFuncName>$<saltHex>$<hashHex>

func (h Hasher) getHash(token string, salt []byte) (hash []byte, err error) {
	v := append([]byte(token), salt...)

	if h.PepperGetter != nil {
		var pepper []byte
		pepper, err = h.PepperGetter()
		if err != nil {
			return
		}
		v = append(v, pepper...)
	}

	hash = h.HashFunc.New().Sum(v)
	return
}

func (h Hasher) Hash(token string) (hash string, err error) {
	if h.SaltSize == 0 {
		err = ErrInvalidSaltSize
		return
	}

	salt, err := random.GetRandByteArray(h.SaltSize)
	if err != nil {
		return
	}

	hashB, err := h.getHash(token, salt)
	if err != nil {
		return
	}

	hash = fmt.Sprintf("%s$%x$%x", h.HashFunc.String(), salt, hashB)

	return
}

func Compare(token, hash string, pepperGetter ...func() ([]byte, error)) (ok bool, err error) {
	split := strings.Split(hash, "$")
	if len(split) != 3 {
		err = ErrInvalidFormat
		return
	}

	hasher := Hasher{}

	if len(pepperGetter) > 0 {
		hasher.PepperGetter = pepperGetter[0]
	}

	hasher.HashFunc, err = GetHashFunc(split[0])
	if err != nil {
		return
	}

	salt := make([]byte, len(split[1])/2)
	if _, err = hex.Decode(salt, []byte(split[1])); err != nil {
		return
	}

	cHash := make([]byte, len(split[2])/2)
	if _, err = hex.Decode(cHash, []byte(split[2])); err != nil {
		return
	}

	hashB, err := hasher.getHash(token, salt)
	if err != nil {
		return
	}

	ok = bytes.Equal(hashB, cHash)

	return
}

func GetHashFunc(name string) (h crypto.Hash, err error) {
	switch name {
	case "MD4":
		h = crypto.MD4
	case "MD5":
		h = crypto.MD5
	case "SHA-1":
		h = crypto.SHA1
	case "SHA-224":
		h = crypto.SHA224
	case "SHA-256":
		h = crypto.SHA256
	case "SHA-384":
		h = crypto.SHA384
	case "SHA-512":
		h = crypto.SHA512
	case "MD5+SHA1":
		h = crypto.MD5SHA1
	case "RIPEMD-160":
		h = crypto.RIPEMD160
	case "SHA3-224":
		h = crypto.SHA3_224
	case "SHA3-256":
		h = crypto.SHA3_256
	case "SHA3-384":
		h = crypto.SHA3_384
	case "SHA3-512":
		h = crypto.SHA3_512
	case "SHA-512/224":
		h = crypto.SHA512_224
	case "SHA-512/256":
		h = crypto.SHA512_256
	case "BLAKE2s-256":
		h = crypto.BLAKE2s_256
	case "BLAKE2b-256":
		h = crypto.BLAKE2b_256
	case "BLAKE2b-384":
		h = crypto.BLAKE2b_384
	case "BLAKE2b-512":
		h = crypto.BLAKE2b_512
	default:
		err = ErrUnsupportedHashFunction
	}

	return
}
