package commands

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/core/database"
	"github.com/zekroTJA/shinpuru/internal/shared/models"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/pkg/acceptmsg"
	"github.com/zekroTJA/shinpuru/pkg/fetch"
	"github.com/zekroTJA/shireikan"
)

type CmdStarboard struct {
	PermLvl int
}

func (c *CmdStarboard) GetInvokes() []string {
	return []string{"starboard", "star", "stb"}
}

func (c *CmdStarboard) GetDescription() string {
	return "Set guild starboard settings."
}

func (c *CmdStarboard) GetHelp() string {
	return "`starboard channel (<channelResolvable>)` - define a starboard channel\n" +
		"`starboard threshold <int>` - define a threshold for reaction count\n" +
		"`starboard emote <emoteName>` - define an emote to be used as starboard reaction\n" +
		"`starboard karma <int>` - define the amount of karma gained\n" +
		"`starboard disable` - disable starboard"
}

func (c *CmdStarboard) GetGroup() string {
	return shireikan.GroupGuildConfig
}

func (c *CmdStarboard) GetDomainName() string {
	return "sp.guild.config.stats"
}

func (c *CmdStarboard) GetSubPermissionRules() []shireikan.SubPermission {
	return nil
}

func (c *CmdStarboard) IsExecutableInDMChannels() bool {
	return false
}

func (c *CmdStarboard) Exec(ctx shireikan.Context) (err error) {

	db := ctx.GetObject("db").(database.Database)

	starboardConfig, err := db.GetStarboardConfig(ctx.GetGuild().ID)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return
	}

	if database.IsErrDatabaseNotFound(err) || starboardConfig == nil {
		starboardConfig = &models.StarboardConfig{
			Threshold: 5,
			EmojiID:   "⭐",
		}
	}
	starboardConfig.GuildID = ctx.GetGuild().ID

	switch ctx.GetArgs().Get(0).AsString() {

	case "channel", "chan", "ch":
		chanResolvable := ctx.GetArgs().Get(1).AsString()
		if chanResolvable == "" {
			_, err = acceptmsg.New().
				WithSession(ctx.GetSession()).
				WithContent("Do you want to set this channel as starboard channel?").
				DeleteAfterAnswer().
				DoOnAccept(func(m *discordgo.Message) {
					starboardConfig.ChannelID = ctx.GetChannel().ID
					c.setConfig(ctx, db, starboardConfig, true)
				}).
				Send(ctx.GetChannel().ID)
			return
		}
		ch, err := fetch.FetchChannel(ctx.GetSession(), ctx.GetGuild().ID, chanResolvable, func(c *discordgo.Channel) bool {
			return c.Type == discordgo.ChannelTypeGuildText || c.Type == discordgo.ChannelTypeGuildNews
		})
		if err == fetch.ErrNotFound {
			return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
				"Failed to find any channel by the specified resolvable.").
				DeleteAfter(15 * time.Second).
				Error()
		}
		starboardConfig.ChannelID = ch.ID

	case "threshold", "margin", "min":
		starboardConfig.Threshold, err = ctx.GetArgs().Get(1).AsInt()
		if err != nil || starboardConfig.Threshold < 1 {
			return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
				"Threshold must be a valid number larger than `0`.").
				DeleteAfter(15 * time.Second).
				Error()
		}

	case "emote", "emoji":
		emote := ctx.GetArgs().Get(1).AsString()
		if emote == "" {
			return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
				"Invalid arguments. Use `help joinmsg` to get help about how to use this command.").
				DeleteAfter(8 * time.Second).
				Error()
		}
		starboardConfig.EmojiID = emote

	case "karma", "karmagain":
		starboardConfig.KarmaGain, err = ctx.GetArgs().Get(1).AsInt()
		if err != nil || starboardConfig.KarmaGain < 0 || starboardConfig.KarmaGain > 100_000 {
			return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
				"Threshold must be a valid number larger or equal `0`.").
				DeleteAfter(15 * time.Second).
				Error()
		}

	case "disable", "reset", "unset", "off":
		starboardConfig.ChannelID = ""

	default:
		return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
			"Invalid arguments. Use `help joinmsg` to get help about how to use this command.").
			DeleteAfter(8 * time.Second).Error()
	}

	return c.setConfig(ctx, db, starboardConfig, false)
}

func (c *CmdStarboard) setConfig(
	ctx shireikan.Context,
	db database.Database,
	cfg *models.StarboardConfig,
	outputError bool,
) (err error) {
	if err = db.SetStarboardConfig(cfg); err != nil {
		if outputError {
			return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
				"Failed settings starboard config to database: ```"+err.Error()+"```").
				Error()
		}
		return err
	}

	msg := fmt.Sprintf(
		"Set starboard config:\n\nChannel: <#%s>\nThreshold: `%d`\nEmote: %s\nKarma Gain: `%d`",
		cfg.ChannelID, cfg.Threshold, cfg.EmojiID, cfg.KarmaGain)
	if cfg.ChannelID == "" {
		msg = "Starboard disabled. Set a channel as starboard channel to enable the starboard."
	}
	return util.SendEmbed(ctx.GetSession(), ctx.GetChannel().ID, msg, "", 0).
		Error()
}
