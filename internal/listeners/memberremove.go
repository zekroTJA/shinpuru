package listeners

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/core"
	"github.com/zekroTJA/shinpuru/internal/util"
)

type ListenerMemberRemove struct {
	db core.Database
}

func NewListenerMemberRemove(db core.Database) *ListenerMemberRemove {
	return &ListenerMemberRemove{
		db: db,
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
