package middleware

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/services/config"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shireikan"
	"github.com/zekroTJA/timedmap"
)

type DeprecationMiddleware struct {
	cfg     config.Provider
	timeout *timedmap.TimedMap
}

var _ shireikan.Middleware = (*DeprecationMiddleware)(nil)

func NewDeprecationMiddleware(ctn di.Container) *DeprecationMiddleware {
	return &DeprecationMiddleware{
		cfg:     ctn.Get(static.DiConfig).(config.Provider),
		timeout: timedmap.New(1 * time.Hour),
	}
}

func (m *DeprecationMiddleware) GetLayer() shireikan.MiddlewareLayer {
	return shireikan.LayerAfterCommand
}

func (m *DeprecationMiddleware) Handle(cmd shireikan.Command, ctx shireikan.Context, layer shireikan.MiddlewareLayer) (next bool, err error) {
	if m.timeout.GetValue(ctx.GetUser().ID) != nil {
		return
	}

	msg := "The legacy command system is now deprecated and will be removed in future versions. " +
		"[Here](https://github.com/zekroTJA/shinpuru/wiki/Legacy-Command-Deprecation) you can read more about that.\n\n" +
		"Just type `/` to a text channel to get a list of available commands."

	if pubURL := m.cfg.Config().WebServer.PublicAddr; pubURL != "" {
		msg += fmt.Sprintf(" [Here](%s/commands) you can also find an interactive list of available slash commands.", pubURL)
	}

	ctx.ReplyEmbed(&discordgo.MessageEmbed{
		Color:       static.ColorEmbedOrange,
		Title:       "⚠️ Deprecation",
		Description: msg,
	})

	m.timeout.Set(ctx.GetUser().ID, true, 24*time.Hour)
	next = true
	return
}
