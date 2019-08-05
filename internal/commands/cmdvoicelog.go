package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/util"
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
	return GroupGuildConfig
}

func (c *CmdVoicelog) GetDomainName() string {
	return "sp.guild.config.voicelog"
}

func (c *CmdVoicelog) Exec(args *CommandArgs) error {
	if len(args.Args) < 1 {
		acceptMsg := &util.AcceptMessage{
			Session: args.Session,
			Embed: &discordgo.MessageEmbed{
				Color:       util.ColorEmbedDefault,
				Description: "Do you want to set this channel as voicelog channel?",
			},
			UserID:         args.User.ID,
			DeleteMsgAfter: true,
			AcceptFunc: func(msg *discordgo.Message) {
				err := args.CmdHandler.db.SetGuildVoiceLog(args.Guild.ID, args.Channel.ID)
				if err != nil {
					util.SendEmbedError(args.Session, args.Channel.ID,
						"Failed setting voicelog channel: ```\n"+err.Error()+"\n```")
				} else {
					msg, _ := util.SendEmbed(args.Session, args.Channel.ID,
						"Set this channel as voicelog channel.", "", util.ColorEmbedUpdated)
					util.DeleteMessageLater(args.Session, msg, 6*time.Second)
				}
			},
		}
		_, err := acceptMsg.Send(args.Channel.ID)
		return err
	}

	if strings.ToLower(args.Args[0]) == "reset" {
		err := args.CmdHandler.db.SetGuildVoiceLog(args.Guild.ID, "")
		if err != nil {
			msg, err := util.SendEmbedError(args.Session, args.Channel.ID,
				"Failed reseting voice log channel: ```\n"+err.Error()+"\n```")
			util.DeleteMessageLater(args.Session, msg, 15*time.Second)
			return err
		}
		msg, err := util.SendEmbed(args.Session, args.Channel.ID,
			"Voicelog channel reset.", "", util.ColorEmbedUpdated)
		util.DeleteMessageLater(args.Session, msg, 5*time.Second)
		return err
	}

	mlChan, err := util.FetchChannel(args.Session, args.Guild.ID, args.Args[0], func(c *discordgo.Channel) bool {
		return c.Type == discordgo.ChannelTypeGuildText
	})
	if err != nil {
		msg, err := util.SendEmbedError(args.Session, args.Channel.ID,
			"Could not find any channel on this guild passing this resolvable.")
		util.DeleteMessageLater(args.Session, msg, 6*time.Second)
		return err
	}
	err = args.CmdHandler.db.SetGuildVoiceLog(args.Guild.ID, mlChan.ID)
	if err != nil {
		return err
	}
	msg, err := util.SendEmbed(args.Session, args.Channel.ID,
		fmt.Sprintf("Set <#%s> as voicelog channel.", mlChan.ID), "", util.ColorEmbedUpdated)
	util.DeleteMessageLater(args.Session, msg, 6*time.Second)
	return err
}
