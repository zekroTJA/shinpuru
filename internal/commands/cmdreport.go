package commands

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/zekroTJA/shinpuru/internal/core/database"
	"github.com/zekroTJA/shinpuru/internal/shared"
	"github.com/zekroTJA/shinpuru/internal/util/imgstore"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/acceptmsg"
	"github.com/zekroTJA/shinpuru/pkg/fetch"

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
	return "report a user"
}

func (c *CmdReport) GetHelp() string {
	repTypes := make([]string, len(static.ReportTypes))
	for i, t := range static.ReportTypes {
		repTypes[i] = fmt.Sprintf("`%d` - %s", i, t)
	}
	return "`report <userResolvable>` - list all reports of a user\n" +
		"`report <userResolvable> [<type>] <reason>` - report a user *(if type is empty, its defaultly 0 = warn)*\n" +
		"`report revoke <caseID> <reason>` - revoke a report\n" +
		"\n**TYPES:**\n" + strings.Join(repTypes, "\n") +
		"\nTypes `BAN`, `KICK` and `MUTE` are reserved for bands and kicks executed with this bot."
}

func (c *CmdReport) GetGroup() string {
	return GroupModeration
}

func (c *CmdReport) GetDomainName() string {
	return "sp.guild.mod.report"
}

func (c *CmdReport) GetSubPermissionRules() []SubPermission {
	return nil
}

func (c *CmdReport) Exec(args *CommandArgs) error {
	if len(args.Args) < 1 {
		return util.SendEmbedError(args.Session, args.Channel.ID,
			"Invalid command arguments. Please use `help report` to see how to use this command.").
			DeleteAfter(8 * time.Second).Error()
	}

	if strings.ToLower(args.Args[0]) == "revoke" {
		return c.revoke(args)
	}

	victim, err := fetch.FetchMember(args.Session, args.Guild.ID, args.Args[0])
	if err != nil || victim == nil {
		return util.SendEmbedError(args.Session, args.Channel.ID,
			"Sorry, could not find any member :cry:").
			DeleteAfter(8 * time.Second).Error()
	}

	if len(args.Args) == 1 {
		emb := &discordgo.MessageEmbed{
			Color: static.ColorEmbedDefault,
			Title: fmt.Sprintf("Reports for %s#%s",
				victim.User.Username, victim.User.Discriminator),
			Description: fmt.Sprintf("[**Here**](%s/guilds/%s/%s) you can find this users reports in the web interface.",
				args.CmdHandler.config.WebServer.PublicAddr, args.Guild.ID, victim.User.ID),
		}
		reps, err := args.CmdHandler.db.GetReportsFiltered(args.Guild.ID, victim.User.ID, -1)
		if err != nil {
			return err
		}
		if len(reps) == 0 {
			emb.Description += "\n\nThis user has a white west. :ok_hand:"
		} else {
			emb.Fields = make([]*discordgo.MessageEmbedField, 0)
			for _, r := range reps {
				emb.Fields = append(emb.Fields, r.AsEmbedField(args.CmdHandler.config.WebServer.PublicAddr))
			}
		}
		_, err = args.Session.ChannelMessageSendEmbed(args.Channel.ID, emb)
		return err
	}

	msgOffset := 1
	repType, err := strconv.Atoi(args.Args[1])
	maxType := len(static.ReportTypes) - 1
	minType := static.ReportTypesReserved
	if repType == 0 {
		repType = minType
	}

	if victim.User.ID == args.User.ID {
		return util.SendEmbedError(args.Session, args.Channel.ID,
			"You can not report yourself...").
			DeleteAfter(8 * time.Second).Error()
	}

	if err == nil {
		if repType < minType || repType > maxType {
			return util.SendEmbedError(args.Session, args.Channel.ID,
				fmt.Sprintf("Report type must be between *(including)* %d and %d.\n", minType, maxType)+
					"Use `help report` to get all types of report which can be used.").
				DeleteAfter(8 * time.Second).Error()
		}
		msgOffset++
	}

	if len(args.Args[msgOffset:]) < 1 {
		return util.SendEmbedError(args.Session, args.Channel.ID,
			"Please enter a valid report description.").
			DeleteAfter(8 * time.Second).Error()
	}
	repMsg := strings.Join(args.Args[msgOffset:], " ")

	var attachment string
	repMsg, attachment = imgstore.ExtractFromMessage(repMsg, args.Message.Attachments)
	if attachment != "" {
		img, err := imgstore.DownloadFromURL(attachment)
		if err == nil && img != nil {
			err = args.CmdHandler.st.PutObject(static.StorageBucketImages, img.ID.String(),
				bytes.NewReader(img.Data), int64(img.Size), img.MimeType)
			if err != nil {
				return err
			}
			attachment = img.ID.String()
		}
	}

	acceptMsg := acceptmsg.AcceptMessage{
		Embed: &discordgo.MessageEmbed{
			Color:       static.ReportColors[repType],
			Title:       "Report Check",
			Description: "Is everything okay so far?",
			Fields: []*discordgo.MessageEmbedField{
				{
					Name: "Victim",
					Value: fmt.Sprintf("<@%s> (%s#%s)",
						victim.User.ID, victim.User.Username, victim.User.Discriminator),
				},
				{
					Name:  "Type",
					Value: static.ReportTypes[repType],
				},
				{
					Name:  "Description",
					Value: repMsg,
				},
			},
			Image: &discordgo.MessageEmbedImage{
				URL: imgstore.GetLink(attachment, args.CmdHandler.config.WebServer.PublicAddr),
			},
		},
		Session:        args.Session,
		UserID:         args.User.ID,
		DeleteMsgAfter: true,
		AcceptFunc: func(msg *discordgo.Message) {
			rep, err := shared.PushReport(
				args.Session,
				args.CmdHandler.db,
				args.CmdHandler.config.WebServer.PublicAddr,
				args.Guild.ID,
				args.User.ID,
				victim.User.ID,
				repMsg,
				attachment,
				repType)

			if err != nil {
				util.SendEmbedError(args.Session, args.Channel.ID,
					"Failed creating report: ```\n"+err.Error()+"\n```")
			} else {
				args.Session.ChannelMessageSendEmbed(args.Channel.ID, rep.AsEmbed(args.CmdHandler.config.WebServer.PublicAddr))
			}
		},
	}

	_, err = acceptMsg.Send(args.Channel.ID)

	return err
}

func (c *CmdReport) revoke(args *CommandArgs) error {
	if len(args.Args) < 3 {
		return util.SendEmbedError(args.Session, args.Channel.ID,
			"Invalid command arguments. Please use `help report` for more information.").
			DeleteAfter(8 * time.Second).Error()
	}

	id, err := strconv.Atoi(args.Args[1])
	if err != nil {
		return err
	}

	reason := strings.Join(args.Args[2:], " ")

	rep, err := args.CmdHandler.db.GetReport(snowflake.ID(id))
	if err != nil {
		if database.IsErrDatabaseNotFound(err) {
			return util.SendEmbedError(args.Session, args.Channel.ID,
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
				rep.AsEmbedField(args.CmdHandler.config.WebServer.PublicAddr),
			},
		},
		Session:        args.Session,
		DeleteMsgAfter: true,
		UserID:         args.User.ID,
		DeclineFunc: func(m *discordgo.Message) {
			util.SendEmbedError(args.Session, args.Channel.ID,
				"Canceled.").
				DeleteAfter(8 * time.Second)
		},
		AcceptFunc: func(m *discordgo.Message) {
			err := args.CmdHandler.db.DeleteReport(rep.ID)
			if err != nil {
				util.SendEmbedError(args.Session, args.Channel.ID,
					fmt.Sprintf("An error occured while deleting report from database: ```\n%s\n```", err.Error()))
				return
			}

			repRevEmb := &discordgo.MessageEmbed{
				Color:       static.ReportRevokedColor,
				Title:       "REPORT REVOCATION",
				Description: "Revoked reports are deleted from the database and no more visible in any commands.",
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:  "Revoke Executor",
						Value: args.User.Mention(),
					},
					{
						Name:  "Revocation Reason",
						Value: reason,
					},
					rep.AsEmbedField(args.CmdHandler.config.WebServer.PublicAddr),
				},
			}

			args.Session.ChannelMessageSendEmbed(args.Channel.ID, repRevEmb)
			if modlogChan, err := args.CmdHandler.db.GetGuildModLog(args.Guild.ID); err == nil {
				args.Session.ChannelMessageSendEmbed(modlogChan, repRevEmb)
			}
			dmChan, err := args.Session.UserChannelCreate(rep.VictimID)
			if err == nil {
				args.Session.ChannelMessageSendEmbed(dmChan.ID, repRevEmb)
			}
		},
	}

	_, err = aceptMsg.Send(args.Channel.ID)
	return err
}
