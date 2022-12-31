package cmdutil

import (
	"bytes"
	"fmt"

	"github.com/zekroTJA/shinpuru/internal/models"
	"github.com/zekroTJA/shinpuru/internal/services/config"
	"github.com/zekroTJA/shinpuru/internal/services/report"
	"github.com/zekroTJA/shinpuru/internal/services/storage"
	"github.com/zekroTJA/shinpuru/internal/services/timeprovider"
	"github.com/zekroTJA/shinpuru/internal/util/imgstore"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/acceptmsg/v2"
	"github.com/zekroTJA/shinpuru/pkg/timeutil"
	"github.com/zekrotja/ken"
)

func CmdReport(ctx ken.Context, typ models.ReportType, tp timeprovider.Provider) (err error) {
	cfg := ctx.Get(static.DiConfig).(config.Provider)
	repSvc := ctx.Get(static.DiReport).(report.Provider)

	victim := ctx.Options().GetByName("user").UserValue(ctx)
	reason := ctx.Options().GetByName("reason").StringValue()

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

	rep := models.Report{
		GuildID:       ctx.GetEvent().GuildID,
		ExecutorID:    ctx.User().ID,
		VictimID:      victim.ID,
		Msg:           reason,
		AttachmentURL: attachment,
		Type:          typ,
	}

	if expire != "" {
		exp, err := timeutil.ParseDuration(expire)
		if err != nil {
			err = ctx.FollowUpError(
				fmt.Sprintf("Invalid duration:\n```\n%s```", err.Error()), "").
				Send().Error
			return err
		}
		expT := tp.Now().Add(exp)
		rep.Timeout = &expT
	}

	emb := rep.AsEmbed(cfg.Config().WebServer.PublicAddr)
	emb.Title = "Report Check"
	emb.Description = "Is everything okay so far?"

	acceptMsg := acceptmsg.AcceptMessage{
		Embed:          emb,
		Ken:            ctx.GetKen(),
		UserID:         ctx.User().ID,
		DeleteMsgAfter: true,
		AcceptFunc: func(cctx ken.ComponentContext) (err error) {
			if err = cctx.Defer(); err != nil {
				return
			}

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

			return cctx.FollowUpEmbed(
				rep.AsEmbed(cfg.Config().WebServer.PublicAddr)).
				Send().Error
		},
	}

	if _, err = acceptMsg.AsFollowUp(ctx); err != nil {
		return
	}
	return acceptMsg.Error()
}
