package inits

import (
	"time"

	"github.com/zekroTJA/shinpuru/pkg/lctimer"
)

func InitLTCTimer() *lctimer.LifeCycleTimer {
	return lctimer.New(10 * time.Second)
}
