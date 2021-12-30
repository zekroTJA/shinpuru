package commands

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/zekroTJA/shinpuru/internal/models"
	"github.com/zekroTJA/shinpuru/internal/services/config"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/services/report"
	"github.com/zekroTJA/shinpuru/internal/services/storage"
	"github.com/zekroTJA/shinpuru/internal/util/imgstore"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/fetch"
	"github.com/zekroTJA/shireikan"
	"github.com/zekrotja/dgrs"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/util"
)

type CmdMute struct {
}

func (c *CmdMute) GetInvokes() []string {
	return []string{"mute", "m", "silence", "unmute", "um", "unsilence"}
}

func (c *CmdMute) GetDescription() string {
	return "Mute members in text channels."
}

func (c *CmdMute) GetHelp() string {
	return "`mute setup (<roleResolvable>)` - creates (or uses given) mute role and sets this role in every channel as muted\n" +
		"`mute <userResolvable> (<timeout duration>)` - mute/unmute a user\n" +
		"`mute list` - display muted users on this guild\n" +
		"`mute` - display currently set mute role"
}

func (c *CmdMute) GetGroup() string {
	return shireikan.GroupModeration
}

func (c *CmdMute) GetDomainName() string {
	return "sp.guild.mod.mute"
}

func (c *CmdMute) GetSubPermissionRules() []shireikan.SubPermission {
	return nil
}

func (c *CmdMute) IsExecutableInDMChannels() bool {
	return false
}

func (c *CmdMute) Exec(ctx shireikan.Context) error {
	if len(ctx.GetArgs()) < 1 {
		return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
			"Wrong command usage. Please use `help mute` to get more information.").
			DeleteAfter(8 * time.Second).Error()
	}

	switch ctx.GetArgs().Get(0).AsString() {
	case "list":
		return c.list(ctx)
	default:
		return c.muteUnmute(ctx)
	}
}

func (c *CmdMute) muteUnmute(ctx shireikan.Context) error {
	victim, err := fetch.FetchMember(ctx.GetSession(), ctx.GetGuild().ID, ctx.GetArgs().Get(0).AsString())
	if err != nil {
		return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
			"Could not fetch any user by the passed resolvable.").
			DeleteAfter(8 * time.Second).Error()
	}

	if victim.User.ID == ctx.GetUser().ID {
		return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
			"You can not mute yourself...").
			DeleteAfter(8 * time.Second).Error()
	}

	cfg, _ := ctx.GetObject(static.DiConfig).(config.Provider)
	repSvc, _ := ctx.GetObject(static.DiReport).(*report.ReportService)

	if victim.CommunicationDisabledUntil != nil {
		emb, err := repSvc.RevokeMute(
			ctx.GetGuild().ID,
			ctx.GetUser().ID,
			victim.User.ID,
			strings.Join(ctx.GetArgs()[1:], " "))
		if err != nil {
			return err
		}

		_, err = ctx.GetSession().ChannelMessageSendEmbed(ctx.GetChannel().ID, emb)

		return err
	}

	repMsgS := ctx.GetArgs()[1:]

	timeout, err := time.ParseDuration(repMsgS[len(repMsgS)-1])
	if err == nil && timeout > 0 {
		repMsgS = repMsgS[:len(repMsgS)-1]
	}
	if err != nil {
		return err
	}
	if timeout == 0 {
		return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
			"Please enter a valid mute timeout.").
			DeleteAfter(8 * time.Second).Error()
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

	rep, err := repSvc.PushMute(&models.Report{
		GuildID:       ctx.GetGuild().ID,
		ExecutorID:    ctx.GetUser().ID,
		VictimID:      victim.User.ID,
		Msg:           strings.Join(ctx.GetArgs()[1:], " "),
		AttachmehtURL: attachment,
	})

	if err != nil {
		err = util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
			"Failed creating report: ```\n"+err.Error()+"\n```").
			Error()
	} else {
		_, err = ctx.GetSession().ChannelMessageSendEmbed(ctx.GetChannel().ID, rep.AsEmbed(cfg.Config().WebServer.PublicAddr))
	}

	return err
}

func (c *CmdMute) list(ctx shireikan.Context) error {
	db, _ := ctx.GetObject(static.DiDatabase).(database.Database)

	emb := &discordgo.MessageEmbed{
		Color:       static.ColorEmbedGray,
		Description: "Fetching muted members...",
		Fields:      make([]*discordgo.MessageEmbedField, 0),
	}

	msg, err := ctx.GetSession().ChannelMessageSendEmbed(ctx.GetChannel().ID, emb)
	if err != nil {
		return err
	}

	muteReports, err := db.GetReportsFiltered(ctx.GetGuild().ID, "",
		int(models.TypeMute), 0, 1000)

	muteReportsMap := make(map[string]*models.Report)
	for _, r := range muteReports {
		muteReportsMap[r.VictimID] = r
	}

	st := ctx.GetObject(static.DiState).(*dgrs.State)
	membs, err := st.Members(ctx.GetGuild().ID)
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

	_, err = ctx.GetSession().ChannelMessageEditEmbed(ctx.GetChannel().ID, msg.ID, emb)
	return err
}
