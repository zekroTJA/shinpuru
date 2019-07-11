package inits

import (
	"time"

	"github.com/zekroTJA/shinpuru/internal/util"
)

func InitLTCTimer() *util.LTCTimer {
	return util.NewLTCTimer(1 * time.Second)
}
