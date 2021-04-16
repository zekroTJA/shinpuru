package inits

import (
	"time"

	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/core/database"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/lctimer"
)

func InitLTCTimer(container di.Container) *lctimer.LifeCycleTimer {
	db := container.Get(static.DiDatabase).(database.Database)

	lct := lctimer.New(10 * time.Second)
	lct.AfterDuration(24*time.Hour, func(now time.Time) {
		n, err := db.CleanupExpiredRefreshTokens()
		if err != nil {
			util.Log.Error("LCT :: failed cleaning up expired refresh tokens:", err)
		} else if n > 0 {
			util.Log.Infof("LCT :: cleaned up %d expired refresh tokens", n)
		}
	})

	return lct
}
