package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/acceptmsg"
	"github.com/zekroTJA/shinpuru/pkg/discordutil"
	"github.com/zekroTJA/shinpuru/pkg/fetch"
	"github.com/zekroTJA/shireikan"
)

type CmdVoicelog struct {
}

func (c *CmdVoicelog) GetInvokes() []string {
	return []string{"voicelog", "setvoicelog", "voicelogchan", "vl"}
}

func (c *CmdVoicelog) GetDescription() string {
	return "Set the mod log channel for a guild."
}

func (c *CmdVoicelog) GetHelp() string {
	return "`voicelog` - set this channel as voicelog channel\n" +
		"`voicelog <chanResolvable>` - set any text channel as voicelog channel\n" +
		"`voicelog reset` - reset voice log channel\n" +
		"`voicelog ignore <chanResolvable>` - add voice channel to ignore list\n" +
		"`voicelog unignore <chanResolvable> - removes a voice channel from the ignore list\n`" +
		"`voicelog ignorelist` - display ignored voice channels"
}

func (c *CmdVoicelog) GetGroup() string {
	return shireikan.GroupGuildConfig
}

func (c *CmdVoicelog) GetDomainName() string {
	return "sp.guild.config.voicelog"
}

func (c *CmdVoicelog) GetSubPermissionRules() []shireikan.SubPermission {
	return nil
}

func (c *CmdVoicelog) IsExecutableInDMChannels() bool {
	return false
}

func (c *CmdVoicelog) Exec(ctx shireikan.Context) error {
	db, _ := ctx.GetObject(static.DiDatabase).(database.Database)

	switch ctx.GetArgs().Get(0).AsString() {

	case "":
		acceptMsg := &acceptmsg.AcceptMessage{
			Session: ctx.GetSession(),
			Embed: &discordgo.MessageEmbed{
				Color:       static.ColorEmbedDefault,
				Description: "Do you want to set this channel as voicelog channel?",
			},
			UserID:         ctx.GetUser().ID,
			DeleteMsgAfter: true,
			AcceptFunc: func(msg *discordgo.Message) {
				err := db.SetGuildVoiceLog(ctx.GetGuild().ID, ctx.GetChannel().ID)
				if err != nil {
					util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
						"Failed setting voicelog channel: ```\n"+err.Error()+"\n```")
				} else {
					util.SendEmbed(ctx.GetSession(), ctx.GetChannel().ID,
						"Set this channel as voicelog channel.", "", static.ColorEmbedUpdated).
						DeleteAfter(8 * time.Second)
				}
			},
		}
		_, err := acceptMsg.Send(ctx.GetChannel().ID)
		return err

	case "reset":
		err := db.SetGuildVoiceLog(ctx.GetGuild().ID, "")
		if err != nil {
			return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
				"Failed reseting voice log channel: ```\n"+err.Error()+"\n```").
				DeleteAfter(15 * time.Second).Error()
		}
		return util.SendEmbed(ctx.GetSession(), ctx.GetChannel().ID,
			"Voicelog channel reset.", "", static.ColorEmbedUpdated).
			DeleteAfter(8 * time.Second).Error()

	case "ignore", "block", "hide":
		vChan, err := c.getVoiceChan(ctx)
		if err != nil {
			return err
		}
		if err = db.SetGuildVoiceLogIngore(vChan.GuildID, vChan.ID); err != nil {
			return err
		}
		return util.SendEmbed(ctx.GetSession(), ctx.GetChannel().ID,
			"Voice channel `"+vChan.Name+"` was added to the ignore list.", "", static.ColorEmbedUpdated).
			DeleteAfter(8 * time.Second).Error()

	case "unignore", "unblock", "unhide", "show":
		vChan, err := c.getVoiceChan(ctx)
		if err != nil {
			return err
		}
		if err = db.RemoveGuildVoiceLogIgnore(vChan.GuildID, vChan.ID); err != nil {
			return err
		}
		return util.SendEmbed(ctx.GetSession(), ctx.GetChannel().ID,
			"Voice channel `"+vChan.Name+"` was removed from the ignore list.", "", static.ColorEmbedUpdated).
			DeleteAfter(8 * time.Second).Error()

	case "ignored", "ignorelist", "blocklist":
		vcIDs, err := db.GetGuildVoiceLogIgnores(ctx.GetGuild().ID)
		if err != nil && !database.IsErrDatabaseNotFound(err) {
			return err
		}
		vcs := make([]string, len(vcIDs))
		i := 0
		for _, id := range vcIDs {
			if c, err := discordutil.GetChannel(ctx.GetSession(), id); err == nil && c != nil {
				vcs[i] = fmt.Sprintf("%s `%s`", c.Name, c.ID)
				i++
			}
		}
		return util.SendEmbed(ctx.GetSession(), ctx.GetChannel().ID,
			strings.Join(vcs, "\n"),
			"Ignored Voice Channels", static.ColorEmbedDefault).
			Error()

	default:
		mlChan, err := fetch.FetchChannel(ctx.GetSession(), ctx.GetGuild().ID, ctx.GetArgs().Get(0).AsString(), func(c *discordgo.Channel) bool {
			return c.Type == discordgo.ChannelTypeGuildText
		})
		if err != nil {
			return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
				"Could not find any channel on this guild passing this resolvable.").
				DeleteAfter(8 * time.Second).Error()
		}
		err = db.SetGuildVoiceLog(ctx.GetGuild().ID, mlChan.ID)
		if err != nil {
			return err
		}

		return util.SendEmbed(ctx.GetSession(), ctx.GetChannel().ID,
			fmt.Sprintf("Set <#%s> as voicelog channel.", mlChan.ID), "", static.ColorEmbedUpdated).
			DeleteAfter(8 * time.Second).Error()
	}

}

func (c *CmdVoicelog) getVoiceChan(ctx shireikan.Context) (vc *discordgo.Channel, err error) {
	chanResolver := ctx.GetArgs().Get(1).AsString()
	if chanResolver == "" {
		return nil, util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
			"Pelase pass a voice channel resover as argument which channel you want to add to the ignore list.").
			DeleteAfter(15 * time.Second).Error()
	}
	vChan, err := fetch.FetchChannel(ctx.GetSession(), ctx.GetGuild().ID, chanResolver, func(c *discordgo.Channel) bool {
		return c.Type == discordgo.ChannelTypeGuildVoice
	})
	if err == fetch.ErrNotFound {
		return nil, util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
			"Could not find any voice channel on this guild with the given resolvable.").
			DeleteAfter(15 * time.Second).Error()
	}
	return vChan, nil
}
