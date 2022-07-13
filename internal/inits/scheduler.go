package inits

import (
	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron/v3"
	"github.com/sarulabs/di/v2"
	"github.com/sirupsen/logrus"
	"github.com/zekroTJA/shinpuru/internal/services/backup"
	"github.com/zekroTJA/shinpuru/internal/services/birthday"
	"github.com/zekroTJA/shinpuru/internal/services/config"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/services/guildlog"
	"github.com/zekroTJA/shinpuru/internal/services/report"
	"github.com/zekroTJA/shinpuru/internal/services/scheduler"
	"github.com/zekroTJA/shinpuru/internal/services/verification"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/antiraid"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/discordutil"
	"github.com/zekroTJA/shinpuru/pkg/twitchnotify"
	"github.com/zekrotja/dgrs"
)

func InitScheduler(container di.Container) scheduler.Provider {
	cfg := container.Get(static.DiConfig).(config.Provider)
	db := container.Get(static.DiDatabase).(database.Database)
	gb := container.Get(static.DiBackupHandler).(*backup.GuildBackups)
	tnw := container.Get(static.DiTwitchNotifyWorker).(*twitchnotify.NotifyWorker)
	rep := container.Get(static.DiReport).(*report.ReportService)
	gl := container.Get(static.DiGuildLog).(guildlog.Logger)
	vs := container.Get(static.DiVerification).(verification.Provider)
	bd := container.Get(static.DiBirthday).(*birthday.BirthdayService)
	s := container.Get(static.DiDiscordSession).(*discordgo.Session)
	st := container.Get(static.DiState).(dgrs.IState)

	shardID, shardTotal := discordutil.GetShardOfSession(s)

	sched := &scheduler.CronScheduler{C: cron.New(cron.WithSeconds())}

	schedule(sched, "refresh token cleanup",
		func() string {
			if shardTotal > 1 && shardID != 0 {
				return ""
			}
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

	schedule(sched, "guild backup",
		func() string {
			return cfg.Config().Schedules.GuildBackups
		},
		func() {
			go gb.BackupAllGuilds()
		})

	schedule(sched, "twitch notify",
		func() string {
			if shardTotal > 1 && shardID != 0 {
				return ""
			}
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

	schedule(sched, "report expiration",
		func() string {
			if shardTotal > 1 && shardID != 0 {
				return ""
			}
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

	schedule(sched, "verification kick routine",
		func() string {
			if shardTotal > 1 && shardID != 0 {
				return ""
			}
			return cfg.Config().Schedules.VerificationKick
		}, vs.KickRoutine)

	schedule(sched, "antiraid joinlog flush",
		func() string {
			if shardTotal > 1 && shardID != 0 {
				return ""
			}
			return "@every 1h"
		}, antiraid.FlushExpired(db, gl))

	schedule(sched, "birthday notifications",
		func() string {
			return "0 0 * * * *"
		}, func() {
			bd.Schedule()
		})

	schedule(sched, "guild membercount refresh",
		staticSpec("@every 24h"),
		func() {
			err := util.UpdateGuildMemberStats(st, s)
			if err != nil {
				logrus.WithError(err).Error("Failed refreshing guild member stats")
			} else {
				logrus.Debug("Refreshed guild member stats")
			}
		})

	return sched
}

func schedule(sched scheduler.Provider, name string, specGetter func() string, job func()) {
	spec := specGetter()
	if spec == "" {
		return
	}
	_, err := sched.Schedule(spec, job)
	if err != nil {
		logrus.WithError(err).WithField("name", name).Fatalf("LCT :: failed scheduling job")
	}
	logrus.WithField("name", name).WithField("spec", spec).Info("LCT :: scheduled job")
}

func staticSpec(v string) func() string {
	return func() string {
		return v
	}
}
