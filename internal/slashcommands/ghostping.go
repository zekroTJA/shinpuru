package slashcommands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/services/permissions"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekrotja/ken"
	"github.com/zekrotja/ken/middlewares/cmdhelp"
)

type Ghostping struct{}

var (
	_ ken.Command             = (*Ghostping)(nil)
	_ permissions.PermCommand = (*Ghostping)(nil)
	_ cmdhelp.HelpProvider    = (*Ghostping)(nil)
)

func (c *Ghostping) Name() string {
	return "ghostping"
}

func (c *Ghostping) Description() string {
	return "Setup the ghost ping system."
}

func (c *Ghostping) Version() string {
	return "1.0.1"
}

func (c *Ghostping) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (c *Ghostping) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "status",
			Description: "Display the current status of the ghost ping settings.",
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "setup",
			Description: "Setup ghostping messages.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "message",
					Description: "The ghost ping message pattern. Use `/ghostping help` to get more info.",
					Required:    true,
				},
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "disable",
			Description: "Disable ghostping messages.",
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "help",
			Description: "Display help about the message pattern which can be used.",
		},
	}
}

func (c *Ghostping) Help(ctx *ken.SubCommandCtx) (emb *discordgo.MessageEmbed, err error) {
	emb = &discordgo.MessageEmbed{
		Color: static.ColorEmbedDefault,
		Description: "You can use patterns in the ghost ping message:\n" +
			"- `{pinger}` - Username of the user who sent the message.\n" +
			"- `{@pinger}` - Mention of the user who sent the message.\n" +
			"- `{pinged}` - Username of the user who has been pinged.\n" +
			"- `{@pinged}` - Mention of the user who has been pinged.\n" +
			"- `{msg}` - The content of the original deleted message.\n" +
			"\n" +
			"If you want to use line breaks, use `\n` in the message parameter.",
	}
	return
}

func (c *Ghostping) Domain() string {
	return "sp.guild.mod.ghostping"
}

func (c *Ghostping) SubDomains() []permissions.SubPermission {
	return nil
}

func (c *Ghostping) Run(ctx *ken.Ctx) (err error) {
	if err = ctx.Defer(); err != nil {
		return
	}

	err = ctx.HandleSubCommands(
		ken.SubCommandHandler{"status", c.status},
		ken.SubCommandHandler{"setup", c.setup},
		ken.SubCommandHandler{"disable", c.disable},
	)

	return
}

func (c *Ghostping) status(ctx *ken.SubCommandCtx) (err error) {
	db := ctx.Get(static.DiDatabase).(database.Database)
	gpMsg, err := db.GetGuildGhostpingMsg(ctx.Event.GuildID)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return err
	}

	msg := "Ghostping is currently not set up.\n\n" +
		"If you want to set up ghostping, use `/ghostping setup`. If you want to " +
		"get info about the message pattern, enter `/ghostping help`."
	if gpMsg != "" {
		msg = "Ghostping is set up with the following message.\n" +
			"```\n" + gpMsg + "\n```\n" +
			"If you want to disable Ghostping, use the `/ghostping disable` command."
	}

	return ctx.FollowUpEmbed(&discordgo.MessageEmbed{
		Description: msg,
	}).Error
}

func (c *Ghostping) setup(ctx *ken.SubCommandCtx) (err error) {
	db := ctx.Get(static.DiDatabase).(database.Database)

	message := ctx.Options().GetByName("message").StringValue()

	if err = db.SetGuildGhostpingMsg(ctx.Event.GuildID, message); err != nil {
		return
	}

	return ctx.FollowUpEmbed(&discordgo.MessageEmbed{
		Description: "Ghostping is now set up with the following message.\n" +
			"```\n" + message + "\n```\n" +
			"If you want to disable Ghostping, use the `/ghostping disable` command.",
	}).Error
}

func (c *Ghostping) disable(ctx *ken.SubCommandCtx) (err error) {
	db := ctx.Get(static.DiDatabase).(database.Database)

	if err = db.SetGuildGhostpingMsg(ctx.Event.GuildID, ""); err != nil {
		return
	}

	return ctx.FollowUpEmbed(&discordgo.MessageEmbed{
		Description: "Ghostping is now disabled.",
	}).Error
}
