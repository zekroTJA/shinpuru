package listeners

import (
	"sync/atomic"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/util"
)

type ListenerStatus struct{}

func NewListenerStatus() *ListenerStatus {
	return &ListenerStatus{}
}

func (t ListenerStatus) ListenerConnect(s *discordgo.Session, e *discordgo.Connect) {
	atomic.StoreInt32(&util.ConnectedState, 1)
}

func (t ListenerStatus) ListenerDisconnect(s *discordgo.Session, e *discordgo.Disconnect) {
	atomic.StoreInt32(&util.ConnectedState, 0)
}
