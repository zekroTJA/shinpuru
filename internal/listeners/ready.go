package listeners

import (
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/zekroTJA/shinpuru/internal/core/config"
	"github.com/zekroTJA/shinpuru/internal/core/database"
	"github.com/zekroTJA/shinpuru/internal/core/lctimer"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/presence"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/internal/util/vote"
)

type ListenerReady struct {
	config *config.Config
	db     database.Database
	lct    *lctimer.LCTimer
}

func NewListenerReady(config *config.Config, db database.Database, lct *lctimer.LCTimer) *ListenerReady {
	return &ListenerReady{
		config: config,
		db:     db,
		lct:    lct,
	}
}

func (l *ListenerReady) Handler(s *discordgo.Session, e *discordgo.Ready) {
	util.Log.Infof("Logged in as %s#%s (%s) - Running on %d servers",
		e.User.Username, e.User.Discriminator, e.User.ID, len(e.Guilds))
	util.Log.Infof("Invite link: https://discordapp.com/api/oauth2/authorize?client_id=%s&scope=bot&permissions=%d",
		e.User.ID, static.InvitePermission)

	s.UpdateStatus(0, static.StdMotd)

	l.lct.Start()

	rawPresence, err := l.db.GetSetting(static.SettingPresence)
	if err == nil {
		pre, err := presence.UnmarshalPresence(rawPresence)
		if err == nil {
			s.UpdateStatusComplex(pre.ToUpdateStatusData())
		}
	}

	for _, g := range e.Guilds {
		if err := s.GuildMemberNickname(g.ID, "@me", static.AutoNick); err != nil {
			util.Log.Errorf("Failed updating nickname on guild %s (%s): %s", g.Name, g.ID, err)
		}
	}

	votes, err := l.db.GetVotes()
	if err != nil {
		util.Log.Error("Failed getting votes from DB: ", err)
	} else {
		vote.VotesRunning = votes
		l.lct.OnTick(func(now time.Time) {
			for _, v := range vote.VotesRunning {
				if (v.Expires != time.Time{}) && v.Expires.Before(now) {
					v.Close(s, vote.VoteStateExpired)
					if err = l.db.DeleteVote(v.ID); err != nil {
						util.Log.Errorf("Failed updating vote with ID %s: %s", v.ID, err.Error())
					}
				}
			}
		})
	}
}
