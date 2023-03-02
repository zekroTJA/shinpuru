package listeners

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekrotja/dgrs"
	"github.com/zekrotja/rogu/log"
)

type ListenerDMSync struct {
	st *dgrs.State
}

func NewListenerDMSync(container di.Container) *ListenerDMSync {
	return &ListenerDMSync{
		st: container.Get(static.DiState).(*dgrs.State),
	}
}

func (l *ListenerDMSync) Handler(s *discordgo.Session, e *discordgo.MessageCreate) {
	ch, _ := l.st.Channel(e.ChannelID)
	if ch == nil || ch.Type != discordgo.ChannelTypeDM && ch.Type != discordgo.ChannelTypeGroupDM {
		return
	}

	err := l.st.Publish("dms", e)
	if err != nil {
		log.Error().Tag("DMSync").Err(err).Msg("Failed publishing DM to state")
	}
}
