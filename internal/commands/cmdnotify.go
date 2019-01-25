package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/zekroTJA/shinpuru/internal/core"
	"github.com/zekroTJA/shinpuru/internal/util"
)

type CmdNotify struct {
	PermLvl int
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

func (c *CmdNotify) GetPermission() int {
	return c.PermLvl
}

func (c *CmdNotify) SetPermission(permLvl int) {
	c.PermLvl = permLvl
}

func (c *CmdNotify) Exec(args *CommandArgs) error {
	if len(args.Args) < 1 {
		notifyRoleID, err := args.CmdHandler.db.GetGuildNotifyRole(args.Guild.ID)
		if core.IsErrDatabaseNotFound(err) || notifyRoleID == "" {
			msg, err := util.SendEmbedError(args.Session, args.Channel.ID,
				"No notify role  was set up for this guild.")
			util.DeleteMessageLater(args.Session, msg, 10*time.Second)
			return err
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
			msg, err := util.SendEmbedError(args.Session, args.Channel.ID,
				"The set notify role does not exist on this guild anymore. Please notify a "+
					"moderator aor admin about this to fix this. ;)")
			util.DeleteMessageLater(args.Session, msg, 10*time.Second)
			return err
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
		msg, err := util.SendEmbed(args.Session, args.Channel.ID, msgStr, "", 0)
		util.DeleteMessageLater(args.Session, msg, 5*time.Second)
		return err
	}

	if strings.ToLower(args.Args[0]) == "setup" {
		permLvl, err := args.CmdHandler.db.GetMemberPermissionLevel(args.Session, args.Guild.ID, args.User.ID)
		if err != nil {
			return err
		}
		if permLvl < 6 && args.User.ID != args.Guild.OwnerID && args.User.ID != args.CmdHandler.config.Discord.OwnerID {
			msg, err := util.SendEmbedError(args.Session, args.Channel.ID,
				"Sorry, but you do'nt have the permission to setup the notify role.")
			util.DeleteMessageLater(args.Session, msg, 10*time.Second)
			return err
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
			am := &util.AcceptMessage{
				Session:        args.Session,
				UserID:         args.User.ID,
				DeleteMsgAfter: true,
				Embed: &discordgo.MessageEmbed{
					Color: util.ColorEmbedDefault,
					Description: fmt.Sprintf("The notify role on this guild is already set to <@&%s>.\n"+
						"Do you want to overwrite this setting? This will also **delete** the role <@&%s>.",
						notifyRoleID, notifyRoleID),
				},
				AcceptFunc: func(m *discordgo.Message) {
					role, err := c.setup(args)
					if err != nil {
						msg, _ := util.SendEmbedError(args.Session, args.Channel.ID,
							"Failed setup: "+err.Error())
						util.DeleteMessageLater(args.Session, msg, 10*time.Second)
						return
					}
					err = args.Session.GuildRoleDelete(args.Guild.ID, notifyRoleID)
					if err != nil {
						msg, _ := util.SendEmbedError(args.Session, args.Channel.ID,
							"Failed deleting old notify role: "+err.Error())
						util.DeleteMessageLater(args.Session, msg, 10*time.Second)
						return
					}
					msg, _ := util.SendEmbed(args.Session, args.Channel.ID,
						fmt.Sprintf("Updated notify role to <@&%s>."+notifiableStr, role.ID), "", 0)
					util.DeleteMessageLater(args.Session, msg, 10*time.Second)
				},
				DeclineFunc: func(m *discordgo.Message) {
					msg, _ := util.SendEmbed(args.Session, args.Channel.ID,
						"Canceled.", "", 0)
					util.DeleteMessageLater(args.Session, msg, 6*time.Second)
				},
			}
			_, err := am.Send(args.Channel.ID)
			return err
		}

		role, err := c.setup(args)
		if err != nil {
			return err
		}
		msg, err := util.SendEmbed(args.Session, args.Channel.ID,
			fmt.Sprintf("Set notify role to <@&%s>.", role.ID)+notifiableStr, "", 0)
		util.DeleteMessageLater(args.Session, msg, 10*time.Second)
		return err
	}

	msg, err := util.SendEmbedError(args.Session, args.Channel.ID,
		"Invalid command arguments. Please use `help notify` to get help about this command.")
	util.DeleteMessageLater(args.Session, msg, 10*time.Second)
	return err
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
