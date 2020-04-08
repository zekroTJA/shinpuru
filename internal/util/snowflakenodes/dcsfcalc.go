package snowflakenodes

import (
	"strconv"
	"time"
)

type DiscordSnowflake struct {
	Snowflake     string
	Time          time.Time
	WorkerID      int
	ProcessID     int
	IncrementalID int
}

func ParseDiscordSnowflake(sf string) (*DiscordSnowflake, error) {
	sfi, err := strconv.Atoi(sf)
	if err != nil {
		return nil, err
	}

	dcsf := new(DiscordSnowflake)
	dcsf.Snowflake = sf

	timestamp := (sfi >> 22) + 1420070400000
	dcsf.Time = ParseUnixTime(timestamp)
	dcsf.WorkerID = (sfi & 0x3E0000) >> 17
	dcsf.ProcessID = (sfi & 0x1F000) >> 12
	dcsf.IncrementalID = sfi & 0xFFF

	return dcsf, nil
}

func ParseUnixTime(t int) time.Time {
	return time.Unix(int64(t/1000), 0)
}
