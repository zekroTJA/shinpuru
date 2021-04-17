package listeners

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/core/database"
	"github.com/zekroTJA/shinpuru/internal/util/static"
)

type ListenerChannelCreate struct {
	db database.Database
}

func NewListenerChannelCreate(container di.Container) *ListenerChannelCreate {
	return &ListenerChannelCreate{
		db: container.Get(static.DiDatabase).(database.Database),
	}
}

func (l *ListenerChannelCreate) Handler(s *discordgo.Session, e *discordgo.ChannelCreate) {
	roleID, err := l.db.GetGuildMuteRole(e.GuildID)
	if err == nil {
		s.ChannelPermissionSet(e.ID, roleID, discordgo.PermissionOverwriteTypeRole, 0, 0x00000800)
	}
}
