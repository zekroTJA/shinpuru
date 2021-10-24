package intutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFromBool(t *testing.T) {
	assert.Equal(t, 1, FromBool(true, 1, 2))
	assert.Equal(t, 2, FromBool(false, 1, 2))
}
