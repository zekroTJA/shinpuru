package listeners

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/sarulabs/di/v2"
	"github.com/sirupsen/logrus"
	"github.com/zekrotja/dgrs"

	"github.com/zekroTJA/shinpuru/internal/config"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/services/guildlog"
	"github.com/zekroTJA/shinpuru/internal/services/lctimer"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/presence"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/internal/util/vote"
)

type ListenerReady struct {
	config *config.Config
	db     database.Database
	gl     guildlog.Logger
	lct    lctimer.LifeCycleTimer
	st     *dgrs.State
}

func NewListenerReady(container di.Container) *ListenerReady {
	return &ListenerReady{
		config: container.Get(static.DiConfig).(*config.Config),
		db:     container.Get(static.DiDatabase).(database.Database),
		gl:     container.Get(static.DiGuildLog).(guildlog.Logger).Section("ready"),
		lct:    container.Get(static.DiLifecycleTimer).(lctimer.LifeCycleTimer),
		st:     container.Get(static.DiState).(*dgrs.State),
	}
}

func (l *ListenerReady) Handler(s *discordgo.Session, e *discordgo.Ready) {
	logrus.WithFields(logrus.Fields{
		"username": e.User.String(),
		"id":       e.User.ID,
		"nGuilds":  len(e.Guilds),
	})
	logrus.Infof("Invite link: %s", util.GetInviteLink(e.User.ID))

	s.UpdateGameStatus(0, static.StdMotd)

	l.lct.Start()

	rawPresence, err := l.db.GetSetting(static.SettingPresence)
	if err == nil {
		pre, err := presence.Unmarshal(rawPresence)
		if err == nil {
			s.UpdateStatusComplex(pre.ToUpdateStatusData())
		}
	}

	votes, err := l.db.GetVotes()
	if err != nil {
		logrus.WithError(err).Fatal("Failed getting votes from DB")
	} else {
		vote.VotesRunning = votes
		_, err = l.lct.Schedule("*/10 * * * * *", func() {
			now := time.Now()
			for _, v := range vote.VotesRunning {
				if (v.Expires != time.Time{}) && v.Expires.Before(now) {
					v.Close(s, vote.VoteStateExpired)
					if err = l.db.DeleteVote(v.ID); err != nil {
						logrus.WithError(err).WithField("gid", v.GuildID).WithField("vid", v.ID).Error("Failed updating vote")
						l.gl.Errorf(v.GuildID, "Failed updating vote (%s): %s", v.ID, err.Error())
					}
				}
			}
		})
		if err != nil {
			logrus.WithError(err).Fatal("LCT :: failed scheduling votes job")
		}
	}

	guilds, err := l.st.Guilds()
	if err != nil {
		logrus.WithError(err).Error("READY :: failed getting guilds")
	} else {
		logrus.WithField("n", len(guilds)).Info("READY :: caching members of guilds ...")
		for _, g := range guilds {
			if _, err := l.st.Members(g.ID, true); err != nil {
				logrus.WithError(err).WithField("gid", g.ID).Error("READY :: failed fetchting members")
			}
		}
		logrus.Info("READY :: caching members finished")
	}
}
