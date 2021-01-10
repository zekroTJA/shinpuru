package timeutil

import (
	"fmt"
	"testing"
	"time"
)

func TestFromUnix(t *testing.T) {
	const unixStamp = 1610292604000
	timeObj, _ := time.Parse(time.UnixDate, "Sun Jan 10 15:30:04 UTC 2021")

	timeRec := FromUnix(unixStamp)

	if !timeRec.Equal(timeObj) {
		t.Error("recovered time unequals actual time")
	}
}

func TestToUnix(t *testing.T) {
	const unixStamp = 1610292604000
	timeObj, _ := time.Parse(time.UnixDate, "Sun Jan 10 15:30:04 UTC 2021")

	unixRec := ToUnix(timeObj)
	fmt.Println(unixRec)

	if unixRec != unixStamp {
		t.Error("recovered stamp unequals actual stamp")
	}
}
