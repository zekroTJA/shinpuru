package commands

import (
	"os"
	"strings"
	"time"

	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shireikan"
	"github.com/zekrotja/dgrs"
)

type CmdMaintenance struct {
}

func (c *CmdMaintenance) GetInvokes() []string {
	return []string{"maintenance", "maintain", "service", "mtn"}
}

func (c *CmdMaintenance) GetDescription() string {
	return "Maintenance utilities."
}

func (c *CmdMaintenance) GetHelp() string {
	return "`mtn flush-state (<subKey> (<subKey> (...)))` - Flush dgrs state.\n" +
		"`mtn kill (<exit Code>)` - Kill the bot process\n" +
		"`mtn reconnect` - Reconnects the Discord session"
}

func (c *CmdMaintenance) GetGroup() string {
	return shireikan.GroupGlobalAdmin
}

func (c *CmdMaintenance) GetDomainName() string {
	return "sp.maintenance"
}

func (c *CmdMaintenance) GetSubPermissionRules() []shireikan.SubPermission {
	return nil
}

func (c *CmdMaintenance) IsExecutableInDMChannels() bool {
	return true
}

func (c *CmdMaintenance) Exec(ctx shireikan.Context) error {

	switch strings.ToLower(ctx.GetArgs().Get(0).AsString()) {
	case "flush-state":
		return c.flushState(ctx)
	case "kill":
		return c.kill(ctx)
	case "reconnect":
		return c.reconnect(ctx)
	default:
		return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
			"Invalid sub command.\nUse `help maintenance` for more information.").DeleteAfter(10 * time.Second).Error()
	}
}

func (c *CmdMaintenance) flushState(ctx shireikan.Context) (err error) {
	st := ctx.GetObject(static.DiState).(*dgrs.State)

	keys := ctx.GetArgs()[1:]

	if err = st.Flush(keys...); err != nil {
		return
	}

	ctx.GetSession().Close()
	ctx.GetSession().Open()

	err = util.SendEmbed(ctx.GetSession(), ctx.GetChannel().ID,
		"✅ State cache flushed.", "", static.ColorEmbedGreen).
		DeleteAfter(5 * time.Second).Error()

	return
}

func (c *CmdMaintenance) kill(ctx shireikan.Context) (err error) {
	code := 1

	if ctx.GetArgs().Get(1).AsString() != "" {
		code, err = ctx.GetArgs().Get(1).AsInt()
		if err != nil {
			return
		}
	}

	err = util.SendEmbed(ctx.GetSession(), ctx.GetChannel().ID,
		"👋 Bye.", "", static.ColorEmbedOrange).
		Error()

	os.Exit(code)

	return
}

func (c *CmdMaintenance) reconnect(ctx shireikan.Context) (err error) {
	if err = ctx.GetSession().Close(); err != nil {
		return
	}

	ctx.GetSession().Open()

	err = util.SendEmbed(ctx.GetSession(), ctx.GetChannel().ID,
		"✅ Successfully reconnected.", "", static.ColorEmbedGreen).
		DeleteAfter(5 * time.Second).Error()

	return
}
