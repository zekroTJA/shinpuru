package commands

import (
	"bytes"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/config"
	"github.com/zekroTJA/shinpuru/internal/models"
	"github.com/zekroTJA/shinpuru/internal/services/report"
	"github.com/zekroTJA/shinpuru/internal/services/storage"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/imgstore"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/acceptmsg"
	"github.com/zekroTJA/shinpuru/pkg/fetch"
	"github.com/zekroTJA/shinpuru/pkg/timeutil"
	"github.com/zekroTJA/shireikan"
)

type CmdBan struct {
}

func (c *CmdBan) GetInvokes() []string {
	return []string{"ban", "userban"}
}

func (c *CmdBan) GetDescription() string {
	return "Ban users with creating a report entry."
}

func (c *CmdBan) GetHelp() string {
	return "`ban <UserResolvable> <Reason> (<timeout duration>)`"
}

func (c *CmdBan) GetGroup() string {
	return shireikan.GroupModeration
}

func (c *CmdBan) GetDomainName() string {
	return "sp.guild.mod.ban"
}

func (c *CmdBan) GetSubPermissionRules() []shireikan.SubPermission {
	return nil
}

func (c *CmdBan) IsExecutableInDMChannels() bool {
	return false
}

func (c *CmdBan) Exec(ctx shireikan.Context) error {
	if len(ctx.GetArgs()) < 2 {
		return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
			"Invalid command arguments. Please use `help ban` to see how to use this command.").
			DeleteAfter(8 * time.Second).Error()
	}
	victim, err := fetch.FetchMember(ctx.GetSession(), ctx.GetGuild().ID, ctx.GetArgs().Get(0).AsString())
	if err != nil || victim == nil {
		return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
			"Sorry, could not find any member :cry:").
			DeleteAfter(10 * time.Second).Error()
	}

	if victim.User.ID == ctx.GetUser().ID {
		return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
			"You can not ban yourself...").
			DeleteAfter(8 * time.Second).Error()
	}

	repMsgS := ctx.GetArgs()[1:]

	timeout, err := time.ParseDuration(repMsgS[len(repMsgS)-1])
	if err == nil && timeout > 0 {
		repMsgS = repMsgS[:len(repMsgS)-1]
	}

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

	cfg, _ := ctx.GetObject(static.DiConfig).(*config.Config)
	repSvc, _ := ctx.GetObject(static.DiReport).(*report.ReportService)

	rep := &models.Report{
		GuildID:       ctx.GetGuild().ID,
		ExecutorID:    ctx.GetUser().ID,
		VictimID:      victim.User.ID,
		Msg:           repMsg,
		AttachmehtURL: attachment,
		Timeout:       timeutil.NowAddPtr(timeout),
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
			rep, err := repSvc.PushBan(rep)

			if err != nil {
				return
			}
			_, err = ctx.GetSession().ChannelMessageSendEmbed(ctx.GetChannel().ID, rep.AsEmbed(cfg.WebServer.PublicAddr))
			return
		},
	}

	if _, err = acceptMsg.Send(ctx.GetChannel().ID); err != nil {
		return err
	}

	return acceptMsg.Error()
}
