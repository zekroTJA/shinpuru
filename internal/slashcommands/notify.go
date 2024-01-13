package slashcommands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/services/permissions"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/acceptmsg/v2"
	"github.com/zekroTJA/shinpuru/pkg/stringutil"
	"github.com/zekrotja/dgrs"
	"github.com/zekrotja/ken"
)

type Notify struct{}

var (
	_ ken.SlashCommand        = (*Notify)(nil)
	_ permissions.PermCommand = (*Notify)(nil)
)

func (c *Notify) Name() string {
	return "notify"
}

func (c *Notify) Description() string {
	return "Get, remove or setup the notify role."
}

func (c *Notify) Version() string {
	return "1.0.0"
}

func (c *Notify) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (c *Notify) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "toggle",
			Description: "Get or remove notify role.",
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "setup",
			Description: "Setup notify role.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionRole,
					Name:        "role",
					Description: "The role to be used as notify role (will be created if not specified).",
				},
			},
		},
	}
}

func (c *Notify) Domain() string {
	return "sp.chat.notify"
}

func (c *Notify) SubDomains() []permissions.SubPermission {
	return []permissions.SubPermission{
		{
			Term:        "setup",
			Explicit:    true,
			Description: "Allows setting up the notify role for this guild.",
		},
	}
}

func (c *Notify) Run(ctx ken.Context) (err error) {
	err = ctx.HandleSubCommands(
		ken.SubCommandHandler{"toggle", c.toggle},
		ken.SubCommandHandler{"setup", c.setup},
	)

	return
}

func (c *Notify) toggle(ctx ken.SubCommandContext) (err error) {
	db := ctx.Get(static.DiDatabase).(database.Database)
	st := ctx.Get(static.DiState).(*dgrs.State)

	ctx.SetEphemeral(true)
	if err = ctx.Defer(); err != nil {
		return
	}

	notifyRoleID, err := db.GetGuildNotifyRole(ctx.GetEvent().GuildID)
	if database.IsErrDatabaseNotFound(err) || notifyRoleID == "" {
		return ctx.FollowUpError(
			"No notify role  was set up for this guild.", "").Send().Error
	}
	if err != nil {
		return err
	}

	roles, err := st.Roles(ctx.GetEvent().GuildID, true)
	if err != nil {
		return
	}
	var roleExists bool
	for _, role := range roles {
		if notifyRoleID == role.ID && !roleExists {
			roleExists = true
		}
	}

	if !roleExists {
		return ctx.FollowUpError(
			"The set notify role does not exist on this guild anymore. Please notify a "+
				"moderator aor admin about this to fix this. ;)", "").Send().Error
	}

	member, err := st.Member(ctx.GetEvent().GuildID, ctx.User().ID)
	if err != nil {
		return err
	}

	msgStr := "Removed notify role."
	if stringutil.IndexOf(notifyRoleID, member.Roles) > -1 {
		err = ctx.GetSession().GuildMemberRoleRemove(ctx.GetEvent().GuildID, ctx.User().ID, notifyRoleID)
		if err != nil {
			return err
		}
	} else {
		err = ctx.GetSession().GuildMemberRoleAdd(ctx.GetEvent().GuildID, ctx.User().ID, notifyRoleID)
		if err != nil {
			return err
		}
		msgStr = "Added notify role."
	}

	return ctx.FollowUpEmbed(&discordgo.MessageEmbed{
		Description: msgStr,
	}).Send().Error
}

func (c *Notify) setup(ctx ken.SubCommandContext) (err error) {
	pmw := ctx.Get(static.DiPermissions).(*permissions.Permissions)
	db := ctx.Get(static.DiDatabase).(database.Database)
	st := ctx.Get(static.DiState).(*dgrs.State)

	if err = ctx.Defer(); err != nil {
		return
	}

	ok, err := pmw.CheckSubPerm(ctx, "setup", true)
	if err != nil {
		return err
	}
	if !ok {
		return ctx.FollowUpError(
			"Sorry, but you do'nt have the permission to setup the notify role.", "").
			Send().Error
	}

	roles, err := st.Roles(ctx.GetEvent().GuildID)
	if err != nil {
		return err
	}

	var notifyRoleExists bool
	notifyRoleID, err := db.GetGuildNotifyRole(ctx.GetEvent().GuildID)
	if err == nil {
		for _, role := range roles {
			if notifyRoleID == role.ID && !notifyRoleExists {
				notifyRoleExists = true
			}
		}
	}
	notifiableStr := "\n*Notify role is defaulty not notifiable. You need to enable this manually by using the " +
		"`ment` command or toggling it manually in the discord settings.*"
	if notifyRoleExists {
		am := &acceptmsg.AcceptMessage{
			Ken:            ctx.GetKen(),
			UserID:         ctx.User().ID,
			DeleteMsgAfter: true,
			Embed: &discordgo.MessageEmbed{
				Color: static.ColorEmbedDefault,
				Description: fmt.Sprintf("The notify role on this guild is already set to <@&%s>.\n"+
					"Do you want to overwrite this setting? This will also **delete** the role <@&%s>.",
					notifyRoleID, notifyRoleID),
			},
			AcceptFunc: func(cctx ken.ComponentContext) (err error) {
				if err = cctx.Defer(); err != nil {
					return
				}
				role, err := c.setupRole(ctx)
				if err != nil {
					return
				}
				err = ctx.GetSession().GuildRoleDelete(ctx.GetEvent().GuildID, notifyRoleID)
				if err != nil {
					return
				}
				return cctx.FollowUpEmbed(&discordgo.MessageEmbed{
					Description: fmt.Sprintf("Updated notify role to <@&%s>."+notifiableStr, role.ID),
				}).Send().Error
			},
			DeclineFunc: func(cctx ken.ComponentContext) (err error) {
				return cctx.RespondEmbed(&discordgo.MessageEmbed{
					Description: "Canceled",
				})
			},
		}

		if _, err := am.AsFollowUp(ctx); err != nil {
			return err
		}

		return am.Error()
	}

	role, err := c.setupRole(ctx)
	if err != nil {
		return err
	}
	err = ctx.FollowUpEmbed(&discordgo.MessageEmbed{
		Description: fmt.Sprintf("Updated notify role to <@&%s>."+notifiableStr, role.ID),
	}).Send().Error
	return
}

func (c *Notify) setupRole(ctx ken.SubCommandContext) (role *discordgo.Role, err error) {
	db, _ := ctx.Get(static.DiDatabase).(database.Database)

	const roleName = "Notify"
	if roleV, ok := ctx.Options().GetByNameOptional("role"); ok {
		role = roleV.RoleValue(ctx)
	} else {
		role, err = ctx.GetSession().GuildRoleCreate(ctx.GetEvent().GuildID, &discordgo.RoleParams{
			Name: roleName,
		})
		if err != nil {
			return
		}
	}

	err = db.SetGuildNotifyRole(ctx.GetEvent().GuildID, role.ID)
	return role, err
}
