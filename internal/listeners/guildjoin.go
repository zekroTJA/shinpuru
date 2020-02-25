package listeners

import (
	"github.com/bwmarrin/discordgo"

	"github.com/zekroTJA/shinpuru/internal/core/config"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/static"
)

type ListenerGuildJoin struct {
	config *config.Config
}

func NewListenerGuildJoin(config *config.Config) *ListenerGuildJoin {
	return &ListenerGuildJoin{
		config: config,
	}
}

func (l *ListenerGuildJoin) Handler(s *discordgo.Session, e *discordgo.GuildCreate) {
	if err := s.GuildMemberNickname(e.Guild.ID, "@me", static.AutoNick); err != nil {
		util.Log.Errorf("Failed updating nickname on guild %s (%s): %s", e.Guild.Name, e.Guild.ID, err)
	}
}
