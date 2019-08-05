package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/zekroTJA/shinpuru/internal/util"
)

type CmdMvall struct {
}

func (c *CmdMvall) GetInvokes() []string {
	return []string{"mvall", "mva"}
}

func (c *CmdMvall) GetDescription() string {
	return "move all members in your current voice channel into another one"
}

func (c *CmdMvall) GetHelp() string {
	return "`mvall <otherChanResolvable>`"
}

func (c *CmdMvall) GetGroup() string {
	return GroupModeration
}

func (c *CmdMvall) GetDomainName() string {
	return "sp.guild.mod.mvall"
}

func (c *CmdMvall) Exec(args *CommandArgs) error {
	if len(args.Args) < 1 {
		msg, err := util.SendEmbedError(args.Session, args.Channel.ID,
			"Please enter a voice channel as argument.")
		util.DeleteMessageLater(args.Session, msg, 10*time.Second)
		return err
	}

	var currVC string
	for _, vs := range args.Guild.VoiceStates {
		if vs.UserID == args.User.ID {
			currVC = vs.ChannelID
		}
	}

	if currVC == "" {
		msg, err := util.SendEmbedError(args.Session, args.Channel.ID,
			"You need to be in a voice channel to use this command.")
		util.DeleteMessageLater(args.Session, msg, 10*time.Second)
		return err
	}

	toVC, err := util.FetchChannel(args.Session, args.Guild.ID, strings.Join(args.Args, " "),
		func(c *discordgo.Channel) bool {
			return c.Type == discordgo.ChannelTypeGuildVoice
		})
	if err != nil {
		msg, err := util.SendEmbedError(args.Session, args.Channel.ID,
			"Could not find any voice channel passing the resolvable.")
		util.DeleteMessageLater(args.Session, msg, 10*time.Second)
		return err
	}

	if toVC.Type != discordgo.ChannelTypeGuildVoice {
		msg, err := util.SendEmbedError(args.Session, args.Channel.ID,
			fmt.Sprintf("The target channel *(`%s`)* is not a voice channel.", toVC.Name))
		util.DeleteMessageLater(args.Session, msg, 10*time.Second)
		return err
	}

	for _, vs := range args.Guild.VoiceStates {
		if vs.ChannelID == currVC {
			err := args.Session.GuildMemberMove(args.Guild.ID, vs.UserID, toVC.ID)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
