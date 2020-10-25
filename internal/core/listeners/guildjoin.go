package listeners

import (
	"github.com/bwmarrin/discordgo"

	"github.com/zekroTJA/shinpuru/internal/core/config"
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

}
