package listeners

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sarulabs/di/v2"

	"github.com/zekroTJA/shinpuru/internal/core/database"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/internal/util/vote"
)

type ListenerVote struct {
	db database.Database
}

func NewListenerVote(container di.Container) *ListenerVote {
	return &ListenerVote{
		db: container.Get(static.DiDatabase).(database.Database),
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
	for _, v := range vote.VotesRunning {
		if v.GuildID != e.GuildID || v.ChannelID != e.ChannelID || v.MsgID != e.MessageID {
			continue
		}
		tick := -1
		for i, ve := range vote.VoteEmotes {
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
