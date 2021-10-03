package vote

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	uid1 = "221905671296253953"
	uid2 = "524847123875889153"
	uid3 = "536916384026722314"
)

func TestHashUserID(t *testing.T) {
	{
		salt := []byte("salt")

		hash1, err := HashUserID(uid1, salt)
		assert.Nil(t, err)

		hash2, err := HashUserID(uid1, salt)
		assert.Nil(t, err)

		assert.Equal(t, hash1, hash2)
	}

	{
		salt := []byte("salt")

		hash1, err := HashUserID(uid1, salt)
		assert.Nil(t, err)

		hash2, err := HashUserID(uid2, salt)
		assert.Nil(t, err)

		assert.NotEqual(t, hash1, hash2)
	}

	{
		salt1 := []byte("salt1")
		salt2 := []byte("salt2")

		hash1, err := HashUserID(uid1, salt1)
		assert.Nil(t, err)

		hash2, err := HashUserID(uid1, salt2)
		assert.Nil(t, err)

		assert.NotEqual(t, hash1, hash2)
	}
}
