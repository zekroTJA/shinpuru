package inline

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestII(t *testing.T) {
	assert.Equal(t,
		2,
		II(false, 1, 2))
	assert.Equal(t,
		1,
		II(true, 1, 2))
}

func TestNC(t *testing.T) {
	assert.Equal(t,
		"a",
		NC("", "a"))
	assert.Equal(t,
		"b",
		NC("b", "a"))

	assert.Equal(t,
		1,
		NC(0, 1))
	assert.Equal(t,
		2,
		NC(2, 1))
}
