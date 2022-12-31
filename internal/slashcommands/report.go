package slashcommands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/bwmarrin/snowflake"
	"github.com/zekroTJA/shinpuru/internal/models"
	"github.com/zekroTJA/shinpuru/internal/services/config"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/services/permissions"
	"github.com/zekroTJA/shinpuru/internal/services/report"
	"github.com/zekroTJA/shinpuru/internal/services/timeprovider"
	"github.com/zekroTJA/shinpuru/internal/util/cmdutil"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/acceptmsg/v2"
	"github.com/zekrotja/ken"
)

type Report struct{}

var (
	_ ken.SlashCommand        = (*Report)(nil)
	_ permissions.PermCommand = (*Report)(nil)
)

func (c *Report) Name() string {
	return "report"
}

func (c *Report) Description() string {
	return "Create, revoke or list user reports."
}

func (c *Report) Version() string {
	return "1.2.0"
}

func (c *Report) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (c *Report) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "create",
			Description: "File a new report.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "type",
					Description: "The type of report.",
					Required:    true,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "warn",
							Value: 3,
						},
						{
							Name:  "ad",
							Value: 4,
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user",
					Description: "The user.",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "reason",
					Description: "A short and concise report reason.",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "imageurl",
					Description: "An image url embedded into the report.",
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "expire",
					Description: "Expire report after given time.",
				},
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "revoke",
			Description: "Revoke a report.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "id",
					Description: "ID of the report to be revoked.",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "reason",
					Description: "Reason of the revoke.",
					Required:    true,
				},
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "list",
			Description: "List the reports of a user.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user",
					Description: "User to list reports of.",
					Required:    true,
				},
			},
		},
	}
}

func (c *Report) Domain() string {
	return "sp.guild.mod.report"
}

func (c *Report) SubDomains() []permissions.SubPermission {
	return []permissions.SubPermission{
		{
			Term:        "list",
			Explicit:    false,
			Description: "List a users reports.",
		},
		{
			Term:        "warn",
			Explicit:    false,
			Description: "Warn a member.",
		},
		{
			Term:        "revoke",
			Explicit:    false,
			Description: "Revoke a report.",
		},
	}
}

func (c *Report) Run(ctx ken.Context) (err error) {
	if err = ctx.Defer(); err != nil {
		return
	}

	err = ctx.HandleSubCommands(
		ken.SubCommandHandler{"create", c.create},
		ken.SubCommandHandler{"revoke", c.revoke},
		ken.SubCommandHandler{"list", c.list},
	)

	return
}

func (c *Report) create(ctx ken.SubCommandContext) (err error) {
	pmw := ctx.Get(static.DiPermissions).(*permissions.Permissions)

	typ := models.ReportType(ctx.Options().GetByName("type").IntValue())

	ok, err := pmw.CheckSubPerm(ctx, "warn", false)
	if err != nil && ok {
		return
	}

	tp := ctx.Get(static.DiTimeProvider).(timeprovider.Provider)

	return cmdutil.CmdReport(ctx, typ, tp)
}

func (c *Report) revoke(ctx ken.SubCommandContext) (err error) {
	db, _ := ctx.Get(static.DiDatabase).(database.Database)
	cfg, _ := ctx.Get(static.DiConfig).(config.Provider)
	repSvc, _ := ctx.Get(static.DiReport).(*report.ReportService)
	pmw := ctx.Get(static.DiPermissions).(*permissions.Permissions)

	ok, err := pmw.CheckSubPerm(ctx, "revoke", false)
	if err != nil && ok {
		return
	}

	idStr := ctx.Options().GetByName("id").StringValue()
	reason := ctx.Options().GetByName("reason").StringValue()

	id, err := snowflake.ParseString(idStr)
	if err != nil {
		return
	}

	rep, err := db.GetReport(id)
	if err != nil {
		if database.IsErrDatabaseNotFound(err) {
			return ctx.FollowUpError(
				fmt.Sprintf("Could not find any report with ID `%d`", id), "").
				Send().Error
		}
		return err
	}

	aceptMsg := acceptmsg.AcceptMessage{
		Embed: &discordgo.MessageEmbed{
			Color: static.ReportRevokedColor,
			Title: "Report Revocation",
			Description: "Do you really want to revoke this report?\n" +
				":warning: **WARNING:** Revoking a report will be displayed in the mod log channel (if set) and " +
				"the revoke will be **deleted** from the database and no more visible again after!",
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:  "Revocation Reason",
					Value: reason,
				},
				rep.AsEmbedField(cfg.Config().WebServer.PublicAddr),
			},
		},
		Ken:            ctx.GetKen(),
		DeleteMsgAfter: true,
		UserID:         ctx.User().ID,
		DeclineFunc: func(cctx ken.ComponentContext) (err error) {
			return cctx.RespondError("Canceled.", "")
		},
		AcceptFunc: func(cctx ken.ComponentContext) (err error) {
			if err = cctx.Defer(); err != nil {
				return err
			}

			emb, err := repSvc.RevokeReport(
				rep,
				ctx.User().ID,
				reason,
				cfg.Config().WebServer.PublicAddr,
				db,
				ctx.GetSession())

			if err != nil {
				return
			}

			return cctx.FollowUpEmbed(emb).Send().Error
		},
	}

	if _, err = aceptMsg.AsFollowUp(ctx); err != nil {
		return err
	}
	return aceptMsg.Error()
}

func (c *Report) list(ctx ken.SubCommandContext) (err error) {
	db, _ := ctx.Get(static.DiDatabase).(database.Database)
	cfg, _ := ctx.Get(static.DiConfig).(config.Provider)
	pmw := ctx.Get(static.DiPermissions).(*permissions.Permissions)

	ok, err := pmw.CheckSubPerm(ctx, "list", false)
	if err != nil && ok {
		return
	}

	victim := ctx.Options().GetByName("user").UserValue(ctx)

	emb := &discordgo.MessageEmbed{
		Color: static.ColorEmbedDefault,
		Title: fmt.Sprintf("Reports for %s#%s",
			victim.Username, victim.Discriminator),
		Description: fmt.Sprintf("[**Here**](%s/guilds/%s/%s) you can find this users reports in the web interface.",
			cfg.Config().WebServer.PublicAddr, ctx.GetEvent().GuildID, victim.ID),
	}
	reps, err := db.GetReportsFiltered(ctx.GetEvent().GuildID, victim.ID, -1, 0, 1000)
	if err != nil {
		return err
	}
	if len(reps) == 0 {
		emb.Description += "\n\nThis user has a white west. :ok_hand:"
	} else {
		emb.Fields = make([]*discordgo.MessageEmbedField, 0)
		for _, r := range reps {
			emb.Fields = append(emb.Fields, r.AsEmbedField(cfg.Config().WebServer.PublicAddr))
		}
	}
	err = ctx.FollowUpEmbed(emb).Send().Error
	return
}
