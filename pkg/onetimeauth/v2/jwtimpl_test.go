package onetimeauth

import (
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	testUserID   = "tester"
	invalidToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MTQ3OTY4ODEsImlhdCI6MTYxNDc5Njg4MCwiaXNzIjoic2hpbnB1cnUgdi5URVNUSU5HX0JVSUxEIiwibmJmIjoxNjE0Nzk2ODgwLCJzdWIiOiJ0ZXN0ZXIifQ.na7K3mAszvtMN9x1VfIEv_QU5ZKWHJUSPONPKABFbCI"
)

var testOptions = &JwtOptions{
	Issuer:           "test issuer",
	Lifetime:         time.Second,
	SigningKeyLength: 128,
	TokenKeyLength:   32,
	SigningMethod:    jwt.SigningMethodHS256,
}

func TestNewJwt(t *testing.T) {
	_, err := NewJwt(nil)
	if err != nil {
		t.Error(err)
	}

	_, err = NewJwt(new(JwtOptions))
	if err != nil {
		t.Error(err)
	}

	_, err = NewJwt(testOptions)
	if err != nil {
		t.Error(err)
	}
}

func TestGetKey(t *testing.T) {
	a, err := NewJwt(testOptions)
	if err != nil {
		t.Error(err)
	}

	_, err = a.GetKey(testUserID)
	if err != nil {
		t.Error(err)
	}
}

func TestValidateKey(t *testing.T) {
	a, err := NewJwt(testOptions)
	if err != nil {
		t.Error(err)
	}

	_, err = a.ValidateKey(invalidToken)
	if err == nil {
		t.Error("invalid token passed falsely")
	}

	token, err := a.GetKey(testUserID)
	if err != nil {
		t.Error(err)
	}

	userID, err := a.ValidateKey(token)
	if err != nil {
		t.Error(err)
	}

	if userID != testUserID {
		t.Error("user id missmatch")
	}

	time.Sleep(2 * time.Second)
	_, err = a.ValidateKey(token)
	if err == nil {
		t.Error("key did not expire")
	}
}
