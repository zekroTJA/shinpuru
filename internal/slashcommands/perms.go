package slashcommands

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/services/config"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	permService "github.com/zekroTJA/shinpuru/internal/services/permissions"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/permissions"
	"github.com/zekroTJA/shinpuru/pkg/roleutil"
	"github.com/zekrotja/ken"
)

type Perms struct{}

type permMode string

const (
	modeAllow    = "+"
	modeDisallow = "-"
)

var (
	_ ken.Command             = (*Bug)(nil)
	_ permService.PermCommand = (*Bug)(nil)
)

func (c *Perms) Name() string {
	return "perms"
}

func (c *Perms) Description() string {
	return "Set the permissions for groups on your guild."
}

func (c *Perms) Version() string {
	return "1.0.0"
}

func (c *Perms) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (c *Perms) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "list",
			Description: "List the current permission definitions.",
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "set",
			Description: "Set a permission rule for specific roles.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "mode",
					Description: "Set the permission as allow or disallow.",
					Required:    true,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "allow",
							Value: modeAllow,
						},
						{
							Name:  "disallow",
							Value: modeDisallow,
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "dns",
					Description: "Permission Domain Name Specifier",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionRole,
					Name:        "role",
					Description: "The role to apply the permission to.",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionRole,
					Name:        "role2",
					Description: "Additional role to apply the permission to.",
				},
				{
					Type:        discordgo.ApplicationCommandOptionRole,
					Name:        "role3",
					Description: "Additional role to apply the permission to.",
				},
				{
					Type:        discordgo.ApplicationCommandOptionRole,
					Name:        "role4",
					Description: "Additional role to apply the permission to.",
				},
				{
					Type:        discordgo.ApplicationCommandOptionRole,
					Name:        "role5",
					Description: "Additional role to apply the permission to.",
				},
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "help",
			Description: "Display help information for this command.",
		},
	}
}

func (c *Bug) DomainName() string {
	return "sp.guild.config.perms"
}

func (c *Perms) SubDomains() []permService.SubPermission {
	return nil
}

func (c *Perms) Run(ctx *ken.Ctx) (err error) {
	if err = ctx.Defer(); err != nil {
		return
	}

	err = ctx.HandleSubCommands(
		ken.SubCommandHandler{"list", c.list},
		ken.SubCommandHandler{"set", c.set},
		ken.SubCommandHandler{"help", c.help},
	)

	return
}

func (c *Perms) list(ctx *ken.SubCommandCtx) (err error) {
	db, _ := ctx.Get(static.DiDatabase).(database.Database)

	perms, err := db.GetGuildPermissions(ctx.Event.GuildID)
	if err != nil {
		return err
	}

	sortedGuildRoles, err := roleutil.GetSortedGuildRoles(ctx.Session, ctx.Event.GuildID, true)
	if err != nil {
		return err
	}

	msgstr := ""

	for _, role := range sortedGuildRoles {
		if pa, ok := perms[role.ID]; ok {
			msgstr += fmt.Sprintf("**<@&%s>**\n%s\n\n", role.ID, strings.Join(pa, "\n"))
		}
	}

	return ctx.FollowUpEmbed(&discordgo.MessageEmbed{
		Description: msgstr + "\n*Guild owners does always have permissions over the domains `sp.guild`, `sp.chat` and `sp.etc` " +
			"and the owner of the bot has everywhere permissions over `sp`.*",
		Title: "Permission settings for this guild",
	}).Error
}

func (c *Perms) set(ctx *ken.SubCommandCtx) (err error) {
	db, _ := ctx.Get(static.DiDatabase).(database.Database)

	mode := ctx.Options().GetByName("mode").StringValue()
	dns := ctx.Options().GetByName("dns").StringValue()

	dns = mode + dns

	roles := []*discordgo.Role{
		ctx.Options().GetByName("role").RoleValue(ctx.Ctx),
	}

	for i := 1; i < 6; i++ {
		if rV, ok := ctx.Options().GetByNameOptional(fmt.Sprintf("role%d", i)); ok {
			roles = append(roles, rV.RoleValue(ctx.Ctx))
		}
	}

	perms, err := db.GetGuildPermissions(ctx.Event.GuildID)
	if err != nil {
		return err
	}

	rolesIds := make([]string, len(roles))
	for i, r := range roles {
		rolesIds[i] = fmt.Sprintf("<@&%s>", r.ID)

		cPerm, ok := perms[r.ID]
		if !ok {
			cPerm = make(permissions.PermissionArray, 0)
		}

		cPerm, changed := cPerm.Update(dns, false)

		if changed {
			err := db.SetGuildRolePermission(ctx.Event.GuildID, r.ID, cPerm)
			if err != nil {
				return err
			}
		}
	}

	multipleRoles := ""
	if len(roles) > 1 {
		multipleRoles = "'s"
	}

	return ctx.FollowUpEmbed(&discordgo.MessageEmbed{
		Description: fmt.Sprintf("Set permission `%s` for role%s %s.",
			dns, multipleRoles, strings.Join(rolesIds, ", ")),
	}).Error
}

func (c *Perms) help(ctx *ken.SubCommandCtx) (err error) {
	cfg := ctx.Get(static.DiConfig).(config.Provider)

	desc := "If you donÄt know how the permissions system works, " +
		"please read [**this**](https://github.com/zekroTJA/shinpuru/wiki/Permissions-Guide) " +
		"wiki article to learn more.\n\n"

	wsc := cfg.Config().WebServer
	if wsc.Enabled {
		desc += fmt.Sprintf("You can also set permissions in the [**web interface**](%s), which "+
			"is way more visual, interactive and easy than doing it via commands. 😉",
			wsc.PublicAddr)
	}

	return ctx.FollowUpEmbed(&discordgo.MessageEmbed{
		Description: desc,
	}).Error
}
