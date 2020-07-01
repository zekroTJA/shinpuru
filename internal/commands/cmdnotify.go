package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/zekroTJA/shinpuru/internal/core/database"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/acceptmsg"
)

type CmdNotify struct {
}

func (c *CmdNotify) GetInvokes() []string {
	return []string{"notify", "n"}
}

func (c *CmdNotify) GetDescription() string {
	return "get, remove or setup the notify rule"
}

func (c *CmdNotify) GetHelp() string {
	return "`notify setup (<roleName>)` - creates the notify role and registers it for this command\n" +
		"`notify` - get or remove the role"
}

func (c *CmdNotify) GetGroup() string {
	return GroupModeration
}

func (c *CmdNotify) GetDomainName() string {
	return "sp.chat.notify"
}

func (c *CmdNotify) GetSubPermissionRules() []SubPermission {
	return []SubPermission{
		{
			Term:        "setup",
			Explicit:    true,
			Description: "Allows setting up the notify role for this guild.",
		},
	}
}

func (c *CmdNotify) Exec(args *CommandArgs) error {
	if len(args.Args) < 1 {
		notifyRoleID, err := args.CmdHandler.db.GetGuildNotifyRole(args.Guild.ID)
		if database.IsErrDatabaseNotFound(err) || notifyRoleID == "" {
			return util.SendEmbedError(args.Session, args.Channel.ID,
				"No notify role  was set up for this guild.").
				DeleteAfter(8 * time.Second).Error()
		}
		if err != nil {
			return err
		}
		var roleExists bool
		for _, role := range args.Guild.Roles {
			if notifyRoleID == role.ID && !roleExists {
				roleExists = true
			}
		}
		if !roleExists {
			return util.SendEmbedError(args.Session, args.Channel.ID,
				"The set notify role does not exist on this guild anymore. Please notify a "+
					"moderator aor admin about this to fix this. ;)").
				DeleteAfter(8 * time.Second).Error()
		}
		member, err := args.Session.GuildMember(args.Guild.ID, args.User.ID)
		if err != nil {
			return err
		}
		msgStr := "Removed notify role."
		if util.IndexOfStrArray(notifyRoleID, member.Roles) > -1 {
			err = args.Session.GuildMemberRoleRemove(args.Guild.ID, args.User.ID, notifyRoleID)
			if err != nil {
				return err
			}
		} else {
			err = args.Session.GuildMemberRoleAdd(args.Guild.ID, args.User.ID, notifyRoleID)
			if err != nil {
				return err
			}
			msgStr = "Added notify role."
		}
		return util.SendEmbed(args.Session, args.Channel.ID, msgStr, "", 0).
			DeleteAfter(8 * time.Second).Error()
	}

	if strings.ToLower(args.Args[0]) == "setup" {
		ok, override, err := args.CmdHandler.CheckPermissions(args.Session, args.Guild.ID, args.User.ID, c.GetDomainName()+".setup")
		if err != nil {
			return err
		}
		if !ok && !override {
			return util.SendEmbedError(args.Session, args.Channel.ID,
				"Sorry, but you do'nt have the permission to setup the notify role.").
				DeleteAfter(8 * time.Second).Error()
		}
		var notifyRoleExists bool
		notifyRoleID, err := args.CmdHandler.db.GetGuildNotifyRole(args.Guild.ID)
		if err == nil {
			for _, role := range args.Guild.Roles {
				if notifyRoleID == role.ID && !notifyRoleExists {
					notifyRoleExists = true
				}
			}
		}
		notifiableStr := "\n*Notify role is defaulty not notifiable. You need to enable this manually by using the " +
			"`ment` command or toggling it manually in the discord settings.*"
		if notifyRoleExists {
			am := &acceptmsg.AcceptMessage{
				Session:        args.Session,
				UserID:         args.User.ID,
				DeleteMsgAfter: true,
				Embed: &discordgo.MessageEmbed{
					Color: static.ColorEmbedDefault,
					Description: fmt.Sprintf("The notify role on this guild is already set to <@&%s>.\n"+
						"Do you want to overwrite this setting? This will also **delete** the role <@&%s>.",
						notifyRoleID, notifyRoleID),
				},
				AcceptFunc: func(m *discordgo.Message) {
					role, err := c.setup(args)
					if err != nil {
						util.SendEmbedError(args.Session, args.Channel.ID,
							"Failed setup: "+err.Error()).
							DeleteAfter(8 * time.Second)
						return
					}
					err = args.Session.GuildRoleDelete(args.Guild.ID, notifyRoleID)
					if err != nil {
						util.SendEmbedError(args.Session, args.Channel.ID,
							"Failed deleting old notify role: "+err.Error()).
							DeleteAfter(8 * time.Second)
						return
					}
					util.SendEmbed(args.Session, args.Channel.ID,
						fmt.Sprintf("Updated notify role to <@&%s>."+notifiableStr, role.ID), "", 0).
						DeleteAfter(8 * time.Second)
				},
				DeclineFunc: func(m *discordgo.Message) {
					util.SendEmbed(args.Session, args.Channel.ID,
						"Canceled.", "", 0).
						DeleteAfter(8 * time.Second)
				},
			}
			_, err := am.Send(args.Channel.ID)
			return err
		}

		role, err := c.setup(args)
		if err != nil {
			return err
		}
		return util.SendEmbed(args.Session, args.Channel.ID,
			fmt.Sprintf("Set notify role to <@&%s>.", role.ID)+notifiableStr, "", 0).
			DeleteAfter(8 * time.Second).Error()
	}

	return util.SendEmbedError(args.Session, args.Channel.ID,
		"Invalid command arguments. Please use `help notify` to get help about this command.").
		DeleteAfter(8 * time.Second).Error()
}

func (c *CmdNotify) setup(args *CommandArgs) (*discordgo.Role, error) {
	name := "Notify"
	if len(args.Args) > 1 {
		name = strings.Join(args.Args[1:], " ")
	}
	role, err := args.Session.GuildRoleCreate(args.Guild.ID)
	if err != nil {
		return nil, err
	}
	role, err = args.Session.GuildRoleEdit(args.Guild.ID, role.ID, name, 0, false, 0, false)
	if err != nil {
		return nil, err
	}
	err = args.CmdHandler.db.SetGuildNotifyRole(args.Guild.ID, role.ID)
	return role, err
}
