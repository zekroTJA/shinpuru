package onetimeauth

import (
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

const (
	testUserID   = "tester"
	invalidToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MTQ3OTY4ODEsImlhdCI6MTYxNDc5Njg4MCwiaXNzIjoic2hpbnB1cnUgdi5URVNUSU5HX0JVSUxEIiwibmJmIjoxNjE0Nzk2ODgwLCJzdWIiOiJ0ZXN0ZXIifQ.na7K3mAszvtMN9x1VfIEv_QU5ZKWHJUSPONPKABFbCI"
)

var testOptions = JwtOptions{
	Issuer:           "test issuer",
	Lifetime:         time.Second,
	SigningKeyLength: 128,
	TokenKeyLength:   32,
	SigningMethod:    jwt.SigningMethodHS256,
}

func TestNewJwt(t *testing.T) {
	_, err := NewJwt(nil)
	assert.Nil(t, err)

	_, err = NewJwt(new(JwtOptions))
	assert.Nil(t, err)

	_, err = NewJwt(&testOptions)
	assert.Nil(t, err)
}

func TestGetKey(t *testing.T) {
	a, err := NewJwt(&testOptions)
	assert.Nil(t, err)

	_, _, err = a.GetKey(testUserID)
	assert.Nil(t, err)
}

func TestValidateKey(t *testing.T) {
	a, err := NewJwt(&testOptions)
	assert.Nil(t, err)

	_, err = a.ValidateKey(invalidToken)
	assert.NotNil(t, err)

	token, _, err := a.GetKey(testUserID)
	assert.Nil(t, err)

	userID, err := a.ValidateKey(token)
	assert.Nil(t, err)
	assert.Equal(t, userID, testUserID)

	time.Sleep(2 * time.Second)
	ident, err := a.ValidateKey(token)
	assert.NotNil(t, err)
	assert.Empty(t, ident)
}

func TestScopes(t *testing.T) {
	opt := testOptions
	opt.Lifetime = 1 * time.Minute
	a, err := NewJwt(&opt)
	assert.Nil(t, err)

	{
		token, _, err := a.GetKey(testUserID, "a", "b", "c")
		assert.Nil(t, err)

		ident, err := a.ValidateKey(token)
		assert.Nil(t, err)
		assert.Equal(t, ident, testUserID)
	}

	{
		token, _, err := a.GetKey(testUserID, "a", "b", "c")
		assert.Nil(t, err)

		ident, err := a.ValidateKey(token, "a", "b")
		assert.Nil(t, err)
		assert.Equal(t, ident, testUserID)
	}

	{
		token, _, err := a.GetKey(testUserID, "a")
		assert.Nil(t, err)

		ident, err := a.ValidateKey(token, "a", "b")
		assert.ErrorIs(t, err, ErrInvalidScopes)
		assert.Empty(t, ident)
	}
}

func TestContains(t *testing.T) {
	assert.True(t, contains("a", []string{"a", "b"}))
	assert.True(t, contains("a", []string{"a"}))
	assert.True(t, contains("", []string{"a", ""}))

	assert.False(t, contains("a", []string{"b", "c"}))
	assert.False(t, contains("a", []string{"b", "c"}))
	assert.False(t, contains("", []string{"b", "c"}))
	assert.False(t, contains("a", []string{}))
	assert.False(t, contains("a", nil))
}

func TestValidateScopes(t *testing.T) {
	assert.True(t, validateScopes([]string{"a", "b"}, []string{"a", "b"}))
	assert.True(t, validateScopes([]string{"a"}, []string{"a"}))
	assert.True(t, validateScopes([]string{"a"}, []string{"a", "b"}))
	assert.True(t, validateScopes(nil, []string{"a", "b"}))
	assert.True(t, validateScopes(nil, nil))

	assert.False(t, validateScopes([]string{"a"}, []string{"b", "c"}))
	assert.False(t, validateScopes([]string{"a", "c"}, []string{"b", "c"}))
	assert.False(t, validateScopes([]string{"a", "c"}, nil))
}
