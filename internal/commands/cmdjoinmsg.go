package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/core/database"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/pkg/fetch"
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

func (c *CmdJoinMsg) GetSubPermissionRules() []SubPermission {
	return nil
}

func (c *CmdJoinMsg) Exec(args *CommandArgs) error {
	db := args.CmdHandler.db

	chanID, msg, err := db.GetGuildJoinMsg(args.Guild.ID)
	if err != nil && err != database.ErrDatabaseNotFound {
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

		return util.SendEmbed(args.Session, args.Channel.ID,
			resTxt, "", 0).
			DeleteAfter(10 * time.Second).Error()
	}

	argsJoined := strings.Join(args.Args[1:], " ")

	switch strings.ToLower(args.Args[0]) {

	case "msg", "message", "text":
		if ok, err := c.checkReqArgs(args, 2); !ok || err != nil {
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
		ch, err := fetch.FetchChannel(args.Session, args.Guild.ID, argsJoined, func(c *discordgo.Channel) bool {
			return c.Type == discordgo.ChannelTypeGuildText
		})
		if err != nil {
			return err
		}
		if ch == nil {
			return util.SendEmbedError(args.Session, args.Channel.ID,
				"text channel could not be found.").
				DeleteAfter(6 * time.Second).Error()
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
		return util.SendEmbedError(args.Session, args.Channel.ID,
			"Invalid arguments. Use `help joinmsg` to get help about how to use this command.").
			DeleteAfter(10 * time.Second).Error()
	}

	return util.SendEmbed(args.Session, args.Channel.ID,
		resTxt, "", 0).
		DeleteAfter(10 * time.Second).Error()
}

func (c *CmdJoinMsg) checkReqArgs(args *CommandArgs, req int) (bool, error) {
	if len(args.Args) < req {
		err := util.SendEmbedError(args.Session, args.Channel.ID,
			"Invalid arguments. Use `help joinmsg` to get help about how to use this command.").
			DeleteAfter(8 * time.Second).Error()
		return false, err
	}
	return true, nil
}
