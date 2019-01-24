package commands

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/util"
)

type CmdPerms struct {
	PermLvl int
}

func (c *CmdPerms) GetInvokes() []string {
	return []string{"perms", "perm", "permlvl", "plvl"}
}

func (c *CmdPerms) GetDescription() string {
	return "Set the permission for specific groups on your server"
}

func (c *CmdPerms) GetHelp() string {
	return "`perms` - get current permission settings\n" +
		"`perms <LvL> <RoleResolvable> (<RoleResolvable> ...)` - set permission level for specific roles"
}

func (c *CmdPerms) GetGroup() string {
	return GroupGuildConfig
}

func (c *CmdPerms) GetPermission() int {
	return c.PermLvl
}

func (c *CmdPerms) SetPermission(permLvl int) {
	c.PermLvl = permLvl
}

func (c *CmdPerms) Exec(args *CommandArgs) error {
	db := args.CmdHandler.db

	if len(args.Args) == 0 {
		msgstr := ""
		perms, err := db.GetGuildPermissions(args.Guild.ID)
		if err != nil {
			return err
		}
		for roleID, permLvl := range perms {
			msgstr += fmt.Sprintf("`%02d` - <@&%s>\n", permLvl, roleID)
		}
		_, err = util.SendEmbed(args.Session, args.Channel.ID,
			msgstr+"\n*Guild owners does always have perm LvL 10 and the owner of the bot has everywhere perm LvL 999.*",
			"Permission Level for this Guild", 0)
		return err
	}

	if len(args.Args) < 2 {
		msg, err := util.SendEmbedError(args.Session, args.Channel.ID,
			"Invalid arguments. Use `help perms` to get information how to use this command.")
		util.DeleteMessageLater(args.Session, msg, 10*time.Second)
		return err
	}

	permLvL, err := strconv.Atoi(args.Args[0])
	if err != nil {
		msg, err := util.SendEmbedError(args.Session, args.Channel.ID,
			"First argument is is the permission level and must be a valid number.")
		util.DeleteMessageLater(args.Session, msg, 10*time.Second)
		return err
	} else if permLvL < 0 || permLvL > 9 {
		msg, err := util.SendEmbedError(args.Session, args.Channel.ID,
			"The permission level must be a number between *(including)* 0 and 9.")
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
		err := db.SetGuildRolePermission(args.Guild.ID, r.ID, permLvL)
		if err != nil {
			return err
		}
	}

	multipleRoles := ""
	if len(roles) > 1 {
		multipleRoles = "'s"
	}
	_, err = util.SendEmbed(args.Session, args.Channel.ID,
		fmt.Sprintf("Set permission level `%d` for role%s %s.",
			permLvL, multipleRoles, strings.Join(rolesIds, ", ")),
		"", util.ColorEmbedUpdated)

	return err
}
