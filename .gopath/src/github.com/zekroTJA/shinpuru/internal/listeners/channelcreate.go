package listeners

import (
	"github.com/bwmarrin/discordgo"

	"github.com/zekroTJA/shinpuru/internal/core"
)

type ListenerChannelCreate struct {
	db core.Database
}

func NewListenerChannelCreate(db core.Database) *ListenerChannelCreate {
	return &ListenerChannelCreate{
		db: db,
	}
}

func (l *ListenerChannelCreate) Handler(s *discordgo.Session, e *discordgo.ChannelCreate) {
	roleID, err := l.db.GetMuteRoleGuild(e.GuildID)
	if err == nil {
		s.ChannelPermissionSet(e.ID, roleID, "role", 0, 0x00000800)
	}
}
