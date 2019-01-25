package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/zekroTJA/shinpuru/internal/util"
)

type CmdMention struct {
	PermLvl int
}

func (c *CmdMention) GetInvokes() []string {
	return []string{"ment", "mnt", "mention", "mentions"}
}

func (c *CmdMention) GetDescription() string {
	return "toggle the mentionability of a role"
}

func (c *CmdMention) GetHelp() string {
	return "`ment` - display currently mentionable roles\n" +
		"`ment <roleResolvable> (g)` - make role mentioanble until you mention the role in a message on the guild. " +
		"By attaching the parameter `g`, the role will be mentionable until this command will be exeuted on the role again."
}

func (c *CmdMention) GetGroup() string {
	return GroupGeneral
}

func (c *CmdMention) GetPermission() int {
	return c.PermLvl
}

func (c *CmdMention) SetPermission(permLvl int) {
	c.PermLvl = permLvl
}

func (c *CmdMention) Exec(args *CommandArgs) error {
	if len(args.Args) < 1 {
		rolesStr := ""
		for _, role := range args.Guild.Roles {
			if role.Mentionable {
				rolesStr += fmt.Sprintf("- <@&%s> (%s)\n", role.ID, role.Name)
			}
		}
		_, err := util.SendEmbed(args.Session, args.Channel.ID, rolesStr, "Currently mentionable roles:", 0)
		return err
	}
	role, err := util.FetchRole(args.Session, args.Guild.ID, args.Args[0])
	if err != nil {
		msg, err := util.SendEmbedError(args.Session, args.Channel.ID, "Could not fetch any message to the passed resolvable.")
		util.DeleteMessageLater(args.Session, msg, 6*time.Second)
		return err
	}

	if role.Mentionable {
		_, err := args.Session.GuildRoleEdit(args.Guild.ID, role.ID, role.Name, role.Color, role.Hoist, role.Permissions, false)
		if err != nil {
			return err
		}
		msg, err := util.SendEmbed(args.Session, args.Channel.ID,
			fmt.Sprintf("Disabled mentionability for role <@&%s>.", role.ID), "", 0)
		util.DeleteMessageLater(args.Session, msg, 10*time.Second)
		return err
	}

	_, err = args.Session.GuildRoleEdit(args.Guild.ID, role.ID, role.Name, role.Color, role.Hoist, role.Permissions, true)
	if err != nil {
		return err
	}
	if len(args.Args) > 1 && strings.ToLower(args.Args[1]) == "g" {
		msg, err := util.SendEmbed(args.Session, args.Channel.ID,
			fmt.Sprintf("Enabled mentionability for role <@&%s> permanently.\n"+
				"Use `ment <roleResolable>` to disable mentionality of this role.", role.ID), "", 0)
		util.DeleteMessageLater(args.Session, msg, 10*time.Second)
		return err
	}

	var handlerRemove func()
	handlerRemove = args.Session.AddHandler(func(s *discordgo.Session, e *discordgo.MessageCreate) {
		if e.GuildID == args.Guild.ID && e.Author.ID == args.User.ID && len(e.MentionRoles) > 0 {
			for _, rID := range e.MentionRoles {
				if rID == role.ID {
					s.GuildRoleEdit(args.Guild.ID, role.ID, role.Name, role.Color, role.Hoist, role.Permissions, false)
					if handlerRemove != nil {
						handlerRemove()
					}
				}
			}
		}
	})

	msg, err := util.SendEmbed(args.Session, args.Channel.ID,
		fmt.Sprintf("Enabled mentionability for role <@&%s>.\n"+
			"This role will be automatically set to unmentionable after you mention this role in any message on this guild.", role.ID), "", 0)
	util.DeleteMessageLater(args.Session, msg, 10*time.Second)
	return err
}
