package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"

	"../util"
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

func (c *CmdMvall) GetPermission() int {
	return 5
}

func (c *CmdMvall) Exec(args *CommandArgs) error {
	args.Session.ChannelMessageDelete(args.Channel.ID, args.Message.ID)

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

	toVC, err := util.FetchChannel(args.Session, args.Guild.ID, strings.Join(args.Args, " "))
	if err != nil {
		msg, err := util.SendEmbedError(args.Session, args.Channel.ID,
			"The channel you have passed could not be found.")
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
