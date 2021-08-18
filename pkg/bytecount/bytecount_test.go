package bytecount

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormat(t *testing.T) {
	assert.Equal(t, Format(615), "615 B")
	assert.Equal(t, Format(5623), "5.491 kiB")
	assert.Equal(t, Format(4425623), "4.221 MiB")
	assert.Equal(t, Format(9134426623), "8.507 GiB")
	assert.Equal(t, Format(4534425386623), "4.124 TiB")
}
