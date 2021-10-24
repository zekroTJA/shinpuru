package middleware

import (
	"time"

	"github.com/zekroTJA/shireikan"
	"github.com/zekroTJA/timedmap"
)

// GhostPingIgnoreMiddleware implements the shireikan.Middleware
// interface to provide message caching for command messages
// containing user mentions to dodge the ghostping trigger.
type GhostPingIgnoreMiddleware struct {
	reg *timedmap.TimedMap
}

// NewGhostPingIgnoreMiddleware returns a new instance of
// GhostPingIgnoreMiddleware.
func NewGhostPingIgnoreMiddleware() *GhostPingIgnoreMiddleware {
	return &GhostPingIgnoreMiddleware{
		timedmap.New(15 * time.Minute),
	}
}

func (m *GhostPingIgnoreMiddleware) Handle(
	cmd shireikan.Command, ctx shireikan.Context, layer shireikan.MiddlewareLayer) (next bool, err error) {

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

// ContainsAndRemove returns true when the passed msgID is
// present in the mentions cache register. If this is true,
// the entry is removed from the register.
func (m *GhostPingIgnoreMiddleware) ContainsAndRemove(msgID string) bool {
	mentions, ok := m.reg.GetValue(msgID).([]string)
	if !ok || mentions == nil {
		return false
	}

	m.reg.Remove(msgID)
	return true
}
