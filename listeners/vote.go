package listeners

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/util"

	"github.com/zekroTJA/shinpuru/core"
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
			v.Tick(s, e.UserID, tick)
		}
		s.MessageReactionRemove(e.ChannelID, e.MessageID, e.Emoji.Name, e.UserID)
	}
	fmt.Println(util.VotesRunning)
}
