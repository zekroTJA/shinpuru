package snowflakenodes

import (
	"strconv"
	"time"

	"github.com/zekroTJA/shinpuru/pkg/timeutil"
)

// DiscordSnowflake wraps detailed informations
// of a Discord snowflake.
type DiscordSnowflake struct {
	Snowflake     int
	Time          time.Time
	WorkerID      int
	ProcessID     int
	IncrementalID int
}

func (dcs *DiscordSnowflake) String() string {
	return strconv.Itoa(dcs.Snowflake)
}

// ParseDiscordSnowflakeStr recovers a DiscordSnowflake
// from the passed snowflake string.
//
// An error is returned when the passed snowflake
// could not be parsed to an integer.
func ParseDiscordSnowflakeStr(sf string) (dsf *DiscordSnowflake, err error) {
	sfi, err := strconv.Atoi(sf)
	if err != nil {
		return nil, err
	}

	dsf = ParseDiscordSnowflake(sfi)
	return
}

// ParseDiscordSnowflake recovers a DiscordSnowflake
// from the passed snowflake integer.
func ParseDiscordSnowflake(sfi int) *DiscordSnowflake {
	dcsf := new(DiscordSnowflake)
	dcsf.Snowflake = sfi

	timestamp := (sfi >> 22) + 1420070400000
	dcsf.Time = timeutil.FromUnix(timestamp)
	dcsf.WorkerID = (sfi & 0x3E0000) >> 17
	dcsf.ProcessID = (sfi & 0x1F000) >> 12
	dcsf.IncrementalID = sfi & 0xFFF

	return dcsf
}
