package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/zekroTJA/shinpuru/internal/core/database"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/static"
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

func (c *CmdInviteBlock) GetSubPermissionRules() []SubPermission {
	return []SubPermission{
		{
			Term:        "send",
			Explicit:    true,
			Description: "Allows sending invites even if invite block is enabled",
		},
	}
}

func (c *CmdInviteBlock) IsExecutableInDMChannels() bool {
	return false
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
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return err
	}

	strStat := "disabled"
	color := static.ColorEmbedOrange
	if status != "" {
		strStat = "enabled"
		color = static.ColorEmbedGreen
	}

	return util.SendEmbed(args.Session, args.Channel.ID,
		fmt.Sprintf("Discord invite link blocking is currently **%s** on this guild.\n\n"+
			"*You can enable or disable this with the command `inv enable` or `inv disable`*.", strStat),
		"", color).
		DeleteAfter(8 * time.Second).Error()
}

func (c *CmdInviteBlock) enable(args *CommandArgs) error {
	err := args.CmdHandler.db.SetGuildInviteBlock(args.Guild.ID, "1")
	if err != nil {
		return err
	}

	return util.SendEmbed(args.Session, args.Channel.ID,
		"Enabled invite link blocking.", "", 0).
		DeleteAfter(8 * time.Second).Error()
}

func (c *CmdInviteBlock) disable(args *CommandArgs) error {
	err := args.CmdHandler.db.SetGuildInviteBlock(args.Guild.ID, "")
	if err != nil {
		return err
	}

	return util.SendEmbed(args.Session, args.Channel.ID,
		"Discord invite links will **no more be blocked** on this guild now.",
		"", 0).
		DeleteAfter(6 * time.Second).Error()
}
