package inits

import (
	"github.com/robfig/cron/v3"
	"github.com/sarulabs/di/v2"
	"github.com/sirupsen/logrus"
	"github.com/zekroTJA/shinpuru/internal/config"
	"github.com/zekroTJA/shinpuru/internal/services/backup"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/services/lctimer"
	"github.com/zekroTJA/shinpuru/internal/services/report"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/twitchnotify"
)

func InitLTCTimer(container di.Container) lctimer.LifeCycleTimer {
	cfg := container.Get(static.DiConfig).(*config.Config)
	db := container.Get(static.DiDatabase).(database.Database)
	guildBackups := container.Get(static.DiBackupHandler).(*backup.GuildBackups)
	tnw := container.Get(static.DiTwitchNotifyWorker).(*twitchnotify.NotifyWorker)
	rep := container.Get(static.DiReport).(*report.ReportService)

	lct := &lctimer.CronLifeCycleTimer{C: cron.New(cron.WithSeconds())}

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
				logrus.WithError(err).Error("LCT :: failed cleaning up expired refresh tokens")
			} else if n > 0 {
				logrus.WithField("n", n).Info("LCT :: cleaned up expired refresh tokens")
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
				logrus.WithError(err).Error("LCT :: failed executing twitch notify handler")
			}
		})

	lctSchedule(lct, "report expiration",
		func() string {
			return "@every 5m"
		},
		func() {
			rep.ExpireExpiredReports().ForEach(func(err error, i int) {
				logrus.WithError(err).Error("LCT :: failed expiring report")
			})
		})

	return lct
}

func lctSchedule(lct lctimer.LifeCycleTimer, name string, specGetter func() string, job func()) {
	spec := specGetter()
	_, err := lct.Schedule(spec, job)
	if err != nil {
		logrus.WithError(err).WithField("name", name).Fatalf("LCT :: failed scheduling job")
	}
	logrus.WithField("name", name).WithField("spec", spec).Info("LCT :: scheduled job")
}
