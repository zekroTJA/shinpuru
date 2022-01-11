package slashcommands

import (
	"bytes"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/bwmarrin/snowflake"
	"github.com/zekroTJA/shinpuru/internal/models"
	"github.com/zekroTJA/shinpuru/internal/services/config"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/services/permissions"
	"github.com/zekroTJA/shinpuru/internal/services/report"
	"github.com/zekroTJA/shinpuru/internal/services/storage"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/imgstore"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/acceptmsg"
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
	return "1.1.0"
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
						{
							Name:  "kick",
							Value: 0,
						},
						{
							Name:  "ban",
							Value: 1,
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
			Term:        "kick",
			Explicit:    false,
			Description: "Kick a member.",
		},
		{
			Term:        "ban",
			Explicit:    false,
			Description: "Ban a member.",
		},
		{
			Term:        "revoke",
			Explicit:    false,
			Description: "Revoke a report.",
		},
	}
}

func (c *Report) Run(ctx *ken.Ctx) (err error) {
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

func (c *Report) create(ctx *ken.SubCommandCtx) (err error) {
	cfg := ctx.Get(static.DiConfig).(config.Provider)
	repSvc := ctx.Get(static.DiReport).(*report.ReportService)
	pmw := ctx.Get(static.DiPermissions).(*permissions.Permissions)

	typ := models.ReportType(ctx.Options().GetByName("type").IntValue())
	victim := ctx.Options().GetByName("user").UserValue(ctx.Ctx)
	reason := ctx.Options().GetByName("reason").StringValue()

	var dn string
	switch typ {
	case models.TypeKick:
		dn = "kick"
	case models.TypeBan:
		dn = "ban"
	case models.TypeWarn, models.TypeAd:
		dn = "warn"
	}

	ok, err := pmw.CheckSubPerm(ctx.Ctx, dn, false)
	if err != nil && ok {
		return
	}

	var attachment, expire string
	if imageurlV, ok := ctx.Options().GetByNameOptional("imageurl"); ok {
		attachment = imageurlV.StringValue()
	}
	if expireV, ok := ctx.Options().GetByNameOptional("expire"); ok {
		expire = expireV.StringValue()
	}

	if attachment != "" {
		img, err := imgstore.DownloadFromURL(attachment)
		if err == nil && img != nil {
			st, _ := ctx.Get(static.DiObjectStorage).(storage.Storage)
			err = st.PutObject(static.StorageBucketImages, img.ID.String(),
				bytes.NewReader(img.Data), int64(img.Size), img.MimeType)
			if err != nil {
				return err
			}
			attachment = img.ID.String()
		}
	}

	rep := &models.Report{
		GuildID:       ctx.Event.GuildID,
		ExecutorID:    ctx.User().ID,
		VictimID:      victim.ID,
		Msg:           reason,
		AttachmentURL: attachment,
		Type:          typ,
	}

	if expire != "" {
		exp, err := time.ParseDuration(expire)
		if err != nil {
			err = ctx.FollowUpError(
				fmt.Sprintf("Invalid duration:\n```\n%s```", err.Error()), "").Error
			return err
		}
		expT := time.Now().Add(exp)
		rep.Timeout = &expT
	}

	emb := rep.AsEmbed(cfg.Config().WebServer.PublicAddr)
	emb.Title = "Report Check"
	emb.Description = "Is everything okay so far?"

	acceptMsg := acceptmsg.AcceptMessage{
		Embed:          emb,
		Session:        ctx.Session,
		UserID:         ctx.User().ID,
		DeleteMsgAfter: true,
		AcceptFunc: func(msg *discordgo.Message) (err error) {
			switch typ {
			case models.TypeKick:
				rep, err = repSvc.PushKick(rep)
			case models.TypeBan:
				rep, err = repSvc.PushBan(rep)
			default:
				rep, err = repSvc.PushReport(rep)
			}

			if err != nil {
				return
			}

			_, err = ctx.Session.ChannelMessageSendEmbed(ctx.Event.ChannelID, rep.AsEmbed(cfg.Config().WebServer.PublicAddr))
			return
		},
	}

	if _, err = acceptMsg.AsFollowUp(ctx.Ctx); err != nil {
		return
	}
	return acceptMsg.Error()
}

func (c *Report) revoke(ctx *ken.SubCommandCtx) (err error) {
	db, _ := ctx.Get(static.DiDatabase).(database.Database)
	cfg, _ := ctx.Get(static.DiConfig).(config.Provider)
	repSvc, _ := ctx.Get(static.DiReport).(*report.ReportService)
	pmw := ctx.Get(static.DiPermissions).(*permissions.Permissions)

	ok, err := pmw.CheckSubPerm(ctx.Ctx, "revoke", false)
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
				Error
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
		Session:        ctx.Session,
		DeleteMsgAfter: true,
		UserID:         ctx.User().ID,
		DeclineFunc: func(m *discordgo.Message) (err error) {
			return util.SendEmbedError(ctx.Session, ctx.Event.ChannelID,
				"Canceled.").
				DeleteAfter(8 * time.Second).Error()
		},
		AcceptFunc: func(m *discordgo.Message) (err error) {
			emb, err := repSvc.RevokeReport(
				rep,
				ctx.User().ID,
				reason,
				cfg.Config().WebServer.PublicAddr,
				db,
				ctx.Session)

			if err != nil {
				return
			}

			_, err = ctx.Session.ChannelMessageSendEmbed(ctx.Event.ChannelID, emb)
			return
		},
	}

	if _, err = aceptMsg.AsFollowUp(ctx.Ctx); err != nil {
		return err
	}
	return aceptMsg.Error()
}

func (c *Report) list(ctx *ken.SubCommandCtx) (err error) {
	db, _ := ctx.Get(static.DiDatabase).(database.Database)
	cfg, _ := ctx.Get(static.DiConfig).(config.Provider)
	pmw := ctx.Get(static.DiPermissions).(*permissions.Permissions)

	ok, err := pmw.CheckSubPerm(ctx.Ctx, "list", false)
	if err != nil && ok {
		return
	}

	victim := ctx.Options().GetByName("user").UserValue(ctx.Ctx)

	emb := &discordgo.MessageEmbed{
		Color: static.ColorEmbedDefault,
		Title: fmt.Sprintf("Reports for %s#%s",
			victim.Username, victim.Discriminator),
		Description: fmt.Sprintf("[**Here**](%s/guilds/%s/%s) you can find this users reports in the web interface.",
			cfg.Config().WebServer.PublicAddr, ctx.Event.GuildID, victim.ID),
	}
	reps, err := db.GetReportsFiltered(ctx.Event.GuildID, victim.ID, -1, 0, 1000)
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
	err = ctx.FollowUpEmbed(emb).Error
	return
}
