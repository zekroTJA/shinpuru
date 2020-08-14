package middleware

import (
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shireikan"
)

// CommandStatsMiddleware implements shireikan.Middleware to
// count executed commands.
type CommandStatsMiddleware struct{}

func (m *CommandStatsMiddleware) Handle(
	cmd shireikan.Command, ctx shireikan.Context, layer shireikan.MiddlewareLayer) (next bool, err error) {

	util.StatsCommandsExecuted++

	return true, nil
}

func (m *CommandStatsMiddleware) GetLayer() shireikan.MiddlewareLayer {
	return shireikan.LayerAfterCommand
}
