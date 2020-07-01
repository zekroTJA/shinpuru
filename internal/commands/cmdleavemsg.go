package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/core/database"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/fetch"
)

type CmdLeaveMsg struct {
}

func (c *CmdLeaveMsg) GetInvokes() []string {
	return []string{"leavemsg", "leavemessage"}
}

func (c *CmdLeaveMsg) GetDescription() string {
	return "Set a message which will be sent into the defined channel when a member leaves."
}

func (c *CmdLeaveMsg) GetHelp() string {
	return "`leavemsg msg <message>` - Set the message of the leave message.\n" +
		"`leavemsg channel <ChannelIdentifier>` - Set the channel where the message will be sent into.\n" +
		"`leavemsg reset` - Reset and disable leave messages.\n\n" +
		"`[user]` will be replaced with the user name and `[ment]` will be replaced with the users mention when used in message text."
}

func (c *CmdLeaveMsg) GetGroup() string {
	return GroupGuildConfig
}

func (c *CmdLeaveMsg) GetDomainName() string {
	return "sp.guild.config.leavemsg"
}

func (c *CmdLeaveMsg) GetSubPermissionRules() []SubPermission {
	return nil
}

func (c *CmdLeaveMsg) Exec(args *CommandArgs) error {
	db := args.CmdHandler.db

	chanID, msg, err := db.GetGuildLeaveMsg(args.Guild.ID)
	if err != nil && err != database.ErrDatabaseNotFound {
		return err
	}

	var resTxt string

	if len(args.Args) < 1 {

		if msg == "" && chanID == "" {
			resTxt = "*Leave message and channel not set.*"
		} else if msg == "" {
			resTxt = "*Leave message not set.*"
		} else if chanID == "" {
			resTxt = "*Leave channel not set.*"
		} else {
			ch, err := args.Session.Channel(chanID)
			if ch == nil || err != nil || ch.GuildID != args.Guild.ID {
				resTxt = "*Leave channel not set.*"
			} else {
				resTxt = fmt.Sprintf("```\n%s\n``` is set as leave message and will be posted "+
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
		if err = db.SetGuildLeaveMsg(args.Guild.ID, chanID, argsJoined); err != nil {
			return err
		}
		resTxt = "Leave message set."

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
				"Text channel could not be found.").
				DeleteAfter(8 * time.Second).Error()
		}

		if err = db.SetGuildLeaveMsg(args.Guild.ID, ch.ID, msg); err != nil {
			return err
		}
		resTxt = "Leave message channel set."

	case "reset", "remove", "rem":
		if err = db.SetGuildLeaveMsg(args.Guild.ID, "", ""); err != nil {
			return err
		}
		resTxt = "Leave message reset and disabled."

	default:
		return util.SendEmbedError(args.Session, args.Channel.ID,
			"Invalid arguments. Use `help leavemsg` to get help about how to use this command.").
			DeleteAfter(8 * time.Second).Error()
	}

	return util.SendEmbed(args.Session, args.Channel.ID,
		resTxt, "", static.ColorEmbedGreen).
		DeleteAfter(10 * time.Second).Error()
}

func (c *CmdLeaveMsg) checkReqArgs(args *CommandArgs, req int) (bool, error) {
	if len(args.Args) < req {
		err := util.SendEmbedError(args.Session, args.Channel.ID,
			"Invalid arguments. Use `help leavemsg` to get help about how to use this command.").
			DeleteAfter(10 * time.Second).Error()
		return false, err
	}
	return true, nil
}
