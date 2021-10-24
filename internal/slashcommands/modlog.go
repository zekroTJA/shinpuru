package slashcommands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/services/permissions"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/acceptmsg"
	"github.com/zekrotja/ken"
)

type Modlog struct{}

var (
	_ ken.Command             = (*Modlog)(nil)
	_ permissions.PermCommand = (*Modlog)(nil)
)

func (c *Modlog) Name() string {
	return "modlog"
}

func (c *Modlog) Description() string {
	return "Set the mod log channel for a guild."
}

func (c *Modlog) Version() string {
	return "1.0.0"
}

func (c *Modlog) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (c *Modlog) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "set",
			Description: "Set this or a specified channel as mod log channel.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionChannel,
					Name:        "channel",
					Description: "A channel to be set as mod log (current channel if not specified).",
				},
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "disable",
			Description: "Disable modlog.",
		},
	}
}

func (c *Modlog) Domain() string {
	return "sp.guild.config.modlog"
}

func (c *Modlog) SubDomains() []permissions.SubPermission {
	return nil
}

func (c *Modlog) Run(ctx *ken.Ctx) (err error) {
	if err = ctx.Defer(); err != nil {
		return
	}

	err = ctx.HandleSubCommands(
		ken.SubCommandHandler{"set", c.set},
		ken.SubCommandHandler{"disable", c.disable},
	)

	return
}

func (c *Modlog) set(ctx *ken.SubCommandCtx) (err error) {
	db := ctx.Get(static.DiDatabase).(database.Database)

	chV, ok := ctx.Options().GetByNameOptional("channel")

	if !ok {
		acceptMsg := &acceptmsg.AcceptMessage{
			Session: ctx.Session,
			Embed: &discordgo.MessageEmbed{
				Color:       static.ColorEmbedDefault,
				Description: "Do you want to set this channel as modlog channel?",
			},
			UserID:         ctx.User().ID,
			DeleteMsgAfter: true,
			AcceptFunc: func(msg *discordgo.Message) (err error) {
				err = db.SetGuildModLog(ctx.Event.GuildID, ctx.Event.ChannelID)
				if err != nil {
					return
				}
				err = ctx.FollowUpEmbed(&discordgo.MessageEmbed{
					Description: "Set this channel as modlog channel.",
				}).Error
				return
			},
		}

		if _, err = acceptMsg.AsFollowUp(ctx.Ctx); err != nil {
			return
		}
		return acceptMsg.Error()
	}

	ch := chV.ChannelValue(ctx.Ctx)
	if ch.Type != discordgo.ChannelTypeGuildText {
		return ctx.FollowUpError("Specified channel is not a text channel.", "").Error
	}

	if err = db.SetGuildModLog(ctx.Event.GuildID, ch.ID); err != nil {
		return
	}

	err = ctx.FollowUpEmbed(&discordgo.MessageEmbed{
		Description: fmt.Sprintf("Set channel <#%s> as modlog channel.", ch.ID),
	}).Error

	return
}

func (c *Modlog) disable(ctx *ken.SubCommandCtx) (err error) {
	db := ctx.Get(static.DiDatabase).(database.Database)

	if err = db.SetGuildModLog(ctx.Event.GuildID, ""); err != nil {
		return
	}

	return ctx.FollowUpEmbed(&discordgo.MessageEmbed{
		Description: "Modloging disabled.",
	}).Error
}
