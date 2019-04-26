package listeners

import (
	"github.com/bwmarrin/discordgo"

	"github.com/zekroTJA/shinpuru/internal/core"
	"github.com/zekroTJA/shinpuru/internal/util"
)

type ListenerVote struct {
	db core.Database
}

func NewListenerVote(db core.Database) *ListenerVote {
	return &ListenerVote{
		db: db,
	}
}

func (l *ListenerVote) Handler(s *discordgo.Session, e *discordgo.MessageReactionAdd) {
	user, err := s.User(e.UserID)
	if err != nil {
		return
	}
	if user == nil || user.Bot || user.ID == s.State.User.ID {
		return
	}
	for _, v := range util.VotesRunning {
		if v.GuildID != e.GuildID || v.ChannelID != e.ChannelID || v.MsgID != e.MessageID {
			continue
		}
		tick := -1
		for i, ve := range util.VoteEmotes {
			if e.Emoji.Name == ve {
				tick = i
			}
		}
		if tick > -1 {
			go func() {
				v.Tick(s, e.UserID, tick)
				l.db.AddUpdateVote(v)
			}()
		}
		s.MessageReactionRemove(e.ChannelID, e.MessageID, e.Emoji.Name, e.UserID)
	}
}
