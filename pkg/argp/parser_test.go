package argp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	{
		args := []string{}
		res, err := New(args).String("-n", "")
		assert.Equal(t, "", res)
		assert.Nil(t, err)
	}
	{
		args := []string{"-n", "heyho"}
		res, err := New(args).String("-n", "")
		assert.Equal(t, "heyho", res)
		assert.Nil(t, err)
	}
	{
		args := []string{"-a", "abc", "-n", "heyho", "was geht ab"}
		res, err := New(args).String("-n", "")
		assert.Equal(t, "heyho", res)
		assert.Nil(t, err)
	}
	{
		args := []string{"-n", `"hey`, "was", "geht", `ab"`}
		res, err := New(args).String("-n", "")
		assert.Equal(t, "hey was geht ab", res)
		assert.Nil(t, err)
	}
	{
		args := []string{"-n"}
		res, err := New(args).String("-n", "")
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
		res, err := New(args).Bool("-n", false)
		assert.Equal(t, false, res)
		assert.Nil(t, err)
	}
	{
		args := []string{"-n"}
		res, err := New(args).Bool("-n", false)
		assert.Equal(t, true, res)
		assert.Nil(t, err)
	}
	{
		args := []string{"-n=true"}
		res, err := New(args).Bool("-n", false)
		assert.Equal(t, true, res)
		assert.Nil(t, err)
	}
	{
		args := []string{"-n=false"}
		res, err := New(args).Bool("-n", false)
		assert.Equal(t, false, res)
		assert.Nil(t, err)
	}
	{
		args := []string{"-n=1"}
		res, err := New(args).Bool("-n", false)
		assert.Equal(t, false, res)
		assert.NotNil(t, err)
	}
}

func TestInt(t *testing.T) {
	{
		args := []string{}
		res, err := New(args).Int("-n", 0)
		assert.Equal(t, 0, res)
		assert.Nil(t, err)
	}
	{
		args := []string{"-n"}
		res, err := New(args).Int("-n", 0)
		assert.Equal(t, 0, res)
		assert.Nil(t, err)
	}
	{
		args := []string{}
		res, err := New(args).Int("-n", 456)
		assert.Equal(t, 456, res)
		assert.Nil(t, err)
	}
	{
		args := []string{"-n=123"}
		res, err := New(args).Int("-n", 0)
		assert.Equal(t, 123, res)
		assert.Nil(t, err)
	}
	{
		args := []string{"-n", "123"}
		res, err := New(args).Int("-n", 0)
		assert.Equal(t, 123, res)
		assert.Nil(t, err)
	}
	{
		args := []string{"-n=ads"}
		res, err := New(args).Int("-n", 0)
		assert.Equal(t, 0, res)
		assert.NotNil(t, err)
	}
	{
		args := []string{"-n", "ads"}
		res, err := New(args).Int("-n", 0)
		assert.Equal(t, 0, res)
		assert.NotNil(t, err)
	}
}

func TestFloat(t *testing.T) {
	{
		args := []string{}
		res, err := New(args).Float("-n", 0)
		assert.Equal(t, 0.0, res)
		assert.Nil(t, err)
	}
	{
		args := []string{"-n"}
		res, err := New(args).Float("-n", 0)
		assert.Equal(t, 0.0, res)
		assert.Nil(t, err)
	}
	{
		args := []string{}
		res, err := New(args).Float("-n", 4.56)
		assert.Equal(t, 4.56, res)
		assert.Nil(t, err)
	}
	{
		args := []string{"-n=1.23"}
		res, err := New(args).Float("-n", 0)
		assert.Equal(t, 1.23, res)
		assert.Nil(t, err)
	}
	{
		args := []string{"-n", "1.23"}
		res, err := New(args).Float("-n", 0)
		assert.Equal(t, 1.23, res)
		assert.Nil(t, err)
	}
	{
		args := []string{"-n=ads"}
		res, err := New(args).Float("-n", 0)
		assert.Equal(t, 0.0, res)
		assert.NotNil(t, err)
	}
	{
		args := []string{"-n", "ads"}
		res, err := New(args).Float("-n", 0)
		assert.Equal(t, 0.0, res)
		assert.NotNil(t, err)
	}
}

func TestArgs(t *testing.T) {
	{
		args := []string{"-a", "abc", "-n", "heyho", "yo"}
		p := New(args)
		res, err := p.String("-n", "")
		assert.Equal(t, "heyho", res)
		assert.Nil(t, err)

		rArgs := p.Args()
		assert.Equal(t,
			[]string{"-a", "abc", "yo"},
			rArgs)
	}
	{
		args := []string{"-a", "abc", "-n", "heyho", "yo"}
		p := New(args)
		res, err := p.Bool("-n", false)
		assert.Equal(t, true, res)
		assert.Nil(t, err)

		rArgs := p.Args()
		assert.Equal(t,
			[]string{"-a", "abc", "heyho", "yo"},
			rArgs)
	}
}
