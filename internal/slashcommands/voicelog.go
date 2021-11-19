package slashcommands

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/services/permissions"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/acceptmsg"
	"github.com/zekrotja/dgrs"
	"github.com/zekrotja/ken"
)

type Voicelog struct{}

var (
	_ ken.Command             = (*Voicelog)(nil)
	_ permissions.PermCommand = (*Voicelog)(nil)
)

func (c *Voicelog) Name() string {
	return "voicelog"
}

func (c *Voicelog) Description() string {
	return "Set the voice log channel for a guild."
}

func (c *Voicelog) Version() string {
	return "1.0.0"
}

func (c *Voicelog) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (c *Voicelog) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "set",
			Description: "Set this or a specified channel as voice log channel.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:         discordgo.ApplicationCommandOptionChannel,
					Name:         "channel",
					Description:  "A channel to be set as voice log (current channel if not specified).",
					ChannelTypes: []discordgo.ChannelType{discordgo.ChannelTypeGuildText},
				},
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "disable",
			Description: "Disable voicelog.",
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "ignore",
			Description: "Add a voice channel to the ignorelist.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:         discordgo.ApplicationCommandOptionChannel,
					Name:         "channel",
					Description:  "A voice channel to be ignored.",
					Required:     true,
					ChannelTypes: []discordgo.ChannelType{discordgo.ChannelTypeGuildVoice},
				},
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "unignore",
			Description: "Remove a voice channel from the ignorelist.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:         discordgo.ApplicationCommandOptionChannel,
					Name:         "channel",
					Description:  "A voice channel to be unset from the ignore list.",
					ChannelTypes: []discordgo.ChannelType{discordgo.ChannelTypeGuildVoice},
					Required:     true,
				},
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "ignorelist",
			Description: "Show all ignored voice channels.",
		},
	}
}

func (c *Voicelog) Domain() string {
	return "sp.guild.config.voicelog"
}

func (c *Voicelog) SubDomains() []permissions.SubPermission {
	return nil
}

func (c *Voicelog) Run(ctx *ken.Ctx) (err error) {
	if err = ctx.Defer(); err != nil {
		return
	}

	err = ctx.HandleSubCommands(
		ken.SubCommandHandler{"set", c.set},
		ken.SubCommandHandler{"disable", c.disable},
		ken.SubCommandHandler{"ignore", c.ignore},
		ken.SubCommandHandler{"unignore", c.unignore},
		ken.SubCommandHandler{"ignorelist", c.ignorelist},
	)

	return
}

func (c *Voicelog) set(ctx *ken.SubCommandCtx) (err error) {
	db := ctx.Get(static.DiDatabase).(database.Database)

	chV, ok := ctx.Options().GetByNameOptional("channel")

	if !ok {
		acceptMsg := &acceptmsg.AcceptMessage{
			Session: ctx.Session,
			Embed: &discordgo.MessageEmbed{
				Color:       static.ColorEmbedDefault,
				Description: "Do you want to set this channel as voicelog channel?",
			},
			UserID:         ctx.User().ID,
			DeleteMsgAfter: true,
			AcceptFunc: func(msg *discordgo.Message) (err error) {
				err = db.SetGuildVoiceLog(ctx.Event.GuildID, ctx.Event.ChannelID)
				if err != nil {
					return
				}
				err = ctx.FollowUpEmbed(&discordgo.MessageEmbed{
					Description: "Set this channel as voicelog channel.",
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

	if err = db.SetGuildVoiceLog(ctx.Event.GuildID, ch.ID); err != nil {
		return
	}

	err = ctx.FollowUpEmbed(&discordgo.MessageEmbed{
		Description: fmt.Sprintf("Set channel <#%s> as voicelog channel.", ch.ID),
	}).Error

	return
}

func (c *Voicelog) disable(ctx *ken.SubCommandCtx) (err error) {
	db := ctx.Get(static.DiDatabase).(database.Database)

	if err = db.SetGuildVoiceLog(ctx.Event.GuildID, ""); err != nil {
		return
	}

	return ctx.FollowUpEmbed(&discordgo.MessageEmbed{
		Description: "Voiceloging disabled.",
	}).Error
}

func (c *Voicelog) ignore(ctx *ken.SubCommandCtx) (err error) {
	db := ctx.Get(static.DiDatabase).(database.Database)

	ch := ctx.Options().GetByName("channel").ChannelValue(ctx.Ctx)

	if err = db.SetGuildVoiceLogIngore(ch.GuildID, ch.ID); err != nil {
		return err
	}

	return ctx.FollowUpEmbed(&discordgo.MessageEmbed{
		Description: fmt.Sprintf("Channel <#%s> is now on the ignore list.", ch.ID),
	}).Error
}

func (c *Voicelog) unignore(ctx *ken.SubCommandCtx) (err error) {
	db := ctx.Get(static.DiDatabase).(database.Database)

	ch := ctx.Options().GetByName("channel").ChannelValue(ctx.Ctx)

	if err = db.RemoveGuildVoiceLogIgnore(ch.GuildID, ch.ID); err != nil {
		return err
	}

	return ctx.FollowUpEmbed(&discordgo.MessageEmbed{
		Description: fmt.Sprintf("Channel <#%s> was removed from the ignore list.", ch.ID),
	}).Error
}

func (c *Voicelog) ignorelist(ctx *ken.SubCommandCtx) (err error) {
	db := ctx.Get(static.DiDatabase).(database.Database)
	st := ctx.Get(static.DiState).(*dgrs.State)

	vcIDs, err := db.GetGuildVoiceLogIgnores(ctx.Event.GuildID)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return err
	}
	vcs := make([]string, len(vcIDs))
	i := 0

	for _, id := range vcIDs {
		if c, err := st.Channel(id); err == nil && c != nil {
			vcs[i] = fmt.Sprintf("%s `%s`", c.Name, c.ID)
			i++
		}
	}

	return ctx.FollowUpEmbed(&discordgo.MessageEmbed{
		Description: strings.Join(vcs, "\n"),
		Title:       "Ignored Voice Channels",
	}).Error
}
