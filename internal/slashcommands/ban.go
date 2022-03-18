package slashcommands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/models"
	"github.com/zekroTJA/shinpuru/internal/services/permissions"
	"github.com/zekroTJA/shinpuru/internal/util/cmdutil"
	"github.com/zekrotja/ken"
)

type Ban struct{}

var (
	_ ken.SlashCommand        = (*Ban)(nil)
	_ permissions.PermCommand = (*Ban)(nil)
)

func (c *Ban) Name() string {
	return "ban"
}

func (c *Ban) Description() string {
	return "ban a member with creating a report."
}

func (c *Ban) Version() string {
	return "1.0.0"
}

func (c *Ban) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (c *Ban) Options() []*discordgo.ApplicationCommandOption {
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
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "expire",
			Description: "Expire report after given time.",
		},
	}
}

func (c *Ban) Domain() string {
	return "sp.guild.mod.ban"
}

func (c *Ban) SubDomains() []permissions.SubPermission {
	return []permissions.SubPermission{}
}

func (c *Ban) Run(ctx *ken.Ctx) (err error) {
	if err = ctx.Defer(); err != nil {
		return
	}
	return cmdutil.CmdReport(ctx, models.TypeBan)
}
