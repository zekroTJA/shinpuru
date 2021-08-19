// Package hammertime provides functionailities to
// format a time.Time into a Discord timestamp mention.
//
// The name was used after the very useful web app
// hammertime.djdavid98.art.
package hammertime

import (
	"fmt"
	"time"
)

type FormatSpec string

const (
	ShortDate      FormatSpec = "d" // 12/12/2020
	LongerDateTime FormatSpec = "f" // December 12, 2020 8:00 AM
	ShortTime      FormatSpec = "t" // 8:00 AM
	LongerDate     FormatSpec = "D" // December 12, 2020
	LongDateTime   FormatSpec = "F" // Saturday, December 12, 2020 8:00 AM
	Span           FormatSpec = "R" // 8 months ago
	LongTime       FormatSpec = "T" // 8:00:00 AM
)

func Format(t time.Time, f FormatSpec) string {
	return fmt.Sprintf("<t:%d:%s>", t.Unix(), f)
}
