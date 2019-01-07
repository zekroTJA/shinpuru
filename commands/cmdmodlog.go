package commands

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/util"
)

type CmdModlog struct {
	PermLvl int
}

func (c *CmdModlog) GetInvokes() []string {
	return []string{"modlog", "setmodlog", "modlogchan", "ml"}
}

func (c *CmdModlog) GetDescription() string {
	return "set the mod log channel for a guild"
}

func (c *CmdModlog) GetHelp() string {
	return "`modlog` - set this channel as modlog channel\n" +
		"`modlog <chanResolvable>` - set any text channel as mod log channel"
}

func (c *CmdModlog) GetGroup() string {
	return GroupGuildConfig
}

func (c *CmdModlog) GetPermission() int {
	return c.PermLvl
}

func (c *CmdModlog) SetPermission(permLvl int) {
	c.PermLvl = permLvl
}

func (c *CmdModlog) Exec(args *CommandArgs) error {
	if len(args.Args) < 1 {
		acceptMsg := &util.AcceptMessage{
			Session: args.Session,
			Embed: &discordgo.MessageEmbed{
				Color:       util.ColorEmbedDefault,
				Description: "Do you want to set this channel as modlog channel?",
			},
			UserID:         args.User.ID,
			DeleteMsgAfter: true,
			AcceptFunc: func(msg *discordgo.Message) {
				err := args.CmdHandler.db.SetGuildModLog(args.Guild.ID, args.Channel.ID)
				if err != nil {
					util.SendEmbedError(args.Session, args.Channel.ID,
						"Failed setting modlog channel: ```\n"+err.Error()+"\n```")
				} else {
					msg, _ := util.SendEmbed(args.Session, args.Channel.ID,
						"Set this channel as modlog channel.", "", util.ColorEmbedUpdated)
					util.DeleteMessageLater(args.Session, msg, 6*time.Second)
				}
			},
		}
		_, err := acceptMsg.Send(args.Channel.ID)
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
	err = args.CmdHandler.db.SetGuildModLog(args.Guild.ID, mlChan.ID)
	if err != nil {
		return err
	}
	msg, err := util.SendEmbed(args.Session, args.Channel.ID,
		fmt.Sprintf("Set <#%s> as modlog channel.", mlChan.ID), "", util.ColorEmbedUpdated)
	util.DeleteMessageLater(args.Session, msg, 6*time.Second)
	return err
}
