package inits

import (
	"time"

	"github.com/zekroTJA/shinpuru/internal/core/lctimer"
)

func InitLTCTimer() *lctimer.LifeCycleTimer {
	return lctimer.New(10 * time.Second)
}
