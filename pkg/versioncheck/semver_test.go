package versioncheck

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseSemver(t *testing.T) {
	var (
		v   Semver
		err error
	)

	v, err = ParseSemver("2")
	assert.Nil(t, err)
	assert.Equal(t, Semver{2, 0, 0, ""}, v)

	v, err = ParseSemver("2.3")
	assert.Nil(t, err)
	assert.Equal(t, Semver{2, 3, 0, ""}, v)

	v, err = ParseSemver("12.3.456")
	assert.Nil(t, err)
	assert.Equal(t, Semver{12, 3, 456, ""}, v)

	v, err = ParseSemver("v12.3.456")
	assert.Nil(t, err)
	assert.Equal(t, Semver{12, 3, 456, ""}, v)

	v, err = ParseSemver("v.12.3.456")
	assert.Nil(t, err)
	assert.Equal(t, Semver{12, 3, 456, ""}, v)

	v, err = ParseSemver("V12.3.456")
	assert.Nil(t, err)
	assert.Equal(t, Semver{12, 3, 456, ""}, v)

	v, err = ParseSemver("V.12.3.456")
	assert.Nil(t, err)
	assert.Equal(t, Semver{12, 3, 456, ""}, v)

	v, err = ParseSemver("12.3.456-beta1")
	assert.Nil(t, err)
	assert.Equal(t, Semver{12, 3, 456, "beta1"}, v)

	v, err = ParseSemver("12.3.456+beta1")
	assert.Nil(t, err)
	assert.Equal(t, Semver{12, 3, 456, "beta1"}, v)

	v, err = ParseSemver("12.3.456+beta_1-test")
	assert.Nil(t, err)
	assert.Equal(t, Semver{12, 3, 456, "beta_1-test"}, v)

	_, err = ParseSemver("invalid semver")
	assert.ErrorIs(t, err, ErrNoMatch)

	_, err = ParseSemver("1.invalidversionnumber.0")
	assert.ErrorIs(t, err, ErrNoMatch)

	_, err = ParseSemver(" 1.0.0 ")
	assert.ErrorIs(t, err, ErrNoMatch)
}

func TestString(t *testing.T) {
	var r string

	r = Semver{}.String()
	assert.Equal(t, "0.0.0", r)

	r = Semver{1, 0, 0, ""}.String()
	assert.Equal(t, "1.0.0", r)

	r = Semver{12, 3, 456, ""}.String()
	assert.Equal(t, "12.3.456", r)

	r = Semver{12, 3, 456, "beta1"}.String()
	assert.Equal(t, "12.3.456-beta1", r)
}

func TestEqual(t *testing.T) {
	assert.True(t, Semver{1, 2, 3, "4"}.Equal(Semver{1, 2, 3, "4"}))

	assert.True(t, Semver{1, 2, 3, "4"}.Equal(Semver{1, 2, 3, "4"}, Exact))
	assert.False(t, Semver{1, 2, 3, ""}.Equal(Semver{1, 2, 3, "4"}, Exact))
	assert.False(t, Semver{1, 2, 2, "4"}.Equal(Semver{1, 2, 3, "4"}, Exact))
	assert.False(t, Semver{1, 1, 3, "4"}.Equal(Semver{1, 2, 3, "4"}, Exact))
	assert.False(t, Semver{3, 2, 3, "4"}.Equal(Semver{1, 2, 3, "4"}, Exact))

	assert.True(t, Semver{1, 2, 3, "4"}.Equal(Semver{1, 2, 3, "4"}, Patch))
	assert.True(t, Semver{1, 2, 3, ""}.Equal(Semver{1, 2, 3, "4"}, Patch))
	assert.False(t, Semver{1, 2, 2, "4"}.Equal(Semver{1, 2, 3, "4"}, Patch))
	assert.False(t, Semver{1, 1, 3, "4"}.Equal(Semver{1, 2, 3, "4"}, Patch))
	assert.False(t, Semver{3, 2, 3, "4"}.Equal(Semver{1, 2, 3, "4"}, Patch))

	assert.True(t, Semver{1, 2, 3, "4"}.Equal(Semver{1, 2, 3, "4"}, Minor))
	assert.True(t, Semver{1, 2, 3, ""}.Equal(Semver{1, 2, 3, "4"}, Minor))
	assert.True(t, Semver{1, 2, 2, "4"}.Equal(Semver{1, 2, 3, "4"}, Minor))
	assert.True(t, Semver{1, 2, 4, "123"}.Equal(Semver{1, 2, 3, "4"}, Minor))
	assert.False(t, Semver{1, 1, 3, "4"}.Equal(Semver{1, 2, 3, "4"}, Minor))
	assert.False(t, Semver{3, 2, 3, "4"}.Equal(Semver{1, 2, 3, "4"}, Minor))

	assert.True(t, Semver{1, 2, 3, "4"}.Equal(Semver{1, 2, 3, "4"}, Major))
	assert.True(t, Semver{1, 2, 3, ""}.Equal(Semver{1, 2, 3, "4"}, Major))
	assert.True(t, Semver{1, 2, 2, "4"}.Equal(Semver{1, 2, 3, "4"}, Major))
	assert.True(t, Semver{1, 1, 3, "4"}.Equal(Semver{1, 2, 3, "4"}, Major))
	assert.True(t, Semver{1, 2, 4, ""}.Equal(Semver{1, 2, 3, "4"}, Major))
	assert.True(t, Semver{1, 3, 6, "123"}.Equal(Semver{1, 2, 3, "4"}, Major))
	assert.False(t, Semver{3, 2, 3, "4"}.Equal(Semver{1, 2, 3, "4"}, Major))
}

func TestOlderThan(t *testing.T) {
	assert.True(t, Semver{1, 2, 3, ""}.OlderThan(Semver{1, 3, 3, ""}, Minor))
	assert.True(t, Semver{0, 3, 3, ""}.OlderThan(Semver{1, 3, 3, ""}, Minor))
	assert.False(t, Semver{1, 3, 2, ""}.OlderThan(Semver{1, 3, 3, ""}, Minor))
	assert.False(t, Semver{1, 4, 4, ""}.OlderThan(Semver{1, 3, 3, ""}, Minor))
}

func TestLaterThan(t *testing.T) {
	assert.True(t, Semver{1, 4, 3, ""}.LaterThan(Semver{1, 3, 3, ""}, Minor))
	assert.True(t, Semver{2, 3, 3, ""}.LaterThan(Semver{1, 3, 3, ""}, Minor))
	assert.False(t, Semver{1, 3, 3, ""}.LaterThan(Semver{1, 3, 3, ""}, Minor))
	assert.False(t, Semver{1, 2, 3, ""}.LaterThan(Semver{1, 3, 3, ""}, Minor))
}
