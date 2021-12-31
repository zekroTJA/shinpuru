package slashcommands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/models"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/services/permissions"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekrotja/ken"
)

type Starboard struct{}

var (
	_ ken.SlashCommand        = (*Starboard)(nil)
	_ permissions.PermCommand = (*Starboard)(nil)
)

func (c *Starboard) Name() string {
	return "starboard"
}

func (c *Starboard) Description() string {
	return "Set guild starboard settings."
}

func (c *Starboard) Version() string {
	return "1.0.0"
}

func (c *Starboard) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (c *Starboard) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "set",
			Description: "Set starboard settings.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:         discordgo.ApplicationCommandOptionChannel,
					Name:         "channel",
					Description:  "The channel where the starboard messages will appear.",
					ChannelTypes: []discordgo.ChannelType{discordgo.ChannelTypeGuildText},
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "threshold",
					Description: "The minimum number of emote votes until a message gets into the starboard.",
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "emote",
					Description: "The name or emote of the emote to be used for staring messages.",
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "karma",
					Description: "The amount of karma gain when a users message gets into the starboard.",
				},
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "disable",
			Description: "Disable the starboard.",
		},
	}
}

func (c *Starboard) Domain() string {
	return "sp.guild.config.starboard"
}

func (c *Starboard) SubDomains() []permissions.SubPermission {
	return nil
}

func (c *Starboard) Run(ctx *ken.Ctx) (err error) {
	if err = ctx.Defer(); err != nil {
		return
	}

	err = ctx.HandleSubCommands(
		ken.SubCommandHandler{"set", c.set},
		ken.SubCommandHandler{"disable", c.disable},
	)

	return
}

func (c *Starboard) set(ctx *ken.SubCommandCtx) (err error) {
	starboardConfig, err := c.getConfig(ctx)

	if v, ok := ctx.Options().GetByNameOptional("channel"); ok {
		ch := v.ChannelValue(ctx.Ctx)
		starboardConfig.ChannelID = ch.ID
	}
	if v, ok := ctx.Options().GetByNameOptional("threshold"); ok {
		starboardConfig.Threshold = int(v.IntValue())
		if starboardConfig.Threshold < 0 {
			return ctx.FollowUpError("Threshold value must be equal or larger than `0`.", "").Error
		}
	}
	if v, ok := ctx.Options().GetByNameOptional("emote"); ok {
		starboardConfig.EmojiID = v.StringValue()
	}
	if v, ok := ctx.Options().GetByNameOptional("karma"); ok {
		starboardConfig.KarmaGain = int(v.IntValue())
		if starboardConfig.KarmaGain < 0 {
			return ctx.FollowUpError("Threshold value must be equal or larger than `0`.", "").Error
		}
	}

	return c.setConfig(ctx, starboardConfig, false)
}

func (c *Starboard) disable(ctx *ken.SubCommandCtx) (err error) {
	starboardConfig, err := c.getConfig(ctx)

	starboardConfig.ChannelID = ""

	return c.setConfig(ctx, starboardConfig, false)
}

func (c *Starboard) getConfig(ctx *ken.SubCommandCtx) (starboardConfig *models.StarboardConfig, err error) {
	db := ctx.Get(static.DiDatabase).(database.Database)

	starboardConfig, err = db.GetStarboardConfig(ctx.Event.GuildID)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return
	}

	if database.IsErrDatabaseNotFound(err) || starboardConfig == nil {
		starboardConfig = &models.StarboardConfig{
			Threshold: 5,
			EmojiID:   "⭐",
			KarmaGain: 3,
		}
	}
	starboardConfig.GuildID = ctx.Event.GuildID

	return
}

func (c *Starboard) setConfig(
	ctx *ken.SubCommandCtx,
	cfg *models.StarboardConfig,
	outputError bool,
) (err error) {
	db := ctx.Get(static.DiDatabase).(database.Database)

	if err = db.SetStarboardConfig(cfg); err != nil {
		if outputError {
			return ctx.FollowUpError(
				"Failed settings starboard config to database: ```"+err.Error()+"```", "").
				Error
		}
		return err
	}

	msg := fmt.Sprintf(
		"Set starboard config:\n\nChannel: <#%s>\nThreshold: `%d`\nEmote: %s\nKarma Gain: `%d`",
		cfg.ChannelID, cfg.Threshold, cfg.EmojiID, cfg.KarmaGain)
	if cfg.ChannelID == "" {
		msg = "Starboard disabled. Set a channel as starboard channel to enable the starboard."
	}

	return ctx.FollowUpEmbed(&discordgo.MessageEmbed{
		Description: msg,
	}).Error
}
