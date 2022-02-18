package listeners

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sarulabs/di/v2"
	"github.com/sirupsen/logrus"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekrotja/dgrs"
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
		logrus.WithError(err).Error("Failed publishing DM to state")
	}
}
