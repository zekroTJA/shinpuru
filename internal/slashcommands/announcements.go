package slashcommands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/services/permissions"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/stringutil"
	"github.com/zekrotja/ken"
)

type announcementType string

const (
	announcementTypeJoin  = "join"
	announcementTypeLeave = "leave"
)

type Announcements struct{}

var (
	_ ken.SlashCommand        = (*Announcements)(nil)
	_ permissions.PermCommand = (*Announcements)(nil)
)

func (c *Announcements) Name() string {
	return "announcements"
}

func (c *Announcements) Description() string {
	return "Set a message which will show up when a user joins or leaves the guild."
}

func (c *Announcements) Version() string {
	return "1.0.0"
}

func (c *Announcements) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (c *Announcements) Options() []*discordgo.ApplicationCommandOption {
	commonOpts := []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "type",
			Description: "The announcement type.",
			Required:    true,
			Choices: []*discordgo.ApplicationCommandOptionChoice{
				{Name: announcementTypeJoin, Value: announcementTypeJoin},
				{Name: announcementTypeLeave, Value: announcementTypeLeave},
			},
		},
	}

	return []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "set",
			Description: "Set a message and channel for an announcement message.",
			Options: append(commonOpts, []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "message",
					Description: "The message. [user] will be replaced with the username and [ment] with the mention.",
				},
				{
					Type:         discordgo.ApplicationCommandOptionChannel,
					Name:         "channel",
					Description:  "A channel to be set.",
					ChannelTypes: []discordgo.ChannelType{discordgo.ChannelTypeGuildText},
				},
			}...),
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "disable",
			Description: "Disable announcements.",
			Options:     commonOpts,
		},
	}
}

func (c *Announcements) Domain() string {
	return "sp.guild.config.announcements"
}

func (c *Announcements) SubDomains() []permissions.SubPermission {
	return nil
}

func (c *Announcements) Run(ctx *ken.Ctx) (err error) {
	if err = ctx.Defer(); err != nil {
		return
	}

	err = ctx.HandleSubCommands(
		ken.SubCommandHandler{"set", c.set},
		ken.SubCommandHandler{"disable", c.disable},
	)

	return
}

func (c *Announcements) set(ctx *ken.SubCommandCtx) (err error) {
	db := ctx.Get(static.DiDatabase).(database.Database)

	typ := announcementType(ctx.Options().GetByName("type").StringValue())

	var currChanID, currMsg string
	if typ == announcementTypeJoin {
		currChanID, currMsg, err = db.GetGuildJoinMsg(ctx.Event.GuildID)
	} else if typ == announcementTypeLeave {
		currChanID, currMsg, err = db.GetGuildLeaveMsg(ctx.Event.GuildID)
	}
	if err != nil {
		return
	}

	if chV, ok := ctx.Options().GetByNameOptional("channel"); ok {
		ch := chV.ChannelValue(ctx.Ctx)
		currChanID = ch.ID
	}
	if msgV, ok := ctx.Options().GetByNameOptional("message"); ok {
		currMsg = msgV.StringValue()
	}

	if typ == announcementTypeJoin {
		err = db.SetGuildJoinMsg(ctx.Event.GuildID, currChanID, currMsg)
	} else if typ == announcementTypeLeave {
		err = db.SetGuildLeaveMsg(ctx.Event.GuildID, currChanID, currMsg)
	}
	if err != nil {
		return
	}

	if currChanID == "" {
		err = ctx.FollowUpEmbed(&discordgo.MessageEmbed{
			Description: fmt.Sprintf(
				"Set %s message to\n```\n%s\n```."+
					"%s messages are still disabled because no channel is set.",
				typ, currMsg, stringutil.Capitalize(string(typ), false)),
			Color: static.ColorEmbedOrange,
		}).Error
	} else if currMsg == "" {
		err = ctx.FollowUpEmbed(&discordgo.MessageEmbed{
			Description: fmt.Sprintf(
				"Set %s message channel to <#%s>."+
					"%s messages are still disabled because no message is set.",
				typ, currMsg, stringutil.Capitalize(string(typ), false)),
			Color: static.ColorEmbedOrange,
		}).Error
	} else {
		err = ctx.FollowUpEmbed(&discordgo.MessageEmbed{
			Description: fmt.Sprintf(
				"Set %s message channel to <#%s> and %s message to\n```\n%s\n```"+
					"%s messages are now enabled.",
				typ, currChanID, typ, currMsg, stringutil.Capitalize(string(typ), false)),
			Color: static.ColorEmbedGreen,
		}).Error
	}

	return
}

func (c *Announcements) disable(ctx *ken.SubCommandCtx) (err error) {
	db := ctx.Get(static.DiDatabase).(database.Database)

	typ := announcementType(ctx.Options().GetByName("type").StringValue())

	if typ == announcementTypeJoin {
		err = db.SetGuildJoinMsg(ctx.Event.GuildID, "", "")
	} else if typ == announcementTypeLeave {
		err = db.SetGuildLeaveMsg(ctx.Event.GuildID, "", "")
	}
	if err != nil {
		return
	}

	return ctx.FollowUpEmbed(&discordgo.MessageEmbed{
		Description: fmt.Sprintf("%s disabled.", stringutil.Capitalize(string(typ), false)),
	}).Error
}
