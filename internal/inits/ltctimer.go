package inits

import (
	"github.com/robfig/cron/v3"
	"github.com/sarulabs/di/v2"
	"github.com/sirupsen/logrus"
	"github.com/zekroTJA/shinpuru/internal/services/backup"
	"github.com/zekroTJA/shinpuru/internal/services/config"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/services/guildlog"
	"github.com/zekroTJA/shinpuru/internal/services/lctimer"
	"github.com/zekroTJA/shinpuru/internal/services/report"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/twitchnotify"
)

func InitLTCTimer(container di.Container) lctimer.LifeCycleTimer {
	cfg := container.Get(static.DiConfig).(config.Provider)
	db := container.Get(static.DiDatabase).(database.Database)
	gb := container.Get(static.DiBackupHandler).(*backup.GuildBackups)
	tnw := container.Get(static.DiTwitchNotifyWorker).(*twitchnotify.NotifyWorker)
	rep := container.Get(static.DiReport).(*report.ReportService)
	gl := container.Get(static.DiGuildLog).(guildlog.Logger)

	lct := &lctimer.CronLifeCycleTimer{C: cron.New(cron.WithSeconds())}

	lctSchedule(lct, "refresh token cleanup",
		func() string {
			return cfg.Config().Schedules.RefreshTokenCleanup
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
			return cfg.Config().Schedules.GuildBackups
		},
		func() {
			go gb.BackupAllGuilds()
		})

	lctSchedule(lct, "twitch notify",
		func() string {
			return "@every 60s"
		},
		func() {
			if tnw == nil {
				return
			}
			if err := tnw.Handle(); err != nil {
				logrus.WithError(err).Error("LCT :: failed executing twitch notify handler")
			}
		})

	lctSchedule(lct, "report expiration",
		func() string {
			return cfg.Config().Schedules.ReportsExpiration
		},
		func() {
			rep.ExpireExpiredReports().ForEach(func(err error, _ int) {
				lentry := logrus.WithError(err)
				if repErr, ok := err.(*report.ReportError); ok {
					lentry = lentry.
						WithField("repID", repErr.ID).
						WithField("gid", repErr.GuildID)
					gl.Section("lct").Errorf(repErr.ID.String(), "Failed expiring report: %s", err.Error())
				}
				lentry.Error("LCT :: failed expiring report")
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
