package listeners

import (
	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/core/database"
)

type ListenerChannelCreate struct {
	db database.Database
}

func NewListenerChannelCreate(db database.Database) *ListenerChannelCreate {
	return &ListenerChannelCreate{
		db: db,
	}
}

func (l *ListenerChannelCreate) Handler(s *discordgo.Session, e *discordgo.ChannelCreate) {
	roleID, err := l.db.GetGuildMuteRole(e.GuildID)
	if err == nil {
		s.ChannelPermissionSet(e.ID, roleID, "role", 0, 0x00000800)
	}
}
