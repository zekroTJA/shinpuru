package timeprovider

import "time"

type Provider interface {
	Now() time.Time
}
