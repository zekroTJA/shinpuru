package slashcommands

import (
	"bytes"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/models"
	"github.com/zekroTJA/shinpuru/internal/services/config"
	"github.com/zekroTJA/shinpuru/internal/services/permissions"
	"github.com/zekroTJA/shinpuru/internal/services/report"
	"github.com/zekroTJA/shinpuru/internal/services/storage"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/imgstore"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/acceptmsg"
	"github.com/zekroTJA/shinpuru/pkg/timeutil"
	"github.com/zekrotja/ken"
)

type Ban struct{}

var (
	_ ken.Command             = (*Ban)(nil)
	_ permissions.PermCommand = (*Ban)(nil)
)

func (c *Ban) Name() string {
	return "ban"
}

func (c *Ban) Description() string {
	return "Ban users with creating a report entry."
}

func (c *Ban) Version() string {
	return ""
}

func (c *Ban) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (c *Ban) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionUser,
			Name:        "user",
			Description: "User to ban",
			Required:    true,
		},
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "reason",
			Description: "Ban reason",
			Required:    true,
		},
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "attachment",
			Description: "Url of an image attachment (e.g. a screenshot)",
			Required:    false,
		},
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "timeout",
			Description: "duration of the ban",
			Required:    false,
		},
	}
}

func (c *Ban) Domain() string {
	return "sp.guild.mod.ban"
}

func (c *Ban) SubDomains() []permissions.SubPermission {
	return nil
}

func (c *Ban) Run(ctx *ken.Ctx) (err error) {
	if err = ctx.Defer(); err != nil {
		return
	}

	victim := ctx.Event.ApplicationCommandData().Options[0].UserValue(nil)

	if victim.ID == ctx.User().ID {
		return util.SendEmbedError(ctx.Session, ctx.Event.ChannelID,
			"You can not ban yourself...").
			DeleteAfter(8 * time.Second).Error()
	}

	reason := ctx.Event.ApplicationCommandData().Options[1].StringValue()

	var attachment string = ""
	var timeout time.Duration = 0

	//parse options 2 and 3 (attachment and timeout)
	for i := 2; i <= 3; i++ {
		if len(ctx.Event.ApplicationCommandData().Options) > i {
			switch ctx.Event.ApplicationCommandData().Options[i].Name {
			case "Attachment":
				attachment = ctx.Event.ApplicationCommandData().Options[i].StringValue()
			case "Timeout":
				timeout, _ = time.ParseDuration(ctx.Event.ApplicationCommandData().Options[i].StringValue())
			}
		}
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

	if strings.TrimSpace(reason) == "" {
		return util.SendEmbedError(ctx.Session, ctx.Event.ChannelID,
			"Please enter a valid report description.").
			DeleteAfter(8 * time.Second).Error()
	}

	cfg, _ := ctx.Get(static.DiConfig).(config.Provider)
	repSvc, _ := ctx.Get(static.DiReport).(*report.ReportService)

	rep := &models.Report{
		GuildID:       ctx.Event.GuildID,
		ExecutorID:    ctx.User().ID,
		VictimID:      victim.ID,
		Msg:           reason,
		AttachmehtURL: attachment,
		Timeout:       timeutil.NowAddPtr(timeout),
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
			rep, err := repSvc.PushBan(rep)

			if err != nil {
				return
			}
			_, err = ctx.Session.ChannelMessageSendEmbed(ctx.Event.ChannelID, rep.AsEmbed(cfg.Config().WebServer.PublicAddr))
			return
		},
	}

	if _, err = acceptMsg.AsFollowUp(ctx); err != nil {
		return err
	}

	return acceptMsg.Error()
}
