package slashcommands

import (
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/services/config"
	"github.com/zekroTJA/shinpuru/internal/services/permissions"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekrotja/dgrs"
	"github.com/zekrotja/ken"
)

type Maintenance struct{}

var (
	_ ken.Command             = (*Maintenance)(nil)
	_ permissions.PermCommand = (*Maintenance)(nil)
	_ ken.DmCapable           = (*Maintenance)(nil)
)

func (c *Maintenance) Name() string {
	return "maintenance"
}

func (c *Maintenance) Description() string {
	return "Maintenance utilities."
}

func (c *Maintenance) Version() string {
	return "1.1.0"
}

func (c *Maintenance) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (c *Maintenance) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "flush-state",
			Description: "Flush dgrs state.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionBoolean,
					Name:        "reconnect",
					Description: "Disconnect and reconnect session after flush.",
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "subkeys",
					Description: "The cache sub keys (comma seperated).",
				},
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "kill",
			Description: "Kill the bot process.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "exitcode",
					Description: "The exit code.",
				},
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "reconnect",
			Description: "Reconnects the Discord session.",
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "reload-config",
			Description: "Reloads the bots config.",
		},
	}
}

func (c *Maintenance) Domain() string {
	return "sp.maintenance"
}

func (c *Maintenance) SubDomains() []permissions.SubPermission {
	return nil
}

func (c *Maintenance) IsDmCapable() bool {
	return true
}

func (c *Maintenance) Run(ctx *ken.Ctx) (err error) {
	if err = ctx.Defer(); err != nil {
		return
	}

	err = ctx.HandleSubCommands(
		ken.SubCommandHandler{"flush-state", c.flushState},
		ken.SubCommandHandler{"kill", c.kill},
		ken.SubCommandHandler{"reconnect", c.reconnect},
		ken.SubCommandHandler{"reload-config", c.reloadConfig},
	)

	return
}

func (c *Maintenance) flushState(ctx *ken.SubCommandCtx) (err error) {
	st := ctx.Get(static.DiState).(*dgrs.State)

	subkeys := ([]string)(nil)
	if subkeysV, ok := ctx.Options().GetByNameOptional("subkeys"); ok {
		subkeys = strings.Split(subkeysV.StringValue(), ",")
		for i, sk := range subkeys {
			subkeys[i] = strings.TrimSpace(sk)
		}
	}

	if err = st.Flush(subkeys...); err != nil {
		return
	}

	if reconnectV, ok := ctx.Options().GetByNameOptional("reconnect"); ok && reconnectV.BoolValue() {
		ctx.Session.Close()
		ctx.Session.Open()
	}

	return ctx.FollowUpEmbed(&discordgo.MessageEmbed{
		Description: "✅ State cache flushed.",
		Color:       static.ColorEmbedGreen,
	}).Error
}

func (c *Maintenance) kill(ctx *ken.SubCommandCtx) (err error) {
	code := 1

	if exitcodeV, ok := ctx.Options().GetByNameOptional("exitcode"); ok {
		code = int(exitcodeV.IntValue())
	}

	err = ctx.FollowUpEmbed(&discordgo.MessageEmbed{
		Description: "👋 Bye.",
		Color:       static.ColorEmbedOrange,
	}).Error
	if err != nil {
		return
	}

	os.Exit(code)

	return
}

func (c *Maintenance) reconnect(ctx *ken.SubCommandCtx) (err error) {
	if err = ctx.Session.Close(); err != nil {
		return
	}

	ctx.Session.Open()

	return ctx.FollowUpEmbed(&discordgo.MessageEmbed{
		Description: "✅ Successfully reconnected.",
		Color:       static.ColorEmbedGreen,
	}).Error
}

func (c *Maintenance) reloadConfig(ctx *ken.SubCommandCtx) (err error) {
	cfg := ctx.Get(static.DiConfig).(config.Provider)

	if err = cfg.Parse(); err != nil {
		return
	}

	return ctx.FollowUpEmbed(&discordgo.MessageEmbed{
		Description: "Config has been reloaded.\n\nSome config changes will only take effect after a restart!",
	}).Error
}
