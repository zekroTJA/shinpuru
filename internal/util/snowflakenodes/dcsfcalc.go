package snowflakenodes

import (
	"strconv"
	"time"

	"github.com/zekroTJA/shinpuru/pkg/timeutil"
)

// DiscordSnowflake wraps detailed informations
// of a Discord snowflake.
type DiscordSnowflake struct {
	Snowflake     string
	Time          time.Time
	WorkerID      int
	ProcessID     int
	IncrementalID int
}

// ParseDiscordSnowflake recovers a DiscordSnowflake
// from the passed snowflake string.
//
// An error is returned when the passed snowflake
// could not be parsed to an integer.
func ParseDiscordSnowflake(sf string) (*DiscordSnowflake, error) {
	sfi, err := strconv.Atoi(sf)
	if err != nil {
		return nil, err
	}

	dcsf := new(DiscordSnowflake)
	dcsf.Snowflake = sf

	timestamp := (sfi >> 22) + 1420070400000
	dcsf.Time = timeutil.FromUnix(timestamp)
	dcsf.WorkerID = (sfi & 0x3E0000) >> 17
	dcsf.ProcessID = (sfi & 0x1F000) >> 12
	dcsf.IncrementalID = sfi & 0xFFF

	return dcsf, nil
}
