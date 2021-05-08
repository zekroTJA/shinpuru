package util

import "time"

var (
	StatsStartupTime             = time.Now()
	StatsCommandsExecuted uint64 = 0
	StatsMessagesAnalysed uint64 = 0
)
