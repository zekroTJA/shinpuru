package commands

import (
	"time"

	"github.com/zekroTJA/shinpuru/internal/core/database"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/static"
)

type CmdPrefix struct {
}

func (c *CmdPrefix) GetInvokes() []string {
	return []string{"prefix", "pre", "guildpre", "guildprefix"}
}

func (c *CmdPrefix) GetDescription() string {
	return "set a custom prefix for your guild"
}

func (c *CmdPrefix) GetHelp() string {
	return "`prefix` - display current guilds prefix\n" +
		"`prefix <newPrefix>` - set the current guilds prefix"
}

func (c *CmdPrefix) GetGroup() string {
	return GroupGuildConfig
}

func (c *CmdPrefix) GetDomainName() string {
	return "sp.guild.config.prefix"
}

func (c *CmdPrefix) GetSubPermissionRules() []SubPermission {
	return nil
}

func (c *CmdPrefix) Exec(args *CommandArgs) error {
	db := args.CmdHandler.db

	if len(args.Args) == 0 {
		prefix, err := db.GetGuildPrefix(args.Guild.ID)
		if !database.IsErrDatabaseNotFound(err) && err != nil {
			return err
		}
		defPrefix := args.CmdHandler.config.Discord.GeneralPrefix
		if prefix == "" || prefix == defPrefix {
			err = util.SendEmbed(args.Session, args.Channel.ID,
				"The current guild prefix is not set, so the default prefix of the bot must be used: ```\n"+defPrefix+"\n```",
				"", 0).DeleteAfter(8 * time.Second).Error()
		} else {
			err = util.SendEmbed(args.Session, args.Channel.ID,
				"The current guild prefix is: ```\n"+prefix+"\n``` "+
					"Surely, you can still use the general prefix (`"+defPrefix+"`)",
				"", 0).DeleteAfter(8 * time.Second).Error()
		}
		return err
	}

	err := db.SetGuildPrefix(args.Guild.ID, args.Args[0])
	if err != nil {
		return err
	}

	return util.SendEmbed(args.Session, args.Channel.ID,
		"Guild prefix is now set to: ```\n"+args.Args[0]+"\n```",
		"", static.ColorEmbedUpdated).
		DeleteAfter(8 * time.Second).Error()
}
