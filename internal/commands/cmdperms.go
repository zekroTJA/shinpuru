package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/zekroTJA/shinpuru/internal/core"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/util"
)

type CmdPerms struct {
}

func (c *CmdPerms) GetInvokes() []string {
	return []string{"perms", "perm", "permlvl", "plvl"}
}

func (c *CmdPerms) GetDescription() string {
	return "Set the permission for specific groups on your server"
}

func (c *CmdPerms) GetHelp() string {
	return "`perms` - get current permission settings\n" +
		"`perms <PDNS> <RoleResolvable> (<RoleResolvable> ...)` - set permission for specific roles\n\n" +
		"PDNS (permission domain name specifier) is used to define permissions to groups by domains. This specifier consists of two parts:\n" +
		"The allow (`+`) / disallow (`-`) part and the domain name (`sp.guilds.config.*` for example).\n\n" +
		"For example, if you want to allow all guild moderation commands for moderators use `+sp.guild.mod.*`. If you want to disallow a role to use a specific command like " +
		"`sp!ban`, you can do this by disallowing the specific domain name `-sp.guild.mod.ban`.\n\n" +
		"Keep in mind:\n" +
		"`-` and `+` of the same domain always results in a disallow.\n" +
		"Higher level rules (like `sp.guild.config.*`) always override lower level rules (like `sp.guild.*`).\n\n" +
		"[**Here**](https://github.com/zekroTJA/shinpuru/blob/master/docs/permissions-guide.md) you can find further information about the permission system."
}

func (c *CmdPerms) GetGroup() string {
	return GroupGuildConfig
}

func (c *CmdPerms) GetDomainName() string {
	return "sp.guild.config.perms"
}

func (c *CmdPerms) Exec(args *CommandArgs) error {
	db := args.CmdHandler.db
	perms, err := db.GetGuildPermissions(args.Guild.ID)
	if err != nil {
		return err
	}

	if len(args.Args) == 0 {
		msgstr := ""

		for roleID, pa := range perms {
			msgstr += fmt.Sprintf("**<@&%s>**\n%s\n\n", roleID, strings.Join(pa, "\n"))
		}

		_, err = util.SendEmbed(args.Session, args.Channel.ID,
			msgstr+"\n*Guild owners does always have permissions over the domains `sp.guild`, `sp.chat` and `sp.etc` "+
				"and the owner of the bot has everywhere permissions over `sp`.*", "Permission settings for this guild", 0)
		return err
	}

	if len(args.Args) < 2 {
		msg, err := util.SendEmbedError(args.Session, args.Channel.ID,
			"Invalid arguments. Use `help perms` to get information how to use this command.")
		util.DeleteMessageLater(args.Session, msg, 10*time.Second)
		return err
	}

	perm := strings.ToLower(args.Args[0])
	sperm := perm[1:]
	if !strings.HasPrefix(sperm, "sp.guild") && !strings.HasPrefix(sperm, "sp.etc") && !strings.HasPrefix(sperm, "sp.chat") {
		msg, err := util.SendEmbedError(args.Session, args.Channel.ID,
			"You are only able to set permissions for the domains `sp.guild`, `sp.etc` and `sp.chat`")
		util.DeleteMessageLater(args.Session, msg, 10*time.Second)
		return err
	}

	roles := make([]*discordgo.Role, 0)
	for _, roleID := range args.Args[1:] {
		if r, err := util.FetchRole(args.Session, args.Guild.ID, roleID); err == nil {
			roles = append(roles, r)
		}
	}

	rolesIds := make([]string, len(roles))
	for i, r := range roles {
		rolesIds[i] = fmt.Sprintf("<@&%s>", r.ID)

		cPerm, ok := perms[r.ID]
		if !ok {
			cPerm = make(core.PermissionArray, 0)
		}

		cPerm = cPerm.Update(perm)

		err := db.SetGuildRolePermission(args.Guild.ID, r.ID, cPerm)
		if err != nil {
			return err
		}
	}

	multipleRoles := ""
	if len(roles) > 1 {
		multipleRoles = "'s"
	}
	_, err = util.SendEmbed(args.Session, args.Channel.ID,
		fmt.Sprintf("Set permission `%s` for role%s %s.",
			perm, multipleRoles, strings.Join(rolesIds, ", ")),
		"", util.ColorEmbedUpdated)

	return err
}
