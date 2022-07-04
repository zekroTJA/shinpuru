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

func TestParseDuration(t *testing.T) {
	var (
		exp, res time.Duration
		err      error
	)

	exp = 5 * time.Second
	res, err = ParseDuration("5s")
	assert.Nil(t, err)
	assert.Equal(t, exp, res)

	exp = 0
	res, err = ParseDuration("0s")
	assert.Nil(t, err)
	assert.Equal(t, exp, res)

	exp = 5 * time.Second
	res, err = ParseDuration("0h0m5s")
	assert.Nil(t, err)
	assert.Equal(t, exp, res)

	exp = 5 * time.Minute
	res, err = ParseDuration("5m")
	assert.Nil(t, err)
	assert.Equal(t, exp, res)

	exp = (9 * 24) * time.Hour
	res, err = ParseDuration("1w2d")
	assert.Nil(t, err)
	assert.Equal(t, exp, res)

	exp = (9*24+3)*time.Hour + 4*time.Minute + 5*time.Second
	res, err = ParseDuration("1w2d3h4m5s")
	assert.Nil(t, err)
	assert.Equal(t, exp, res)

	exp = 2*time.Millisecond + 3*time.Microsecond + 4*time.Nanosecond
	res, err = ParseDuration("2ms3us4ns")
	assert.Nil(t, err)
	assert.Equal(t, exp, res)

	exp = 2*time.Millisecond + 3*time.Microsecond + 4*time.Nanosecond
	res, err = ParseDuration("2ms3µs4ns")
	assert.Nil(t, err)
	assert.Equal(t, exp, res)

	exp = (9*24+3)*time.Hour + 4*time.Minute + 5*time.Second
	res, err = ParseDuration(" 		1w 2d	3h  4m 5s")
	assert.Nil(t, err)
	assert.Equal(t, exp, res)

	exp = 23 * time.Hour
	res, err = ParseDuration("1d-1h")
	assert.Nil(t, err)
	assert.Equal(t, exp, res)

	exp = 23 * time.Hour
	res, err = ParseDuration("1d  -1h")
	assert.Nil(t, err)
	assert.Equal(t, exp, res)

	exp = -3 * time.Second
	res, err = ParseDuration("-3s")
	assert.Nil(t, err)
	assert.Equal(t, exp, res)

	_, err = ParseDuration("")
	assert.EqualError(t, err, ErrInvalidDurationFormat.Error())

	_, err = ParseDuration("jskahdfjkas")
	assert.EqualError(t, err, ErrInvalidDurationFormat.Error())

	_, err = ParseDuration("invalid 4h 3s")
	assert.EqualError(t, err, ErrInvalidDurationFormat.Error())

	_, err = ParseDuration("4h invalid 3s")
	assert.EqualError(t, err, ErrInvalidDurationFormat.Error())
}
