package timeutil

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFromUnix(t *testing.T) {
	const unixStamp = 1610292604000
	timeObj, err := time.Parse(time.UnixDate, "Sun Jan 10 15:30:04 UTC 2021")
	assert.Nil(t, err)

	timeRec := FromUnix(unixStamp)
	assert.Zero(t, timeRec.Sub(timeObj))
}

func TestToUnix(t *testing.T) {
	const unixStamp = 1610292604000
	timeObj, err := time.Parse(time.UnixDate, "Sun Jan 10 15:30:04 UTC 2021")
	assert.Nil(t, err)

	unixRec := ToUnix(timeObj)
	assert.Equal(t, unixRec, unixStamp)
}

func TestNowAddPtr(t *testing.T) {
	res := NowAddPtr(0)
	assert.Nil(t, res)

	res = NowAddPtr(-1)
	assert.Nil(t, res)

	now := time.Now().Add(5 * time.Second)
	res = NowAddPtr(5 * time.Second)
	assert.NotNil(t, res)
	assert.InDelta(t, 0, now.Sub(*res), float64(100*time.Millisecond))
}
