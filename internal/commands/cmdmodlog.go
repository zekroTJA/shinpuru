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

type CmdModlog struct {
}

func (c *CmdModlog) GetInvokes() []string {
	return []string{"modlog", "setmodlog", "modlogchan", "ml"}
}

func (c *CmdModlog) GetDescription() string {
	return "set the mod log channel for a guild"
}

func (c *CmdModlog) GetHelp() string {
	return "`modlog` - set this channel as modlog channel\n" +
		"`modlog <chanResolvable>` - set any text channel as mod log channel\n" +
		"`modlog reset` - reset mod log channel"
}

func (c *CmdModlog) GetGroup() string {
	return shireikan.GroupGuildConfig
}

func (c *CmdModlog) GetDomainName() string {
	return "sp.guild.config.modlog"
}

func (c *CmdModlog) GetSubPermissionRules() []shireikan.SubPermission {
	return nil
}

func (c *CmdModlog) IsExecutableInDMChannels() bool {
	return false
}

func (c *CmdModlog) Exec(ctx shireikan.Context) error {
	db, _ := ctx.GetObject("db").(database.Database)

	if len(ctx.GetArgs()) < 1 {
		acceptMsg := &acceptmsg.AcceptMessage{
			Session: ctx.GetSession(),
			Embed: &discordgo.MessageEmbed{
				Color:       static.ColorEmbedDefault,
				Description: "Do you want to set this channel as modlog channel?",
			},
			UserID:         ctx.GetUser().ID,
			DeleteMsgAfter: true,
			AcceptFunc: func(msg *discordgo.Message) {
				err := db.SetGuildModLog(ctx.GetGuild().ID, ctx.GetChannel().ID)
				if err != nil {
					util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
						"Failed setting modlog channel: ```\n"+err.Error()+"\n```")
				} else {
					util.SendEmbed(ctx.GetSession(), ctx.GetChannel().ID,
						"Set this channel as modlog channel.", "", static.ColorEmbedUpdated).
						DeleteAfter(6 * time.Second)
				}
			},
		}
		_, err := acceptMsg.Send(ctx.GetChannel().ID)
		return err
	}

	if strings.ToLower(ctx.GetArgs().Get(0).AsString()) == "reset" {
		err := db.SetGuildModLog(ctx.GetGuild().ID, "")
		if err != nil {
			return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
				"Failed reseting mod log channel: ```\n"+err.Error()+"\n```").
				DeleteAfter(15 * time.Second).Error()
		}
		return util.SendEmbed(ctx.GetSession(), ctx.GetChannel().ID,
			"Modlog channel reset.", "", static.ColorEmbedUpdated).
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
	err = db.SetGuildModLog(ctx.GetGuild().ID, mlChan.ID)
	if err != nil {
		return err
	}
	return util.SendEmbed(ctx.GetSession(), ctx.GetChannel().ID,
		fmt.Sprintf("Set <#%s> as modlog channel.", mlChan.ID), "", static.ColorEmbedUpdated).
		DeleteAfter(8 * time.Second).Error()
}
