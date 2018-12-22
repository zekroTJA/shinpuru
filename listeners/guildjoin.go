package listeners

import (
	"github.com/bwmarrin/discordgo"

	"github.com/zekroTJA/shinpuru/core"
	"github.com/zekroTJA/shinpuru/util"
)

type ListenerGuildJoin struct {
	config *core.Config
}

func NewListenerGuildJoin(config *core.Config) *ListenerGuildJoin {
	return &ListenerGuildJoin{
		config: config,
	}
}

func (l *ListenerGuildJoin) Handler(s *discordgo.Session, e *discordgo.GuildCreate) {
	if err := s.GuildMemberNickname(e.Guild.ID, "@me", util.AutoNick); err != nil {
		util.Log.Errorf("Failed updating nickname on guild %s (%s): %s", e.Guild.Name, e.Guild.ID, err)
	}
}
