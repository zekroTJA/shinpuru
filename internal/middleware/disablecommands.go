package middleware

import (
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/services/config"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/stringutil"
	"github.com/zekrotja/ken"
)

type DisableCommandsMiddleware struct {
	cfg config.Provider
}

var (
	_ ken.MiddlewareBefore = (*DisableCommandsMiddleware)(nil)
)

func NewDisableCommandsMiddleware(ctn di.Container) *DisableCommandsMiddleware {
	return &DisableCommandsMiddleware{
		cfg: ctn.Get(static.DiConfig).(config.Provider),
	}
}

func (m *DisableCommandsMiddleware) Before(ctx *ken.Ctx) (next bool, err error) {
	next = true

	if m.isDisabled(ctx.Command.Name()) {
		next = false
		err = ctx.RespondError("This command is disabled by config.", "")
	}

	return
}

func (m *DisableCommandsMiddleware) isDisabled(invoke string) bool {
	disabledCmds := m.cfg.Config().Discord.DisabledCommands
	return len(disabledCmds) != 0 && stringutil.ContainsAny(invoke, disabledCmds)
}
