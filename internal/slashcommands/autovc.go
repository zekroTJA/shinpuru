package slashcommands

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/services/permissions"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/stringutil"
	"github.com/zekrotja/ken"
)

type Autovc struct{}

var (
	_ ken.SlashCommand        = (*Autovc)(nil)
	_ permissions.PermCommand = (*Autovc)(nil)
)

func (c *Autovc) Name() string {
	return "autovc"
}

func (c *Autovc) Description() string {
	return "Manage guild auto voicechannels."
}

func (c *Autovc) Version() string {
	return "1.0.0"
}

func (c *Autovc) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (c *Autovc) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "show",
			Description: "Display the currently set auto voicechannels.",
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "add",
			Description: "Add an auto voicechannel.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:         discordgo.ApplicationCommandOptionChannel,
					Name:         "voicechannel",
					ChannelTypes: []discordgo.ChannelType{discordgo.ChannelTypeGuildVoice},
					Description:  "The voicechannel to be set.",
					Required:     true,
				},
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "remove",
			Description: "Remove an auto voicechannel.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:         discordgo.ApplicationCommandOptionChannel,
					ChannelTypes: []discordgo.ChannelType{discordgo.ChannelTypeGuildVoice},
					Name:         "voicechannel",
					Description:  "The voicechannel to be set.",
					Required:     true,
				},
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "purge",
			Description: "Unset all auto voicechannels.",
		},
	}
}

func (c *Autovc) Domain() string {
	return "sp.guild.config.autovc"
}

func (c *Autovc) SubDomains() []permissions.SubPermission {
	return nil
}

func (c *Autovc) Run(ctx ken.Context) (err error) {
	if err = ctx.Defer(); err != nil {
		return
	}

	err = ctx.HandleSubCommands(
		ken.SubCommandHandler{"show", c.show},
		ken.SubCommandHandler{"add", c.add},
		ken.SubCommandHandler{"remove", c.remove},
		ken.SubCommandHandler{"purge", c.purge},
	)

	return
}

func (c *Autovc) show(ctx ken.SubCommandContext) (err error) {
	db := ctx.Get(static.DiDatabase).(database.Database)

	autovcs, err := db.GetGuildAutoVC(ctx.GetEvent().GuildID)
	if err != nil {
		return
	}

	if len(autovcs) == 0 {
		err = ctx.FollowUpEmbed(&discordgo.MessageEmbed{
			Description: "Currently, no auto voicechannels are defined.",
		}).Error
		return
	}

	var res strings.Builder
	for _, id := range autovcs {
		res.WriteString(fmt.Sprintf("- <#%s>\n", id))
	}

	err = ctx.FollowUpEmbed(&discordgo.MessageEmbed{
		Description: "Currently, following channels are set as auto voicechannels:\n" + res.String(),
	}).Error

	return
}

func (c *Autovc) add(ctx ken.SubCommandContext) (err error) {
	db := ctx.Get(static.DiDatabase).(database.Database)

	vc := ctx.Options().Get(0).
		ChannelValue(ctx)

	autovcs, err := db.GetGuildAutoVC(ctx.GetEvent().GuildID)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return
	}

	if stringutil.ContainsAny(vc.ID, autovcs) {
		err = ctx.FollowUpError("The given voicechannel is already assigned.", "").Error
		return
	}

	if err = db.SetGuildAutoVC(ctx.GetEvent().GuildID, append(autovcs, vc.ID)); err != nil {
		return
	}

	err = ctx.FollowUpEmbed(&discordgo.MessageEmbed{
		Color:       static.ColorEmbedGreen,
		Description: "Voicechannel was successfully assigned as auto voicechannel.",
	}).Error

	return
}

func (c *Autovc) remove(ctx ken.SubCommandContext) (err error) {
	db := ctx.Get(static.DiDatabase).(database.Database)

	vc := ctx.Options().Get(0).
		ChannelValue(ctx)

	autovcs, err := db.GetGuildAutoVC(ctx.GetEvent().GuildID)
	if err != nil {
		return
	}

	if !stringutil.ContainsAny(vc.ID, autovcs) {
		err = ctx.FollowUpError("The given voicechannel is not assigned as auto voicechannel.", "").Error
		return
	}

	autovcs = stringutil.Splice(autovcs, stringutil.IndexOf(vc.ID, autovcs))
	if err = db.SetGuildAutoVC(ctx.GetEvent().GuildID, autovcs); err != nil {
		return
	}

	err = ctx.FollowUpEmbed(&discordgo.MessageEmbed{
		Color:       static.ColorEmbedGreen,
		Description: "Channel was successfully removed as autochannel.",
	}).Error

	return
}

func (c *Autovc) purge(ctx ken.SubCommandContext) (err error) {
	db := ctx.Get(static.DiDatabase).(database.Database)

	if err = db.SetGuildAutoVC(ctx.GetEvent().GuildID, []string{}); err != nil {
		return
	}

	err = ctx.FollowUpEmbed(&discordgo.MessageEmbed{
		Color:       static.ColorEmbedGreen,
		Description: "All auto voicechannels were successfully removed.",
	}).Error

	return
}
