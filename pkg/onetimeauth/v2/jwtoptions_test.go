package onetimeauth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComplete(t *testing.T) {
	options := &JwtOptions{}
	options.complete()

	assert.Equal(t, options.Issuer, DefaultOptions.Issuer)
	assert.Equal(t, options.Lifetime, DefaultOptions.Lifetime)
	assert.Equal(t, options.SigningKeyLength, DefaultOptions.SigningKeyLength)
	assert.Equal(t, options.TokenKeyLength, DefaultOptions.TokenKeyLength)
	assert.Equal(t, options.SigningMethod, DefaultOptions.SigningMethod)
}
