package inits

import (
	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron/v3"
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/services/backup"
	"github.com/zekroTJA/shinpuru/internal/services/birthday"
	"github.com/zekroTJA/shinpuru/internal/services/config"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/services/guildlog"
	"github.com/zekroTJA/shinpuru/internal/services/report"
	"github.com/zekroTJA/shinpuru/internal/services/scheduler"
	"github.com/zekroTJA/shinpuru/internal/services/timeprovider"
	"github.com/zekroTJA/shinpuru/internal/services/verification"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/antiraid"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/discordutil"
	"github.com/zekroTJA/shinpuru/pkg/twitchnotify"
	"github.com/zekrotja/dgrs"
	"github.com/zekrotja/rogu"
	"github.com/zekrotja/rogu/log"
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
	tp := container.Get(static.DiTimeProvider).(timeprovider.Provider)

	shardID, shardTotal := discordutil.GetShardOfSession(s)

	sched := &scheduler.CronScheduler{C: cron.New(cron.WithSeconds())}

	log := log.Tagged("LCT")
	log.Info().Msg("Initializing lifecycle timer ...")

	schedule(log, sched, "refresh token cleanup",
		func() string {
			if shardTotal > 1 && shardID != 0 {
				return ""
			}
			return cfg.Config().Schedules.RefreshTokenCleanup
		},
		func() {
			n, err := db.CleanupExpiredRefreshTokens()
			if err != nil {
				log.Error().Err(err).Msg("Failed cleaning up expired refresh tokens")
			} else if n > 0 {
				log.Info().Field("n", n).Msg("Cleaned up expired refresh tokens")
			}
		})

	schedule(log, sched, "guild backup",
		func() string {
			return cfg.Config().Schedules.GuildBackups
		},
		func() {
			go gb.BackupAllGuilds()
		})

	schedule(log, sched, "twitch notify",
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
				log.Error().Err(err).Msg("Failed executing twitch notify handler")
			}
		})

	schedule(log, sched, "report expiration",
		func() string {
			if shardTotal > 1 && shardID != 0 {
				return ""
			}
			return cfg.Config().Schedules.ReportsExpiration
		},
		func() {
			rep.ExpireExpiredReports().ForEach(func(err error, _ int) {
				entry := log.Error().Err(err)
				if repErr, ok := err.(*report.ReportError); ok {
					entry.Fields(
						"repID", repErr.ID,
						"gid", repErr.GuildID,
					)
					gl.Section("lct").Errorf(repErr.ID.String(), "Failed expiring report: %s", err.Error())
				}
				entry.Msg("Failed expiring report")
			})
		})

	schedule(log, sched, "verification kick routine",
		func() string {
			if shardTotal > 1 && shardID != 0 {
				return ""
			}
			return cfg.Config().Schedules.VerificationKick
		}, vs.KickRoutine)

	schedule(log, sched, "antiraid joinlog flush",
		func() string {
			if shardTotal > 1 && shardID != 0 {
				return ""
			}
			return "@every 1h"
		}, antiraid.FlushExpired(db, gl, tp))

	schedule(log, sched, "birthday notifications",
		func() string {
			return "0 0 * * * *"
		}, func() {
			bd.Schedule()
		})

	schedule(log, sched, "guild membercount refresh",
		staticSpec("@every 24h"),
		func() {
			err := util.UpdateGuildMemberStats(st, s)
			if err != nil {
				log.Error().Err(err).Msg("Failed refreshing guild member stats")
			} else {
				log.Debug().Msg("Refreshed guild member stats")
			}
		})

	return sched
}

func schedule(log rogu.Logger, sched scheduler.Provider, name string, specGetter func() string, job func()) {
	spec := specGetter()
	if spec == "" {
		return
	}
	_, err := sched.Schedule(spec, job)
	if err != nil {
		log.Fatal().Err(err).Field("name", name).Msg("Failed scheduling job")
	}
	log.Info().Fields("name", name, "spec", spec).Msg("Scheduled job")
}

func staticSpec(v string) func() string {
	return func() string {
		return v
	}
}
