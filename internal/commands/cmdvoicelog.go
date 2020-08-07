package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/core/database"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/acceptmsg"
	"github.com/zekroTJA/shinpuru/pkg/fetch"
	"github.com/zekroTJA/shireikan"
)

type CmdVoicelog struct {
}

func (c *CmdVoicelog) GetInvokes() []string {
	return []string{"voicelog", "setvoicelog", "voicelogchan", "vl"}
}

func (c *CmdVoicelog) GetDescription() string {
	return "set the mod log channel for a guild"
}

func (c *CmdVoicelog) GetHelp() string {
	return "`voicelog` - set this channel as voicelog channel\n" +
		"`voicelog <chanResolvable>` - set any text channel as voicelog channel\n" +
		"`voicelog reset` - reset voice log channel"
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
	db, _ := ctx.GetObject("dbtnw").(database.Database)

	if len(ctx.GetArgs()) < 1 {
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
	}

	if strings.ToLower(ctx.GetArgs().Get(0).AsString()) == "reset" {
		err := db.SetGuildVoiceLog(ctx.GetGuild().ID, "")
		if err != nil {
			return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
				"Failed reseting voice log channel: ```\n"+err.Error()+"\n```").
				DeleteAfter(15 * time.Second).Error()
		}
		return util.SendEmbed(ctx.GetSession(), ctx.GetChannel().ID,
			"Voicelog channel reset.", "", static.ColorEmbedUpdated).
			DeleteAfter(8 * time.Second).Error()
	}

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
