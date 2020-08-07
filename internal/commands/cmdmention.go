package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/pkg/fetch"
	"github.com/zekroTJA/shireikan"
)

type CmdMention struct {
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
	return shireikan.GroupModeration
}

func (c *CmdMention) GetDomainName() string {
	return "sp.guild.mod.ment"
}

func (c *CmdMention) GetSubPermissionRules() []shireikan.SubPermission {
	return nil
}

func (c *CmdMention) IsExecutableInDMChannels() bool {
	return false
}

func (c *CmdMention) Exec(ctx shireikan.Context) error {
	if len(ctx.GetArgs()) < 1 {
		rolesStr := ""
		for _, role := range ctx.GetGuild().Roles {
			if role.Mentionable {
				rolesStr += fmt.Sprintf("- <@&%s> (%s)\n", role.ID, role.Name)
			}
		}
		return util.SendEmbed(ctx.GetSession(), ctx.GetChannel().ID, rolesStr, "Currently mentionable roles:", 0).
			Error()
	}
	role, err := fetch.FetchRole(ctx.GetSession(), ctx.GetGuild().ID, ctx.GetArgs().Get(0).AsString())
	if err != nil {
		return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID, "Could not fetch any message to the passed resolvable.").
			DeleteAfter(8 * time.Second).Error()
	}

	if role.Mentionable {
		_, err := ctx.GetSession().GuildRoleEdit(ctx.GetGuild().ID, role.ID, role.Name, role.Color, role.Hoist, role.Permissions, false)
		if err != nil {
			return err
		}
		return util.SendEmbed(ctx.GetSession(), ctx.GetChannel().ID,
			fmt.Sprintf("Disabled mentionability for role <@&%s>.", role.ID), "", 0).
			DeleteAfter(10 * time.Second).Error()
	}

	_, err = ctx.GetSession().GuildRoleEdit(ctx.GetGuild().ID, role.ID, role.Name, role.Color, role.Hoist, role.Permissions, true)
	if err != nil {
		return err
	}
	if len(ctx.GetArgs()) > 1 && strings.ToLower(ctx.GetArgs().Get(1).AsString()) == "g" {
		return util.SendEmbed(ctx.GetSession(), ctx.GetChannel().ID,
			fmt.Sprintf("Enabled mentionability for role <@&%s> permanently.\n"+
				"Use `ment <roleResolable>` to disable mentionality of this role.", role.ID), "", 0).
			DeleteAfter(10 * time.Second).Error()
	}

	var handlerRemove func()
	handlerRemove = ctx.GetSession().AddHandler(func(s *discordgo.Session, e *discordgo.MessageCreate) {
		if e.GuildID == ctx.GetGuild().ID && e.Author.ID == ctx.GetUser().ID && len(e.MentionRoles) > 0 {
			for _, rID := range e.MentionRoles {
				if rID == role.ID {
					s.GuildRoleEdit(ctx.GetGuild().ID, role.ID, role.Name, role.Color, role.Hoist, role.Permissions, false)
					if handlerRemove != nil {
						handlerRemove()
					}
				}
			}
		}
	})

	return util.SendEmbed(ctx.GetSession(), ctx.GetChannel().ID,
		fmt.Sprintf("Enabled mentionability for role <@&%s>.\n"+
			"This role will be automatically set to unmentionable after you mention this role in any message on this guild.", role.ID), "", 0).
		DeleteAfter(10 * time.Second).Error()
}
