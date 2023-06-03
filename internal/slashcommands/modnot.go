package slashcommands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/services/permissions"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/acceptmsg/v2"
	"github.com/zekrotja/ken"
)

type Modnot struct{}

var (
	_ ken.SlashCommand        = (*Modnot)(nil)
	_ permissions.PermCommand = (*Modnot)(nil)
)

func (c *Modnot) Name() string {
	return "modnot"
}

func (c *Modnot) Description() string {
	return "Set the mod notification channel for a guild."
}

func (c *Modnot) Version() string {
	return "1.0.0"
}

func (c *Modnot) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (c *Modnot) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "set",
			Description: "Set this or a specified channel as mod notification channel.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:         discordgo.ApplicationCommandOptionChannel,
					Name:         "channel",
					Description:  "A channel to be set as mod notification channel (current channel if not specified).",
					ChannelTypes: []discordgo.ChannelType{discordgo.ChannelTypeGuildText},
				},
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "disable",
			Description: "Disable mod notifications.",
		},
	}
}

func (c *Modnot) Domain() string {
	return "sp.guild.config.modnot"
}

func (c *Modnot) SubDomains() []permissions.SubPermission {
	return nil
}

func (c *Modnot) Run(ctx ken.Context) (err error) {
	if err = ctx.Defer(); err != nil {
		return
	}

	err = ctx.HandleSubCommands(
		ken.SubCommandHandler{"set", c.set},
		ken.SubCommandHandler{"disable", c.disable},
	)

	return
}

func (c *Modnot) set(ctx ken.SubCommandContext) (err error) {
	db := ctx.Get(static.DiDatabase).(database.Database)

	chV, ok := ctx.Options().GetByNameOptional("channel")

	if !ok {
		acceptMsg := &acceptmsg.AcceptMessage{
			Ken: ctx.GetKen(),
			Embed: &discordgo.MessageEmbed{
				Color:       static.ColorEmbedDefault,
				Description: "Do you want to set this channel as mod notification channel?",
			},
			UserID:         ctx.User().ID,
			DeleteMsgAfter: true,
			AcceptFunc: func(cctx ken.ComponentContext) (err error) {
				if err = cctx.Defer(); err != nil {
					return
				}
				err = db.SetGuildModNot(ctx.GetEvent().GuildID, ctx.GetEvent().ChannelID)
				if err != nil {
					return
				}
				err = cctx.FollowUpEmbed(&discordgo.MessageEmbed{
					Description: "Set this channel as mod notification channel.",
				}).Send().Error
				return
			},
		}

		if _, err = acceptMsg.AsFollowUp(ctx); err != nil {
			return
		}
		return acceptMsg.Error()
	}

	ch := chV.ChannelValue(ctx)

	if err = db.SetGuildModNot(ctx.GetEvent().GuildID, ch.ID); err != nil {
		return
	}

	err = ctx.FollowUpEmbed(&discordgo.MessageEmbed{
		Description: fmt.Sprintf("Set channel <#%s> as mod notification channel.", ch.ID),
	}).Send().Error

	return
}

func (c *Modnot) disable(ctx ken.SubCommandContext) (err error) {
	db := ctx.Get(static.DiDatabase).(database.Database)

	if err = db.SetGuildModNot(ctx.GetEvent().GuildID, ""); err != nil {
		return
	}

	return ctx.FollowUpEmbed(&discordgo.MessageEmbed{
		Description: "Mod notifications disabled.",
	}).Send().Error
}
