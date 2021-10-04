package wrappers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/zekrotja/dgrs"
	"github.com/zekrotja/ken/state"
)

type StateWrapper struct {
	*dgrs.State
}

var _ state.State = (*StateWrapper)(nil)

func (w *StateWrapper) SelfUser(s *discordgo.Session) (*discordgo.User, error) {
	return w.State.SelfUser()
}

func (w *StateWrapper) Channel(s *discordgo.Session, id string) (*discordgo.Channel, error) {
	return w.State.Channel(id)
}

func (w *StateWrapper) Guild(s *discordgo.Session, id string) (*discordgo.Guild, error) {
	return w.State.Guild(id)
}
