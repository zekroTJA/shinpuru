package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/fetch"
	"github.com/zekroTJA/shireikan"
	"github.com/zekrotja/discordgo"
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
	return shireikan.GroupGuildConfig
}

func (c *CmdLeaveMsg) GetDomainName() string {
	return "sp.guild.config.leavemsg"
}

func (c *CmdLeaveMsg) GetSubPermissionRules() []shireikan.SubPermission {
	return nil
}

func (c *CmdLeaveMsg) IsExecutableInDMChannels() bool {
	return false
}

func (c *CmdLeaveMsg) Exec(ctx shireikan.Context) error {
	db, _ := ctx.GetObject(static.DiDatabase).(database.Database)

	chanID, msg, err := db.GetGuildLeaveMsg(ctx.GetGuild().ID)
	if err != nil && err != database.ErrDatabaseNotFound {
		return err
	}

	var resTxt string

	if len(ctx.GetArgs()) < 1 {

		if msg == "" && chanID == "" {
			resTxt = "*Leave message and channel not set.*"
		} else if msg == "" {
			resTxt = "*Leave message not set.*"
		} else if chanID == "" {
			resTxt = "*Leave channel not set.*"
		} else {
			ch, err := ctx.GetSession().Channel(chanID)
			if ch == nil || err != nil || ch.GuildID != ctx.GetGuild().ID {
				resTxt = "*Leave channel not set.*"
			} else {
				resTxt = fmt.Sprintf("```\n%s\n``` is set as leave message and will be posted "+
					"into channel %s (%s).", msg, ch.Mention(), ch.ID)
			}
		}

		return util.SendEmbed(ctx.GetSession(), ctx.GetChannel().ID,
			resTxt, "", 0).
			DeleteAfter(10 * time.Second).Error()
	}

	argsJoined := strings.Join(ctx.GetArgs()[1:], " ")

	switch strings.ToLower(ctx.GetArgs().Get(0).AsString()) {

	case "msg", "message", "text":
		if ok, err := c.checkReqArgs(ctx, 2); !ok || err != nil {
			return err
		}
		if err = db.SetGuildLeaveMsg(ctx.GetGuild().ID, chanID, argsJoined); err != nil {
			return err
		}
		resTxt = "Leave message set."

	case "chan", "channel", "ch":
		if ok, err := c.checkReqArgs(ctx, 2); !ok || err != nil {
			return err
		}
		ch, err := fetch.FetchChannel(ctx.GetSession(), ctx.GetGuild().ID, argsJoined, func(c *discordgo.Channel) bool {
			return c.Type == discordgo.ChannelTypeGuildText
		})
		if err != nil {
			return err
		}
		if ch == nil {
			return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
				"Text channel could not be found.").
				DeleteAfter(8 * time.Second).Error()
		}

		if err = db.SetGuildLeaveMsg(ctx.GetGuild().ID, ch.ID, msg); err != nil {
			return err
		}
		resTxt = "Leave message channel set."

	case "reset", "remove", "rem":
		if err = db.SetGuildLeaveMsg(ctx.GetGuild().ID, "", ""); err != nil {
			return err
		}
		resTxt = "Leave message reset and disabled."

	default:
		return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
			"Invalid arguments. Use `help leavemsg` to get help about how to use this command.").
			DeleteAfter(8 * time.Second).Error()
	}

	return util.SendEmbed(ctx.GetSession(), ctx.GetChannel().ID,
		resTxt, "", static.ColorEmbedGreen).
		DeleteAfter(10 * time.Second).Error()
}

func (c *CmdLeaveMsg) checkReqArgs(ctx shireikan.Context, req int) (bool, error) {
	if len(ctx.GetArgs()) < req {
		err := util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
			"Invalid arguments. Use `help leavemsg` to get help about how to use this command.").
			DeleteAfter(10 * time.Second).Error()
		return false, err
	}
	return true, nil
}
