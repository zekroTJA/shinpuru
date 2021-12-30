package slashcommands

import (
	"bytes"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/models"
	"github.com/zekroTJA/shinpuru/internal/services/config"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/services/permissions"
	"github.com/zekroTJA/shinpuru/internal/services/report"
	"github.com/zekroTJA/shinpuru/internal/services/storage"
	"github.com/zekroTJA/shinpuru/internal/util/imgstore"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekrotja/dgrs"
	"github.com/zekrotja/ken"
)

type Mute struct{}

var (
	_ ken.SlashCommand        = (*Mute)(nil)
	_ permissions.PermCommand = (*Mute)(nil)
)

func (c *Mute) Name() string {
	return "mute"
}

func (c *Mute) Description() string {
	return "Mute members or setup mute."
}

func (c *Mute) Version() string {
	return "2.0.0"
}

func (c *Mute) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (c *Mute) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "toggle",
			Description: "Toggle mute/unmute state of a member.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user",
					Description: "The user to be muted/unmuted.",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "reason",
					Description: "The mute reason.",
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "imageurl",
					Description: "Image attachment URL.",
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "expire",
					Description: "Expiration time.",
				},
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "list",
			Description: "List muted members.",
		},
	}
}

func (c *Mute) Domain() string {
	return "sp.guild.mod.mute"
}

func (c *Mute) SubDomains() []permissions.SubPermission {
	return nil
}

func (c *Mute) Run(ctx *ken.Ctx) (err error) {
	if err = ctx.Defer(); err != nil {
		return
	}

	err = ctx.HandleSubCommands(
		ken.SubCommandHandler{"toggle", c.toggle},
		ken.SubCommandHandler{"list", c.list},
	)

	return
}

func (c *Mute) toggle(ctx *ken.SubCommandCtx) (err error) {
	victim := ctx.Options().GetByName("user").UserValue(ctx.Ctx)

	var reason string
	if reasonV, ok := ctx.Options().GetByNameOptional("reason"); ok {
		reason = reasonV.StringValue()
	}

	if victim.ID == ctx.User().ID {
		return ctx.FollowUpError(
			"You can not mute yourself...", "").
			Error
	}

	st := ctx.Get(static.DiState).(*dgrs.State)

	// TODO: forcefetch is set to true because dgrs does not
	//       track member timeout states at the moment.
	member, err := st.Member(ctx.Event.GuildID, victim.ID, true)
	if err != nil {
		return
	}

	cfg := ctx.Get(static.DiConfig).(config.Provider)
	repSvc := ctx.Get(static.DiReport).(*report.ReportService)

	if member.CommunicationDisabledUntil != nil {
		emb, err := repSvc.RevokeMute(
			ctx.Event.GuildID,
			ctx.User().ID,
			victim.ID,
			reason)
		if err != nil {
			return err
		}

		return ctx.FollowUpEmbed(emb).Error
	}

	if len(reason) == 0 {
		return ctx.FollowUpError(
			"Please enter a valid report description.", "").
			Error
	}

	var attachment string
	if imageurl, ok := ctx.Options().GetByNameOptional("imageurl"); ok {
		img, err := imgstore.DownloadFromURL(imageurl.StringValue())
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
	}

	expireV, ok := ctx.Options().GetByNameOptional("expire")
	if !ok {
		return ctx.FollowUpError(
			"Please enter a valid timeout.", "").
			Error
	}
	expire, err := time.ParseDuration(expireV.StringValue())
	if err != nil {
		return ctx.FollowUpError(
			fmt.Sprintf("Invalid expire value:\n```\n%s```", err.Error()), "").Error
	}
	expireTime := time.Now().Add(expire)
	rep.Timeout = &expireTime

	rep, err = repSvc.PushMute(rep)
	if err != nil {
		err = ctx.FollowUpError(
			"Failed creating report: ```\n"+err.Error()+"\n```", "").
			Error
	} else {
		err = ctx.FollowUpEmbed(rep.AsEmbed(cfg.Config().WebServer.PublicAddr)).Error
	}

	return err
}

func (c *Mute) list(ctx *ken.SubCommandCtx) (err error) {
	db := ctx.Get(static.DiDatabase).(database.Database)

	emb := &discordgo.MessageEmbed{
		Color:       static.ColorEmbedGray,
		Description: "Fetching muted members...",
		Fields:      make([]*discordgo.MessageEmbedField, 0),
	}

	fum := ctx.FollowUpEmbed(emb)
	err = fum.Error
	if err != nil {
		return err
	}

	muteReports, err := db.GetReportsFiltered(ctx.Event.GuildID, "",
		int(models.TypeMute), 0, 1000)

	muteReportsMap := make(map[string]*models.Report)
	for _, r := range muteReports {
		muteReportsMap[r.VictimID] = r
	}

	st := ctx.Get(static.DiState).(*dgrs.State)
	membs, err := st.Members(ctx.Event.GuildID)
	if err != nil {
		return err
	}
	for _, m := range membs {
		if m.CommunicationDisabledUntil != nil {
			if r, ok := muteReportsMap[m.User.ID]; ok {
				emb.Fields = append(emb.Fields, &discordgo.MessageEmbedField{
					Name: fmt.Sprintf("CaseID: %d", r.ID),
					Value: fmt.Sprintf("<@%s> since `%s` with reason:\n%s",
						m.User.ID, r.GetTimestamp().Format(time.RFC1123), r.Msg),
				})
			}
		}
	}

	emb.Color = static.ColorEmbedDefault
	emb.Description = ""

	err = fum.Edit(&discordgo.WebhookEdit{
		Embeds: []*discordgo.MessageEmbed{emb},
	})
	return
}
