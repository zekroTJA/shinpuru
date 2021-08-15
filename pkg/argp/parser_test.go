package argp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	{
		args := []string{}
		res, err := New(args).String("-n")
		assert.Equal(t, "", res)
		assert.Nil(t, err)
	}
	{
		args := []string{"-n", "heyho"}
		res, err := New(args).String("-n")
		assert.Equal(t, "heyho", res)
		assert.Nil(t, err)
	}
	{
		args := []string{"-a", "abc", "-n", "heyho", "was geht ab"}
		res, err := New(args).String("-n")
		assert.Equal(t, "heyho", res)
		assert.Nil(t, err)
	}
	{
		args := []string{"-n", `"hey`, "was", "geht", `ab"`}
		res, err := New(args).String("-n")
		assert.Equal(t, "hey was geht ab", res)
		assert.Nil(t, err)
	}
	{
		args := []string{"-n"}
		res, err := New(args).String("-n")
		assert.Equal(t, "", res)
		assert.Nil(t, err)
	}
	{
		args := []string{}
		res, err := New(args).String("-n", "def")
		assert.Equal(t, "def", res)
		assert.Nil(t, err)
	}
	{
		args := []string{"-n"}
		res, err := New(args).String("-n", "def")
		assert.Equal(t, "def", res)
		assert.Nil(t, err)
	}
	{
		args := []string{"-n", "notdef"}
		res, err := New(args).String("-n", "def")
		assert.Equal(t, "notdef", res)
		assert.Nil(t, err)
	}
}

func TestBool(t *testing.T) {
	{
		args := []string{}
		res, err := New(args).Bool("-n")
		assert.Equal(t, false, res)
		assert.Nil(t, err)
	}
	{
		args := []string{"-n"}
		res, err := New(args).Bool("-n")
		assert.Equal(t, true, res)
		assert.Nil(t, err)
	}
	{
		args := []string{"-n=true"}
		res, err := New(args).Bool("-n")
		assert.Equal(t, true, res)
		assert.Nil(t, err)
	}
	{
		args := []string{"-n=false"}
		res, err := New(args).Bool("-n")
		assert.Equal(t, false, res)
		assert.Nil(t, err)
	}
	{
		args := []string{"-n=1"}
		res, err := New(args).Bool("-n")
		assert.Equal(t, false, res)
		assert.NotNil(t, err)
	}
}
