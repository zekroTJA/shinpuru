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
	return shireikan.GroupGuildConfig
}

func (c *CmdJoinMsg) GetDomainName() string {
	return "sp.guild.config.joinmsg"
}

func (c *CmdJoinMsg) GetSubPermissionRules() []shireikan.SubPermission {
	return nil
}

func (c *CmdJoinMsg) IsExecutableInDMChannels() bool {
	return false
}

func (c *CmdJoinMsg) Exec(ctx shireikan.Context) error {
	db, _ := ctx.GetObject(static.DiDatabase).(database.Database)

	chanID, msg, err := db.GetGuildJoinMsg(ctx.GetGuild().ID)
	if err != nil && err != database.ErrDatabaseNotFound {
		return err
	}

	var resTxt string

	if len(ctx.GetArgs()) < 1 {

		if msg == "" && chanID == "" {
			resTxt = "*Join message and channel not set.*"
		} else if msg == "" {
			resTxt = "*Join message not set.*"
		} else if chanID == "" {
			resTxt = "*Join channel not set.*"
		} else {
			ch, err := ctx.GetSession().Channel(chanID)
			if ch == nil || err != nil || ch.GuildID != ctx.GetGuild().ID {
				resTxt = "*Join channel not set.*"
			} else {
				resTxt = fmt.Sprintf("```\n%s\n``` is set as join message and will be posted "+
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
		if err = db.SetGuildJoinMsg(ctx.GetGuild().ID, chanID, argsJoined); err != nil {
			return err
		}
		resTxt = "Join message set."

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
				"text channel could not be found.").
				DeleteAfter(6 * time.Second).Error()
		}

		if err = db.SetGuildJoinMsg(ctx.GetGuild().ID, ch.ID, msg); err != nil {
			return err
		}
		resTxt = "Join message channel set."

	case "reset", "remove", "rem":
		if err = db.SetGuildJoinMsg(ctx.GetGuild().ID, "", ""); err != nil {
			return err
		}
		resTxt = "Join message reset and disabled."

	default:
		return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
			"Invalid arguments. Use `help joinmsg` to get help about how to use this command.").
			DeleteAfter(10 * time.Second).Error()
	}

	return util.SendEmbed(ctx.GetSession(), ctx.GetChannel().ID,
		resTxt, "", 0).
		DeleteAfter(10 * time.Second).Error()
}

func (c *CmdJoinMsg) checkReqArgs(ctx shireikan.Context, req int) (bool, error) {
	if len(ctx.GetArgs()) < req {
		err := util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
			"Invalid arguments. Use `help joinmsg` to get help about how to use this command.").
			DeleteAfter(8 * time.Second).Error()
		return false, err
	}
	return true, nil
}
