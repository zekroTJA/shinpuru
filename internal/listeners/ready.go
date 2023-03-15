package listeners

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/sarulabs/di/v2"
	"github.com/zekrotja/dgrs"
	"github.com/zekrotja/rogu"
	"github.com/zekrotja/rogu/log"

	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/services/guildlog"
	"github.com/zekroTJA/shinpuru/internal/services/scheduler"
	"github.com/zekroTJA/shinpuru/internal/services/timeprovider"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/presence"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/internal/util/vote"
)

type ListenerReady struct {
	db    database.Database
	gl    guildlog.Logger
	sched scheduler.Provider
	st    *dgrs.State
	tp    timeprovider.Provider
	log   rogu.Logger
}

func NewListenerReady(container di.Container) *ListenerReady {
	return &ListenerReady{
		db:    container.Get(static.DiDatabase).(database.Database),
		gl:    container.Get(static.DiGuildLog).(guildlog.Logger).Section("ready"),
		sched: container.Get(static.DiScheduler).(scheduler.Provider),
		st:    container.Get(static.DiState).(*dgrs.State),
		tp:    container.Get(static.DiTimeProvider).(timeprovider.Provider),
		log:   log.Tagged("Ready"),
	}
}

func (l *ListenerReady) Handler(s *discordgo.Session, e *discordgo.Ready) {
	l.log.Info().Fields(
		"username", e.User.String(),
		"id", e.User.ID,
		"nGuilds", len(e.Guilds),
	).Msg("Discord Connection ready")
	l.log.Info().Msgf("Invite link: %s", util.GetInviteLink(e.User.ID))

	s.UpdateGameStatus(0, static.StdMotd)

	l.sched.Start()

	rawPresence, err := l.db.GetSetting(static.SettingPresence)
	if err == nil {
		pre, err := presence.Unmarshal(rawPresence)
		if err == nil {
			s.UpdateStatusComplex(pre.ToUpdateStatusData())
		}
	}

	votes, err := l.db.GetVotes()
	if err != nil {
		l.log.Error().Err(err).Msg("Failed getting votes from DB")
	} else {
		vote.VotesRunning = votes
		_, err = l.sched.Schedule("*/10 * * * * *", func() {
			now := l.tp.Now()
			for _, v := range vote.VotesRunning {
				if (v.Expires != time.Time{}) && v.Expires.Before(now) {
					v.Close(s, vote.VoteStateExpired)
					if err = l.db.DeleteVote(v.ID); err != nil {
						log.Error().Tag("LCT").Err(err).Fields("gid", v.GuildID, "vid", v.ID).Msg("Failed updating vote")
						l.gl.Errorf(v.GuildID, "Failed updating vote (%s): %s", v.ID, err.Error())
					}
				}
			}
		})
		if err != nil {
			log.Error().Tag("LCT").Err(err).Msg("Failed scheduling votes job")
		}
	}

	time.Sleep(1 * time.Second)

	l.log.Info().Field("n", len(e.Guilds)).Msg("Start caching members of guilds ...")
	for _, g := range e.Guilds {
		gs, _ := l.st.Guild(g.ID)
		if gs != nil && gs.MemberCount > 0 {
			membs, _ := l.st.Members(g.ID)
			if len(membs) >= gs.MemberCount {
				l.log.Debug().Field("gid", g.ID).Msg("Skip fetching members because state is hydrated")
				continue
			}
		}

		if _, err := l.st.Members(g.ID, true); err != nil {
			l.log.Error().Err(err).Field("gid", g.ID).Msg("Failed fetchting members")
		} else {
			l.log.Debug().Field("gid", g.ID).Msg("Fetched members")
		}
	}
	l.log.Info().Msg("Caching members finished")
}
