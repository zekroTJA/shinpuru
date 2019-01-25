package listeners

import (
	"github.com/bwmarrin/discordgo"

	"github.com/zekroTJA/shinpuru/internal/core"
	"github.com/zekroTJA/shinpuru/internal/util"
)

type ListenerReady struct {
	config *core.Config
	db     core.Database
}

func NewListenerReady(config *core.Config, db core.Database) *ListenerReady {
	return &ListenerReady{
		config: config,
		db:     db,
	}
}

func (l *ListenerReady) Handler(s *discordgo.Session, e *discordgo.Ready) {
	util.Log.Infof("Logged in as %s#%s (%s) - Running on %d servers",
		e.User.Username, e.User.Discriminator, e.User.ID, len(e.Guilds))
	util.Log.Infof("Invite link: https://discordapp.com/api/oauth2/authorize?client_id=%s&scope=bot&permissions=%d",
		e.User.ID, util.InvitePermission)

	s.UpdateStatus(0, util.StdMotd)

	rawPresence, err := l.db.GetSetting(util.SettingPresence)
	if err == nil {
		presence, err := util.UnmarshalPresence(rawPresence)
		if err == nil {
			s.UpdateStatusComplex(presence.ToUpdateStatusData())
		}
	}

	for _, g := range e.Guilds {
		if err := s.GuildMemberNickname(g.ID, "@me", util.AutoNick); err != nil {
			util.Log.Errorf("Failed updating nickname on guild %s (%s): %s", g.Name, g.ID, err)
		}
	}

	votes, err := l.db.GetVotes()
	if err != nil {
		util.Log.Error("Failed getting votes from DB: ", err)
	} else {
		util.VotesRunning = votes
	}
}
