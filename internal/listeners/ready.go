package listeners

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/sarulabs/di/v2"
	"github.com/sirupsen/logrus"
	"github.com/zekrotja/dgrs"

	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/services/guildlog"
	"github.com/zekroTJA/shinpuru/internal/services/lctimer"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/presence"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/internal/util/vote"
)

type ListenerReady struct {
	db  database.Database
	gl  guildlog.Logger
	lct lctimer.LifeCycleTimer
	st  *dgrs.State
}

func NewListenerReady(container di.Container) *ListenerReady {
	return &ListenerReady{
		db:  container.Get(static.DiDatabase).(database.Database),
		gl:  container.Get(static.DiGuildLog).(guildlog.Logger).Section("ready"),
		lct: container.Get(static.DiLifecycleTimer).(lctimer.LifeCycleTimer),
		st:  container.Get(static.DiState).(*dgrs.State),
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
		logrus.WithError(err).Error("Failed getting votes from DB")
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
			logrus.WithError(err).Error("LCT :: failed scheduling votes job")
		}
	}

	time.Sleep(1 * time.Second)

	logrus.WithField("n", len(e.Guilds)).Info("READY :: caching members of guilds ...")
	for _, g := range e.Guilds {
		gs, _ := l.st.Guild(g.ID)
		if gs != nil && gs.MemberCount > 0 {
			membs, _ := l.st.Members(g.ID)
			if len(membs) >= gs.MemberCount {
				logrus.WithField("gid", g.ID).Debug("READY :: skip fetching members because state is hydrated")
				continue
			}
		}

		if _, err := l.st.Members(g.ID, true); err != nil {
			logrus.WithError(err).WithField("gid", g.ID).Error("READY :: failed fetchting members")
		} else {
			logrus.WithField("gid", g.ID).Debug("READY :: fetched members")
		}
	}
	logrus.Info("READY :: caching members finished")
}
