package slashcommands

import (
	"fmt"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/services/config"
	"github.com/zekroTJA/shinpuru/internal/services/permissions"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/mody"
	"github.com/zekrotja/dgrs"
	"github.com/zekrotja/ken"
)

type Maintenance struct {
	ken.EphemeralCommand
}

var (
	_ ken.SlashCommand        = (*Maintenance)(nil)
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
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "set-config-value",
			Description: "Set a specific config value.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "field",
					Description: "The config fild path and name.",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "jsonvalue",
					Description: "The value as JSON representation.",
					Required:    true,
				},
			},
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

func (c *Maintenance) Run(ctx ken.Context) (err error) {
	if err = ctx.Defer(); err != nil {
		return
	}

	err = ctx.HandleSubCommands(
		ken.SubCommandHandler{"flush-state", c.flushState},
		ken.SubCommandHandler{"kill", c.kill},
		ken.SubCommandHandler{"reconnect", c.reconnect},
		ken.SubCommandHandler{"reload-config", c.reloadConfig},
		ken.SubCommandHandler{"set-config-value", c.setConfigValue},
	)

	return
}

func (c *Maintenance) flushState(ctx ken.SubCommandContext) (err error) {
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
		ctx.GetSession().Close()
		ctx.GetSession().Open()
	}

	return ctx.FollowUpEmbed(&discordgo.MessageEmbed{
		Description: "✅ State cache flushed.",
		Color:       static.ColorEmbedGreen,
	}).Send().Error
}

func (c *Maintenance) kill(ctx ken.SubCommandContext) (err error) {
	code := 1

	if exitcodeV, ok := ctx.Options().GetByNameOptional("exitcode"); ok {
		code = int(exitcodeV.IntValue())
	}

	err = ctx.FollowUpEmbed(&discordgo.MessageEmbed{
		Description: "👋 Bye.",
		Color:       static.ColorEmbedOrange,
	}).Send().Error
	if err != nil {
		return
	}

	os.Exit(code)

	return
}

func (c *Maintenance) reconnect(ctx ken.SubCommandContext) (err error) {
	if err = ctx.GetSession().Close(); err != nil {
		return
	}

	ctx.GetSession().Open()

	return ctx.FollowUpEmbed(&discordgo.MessageEmbed{
		Description: "✅ Successfully reconnected.",
		Color:       static.ColorEmbedGreen,
	}).Send().Error
}

func (c *Maintenance) reloadConfig(ctx ken.SubCommandContext) (err error) {
	cfg := ctx.Get(static.DiConfig).(config.Provider)

	if err = cfg.Parse(); err != nil {
		return
	}

	return ctx.FollowUpEmbed(&discordgo.MessageEmbed{
		Description: "Config has been reloaded.\n\nSome config changes will only take effect after a restart!",
	}).Send().Error
}

func (c *Maintenance) setConfigValue(ctx ken.SubCommandContext) (err error) {
	cfg := ctx.Get(static.DiConfig).(config.Provider)

	field := ctx.Options().GetByName("field").StringValue()
	jsonvalue := ctx.Options().GetByName("jsonvalue").StringValue()

	var errInner error
	err = mody.Catch(func() {
		errInner = mody.UpdateJson(cfg.Config(), field, jsonvalue)
	})
	if err != nil {
		return
	}
	if err = errInner; err != nil {
		return
	}

	return ctx.FollowUpEmbed(&discordgo.MessageEmbed{
		Description: fmt.Sprintf("Config value `%s` has been updated to `%s`.\n\n"+
			"Keep in mind that not all config value changes will be effective.",
			field, jsonvalue),
	}).Send().Error
}
