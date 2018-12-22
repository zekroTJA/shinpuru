package listeners

import (
	"github.com/bwmarrin/discordgo"

	"github.com/zekroTJA/shinpuru/core"
	"github.com/zekroTJA/shinpuru/util"
)

type ListenerReady struct {
	config *core.Config
}

func NewListenerReady(config *core.Config) *ListenerReady {
	return &ListenerReady{
		config: config,
	}
}

func (l *ListenerReady) Handler(s *discordgo.Session, e *discordgo.Ready) {
	util.Log.Infof("Logged in as %s#%s (%s) - Running on %d servers",
		e.User.Username, e.User.Discriminator, e.User.ID, len(e.Guilds))
	util.Log.Infof("Invite link: https://discordapp.com/api/oauth2/authorize?client_id=%s&scope=bot&permissions=%d",
		e.User.ID, util.InvitePermission)

	s.UpdateStatus(0, util.StdMotd)
	for _, g := range e.Guilds {
		if err := s.GuildMemberNickname(g.ID, "@me", util.AutoNick); err != nil {
			util.Log.Errorf("Failed updating nickname on guild %s (%s): %s", g.Name, g.ID, err)
		}
	}
}
