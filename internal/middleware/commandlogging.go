package middleware

import (
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/services/config"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekrotja/ken"
	"github.com/zekrotja/rogu"
	"github.com/zekrotja/rogu/log"
)

type CommandLoggingMiddleware struct {
	cfg config.Provider
	log rogu.Logger
}

var _ ken.MiddlewareAfter = (*CommandLoggingMiddleware)(nil)

func NewCommandLoggingMiddleware(ct di.Container) *CommandLoggingMiddleware {
	return &CommandLoggingMiddleware{
		log: log.Tagged("Command"),
		cfg: ct.Get(static.DiConfig).(config.Provider),
	}
}

func (t *CommandLoggingMiddleware) After(ctx *ken.Ctx, cmdError error) (err error) {
	if !t.cfg.Config().Logging.CommandLogging {
		return nil
	}

	var e *rogu.Event

	if cmdError != nil {
		e = t.log.Error().
			Err(cmdError)
	} else {
		e = t.log.Info()
	}

	e.Fields(
		"byId", ctx.User().ID,
		"byName", ctx.User().String(),
	)

	if ctx.GetEvent().GuildID != "" {
		if guild, err := ctx.Guild(); err == nil {
			e.Fields(
				"guildId", guild.ID,
				"guildName", guild.Name,
			)
		}
	}

	e.Msgf("/%s", ctx.Command.Name())

	return nil
}
