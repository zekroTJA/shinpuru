package slashcommands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/models"
	"github.com/zekroTJA/shinpuru/internal/services/permissions"
	"github.com/zekroTJA/shinpuru/internal/services/timeprovider"
	"github.com/zekroTJA/shinpuru/internal/util/cmdutil"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekrotja/ken"
)

type Kick struct{}

var (
	_ ken.SlashCommand        = (*Kick)(nil)
	_ permissions.PermCommand = (*Kick)(nil)
)

func (c *Kick) Name() string {
	return "kick"
}

func (c *Kick) Description() string {
	return "Kick a member with creating a report."
}

func (c *Kick) Version() string {
	return "1.0.0"
}

func (c *Kick) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (c *Kick) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionUser,
			Name:        "user",
			Description: "The user.",
			Required:    true,
		},
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "reason",
			Description: "A short and concise report reason.",
			Required:    true,
		},
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "imageurl",
			Description: "An image url embedded into the report.",
		},
	}
}

func (c *Kick) Domain() string {
	return "sp.guild.mod.kick"
}

func (c *Kick) SubDomains() []permissions.SubPermission {
	return []permissions.SubPermission{}
}

func (c *Kick) Run(ctx *ken.Ctx) (err error) {
	if err = ctx.Defer(); err != nil {
		return
	}

	tp := ctx.Get(static.DiTimeProvider).(timeprovider.Provider)

	return cmdutil.CmdReport(ctx, models.TypeKick, tp)
}
