package commands

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/core/config"
	"github.com/zekroTJA/shinpuru/internal/core/database"
	"github.com/zekroTJA/shinpuru/internal/core/storage"
	"github.com/zekroTJA/shinpuru/internal/shared"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/imgstore"
	"github.com/zekroTJA/shinpuru/internal/util/report"
	"github.com/zekroTJA/shinpuru/internal/util/snowflakenodes"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/acceptmsg"
	"github.com/zekroTJA/shinpuru/pkg/fetch"
	"github.com/zekroTJA/shinpuru/pkg/roleutil"
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
	return "`ban <UserResolvable> <Reason>`"
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

	authorMemb, err := ctx.GetSession().GuildMember(ctx.GetGuild().ID, ctx.GetUser().ID)
	if err != nil {
		return err
	}

	if roleutil.PositionDiff(victim, authorMemb, ctx.GetGuild()) >= 0 {
		return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
			"You can only ban members with lower permissions than yours.").
			DeleteAfter(8 * time.Second).Error()
	}

	repMsg := strings.Join(ctx.GetArgs()[1:], " ")
	var repType int
	for i, v := range report.ReportTypes {
		if v == "BAN" {
			repType = i
		}
	}
	repID := snowflakenodes.NodesReport[repType].Generate()

	var attachment string
	repMsg, attachment = imgstore.ExtractFromMessage(repMsg, ctx.GetMessage().Attachments)
	if attachment != "" {
		img, err := imgstore.DownloadFromURL(attachment)
		if err == nil && img != nil {
			st, _ := ctx.GetObject("storage").(storage.Storage)
			err = st.PutObject(static.StorageBucketImages, img.ID.String(),
				bytes.NewReader(img.Data), int64(img.Size), img.MimeType)
			if err != nil {
				return err
			}
			attachment = img.ID.String()
		}
	}

	cfg, _ := ctx.GetObject("config").(*config.Config)
	db, _ := ctx.GetObject("db").(database.Database)

	acceptMsg := acceptmsg.AcceptMessage{
		Embed: &discordgo.MessageEmbed{
			Color:       report.ReportColors[repType],
			Title:       "Ban Check",
			Description: "Is everything okay so far?",
			Fields: []*discordgo.MessageEmbedField{
				{
					Name: "Victim",
					Value: fmt.Sprintf("<@%s> (%s#%s)",
						victim.User.ID, victim.User.Username, victim.User.Discriminator),
				},
				{
					Name:  "ID",
					Value: repID.String(),
				},
				{
					Name:  "Type",
					Value: report.ReportTypes[repType],
				},
				{
					Name:  "Description",
					Value: repMsg,
				},
			},
			Image: &discordgo.MessageEmbedImage{
				URL: imgstore.GetLink(attachment, cfg.WebServer.PublicAddr),
			},
		},
		Session:        ctx.GetSession(),
		UserID:         ctx.GetUser().ID,
		DeleteMsgAfter: true,
		AcceptFunc: func(msg *discordgo.Message) {
			rep, err := shared.PushBan(
				ctx.GetSession(),
				db,
				cfg.WebServer.PublicAddr,
				ctx.GetGuild().ID,
				ctx.GetUser().ID,
				victim.User.ID,
				repMsg,
				attachment)

			if err != nil {
				util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
					"Failed banning member: ```\n"+err.Error()+"\n```")
			} else {
				ctx.GetSession().ChannelMessageSendEmbed(ctx.GetChannel().ID, rep.AsEmbed(cfg.WebServer.PublicAddr))
			}
		},
	}

	_, err = acceptMsg.Send(ctx.GetChannel().ID)

	return err
}
