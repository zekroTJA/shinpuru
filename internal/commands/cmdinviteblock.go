package commands

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/zekroTJA/shinpuru/internal/core"
	"github.com/zekroTJA/shinpuru/internal/util"
)

type CmdInviteBlock struct {
}

func (c *CmdInviteBlock) GetInvokes() []string {
	return []string{"inv", "invblock"}
}

func (c *CmdInviteBlock) GetDescription() string {
	return "manage Discord invite blocking in chat"
}

func (c *CmdInviteBlock) GetHelp() string {
	return "`inv enable` - enable invite link blocking\n" +
		"`inv disable` - disable link blocking"
}

func (c *CmdInviteBlock) GetGroup() string {
	return GroupModeration
}

func (c *CmdInviteBlock) GetDomainName() string {
	return "sp.guild.mod.inviteblock"
}

func (c *CmdInviteBlock) Exec(args *CommandArgs) error {
	if len(args.Args) < 1 {
		return c.printStatus(args)
	}

	switch strings.ToLower(args.Args[0]) {
	case "enable", "e", "on":
		return c.enable(args)
	case "disable", "d", "off":
		return c.disable(args)
	default:
		return c.printStatus(args)
	}
}

func (c *CmdInviteBlock) printStatus(args *CommandArgs) error {
	status, err := args.CmdHandler.db.GetGuildInviteBlock(args.Guild.ID)
	if err != nil && !core.IsErrDatabaseNotFound(err) {
		return err
	}

	strStat := "disabled"
	color := util.ColorEmbedOrange
	if status != "" {
		strStat = "enabled (*for members with permission level < " + status + "*)"
		color = util.ColorEmbedGreen
	}

	msg, err := util.SendEmbed(args.Session, args.Channel.ID,
		fmt.Sprintf("Discord invite link blocking is currently **%s** on this guild.\n\n"+
			"*You can enable or disable this with the command `inv enable` or `inv disable`*.", strStat),
		"", color)
	util.DeleteMessageLater(args.Session, msg, 8*time.Second)
	return err
}

func (c *CmdInviteBlock) enable(args *CommandArgs) error {
	if len(args.Args) < 2 {
		msg, err := util.SendEmbedError(args.Session, args.Channel.ID,
			"Please enter a permission level, which members must have to be allowed to send guild invites.")
		util.DeleteMessageLater(args.Session, msg, 8*time.Second)
		return err
	}

	lvl := args.Args[1]

	if i, err := strconv.Atoi(lvl); err != nil {
		msg, err := util.SendEmbedError(args.Session, args.Channel.ID,
			"Please enter a valid number as permission level.")
		util.DeleteMessageLater(args.Session, msg, 6*time.Second)
		return err
	} else if i < 1 {
		msg, err := util.SendEmbedError(args.Session, args.Channel.ID,
			"Permission level should be larger than 0.")
		util.DeleteMessageLater(args.Session, msg, 6*time.Second)
		return err
	}

	err := args.CmdHandler.db.SetGuildInviteBlock(args.Guild.ID, lvl)
	if err != nil {
		return err
	}

	msg, err := util.SendEmbed(args.Session, args.Channel.ID,
		fmt.Sprintf("Enabled invite link blocking for members with a permission "+
			"level below `%s`.", lvl), "", 0)
	util.DeleteMessageLater(args.Session, msg, 8*time.Second)
	return err
}

func (c *CmdInviteBlock) disable(args *CommandArgs) error {
	err := args.CmdHandler.db.SetGuildInviteBlock(args.Guild.ID, "")
	if err != nil {
		return err
	}

	msg, err := util.SendEmbed(args.Session, args.Channel.ID,
		"Discord invite links will **no more be blocked** on this guild now.",
		"", 0)
	util.DeleteMessageLater(args.Session, msg, 6*time.Second)
	return err
}
