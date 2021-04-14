package commands

import (
	"time"

	"github.com/zekroTJA/shinpuru/internal/core/config"
	"github.com/zekroTJA/shinpuru/internal/core/database"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shireikan"
)

type CmdPrefix struct {
}

func (c *CmdPrefix) GetInvokes() []string {
	return []string{"prefix", "pre", "guildpre", "guildprefix"}
}

func (c *CmdPrefix) GetDescription() string {
	return "Set a custom prefix for your guild."
}

func (c *CmdPrefix) GetHelp() string {
	return "`prefix` - display current guilds prefix\n" +
		"`prefix <newPrefix>` - set the current guilds prefix"
}

func (c *CmdPrefix) GetGroup() string {
	return shireikan.GroupGuildConfig
}

func (c *CmdPrefix) GetDomainName() string {
	return "sp.guild.config.prefix"
}

func (c *CmdPrefix) GetSubPermissionRules() []shireikan.SubPermission {
	return nil
}

func (c *CmdPrefix) IsExecutableInDMChannels() bool {
	return false
}

func (c *CmdPrefix) Exec(ctx shireikan.Context) error {
	db, _ := ctx.GetObject(static.DiDatabase).(database.Database)
	cfg, _ := ctx.GetObject(static.DiConfig).(*config.Config)

	if len(ctx.GetArgs()) == 0 {
		prefix, err := db.GetGuildPrefix(ctx.GetGuild().ID)
		if !database.IsErrDatabaseNotFound(err) && err != nil {
			return err
		}
		defPrefix := cfg.Discord.GeneralPrefix
		if prefix == "" || prefix == defPrefix {
			err = util.SendEmbed(ctx.GetSession(), ctx.GetChannel().ID,
				"The current guild prefix is not set, so the default prefix of the bot must be used: ```\n"+defPrefix+"\n```",
				"", 0).DeleteAfter(8 * time.Second).Error()
		} else {
			err = util.SendEmbed(ctx.GetSession(), ctx.GetChannel().ID,
				"The current guild prefix is: ```\n"+prefix+"\n``` "+
					"Surely, you can still use the general prefix (`"+defPrefix+"`)",
				"", 0).DeleteAfter(8 * time.Second).Error()
		}
		return err
	}

	err := db.SetGuildPrefix(ctx.GetGuild().ID, ctx.GetArgs().Get(0).AsString())
	if err != nil {
		return err
	}

	return util.SendEmbed(ctx.GetSession(), ctx.GetChannel().ID,
		"Guild prefix is now set to: ```\n"+ctx.GetArgs().Get(0).AsString()+"\n```",
		"", static.ColorEmbedUpdated).
		DeleteAfter(8 * time.Second).Error()
}
