package timeprovider

import "time"

type Time struct{}

func (Time) Now() time.Time {
	return time.Now()
}
