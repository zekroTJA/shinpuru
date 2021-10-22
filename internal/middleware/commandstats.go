package middleware

import (
	"sync/atomic"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/zekroTJA/shinpuru/internal/services/metrics"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekrotja/ken"
)

// CommandStatsMiddleware implements ken.MiddlewareAfter to
// count command exeuction stats.
type CommandStatsMiddleware struct{}

var _ ken.MiddlewareAfter = (*CommandStatsMiddleware)(nil)

// NewCommandStatsMiddleware returns a new instance of
// CommandStatsMiddleware.
func NewCommandStatsMiddleware() *CommandStatsMiddleware {
	return &CommandStatsMiddleware{}
}

func (m *CommandStatsMiddleware) After(ctx *ken.Ctx, cmdError error) (err error) {
	name := ctx.Command.Name()

	metrics.DiscordCommandsProcessed.
		With(prometheus.Labels{"command": name}).
		Add(1)

	atomic.AddUint64(&util.StatsCommandsExecuted, 1)

	return
}
