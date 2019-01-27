package listeners

import (
	"github.com/zekroTJA/shinpuru/internal/core"
)

type ListenerTwitchNotify struct {
	config *core.Config
	db     core.Database
}

func NewListenerTwitchNotify(config *core.Config, db core.Database) *ListenerTwitchNotify {
	return &ListenerTwitchNotify{
		config: config,
		db:     db,
	}
}

func (l *ListenerTwitchNotify) Handler(d *core.TwitchNotifyData, u *core.TwitchNotifyUser) {

}
