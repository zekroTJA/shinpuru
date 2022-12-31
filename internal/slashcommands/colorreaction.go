package slashcommands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/services/permissions"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/intutil"
	"github.com/zekroTJA/shinpuru/pkg/stringutil"
	"github.com/zekrotja/ken"
)

type Colorreation struct {
	ken.EphemeralCommand
}

var (
	_ ken.SlashCommand        = (*Colorreation)(nil)
	_ permissions.PermCommand = (*Colorreation)(nil)
)

func (c *Colorreation) Name() string {
	return "colorreaction"
}

func (c *Colorreation) Description() string {
	return "Toggle color reactions enable or disable."
}

func (c *Colorreation) Version() string {
	return "1.0.0"
}

func (c *Colorreation) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (c *Colorreation) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionBoolean,
			Name:        "enable",
			Description: "Set the enabled state of color reactions.",
		},
	}
}

func (c *Colorreation) Domain() string {
	return "sp.guild.config.color"
}

func (c *Colorreation) SubDomains() []permissions.SubPermission {
	return []permissions.SubPermission{
		{
			Term:        "/sp.chat.colorreactions",
			Explicit:    false,
			Description: "Allows executing color reactions in chat by reaction",
		},
	}
}

func (c *Colorreation) Run(ctx ken.Context) (err error) {
	if err = ctx.Defer(); err != nil {
		return
	}

	db := ctx.Get(static.DiDatabase).(database.Database)

	var enable bool
	enableV, ok := ctx.Options().GetByNameOptional("enable")
	if ok {
		enable = enableV.BoolValue()
		if err = db.SetGuildColorReaction(ctx.GetEvent().GuildID, enable); err != nil {
			return
		}
		err = ctx.FollowUpEmbed(&discordgo.MessageEmbed{
			Description: fmt.Sprintf("Color reaction has been %s.",
				stringutil.FromBool(enable, "enabled", "disabled")),
			Color: intutil.FromBool(enable, static.ColorEmbedGreen, static.ColorEmbedOrange),
		}).Send().Error
	} else {
		enable, err = db.GetGuildColorReaction(ctx.GetEvent().GuildID)
		if err != nil {
			return
		}
		err = ctx.FollowUpEmbed(&discordgo.MessageEmbed{
			Description: fmt.Sprintf("Color reaction is currently %s.",
				stringutil.FromBool(enable, "enabled", "disabled")),
			Color: intutil.FromBool(enable, static.ColorEmbedGreen, static.ColorEmbedOrange),
		}).Send().Error
	}

	return
}
