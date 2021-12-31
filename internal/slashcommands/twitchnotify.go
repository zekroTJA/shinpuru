package slashcommands

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/services/permissions"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/twitchnotify"
	"github.com/zekrotja/ken"
)

type Twitchnotify struct{}

var (
	_ ken.SlashCommand        = (*Twitchnotify)(nil)
	_ permissions.PermCommand = (*Twitchnotify)(nil)
)

func (c *Twitchnotify) Name() string {
	return "twitchnotify"
}

func (c *Twitchnotify) Description() string {
	return "Get notifications in channels when someone goes live on Twitch."
}

func (c *Twitchnotify) Version() string {
	return "1.0.0"
}

func (c *Twitchnotify) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (c *Twitchnotify) Options() []*discordgo.ApplicationCommandOption {
	commonOpts := []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "twitchname",
			Description: "The username of the twitch user.",
			Required:    true,
		},
	}
	return []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "list",
			Description: "List alöl registered notifies for the guild.",
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "add",
			Description: "Add a twitch user to be watched.",
			Options: append(commonOpts, []*discordgo.ApplicationCommandOption{
				{
					Type:         discordgo.ApplicationCommandOptionChannel,
					Name:         "channel",
					Description:  "The channel where the notifications are sent into (defaultly current channel).",
					ChannelTypes: []discordgo.ChannelType{discordgo.ChannelTypeGuildText},
				},
			}...),
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "remove",
			Description: "Remove a twitch user from the watch list.",
			Options:     commonOpts,
		},
	}
}

func (c *Twitchnotify) Domain() string {
	return "sp.chat.twitch"
}

func (c *Twitchnotify) SubDomains() []permissions.SubPermission {
	return nil
}

func (c *Twitchnotify) Run(ctx *ken.Ctx) (err error) {
	tnw := ctx.Get(static.DiTwitchNotifyWorker).(*twitchnotify.NotifyWorker)
	if tnw == nil {
		return ctx.Respond(&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Color:       static.ColorEmbedError,
						Description: "This feature is currently disabled.",
					},
				},
			},
		})
	}

	if err = ctx.Defer(); err != nil {
		return
	}

	err = ctx.HandleSubCommands(
		ken.SubCommandHandler{"list", c.list},
		ken.SubCommandHandler{"add", c.add},
		ken.SubCommandHandler{"remove", c.remove},
	)

	return
}

func (c *Twitchnotify) list(ctx *ken.SubCommandCtx) (err error) {
	tnw := ctx.Get(static.DiTwitchNotifyWorker).(*twitchnotify.NotifyWorker)
	db := ctx.Get(static.DiDatabase).(database.Database)

	nots, err := db.GetAllTwitchNotifies("")
	if err != nil {
		return err
	}

	var notsStr strings.Builder

	for _, not := range nots {
		if not.GuildID == ctx.Event.GuildID {
			if tUser, err := tnw.GetUser(not.TwitchUserID, twitchnotify.IdentID); err == nil {
				fmt.Fprintf(&notsStr, ":white_small_square:  **%s** in <#%s>\n",
					tUser.DisplayName, not.ChannelID)
			}
		}
	}

	return ctx.FollowUpEmbed(&discordgo.MessageEmbed{
		Title:       "Watched Twitch Channels",
		Description: notsStr.String(),
	}).Error
}

func (c *Twitchnotify) add(ctx *ken.SubCommandCtx) (err error) {
	tnw := ctx.Get(static.DiTwitchNotifyWorker).(*twitchnotify.NotifyWorker)
	db := ctx.Get(static.DiDatabase).(database.Database)

	twitchname := ctx.Options().GetByName("twitchname").StringValue()

	channelID := ctx.Event.ChannelID
	if channelV, ok := ctx.Options().GetByNameOptional("channel"); ok {
		ch := channelV.ChannelValue(ctx.Ctx)
		channelID = ch.ID
	}

	twitchuser, err := tnw.GetUser(twitchname, twitchnotify.IdentLogin)
	if err != nil {
		if err.Error() == "not found" {
			return ctx.FollowUpError("Twitch user with this name could not be found.", "").Error
		}
		return
	}

	err = tnw.AddUser(twitchuser)
	if err != nil {
		err = ctx.FollowUpError("Maximum count of registered Twitch accounts has been reached.", "").Error
		return
	}

	err = db.SetTwitchNotify(&twitchnotify.DBEntry{
		ChannelID:    channelID,
		GuildID:      ctx.Event.GuildID,
		TwitchUserID: twitchuser.ID,
	})
	if err != nil {
		return
	}

	return ctx.FollowUpEmbed(&discordgo.MessageEmbed{
		Description: fmt.Sprintf("You will now get a notification in channel <#%s> when `%s` goes live on twitch!",
			channelID, twitchuser.DisplayName),
	}).Error
}

func (c *Twitchnotify) remove(ctx *ken.SubCommandCtx) (err error) {
	tnw := ctx.Get(static.DiTwitchNotifyWorker).(*twitchnotify.NotifyWorker)
	db := ctx.Get(static.DiDatabase).(database.Database)

	twitchname := ctx.Options().GetByName("twitchname").StringValue()

	twitchuser, err := tnw.GetUser(twitchname, twitchnotify.IdentLogin)
	if err != nil {
		if err.Error() == "not found" {
			return ctx.FollowUpError("Twitch user with this name could not be found.", "").Error
		}
		return
	}

	nots, err := db.GetAllTwitchNotifies(twitchuser.ID)
	if err != nil {
		return err
	}

	var notify *twitchnotify.DBEntry
	for _, not := range nots {
		if not.GuildID == ctx.Event.GuildID {
			notify = not
		}
	}

	if notify == nil {
		return ctx.FollowUpError("Twitch user was nto set to be monitored on this guild.", "").Error
	}

	err = db.DeleteTwitchNotify(notify.TwitchUserID, notify.GuildID)
	if err != nil {
		return err
	}

	return ctx.FollowUpEmbed(&discordgo.MessageEmbed{
		Description: fmt.Sprintf("Notifications for twitch user `%s` in channel <#%s> have been removed.",
			twitchuser.DisplayName, notify.ChannelID),
	}).Error
}
