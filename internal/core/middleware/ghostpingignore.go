package middleware

import (
	"time"

	"github.com/zekroTJA/shireikan"
	"github.com/zekroTJA/timedmap"
)

// TODO: docs

type GhostPingIgnoreMiddleware struct {
	reg *timedmap.TimedMap
}

func NewGhostPingIgnoreMiddleware() *GhostPingIgnoreMiddleware {
	return &GhostPingIgnoreMiddleware{
		timedmap.New(15 * time.Minute),
	}
}

func (m *GhostPingIgnoreMiddleware) Handle(cmd shireikan.Command, ctx shireikan.Context) (next bool, err error) {
	next = true

	mentions := ctx.GetMessage().Mentions

	if mentions == nil || len(mentions) < 1 {
		return
	}

	mentionsStr := make([]string, len(mentions))
	for i, u := range mentions {
		mentionsStr[i] = u.ID
	}

	m.reg.Set(ctx.GetMessage().ID, mentionsStr, 10*time.Minute)

	return
}

func (m *GhostPingIgnoreMiddleware) GetLayer() shireikan.MiddlewareLayer {
	return shireikan.LayerBeforeCommand
}

func (m *GhostPingIgnoreMiddleware) ContainsAndRemove(msgID string) bool {
	mentions, ok := m.reg.GetValue(msgID).([]string)
	if !ok || mentions == nil {
		return false
	}

	m.reg.Remove(msgID)
	return true
}
