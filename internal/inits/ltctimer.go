package inits

import (
	"time"

	"github.com/zekroTJA/shinpuru/internal/core"
)

func InitLTCTimer() *core.LCTimer {
	return core.NewLTCTimer(10 * time.Second)
}
