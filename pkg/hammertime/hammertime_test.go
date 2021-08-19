package hammertime

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFormat(t *testing.T) {
	ts, _ := time.Parse(time.RFC3339, "2020-12-12T08:00:00Z")

	assert.Equal(t, "<t:1607760000:d>", Format(ts, ShortDate))
	assert.Equal(t, "<t:1607760000:f>", Format(ts, LongerDateTime))
	assert.Equal(t, "<t:1607760000:t>", Format(ts, ShortTime))
	assert.Equal(t, "<t:1607760000:D>", Format(ts, LongerDate))
	assert.Equal(t, "<t:1607760000:F>", Format(ts, LongDateTime))
	assert.Equal(t, "<t:1607760000:R>", Format(ts, Span))
	assert.Equal(t, "<t:1607760000:T>", Format(ts, LongTime))
}
