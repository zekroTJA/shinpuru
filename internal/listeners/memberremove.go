package listeners

import (
	"strings"

	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekrotja/discordgo"
)

type ListenerMemberRemove struct {
	db database.Database
}

func NewListenerMemberRemove(container di.Container) *ListenerMemberRemove {
	return &ListenerMemberRemove{
		db: container.Get(static.DiDatabase).(database.Database),
	}
}

func (l *ListenerMemberRemove) Handler(s *discordgo.Session, e *discordgo.GuildMemberRemove) {
	chanID, msg, err := l.db.GetGuildLeaveMsg(e.GuildID)
	if err == nil && msg != "" && chanID != "" {
		msg = strings.Replace(msg, "[user]", e.User.Username, -1)
		msg = strings.Replace(msg, "[ment]", e.User.Mention(), -1)

		util.SendEmbed(s, chanID, msg, "", 0)
	}
}
