package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/zekroTJA/shinpuru/internal/middleware"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/acceptmsg"
	"github.com/zekroTJA/shinpuru/pkg/stringutil"
	"github.com/zekroTJA/shireikan"
)

type CmdNotify struct {
}

func (c *CmdNotify) GetInvokes() []string {
	return []string{"notify", "n"}
}

func (c *CmdNotify) GetDescription() string {
	return "Get, remove or setup the notify rule."
}

func (c *CmdNotify) GetHelp() string {
	return "`notify setup (<roleName>)` - creates the notify role and registers it for this command\n" +
		"`notify` - get or remove the role"
}

func (c *CmdNotify) GetGroup() string {
	return shireikan.GroupModeration
}

func (c *CmdNotify) GetDomainName() string {
	return "sp.chat.notify"
}

func (c *CmdNotify) GetSubPermissionRules() []shireikan.SubPermission {
	return []shireikan.SubPermission{
		{
			Term:        "setup",
			Explicit:    true,
			Description: "Allows setting up the notify role for this guild.",
		},
	}
}

func (c *CmdNotify) IsExecutableInDMChannels() bool {
	return false
}

func (c *CmdNotify) Exec(ctx shireikan.Context) error {
	db, _ := ctx.GetObject(static.DiDatabase).(database.Database)

	if len(ctx.GetArgs()) < 1 {
		notifyRoleID, err := db.GetGuildNotifyRole(ctx.GetGuild().ID)
		if database.IsErrDatabaseNotFound(err) || notifyRoleID == "" {
			return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
				"No notify role  was set up for this guild.").
				DeleteAfter(8 * time.Second).Error()
		}
		if err != nil {
			return err
		}
		var roleExists bool
		for _, role := range ctx.GetGuild().Roles {
			if notifyRoleID == role.ID && !roleExists {
				roleExists = true
			}
		}
		if !roleExists {
			return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
				"The set notify role does not exist on this guild anymore. Please notify a "+
					"moderator aor admin about this to fix this. ;)").
				DeleteAfter(8 * time.Second).Error()
		}
		member, err := ctx.GetSession().GuildMember(ctx.GetGuild().ID, ctx.GetUser().ID)
		if err != nil {
			return err
		}
		msgStr := "Removed notify role."
		if stringutil.IndexOf(notifyRoleID, member.Roles) > -1 {
			err = ctx.GetSession().GuildMemberRoleRemove(ctx.GetGuild().ID, ctx.GetUser().ID, notifyRoleID)
			if err != nil {
				return err
			}
		} else {
			err = ctx.GetSession().GuildMemberRoleAdd(ctx.GetGuild().ID, ctx.GetUser().ID, notifyRoleID)
			if err != nil {
				return err
			}
			msgStr = "Added notify role."
		}
		return util.SendEmbed(ctx.GetSession(), ctx.GetChannel().ID, msgStr, "", 0).
			DeleteAfter(8 * time.Second).Error()
	}

	if strings.ToLower(ctx.GetArgs().Get(0).AsString()) == "setup" {
		pmw, _ := ctx.GetObject(static.DiPermissionMiddleware).(*middleware.PermissionsMiddleware)
		ok, override, err := pmw.CheckPermissions(ctx.GetSession(), ctx.GetGuild().ID, ctx.GetUser().ID, c.GetDomainName()+".setup")
		if err != nil {
			return err
		}
		if !ok && !override {
			return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
				"Sorry, but you do'nt have the permission to setup the notify role.").
				DeleteAfter(8 * time.Second).Error()
		}
		var notifyRoleExists bool
		notifyRoleID, err := db.GetGuildNotifyRole(ctx.GetGuild().ID)
		if err == nil {
			for _, role := range ctx.GetGuild().Roles {
				if notifyRoleID == role.ID && !notifyRoleExists {
					notifyRoleExists = true
				}
			}
		}
		notifiableStr := "\n*Notify role is defaulty not notifiable. You need to enable this manually by using the " +
			"`ment` command or toggling it manually in the discord settings.*"
		if notifyRoleExists {
			am := &acceptmsg.AcceptMessage{
				Session:        ctx.GetSession(),
				UserID:         ctx.GetUser().ID,
				DeleteMsgAfter: true,
				Embed: &discordgo.MessageEmbed{
					Color: static.ColorEmbedDefault,
					Description: fmt.Sprintf("The notify role on this guild is already set to <@&%s>.\n"+
						"Do you want to overwrite this setting? This will also **delete** the role <@&%s>.",
						notifyRoleID, notifyRoleID),
				},
				AcceptFunc: func(m *discordgo.Message) (err error) {
					role, err := c.setup(ctx)
					if err != nil {
						return
					}
					err = ctx.GetSession().GuildRoleDelete(ctx.GetGuild().ID, notifyRoleID)
					if err != nil {
						return
					}
					err = util.SendEmbed(ctx.GetSession(), ctx.GetChannel().ID,
						fmt.Sprintf("Updated notify role to <@&%s>."+notifiableStr, role.ID), "", 0).
						DeleteAfter(8 * time.Second).Error()
					return
				},
				DeclineFunc: func(m *discordgo.Message) (err error) {
					err = util.SendEmbed(ctx.GetSession(), ctx.GetChannel().ID,
						"Canceled.", "", 0).
						DeleteAfter(8 * time.Second).Error()
					return
				},
			}

			if _, err := am.Send(ctx.GetChannel().ID); err != nil {
				return err
			}

			return am.Error()
		}

		role, err := c.setup(ctx)
		if err != nil {
			return err
		}
		return util.SendEmbed(ctx.GetSession(), ctx.GetChannel().ID,
			fmt.Sprintf("Set notify role to <@&%s>.", role.ID)+notifiableStr, "", 0).
			DeleteAfter(8 * time.Second).Error()
	}

	return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
		"Invalid command arguments. Please use `help notify` to get help about this command.").
		DeleteAfter(8 * time.Second).Error()
}

func (c *CmdNotify) setup(ctx shireikan.Context) (*discordgo.Role, error) {
	db, _ := ctx.GetObject(static.DiDatabase).(database.Database)

	name := "Notify"
	if len(ctx.GetArgs()) > 1 {
		name = strings.Join(ctx.GetArgs()[1:], " ")
	}
	role, err := ctx.GetSession().GuildRoleCreate(ctx.GetGuild().ID)
	if err != nil {
		return nil, err
	}
	role, err = ctx.GetSession().GuildRoleEdit(ctx.GetGuild().ID, role.ID, name, 0, false, 0, false)
	if err != nil {
		return nil, err
	}
	err = db.SetGuildNotifyRole(ctx.GetGuild().ID, role.ID)
	return role, err
}
