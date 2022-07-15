package slashcommands

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/services/permissions"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/discordutil"
	"github.com/zekrotja/ken"
)

type Say struct{}

var (
	_ ken.SlashCommand        = (*Say)(nil)
	_ permissions.PermCommand = (*Say)(nil)
)

func (c *Say) Name() string {
	return "say"
}

func (c *Say) Description() string {
	return "Send an embedded message with the bot."
}

func (c *Say) Version() string {
	return "1.1.0"
}

func (c *Say) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (c *Say) Options() []*discordgo.ApplicationCommandOption {
	commonOpts := []*discordgo.ApplicationCommandOption{
		{
			Type:         discordgo.ApplicationCommandOptionChannel,
			Name:         "channel",
			Description:  "The channel to send the message into (or to edit a message in).",
			ChannelTypes: []discordgo.ChannelType{discordgo.ChannelTypeGuildText},
		},
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "editmessage",
			Description: "The ID of the message to be edited.",
		},
	}

	return []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "embed",
			Description: "Send an embed message.",
			Options: append([]*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "message",
					Description: "The message content.",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "color",
					Description: "The color.",
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{Name: "default", Value: static.ColorEmbedDefault},
						{Name: "cyan", Value: static.ColorEmbedCyan},
						{Name: "red", Value: static.ColorEmbedError},
						{Name: "gray", Value: static.ColorEmbedGray},
						{Name: "green", Value: static.ColorEmbedGreen},
						{Name: "lime", Value: static.ColorEmbedUpdated},
						{Name: "orange", Value: static.ColorEmbedOrange},
						{Name: "violett", Value: static.ColorEmbedViolett},
						{Name: "yellow", Value: static.ColorEmbedYellow},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "title",
					Description: "The title content.",
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "footer",
					Description: "The footer content.",
				},
			}, commonOpts...),
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "raw",
			Description: "Send raw embed message.",
			Options: append([]*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "json",
					Description: "The raw JSON data of the embed to be sent.",
					Required:    true,
				},
			}, commonOpts...),
		},
	}
}

func (c *Say) Domain() string {
	return "sp.chat.say"
}

func (c *Say) SubDomains() []permissions.SubPermission {
	return nil
}

func (c *Say) Run(ctx *ken.Ctx) (err error) {
	if err = ctx.Defer(); err != nil {
		return
	}

	err = ctx.HandleSubCommands(
		ken.SubCommandHandler{"embed", c.embed},
		ken.SubCommandHandler{"raw", c.raw},
	)

	return
}

func (c *Say) embed(ctx *ken.SubCommandCtx) (err error) {
	emb := &discordgo.MessageEmbed{
		Color: static.ColorEmbedDefault,
	}

	emb.Description = ctx.Options().GetByName("message").StringValue()

	if titleV, ok := ctx.Options().GetByNameOptional("title"); ok {
		emb.Title = titleV.StringValue()
	}
	if footerV, ok := ctx.Options().GetByNameOptional("footer"); ok {
		emb.Footer = &discordgo.MessageEmbedFooter{
			Text: footerV.StringValue(),
		}
	}
	if colorV, ok := ctx.Options().GetByNameOptional("color"); ok {
		emb.Color = int(colorV.IntValue())
	}

	err = c.sendMessage(ctx, emb)
	return
}

func (c *Say) raw(ctx *ken.SubCommandCtx) (err error) {
	raw := ctx.Options().GetByName("json").StringValue()
	var emb discordgo.MessageEmbed
	if err = json.Unmarshal([]byte(raw), &emb); err != nil {
		return ctx.FollowUpError(
			fmt.Sprintf("Failed parsing JSON data:\n```\n%s\n```",
				err.Error()), "").Error
	}
	err = c.sendMessage(ctx, &emb)
	return
}

func (c *Say) sendMessage(ctx *ken.SubCommandCtx, emb *discordgo.MessageEmbed) (err error) {
	emb.Author = &discordgo.MessageEmbedAuthor{
		Name:    ctx.User().String(),
		IconURL: ctx.User().AvatarURL("16x16"),
	}

	chanID := ctx.Event.ChannelID
	if chanV, ok := ctx.Options().GetByNameOptional("channel"); ok {
		ch := chanV.ChannelValue(ctx.Ctx)
		chanID = ch.ID
	}

	messageID := ""
	if msgV, ok := ctx.Options().GetByNameOptional("editmessage"); ok {
		messageID = msgV.StringValue()
	}

	var msg *discordgo.Message
	var status string
	if messageID != "" {
		msg, err = ctx.Session.ChannelMessageEditEmbed(chanID, messageID, emb)
		status = "edited"
	} else {
		msg, err = ctx.Session.ChannelMessageSendEmbed(chanID, emb)
		status = "sent"
	}
	if err != nil {
		return
	}

	fum := ctx.FollowUpEmbed(&discordgo.MessageEmbed{
		Description: fmt.Sprintf("Message has been %s. [Here](%s) you can find the message.",
			status, discordutil.GetMessageLink(msg, ctx.Event.GuildID)),
	})

	if chanID == ctx.Event.ChannelID {
		fum.DeleteAfter(5 * time.Second)
	}

	return fum.Error
}
