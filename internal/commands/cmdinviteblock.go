package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shireikan"
)

type CmdInviteBlock struct {
}

func (c *CmdInviteBlock) GetInvokes() []string {
	return []string{"inv", "invblock"}
}

func (c *CmdInviteBlock) GetDescription() string {
	return "Manage Discord invite blocking in chat."
}

func (c *CmdInviteBlock) GetHelp() string {
	return "`inv enable` - enable invite link blocking\n" +
		"`inv disable` - disable link blocking"
}

func (c *CmdInviteBlock) GetGroup() string {
	return shireikan.GroupModeration
}

func (c *CmdInviteBlock) GetDomainName() string {
	return "sp.guild.mod.inviteblock"
}

func (c *CmdInviteBlock) GetSubPermissionRules() []shireikan.SubPermission {
	return []shireikan.SubPermission{
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

func (c *CmdInviteBlock) Exec(ctx shireikan.Context) error {
	if len(ctx.GetArgs()) < 1 {
		return c.printStatus(ctx)
	}

	switch strings.ToLower(ctx.GetArgs().Get(0).AsString()) {
	case "enable", "e", "on":
		return c.enable(ctx)
	case "disable", "d", "off":
		return c.disable(ctx)
	default:
		return c.printStatus(ctx)
	}
}

func (c *CmdInviteBlock) printStatus(ctx shireikan.Context) error {
	db, _ := ctx.GetObject(static.DiDatabase).(database.Database)

	status, err := db.GetGuildInviteBlock(ctx.GetGuild().ID)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return err
	}

	strStat := "disabled"
	color := static.ColorEmbedOrange
	if status != "" {
		strStat = "enabled"
		color = static.ColorEmbedGreen
	}

	return util.SendEmbed(ctx.GetSession(), ctx.GetChannel().ID,
		fmt.Sprintf("Discord invite link blocking is currently **%s** on this guild.\n\n"+
			"*You can enable or disable this with the command `inv enable` or `inv disable`*.", strStat),
		"", color).
		DeleteAfter(8 * time.Second).Error()
}

func (c *CmdInviteBlock) enable(ctx shireikan.Context) error {
	db, _ := ctx.GetObject(static.DiDatabase).(database.Database)

	err := db.SetGuildInviteBlock(ctx.GetGuild().ID, "1")
	if err != nil {
		return err
	}

	return util.SendEmbed(ctx.GetSession(), ctx.GetChannel().ID,
		"Enabled invite link blocking.", "", 0).
		DeleteAfter(8 * time.Second).Error()
}

func (c *CmdInviteBlock) disable(ctx shireikan.Context) error {
	db, _ := ctx.GetObject(static.DiDatabase).(database.Database)

	err := db.SetGuildInviteBlock(ctx.GetGuild().ID, "")
	if err != nil {
		return err
	}

	return util.SendEmbed(ctx.GetSession(), ctx.GetChannel().ID,
		"Discord invite links will **no more be blocked** on this guild now.",
		"", 0).
		DeleteAfter(6 * time.Second).Error()
}
