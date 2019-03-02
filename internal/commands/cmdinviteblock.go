package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/zekroTJA/shinpuru/internal/core"
	"github.com/zekroTJA/shinpuru/internal/util"
)

type CmdInviteBlock struct {
	PermLvl int
}

func (c *CmdInviteBlock) GetInvokes() []string {
	return []string{"inv", "invblock"}
}

func (c *CmdInviteBlock) GetDescription() string {
	return "manage Discord invite blocking in chat"
}

func (c *CmdInviteBlock) GetHelp() string {
	return "`inv <enable|disable>`"
}

func (c *CmdInviteBlock) GetGroup() string {
	return GroupModeration
}

func (c *CmdInviteBlock) GetPermission() int {
	return c.PermLvl
}

func (c *CmdInviteBlock) SetPermission(permLvl int) {
	c.PermLvl = permLvl
}

func (c *CmdInviteBlock) Exec(args *CommandArgs) error {
	if len(args.Args) < 1 {
		return c.printStatus(args)
	}

	switch strings.ToLower(args.Args[0]) {
	case "enable", "e", "on":
		return c.swtitchStatus(args, true)
	case "disable", "d", "off":
		return c.swtitchStatus(args, false)
	default:
		return c.printStatus(args)
	}
}

func (c *CmdInviteBlock) printStatus(args *CommandArgs) error {
	enabled, err := args.CmdHandler.db.GetGuildInviteBlock(args.Guild.ID)
	if err != nil && !core.IsErrDatabaseNotFound(err) {
		return err
	}

	strStat := "disabled"
	color := util.ColorEmbedOrange
	if enabled {
		strStat = "enabled"
		color = util.ColorEmbedGreen
	}

	msg, err := util.SendEmbed(args.Session, args.Channel.ID,
		fmt.Sprintf("Discord invite link blocing is currently **%s** on this guild.\n\n"+
			"*You can enable or disable this with the command `inv enable` or `inv disable`*.", strStat),
		"", color)
	util.DeleteMessageLater(args.Session, msg, 8*time.Second)
	return err
}

func (c *CmdInviteBlock) swtitchStatus(args *CommandArgs, enable bool) error {
	err := args.CmdHandler.db.SetGuildInviteBlock(args.Guild.ID, enable)
	if err != nil {
		return err
	}

	strStat := "will **no more be blocked** now"
	color := util.ColorEmbedOrange
	if enable {
		strStat = "will now **be blocked**"
		color = util.ColorEmbedGreen
	}

	msg, err := util.SendEmbed(args.Session, args.Channel.ID,
		fmt.Sprintf("Discord invite links %s on this guild.", strStat),
		"", color)
	util.DeleteMessageLater(args.Session, msg, 6*time.Second)
	return err
}
