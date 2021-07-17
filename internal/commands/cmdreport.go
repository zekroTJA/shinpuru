package commands

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/zekroTJA/shinpuru/internal/config"
	"github.com/zekroTJA/shinpuru/internal/models"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/services/report"
	"github.com/zekroTJA/shinpuru/internal/services/storage"
	"github.com/zekroTJA/shinpuru/internal/util/imgstore"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/acceptmsg"
	"github.com/zekroTJA/shinpuru/pkg/fetch"
	"github.com/zekroTJA/shireikan"

	"github.com/bwmarrin/discordgo"
	"github.com/bwmarrin/snowflake"

	"github.com/zekroTJA/shinpuru/internal/util"
)

type CmdReport struct {
}

func (c *CmdReport) GetInvokes() []string {
	return []string{"report", "rep", "warn"}
}

func (c *CmdReport) GetDescription() string {
	return "Report a user."
}

func (c *CmdReport) GetHelp() string {
	repTypes := make([]string, len(models.ReportTypes))
	for i, t := range models.ReportTypes {
		repTypes[i] = fmt.Sprintf("`%d` - %s", i, t)
	}
	return "`report <userResolvable>` - list all reports of a user\n" +
		"`report <userResolvable> [<type>] <reason>` - report a user *(if type is empty, its defaultly 0 = warn)*\n" +
		"`report revoke <caseID> <reason>` - revoke a report\n" +
		"\n**TYPES:**\n" + strings.Join(repTypes, "\n") +
		"\nTypes `BAN`, `KICK` and `MUTE` are reserved for bans and kicks executed with this bot."
}

func (c *CmdReport) GetGroup() string {
	return shireikan.GroupModeration
}

func (c *CmdReport) GetDomainName() string {
	return "sp.guild.mod.report"
}

func (c *CmdReport) GetSubPermissionRules() []shireikan.SubPermission {
	return nil
}

func (c *CmdReport) IsExecutableInDMChannels() bool {
	return false
}

func (c *CmdReport) Exec(ctx shireikan.Context) error {
	db, _ := ctx.GetObject(static.DiDatabase).(database.Database)
	cfg, _ := ctx.GetObject(static.DiConfig).(*config.Config)
	repSvc, _ := ctx.GetObject(static.DiReport).(*report.ReportService)

	if len(ctx.GetArgs()) < 1 {
		return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
			"Invalid command arguments. Please use `help report` to see how to use this command.").
			DeleteAfter(8 * time.Second).Error()
	}

	if strings.ToLower(ctx.GetArgs().Get(0).AsString()) == "revoke" {
		return c.revoke(ctx)
	}

	victim, err := fetch.FetchMember(ctx.GetSession(), ctx.GetGuild().ID, ctx.GetArgs().Get(0).AsString())
	if err != nil || victim == nil {
		return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
			"Sorry, could not find any member :cry:").
			DeleteAfter(8 * time.Second).Error()
	}

	if len(ctx.GetArgs()) == 1 {
		emb := &discordgo.MessageEmbed{
			Color: static.ColorEmbedDefault,
			Title: fmt.Sprintf("Reports for %s#%s",
				victim.User.Username, victim.User.Discriminator),
			Description: fmt.Sprintf("[**Here**](%s/guilds/%s/%s) you can find this users reports in the web interface.",
				cfg.WebServer.PublicAddr, ctx.GetGuild().ID, victim.User.ID),
		}
		reps, err := db.GetReportsFiltered(ctx.GetGuild().ID, victim.User.ID, -1, 0, 1000)
		if err != nil {
			return err
		}
		if len(reps) == 0 {
			emb.Description += "\n\nThis user has a white west. :ok_hand:"
		} else {
			emb.Fields = make([]*discordgo.MessageEmbedField, 0)
			for _, r := range reps {
				emb.Fields = append(emb.Fields, r.AsEmbedField(cfg.WebServer.PublicAddr))
			}
		}
		_, err = ctx.GetSession().ChannelMessageSendEmbed(ctx.GetChannel().ID, emb)
		return err
	}

	msgOffset := 1
	repType, err := models.TypeFromString(ctx.GetArgs().Get(1).AsString())
	if repType == 0 {
		repType = models.TypesReserved
	}

	if victim.User.ID == ctx.GetUser().ID {
		return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
			"You can not report yourself...").
			DeleteAfter(8 * time.Second).Error()
	}

	if err == nil {
		if repType < models.TypesReserved || repType > models.TypeMax {
			return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
				fmt.Sprintf("Report type must be between *(including)* %d and %d.\n", models.TypesReserved, models.TypeMax)+
					"Use `help report` to get all types of report which can be used.").
				DeleteAfter(8 * time.Second).Error()
		}
		msgOffset++
	}

	repMsgS := ctx.GetArgs()[msgOffset:]

	if len(repMsgS) < 1 {
		return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
			"Please enter a valid report description.").
			DeleteAfter(8 * time.Second).Error()
	}

	repMsg := strings.Join(repMsgS, " ")

	var attachment string
	repMsg, attachment = imgstore.ExtractFromMessage(repMsg, ctx.GetMessage().Attachments)
	if attachment != "" {
		img, err := imgstore.DownloadFromURL(attachment)
		if err == nil && img != nil {
			st, _ := ctx.GetObject(static.DiObjectStorage).(storage.Storage)
			err = st.PutObject(static.StorageBucketImages, img.ID.String(),
				bytes.NewReader(img.Data), int64(img.Size), img.MimeType)
			if err != nil {
				return err
			}
			attachment = img.ID.String()
		}
	}

	rep := &models.Report{
		GuildID:       ctx.GetGuild().ID,
		ExecutorID:    ctx.GetUser().ID,
		VictimID:      victim.User.ID,
		Msg:           repMsg,
		AttachmehtURL: attachment,
		Type:          repType,
	}

	emb := rep.AsEmbed(cfg.WebServer.PublicAddr)
	emb.Title = "Report Check"
	emb.Description = "Is everything okay so far?"

	acceptMsg := acceptmsg.AcceptMessage{
		Embed:          emb,
		Session:        ctx.GetSession(),
		UserID:         ctx.GetUser().ID,
		DeleteMsgAfter: true,
		AcceptFunc: func(msg *discordgo.Message) (err error) {
			rep, err := repSvc.PushReport(rep)

			if err != nil {
				return
			}

			_, err = ctx.GetSession().ChannelMessageSendEmbed(ctx.GetChannel().ID, rep.AsEmbed(cfg.WebServer.PublicAddr))
			return
		},
	}

	_, err = acceptMsg.Send(ctx.GetChannel().ID)

	return err
}

func (c *CmdReport) revoke(ctx shireikan.Context) error {
	db, _ := ctx.GetObject(static.DiDatabase).(database.Database)
	cfg, _ := ctx.GetObject(static.DiConfig).(*config.Config)
	repSvc, _ := ctx.GetObject(static.DiReport).(*report.ReportService)

	if len(ctx.GetArgs()) < 3 {
		return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
			"Invalid command arguments. Please use `help report` for more information.").
			DeleteAfter(8 * time.Second).Error()
	}

	id, err := strconv.Atoi(ctx.GetArgs().Get(1).AsString())
	if err != nil {
		return err
	}

	reason := strings.Join(ctx.GetArgs()[2:], " ")

	rep, err := db.GetReport(snowflake.ID(id))
	if err != nil {
		if database.IsErrDatabaseNotFound(err) {
			return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
				fmt.Sprintf("Could not find any report with ID `%d`", id)).
				DeleteAfter(8 * time.Second).Error()
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
				rep.AsEmbedField(cfg.WebServer.PublicAddr),
			},
		},
		Session:        ctx.GetSession(),
		DeleteMsgAfter: true,
		UserID:         ctx.GetUser().ID,
		DeclineFunc: func(m *discordgo.Message) (err error) {
			return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
				"Canceled.").
				DeleteAfter(8 * time.Second).Error()
		},
		AcceptFunc: func(m *discordgo.Message) (err error) {
			emb, err := repSvc.RevokeReport(
				rep,
				ctx.GetUser().ID,
				reason,
				cfg.WebServer.PublicAddr,
				db,
				ctx.GetSession())

			if err != nil {
				return
			}

			_, err = ctx.GetSession().ChannelMessageSendEmbed(ctx.GetChannel().ID, emb)
			return
		},
	}

	if _, err = aceptMsg.Send(ctx.GetChannel().ID); err != nil {
		return err
	}

	return aceptMsg.Error()
}
