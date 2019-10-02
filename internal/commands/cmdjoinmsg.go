package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/util"

	"github.com/zekroTJA/shinpuru/internal/core"
)

type CmdJoinMsg struct {
}

func (c *CmdJoinMsg) GetInvokes() []string {
	return []string{"joinmsg", "joinmessage"}
}

func (c *CmdJoinMsg) GetDescription() string {
	return "Set a message which will be sent into the defined channel when a member joins."
}

func (c *CmdJoinMsg) GetHelp() string {
	return "`joinmsg msg <message>` - Set the message of the join message.\n" +
		"`joinmsg channel <ChannelIdentifier>` - Set the channel where the message will be sent into.\n" +
		"`joinmsg reset` - Reset and disable join messages.\n\n" +
		"`[user]` will be replaced with the user name and `[ment]` will be replaced with the users mention when used in message text."
}

func (c *CmdJoinMsg) GetGroup() string {
	return GroupGuildConfig
}

func (c *CmdJoinMsg) GetDomainName() string {
	return "sp.guild.config.joinmsg"
}

func (c *CmdJoinMsg) Exec(args *CommandArgs) error {
	db := args.CmdHandler.db

	chanID, msg, err := db.GetGuildJoinMsg(args.Guild.ID)
	if err != nil && err != core.ErrDatabaseNotFound {
		return err
	}

	var resTxt string

	if len(args.Args) < 1 {

		if msg == "" && chanID == "" {
			resTxt = "*Join message and channel not set.*"
		} else if msg == "" {
			resTxt = "*Join message not set.*"
		} else if chanID == "" {
			resTxt = "*Join channel not set.*"
		} else {
			ch, err := args.Session.Channel(chanID)
			if ch == nil || err != nil || ch.GuildID != args.Guild.ID {
				resTxt = "*Join channel not set.*"
			} else {
				resTxt = fmt.Sprintf("```\n%s\n``` is set as join message and will be posted "+
					"into channel %s (%s).", msg, ch.Mention(), ch.ID)
			}
		}

		rmsg, err := util.SendEmbed(args.Session, args.Channel.ID,
			resTxt, "", 0)
		util.DeleteMessageLater(args.Session, rmsg, 10*time.Second)
		return err
	}

	argsJoined := strings.Join(args.Args[1:], " ")

	switch strings.ToLower(args.Args[0]) {

	case "msg", "message", "text":
		if ok, err := c.checkReqArgs(args, 2); !ok || err != nil {
			fmt.Println(ok, err)
			return err
		}
		if err = db.SetGuildJoinMsg(args.Guild.ID, chanID, argsJoined); err != nil {
			return err
		}
		resTxt = "Join message set."

	case "chan", "channel", "ch":
		if ok, err := c.checkReqArgs(args, 2); !ok || err != nil {
			return err
		}
		ch, err := util.FetchChannel(args.Session, args.Guild.ID, argsJoined, func(c *discordgo.Channel) bool {
			return c.Type == discordgo.ChannelTypeGuildText
		})
		if err != nil {
			return err
		}
		if ch == nil {
			rmsg, err := util.SendEmbedError(args.Session, args.Channel.ID,
				"text channel could not be found.")
			util.DeleteMessageLater(args.Session, rmsg, 6*time.Second)
			return err
		}

		if err = db.SetGuildJoinMsg(args.Guild.ID, ch.ID, msg); err != nil {
			return err
		}
		resTxt = "Join message channel set."

	case "reset", "remove", "rem":
		if err = db.SetGuildJoinMsg(args.Guild.ID, "", ""); err != nil {
			return err
		}
		resTxt = "Join message reset and disabled."

	default:
		rmsg, err := util.SendEmbedError(args.Session, args.Channel.ID,
			"Invalid arguments. Use `help joinmsg` to get help about how to use this command.")
		util.DeleteMessageLater(args.Session, rmsg, 10*time.Second)
		return err
	}

	rmsg, err := util.SendEmbed(args.Session, args.Channel.ID,
		resTxt, "", 0)
	util.DeleteMessageLater(args.Session, rmsg, 10*time.Second)
	return err
}

func (c *CmdJoinMsg) checkReqArgs(args *CommandArgs, req int) (bool, error) {
	if len(args.Args) < req {
		rmsg, err := util.SendEmbedError(args.Session, args.Channel.ID,
			"Invalid arguments. Use `help joinmsg` to get help about how to use this command.")
		util.DeleteMessageLater(args.Session, rmsg, 10*time.Second)
		return false, err
	}
	return true, nil
}
