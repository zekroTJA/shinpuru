package inits

import (
	"github.com/robfig/cron/v3"
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/core/backup"
	"github.com/zekroTJA/shinpuru/internal/core/config"
	"github.com/zekroTJA/shinpuru/internal/core/database"
	"github.com/zekroTJA/shinpuru/internal/shared"
	"github.com/zekroTJA/shinpuru/internal/shared/wrappers"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/twitchnotify"
)

func InitLTCTimer(container di.Container) shared.LifeCycleTimer {
	cfg := container.Get(static.DiConfig).(*config.Config)
	db := container.Get(static.DiDatabase).(database.Database)
	guildBackups := container.Get(static.DiBackupHandler).(*backup.GuildBackups)
	tnw := container.Get(static.DiTwitchNotifyWorker).(*twitchnotify.NotifyWorker)

	lct := &wrappers.CronLifeCycleTimer{C: cron.New(cron.WithSeconds())}

	lctSchedule(lct, "refresh token cleanup",
		func() string {
			spec := cfg.Defaults.Schedules.RefreshTokenCleanup
			if cfg.Schedules != nil && cfg.Schedules.RefreshTokenCleanup != "" {
				spec = cfg.Schedules.RefreshTokenCleanup
			}
			return spec
		},
		func() {
			n, err := db.CleanupExpiredRefreshTokens()
			if err != nil {
				util.Log.Error("LCT :: failed cleaning up expired refresh tokens:", err)
			} else if n > 0 {
				util.Log.Infof("LCT :: cleaned up %d expired refresh tokens", n)
			}
		})

	lctSchedule(lct, "guild backup",
		func() string {
			spec := cfg.Defaults.Schedules.GuildBackups
			if cfg.Schedules != nil && cfg.Schedules.GuildBackups != "" {
				spec = cfg.Schedules.GuildBackups
			}
			return spec
		},
		func() {
			go guildBackups.BackupAllGuilds()
		})

	lctSchedule(lct, "twitch notify",
		func() string {
			return "@every 60s"
		},
		func() {
			if err := tnw.Handle(); err != nil {
				util.Log.Errorf("LCT :: failed executing twitch notify handler: %s", err.Error())
			}
		})

	return lct
}

func lctSchedule(lct shared.LifeCycleTimer, name string, specGetter func() string, job func()) {
	spec := specGetter()
	_, err := lct.Schedule(spec, job)
	if err != nil {
		util.Log.Fatalf("LCT :: failed scheduling %s job: %s", name, err.Error())
	}
	util.Log.Infof("LCT :: scheduled %s job for '%s'", name, spec)
}
