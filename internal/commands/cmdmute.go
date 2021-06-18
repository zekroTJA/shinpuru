package commands

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/zekroTJA/shinpuru/internal/config"
	"github.com/zekroTJA/shinpuru/internal/models"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/services/report"
	"github.com/zekroTJA/shinpuru/internal/services/storage"
	"github.com/zekroTJA/shinpuru/internal/util/imgstore"
	"github.com/zekroTJA/shinpuru/internal/util/mute"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/acceptmsg"
	"github.com/zekroTJA/shinpuru/pkg/fetch"
	"github.com/zekroTJA/shinpuru/pkg/stringutil"
	"github.com/zekroTJA/shireikan"

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
		"`mute <userResolvable>` - mute/unmute a user\n" +
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
		return c.displayMuteRole(ctx)
	}

	switch ctx.GetArgs().Get(0).AsString() {
	case "setup":
		return c.setup(ctx)
	case "list":
		return c.list(ctx)
	default:
		return c.muteUnmute(ctx)
	}
}

func (c *CmdMute) setup(ctx shireikan.Context) error {
	var muteRole *discordgo.Role
	var err error

	desc := "Following, a rolen with the name `shinpuru-muted` will be created *(if not existend yet)* and set as mute role."

	if len(ctx.GetArgs()) > 1 {
		muteRole, err = fetch.FetchRole(ctx.GetSession(), ctx.GetGuild().ID, ctx.GetArgs().Get(1).AsString())
		if err != nil {
			return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
				"Role could not be fetched by passed identifier.").
				DeleteAfter(8 * time.Second).Error()
		}

		desc = fmt.Sprintf("Follwoing, the role %s will be set as mute role.", muteRole.Mention())
	}

	db, _ := ctx.GetObject(static.DiDatabase).(database.Database)

	acmsg := &acceptmsg.AcceptMessage{
		Session: ctx.GetSession(),
		Embed: &discordgo.MessageEmbed{
			Color: static.ColorEmbedDefault,
			Title: "Warning",
			Description: desc + " Also, all channels *(which the bot has access to)* will be permission-overwritten that " +
				"members with this role will not be able to write in these channels anymore.",
		},
		UserID:         ctx.GetUser().ID,
		DeleteMsgAfter: true,
		AcceptFunc: func(msg *discordgo.Message) {
			if muteRole == nil {
				for _, r := range ctx.GetGuild().Roles {
					if r.Name == static.MutedRoleName {
						muteRole = r
					}
				}
			}

			if muteRole == nil {
				muteRole, err = ctx.GetSession().GuildRoleCreate(ctx.GetGuild().ID)
				if err != nil {
					util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
						"Failed creating mute role: ```\n"+err.Error()+"\n```").
						DeleteAfter(15 * time.Second).Error()
					return
				}

				muteRole, err = ctx.GetSession().GuildRoleEdit(ctx.GetGuild().ID, muteRole.ID,
					static.MutedRoleName, 0, false, 0, false)
				if err != nil {
					util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
						"Failed editing mute role: ```\n"+err.Error()+"\n```").
						DeleteAfter(15 * time.Second).Error()
					return
				}
			}

			err := db.SetGuildMuteRole(ctx.GetGuild().ID, muteRole.ID)
			if err != nil {
				util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
					"Failed setting mute role in database: ```\n"+err.Error()+"\n```").
					DeleteAfter(15 * time.Second).Error()
				return
			}

			err = mute.SetupChannels(ctx.GetSession(), ctx.GetGuild().ID, muteRole.ID)
			if err != nil {
				util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
					"Failed updating channels: ```\n"+err.Error()+"\n```").
					DeleteAfter(15 * time.Second).Error()
				return
			}

			util.SendEmbed(ctx.GetSession(), ctx.GetChannel().ID,
				"Set up mute role and edited channel permissions.\nMaybe you need to increase the "+
					"position of the role to override other roles permission settings.",
				"", static.ColorEmbedUpdated).
				DeleteAfter(15 * time.Second).Error()
		},
		DeclineFunc: func(msg *discordgo.Message) {
			util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
				"Setup canceled.").
				DeleteAfter(8 * time.Second).Error()
		},
	}

	_, err = acmsg.Send(ctx.GetChannel().ID)
	return err
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

	db, _ := ctx.GetObject(static.DiDatabase).(database.Database)

	muteRoleID, err := db.GetGuildMuteRole(ctx.GetGuild().ID)
	if database.IsErrDatabaseNotFound(err) {
		return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
			"Mute command is not set up. Please enter the command `mute setup`.").
			DeleteAfter(8 * time.Second).Error()
	} else if err != nil {
		return err
	}

	var roleExists bool
	for _, r := range ctx.GetGuild().Roles {
		if r.ID == muteRoleID && !roleExists {
			roleExists = true
		}
	}
	if !roleExists {
		return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
			"Mute role does not exist on this guild. Please enter `mute setup`.").
			DeleteAfter(8 * time.Second).Error()
	}

	var victimIsMuted bool
	for _, rID := range victim.Roles {
		if rID == muteRoleID && !victimIsMuted {
			victimIsMuted = true
		}
	}

	cfg, _ := ctx.GetObject(static.DiConfig).(*config.Config)
	repSvc, _ := ctx.GetObject(static.DiReport).(*report.ReportService)

	if victimIsMuted {
		emb, err := repSvc.RevokeMute(
			ctx.GetGuild().ID,
			ctx.GetUser().ID,
			victim.User.ID,
			strings.Join(ctx.GetArgs()[1:], " "),
			muteRoleID)
		if err != nil {
			return err
		}

		_, err = ctx.GetSession().ChannelMessageSendEmbed(ctx.GetChannel().ID, emb)

		return err
	}

	err = ctx.GetSession().GuildMemberRoleAdd(ctx.GetGuild().ID, victim.User.ID, muteRoleID)
	if err != nil {
		return err
	}

	repMsg := strings.Join(ctx.GetArgs()[1:], " ")

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

	rep, err := repSvc.PushMute(
		ctx.GetGuild().ID,
		ctx.GetUser().ID,
		victim.User.ID,
		strings.Join(ctx.GetArgs()[1:], " "),
		attachment,
		muteRoleID)

	if err != nil {
		err = util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
			"Failed creating report: ```\n"+err.Error()+"\n```").
			Error()
	} else {
		_, err = ctx.GetSession().ChannelMessageSendEmbed(ctx.GetChannel().ID, rep.AsEmbed(cfg.WebServer.PublicAddr))
	}

	return err
}

func (c *CmdMute) list(ctx shireikan.Context) error {
	db, _ := ctx.GetObject(static.DiDatabase).(database.Database)

	muteRoleID, err := db.GetGuildMuteRole(ctx.GetGuild().ID)
	if err != nil {
		return err
	}

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
		stringutil.IndexOf("MUTE", models.ReportTypes))

	muteReportsMap := make(map[string]*models.Report)
	for _, r := range muteReports {
		muteReportsMap[r.VictimID] = r
	}

	for _, m := range ctx.GetGuild().Members {
		if stringutil.IndexOf(muteRoleID, m.Roles) > -1 {
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

func (c *CmdMute) displayMuteRole(ctx shireikan.Context) error {
	db, _ := ctx.GetObject(static.DiDatabase).(database.Database)

	roleID, err := db.GetGuildMuteRole(ctx.GetGuild().ID)
	if err != nil {
		return err
	}

	if roleID == "" {
		return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
			"Mute role is currently unset.").
			DeleteAfter(8 * time.Second).Error()
	}

	return util.SendEmbed(ctx.GetSession(), ctx.GetChannel().ID,
		fmt.Sprintf("Role <@&%s> is currently set as mute role.", roleID), "", 0).
		DeleteAfter(8 * time.Second).Error()
}
