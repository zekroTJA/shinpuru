package listeners

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/sarulabs/di/v2"

	"github.com/zekroTJA/shinpuru/internal/config"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/services/lctimer"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/presence"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/internal/util/vote"
)

type ListenerReady struct {
	config *config.Config
	db     database.Database
	lct    lctimer.LifeCycleTimer
}

func NewListenerReady(container di.Container) *ListenerReady {
	return &ListenerReady{
		config: container.Get(static.DiConfig).(*config.Config),
		db:     container.Get(static.DiDatabase).(database.Database),
		lct:    container.Get(static.DiLifecycleTimer).(lctimer.LifeCycleTimer),
	}
}

func (l *ListenerReady) Handler(s *discordgo.Session, e *discordgo.Ready) {
	util.Log.Infof("Logged in as %s#%s (%s) - Running on %d servers",
		e.User.Username, e.User.Discriminator, e.User.ID, len(e.Guilds))
	util.Log.Info("Invite Link: " + util.GetInviteLink(s))

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
		util.Log.Error("Failed getting votes from DB: ", err)
	} else {
		vote.VotesRunning = votes
		_, err = l.lct.Schedule("*/10 * * * * *", func() {
			now := time.Now()
			for _, v := range vote.VotesRunning {
				if (v.Expires != time.Time{}) && v.Expires.Before(now) {
					v.Close(s, vote.VoteStateExpired)
					if err = l.db.DeleteVote(v.ID); err != nil {
						util.Log.Errorf("Failed updating vote with ID %s: %s", v.ID, err.Error())
					}
				}
			}
		})
		if err != nil {
			util.Log.Errorf("LCT :: failed scheduling votes job: %s", err.Error())
		}
	}
}
