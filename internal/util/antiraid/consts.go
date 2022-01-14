package antiraid

import "time"

const (
	TriggerRecordLifetime = 24 * time.Hour
	TriggerLifetime       = 2 * TriggerRecordLifetime
)
