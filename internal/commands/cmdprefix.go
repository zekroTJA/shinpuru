package commands

import (
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/zekroTJA/shinpuru/internal/core"
	"github.com/zekroTJA/shinpuru/internal/util"
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

func (c *CmdPrefix) Exec(args *CommandArgs) error {
	db := args.CmdHandler.db

	if len(args.Args) == 0 {
		prefix, err := db.GetGuildPrefix(args.Guild.ID)
		if !core.IsErrDatabaseNotFound(err) && err != nil {
			return err
		}
		defPrefix := args.CmdHandler.config.Discord.GeneralPrefix
		var msg *discordgo.Message
		if prefix == "" || prefix == defPrefix {
			msg, err = util.SendEmbed(args.Session, args.Channel.ID,
				"The current guild prefix is not set, so the default prefix of the bot must be used: ```\n"+defPrefix+"\n```",
				"", 0)
		} else {
			msg, err = util.SendEmbed(args.Session, args.Channel.ID,
				"The current guild prefix is: ```\n"+prefix+"\n``` "+
					"Surely, you can still use the general prefix (`"+defPrefix+"`)",
				"", 0)
		}
		util.DeleteMessageLater(args.Session, msg, 10*time.Second)
		return err
	}

	err := db.SetGuildPrefix(args.Guild.ID, args.Args[0])
	if err != nil {
		return err
	}

	msg, err := util.SendEmbed(args.Session, args.Channel.ID,
		"Guild prefix is now set to: ```\n"+args.Args[0]+"\n```",
		"", util.ColorEmbedUpdated)
	util.DeleteMessageLater(args.Session, msg, 10*time.Second)

	return err
}
