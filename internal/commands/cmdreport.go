package commands

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/zekroTJA/shinpuru/internal/core"

	"github.com/bwmarrin/discordgo"
	"github.com/bwmarrin/snowflake"

	"github.com/zekroTJA/shinpuru/internal/util"
)

type CmdReport struct {
	PermLvl int
}

func (c *CmdReport) GetInvokes() []string {
	return []string{"report", "rep", "warn"}
}

func (c *CmdReport) GetDescription() string {
	return "report a user"
}

func (c *CmdReport) GetHelp() string {
	repTypes := make([]string, len(util.ReportTypes))
	for i, t := range util.ReportTypes {
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

func (c *CmdReport) GetPermission() int {
	return c.PermLvl
}

func (c *CmdReport) SetPermission(permLvl int) {
	c.PermLvl = permLvl
}

func (c *CmdReport) Exec(args *CommandArgs) error {
	if len(args.Args) < 1 {
		msg, err := util.SendEmbedError(args.Session, args.Channel.ID,
			"Invalid command arguments. Please use `help report` to see how to use this command.")
		util.DeleteMessageLater(args.Session, msg, 8*time.Second)
		return err
	}

	if strings.ToLower(args.Args[0]) == "revoke" {
		return c.revoke(args)
	}

	victim, err := util.FetchMember(args.Session, args.Guild.ID, args.Args[0])
	if err != nil || victim == nil {
		msg, err := util.SendEmbedError(args.Session, args.Channel.ID,
			"Sorry, could not find any member :cry:")
		util.DeleteMessageLater(args.Session, msg, 10*time.Second)
		return err
	}

	if len(args.Args) == 1 {
		emb := &discordgo.MessageEmbed{
			Color: util.ColorEmbedDefault,
			Title: fmt.Sprintf("Reports for %s#%s",
				victim.User.Username, victim.User.Discriminator),
		}
		reps, err := args.CmdHandler.db.GetReportsGuild(args.Guild.ID)
		if err != nil {
			return err
		}
		if len(reps) == 0 {
			emb.Description = "This user has a white west. :ok_hand:"
		} else {
			emb.Fields = make([]*discordgo.MessageEmbedField, 0)
			for _, r := range reps {
				if r.VictimID == victim.User.ID {
					emb.Fields = append(emb.Fields, r.AsEmbedField())
				}
			}
		}
		_, err = args.Session.ChannelMessageSendEmbed(args.Channel.ID, emb)
		return err
	}

	msgOffset := 1
	repType, err := strconv.Atoi(args.Args[1])
	maxType := len(util.ReportTypes) - 1
	minType := util.ReportTypesReserved
	if repType == 0 {
		repType = minType
	}

	if victim.User.ID == args.User.ID {
		msg, err := util.SendEmbedError(args.Session, args.Channel.ID,
			"You can not report yourself...")
		util.DeleteMessageLater(args.Session, msg, 6*time.Second)
		return err
	}

	if err == nil {
		if repType < minType || repType > maxType {
			msg, err := util.SendEmbedError(args.Session, args.Channel.ID,
				fmt.Sprintf("Report type must be between *(including)* %d and %d.\n", minType, maxType)+
					"Use `help report` to get all types of report which can be used.")
			util.DeleteMessageLater(args.Session, msg, 10*time.Second)
			return err
		}
		msgOffset++
	}
	if len(args.Args[msgOffset:]) < 1 {
		msg, err := util.SendEmbedError(args.Session, args.Channel.ID,
			"Please enter a valid report description.")
		util.DeleteMessageLater(args.Session, msg, 6*time.Second)
		return err
	}
	repMsg := strings.Join(args.Args[msgOffset:], " ")
	repID := util.NodesReport[repType].Generate()

	var attachment string
	repMsg, attachment = util.ExtractImageURLFromMessage(repMsg, args.Message.Attachments)

	acceptMsg := util.AcceptMessage{
		Embed: &discordgo.MessageEmbed{
			Color:       util.ReportColors[repType],
			Title:       "Report Check",
			Description: "Is everything okay so far?",
			Fields: []*discordgo.MessageEmbedField{
				&discordgo.MessageEmbedField{
					Name: "Victim",
					Value: fmt.Sprintf("<@%s> (%s#%s)",
						victim.User.ID, victim.User.Username, victim.User.Discriminator),
				},
				&discordgo.MessageEmbedField{
					Name:  "ID",
					Value: repID.String(),
				},
				&discordgo.MessageEmbedField{
					Name:  "Type",
					Value: util.ReportTypes[repType],
				},
				&discordgo.MessageEmbedField{
					Name:  "Description",
					Value: repMsg,
				},
			},
			Image: &discordgo.MessageEmbedImage{
				URL: attachment,
			},
		},
		Session:        args.Session,
		UserID:         args.User.ID,
		DeleteMsgAfter: true,
		AcceptFunc: func(msg *discordgo.Message) {
			rep := &util.Report{
				ID:            repID,
				Type:          repType,
				GuildID:       args.Guild.ID,
				ExecutorID:    args.User.ID,
				VictimID:      victim.User.ID,
				Msg:           repMsg,
				AttachmehtURL: attachment,
			}
			err = args.CmdHandler.db.AddReport(rep)
			if err != nil {
				util.SendEmbedError(args.Session, args.Channel.ID,
					"Failed creating report: ```\n"+err.Error()+"\n```")
			} else {
				args.Session.ChannelMessageSendEmbed(args.Channel.ID, rep.AsEmbed())
				if modlogChan, err := args.CmdHandler.db.GetGuildModLog(args.Guild.ID); err == nil {
					args.Session.ChannelMessageSendEmbed(modlogChan, rep.AsEmbed())
				}
				dmChan, err := args.Session.UserChannelCreate(victim.User.ID)
				if err == nil {
					args.Session.ChannelMessageSendEmbed(dmChan.ID, rep.AsEmbed())
				}
			}
		},
	}

	_, err = acceptMsg.Send(args.Channel.ID)

	return err
}

func (c *CmdReport) revoke(args *CommandArgs) error {
	if len(args.Args) < 3 {
		msg, err := util.SendEmbedError(args.Session, args.Channel.ID,
			"Invalid command arguments. Please use `help report` for more information.")
		util.DeleteMessageLater(args.Session, msg, 6*time.Second)
		return err
	}

	id, err := strconv.Atoi(args.Args[1])
	if err != nil {
		return err
	}

	reason := strings.Join(args.Args[2:], " ")

	rep, err := args.CmdHandler.db.GetReport(snowflake.ID(id))
	if err != nil {
		if core.IsErrDatabaseNotFound(err) {
			msg, err := util.SendEmbedError(args.Session, args.Channel.ID,
				fmt.Sprintf("Could not find any report with ID `%d`", id))
			util.DeleteMessageLater(args.Session, msg, 6*time.Second)
			return err
		}
		return err
	}

	aceptMsg := util.AcceptMessage{
		Embed: &discordgo.MessageEmbed{
			Color: util.ReportRevokedColor,
			Title: "Report Revocation",
			Description: "Do you really want to revoke this report?\n" +
				":warning: **WARNING:** Revoking a report will be displayed in the mod log channel (if set) and " +
				"the revoke will be **deleted** from the database and no more visible again after!",
			Fields: []*discordgo.MessageEmbedField{
				&discordgo.MessageEmbedField{
					Name:  "Revocation Reason",
					Value: reason,
				},
				rep.AsEmbedField(),
			},
		},
		Session:        args.Session,
		DeleteMsgAfter: true,
		UserID:         args.User.ID,
		DeclineFunc: func(m *discordgo.Message) {
			msg, _ := util.SendEmbedError(args.Session, args.Channel.ID,
				"Canceled.")
			util.DeleteMessageLater(args.Session, msg, 6*time.Second)
		},
		AcceptFunc: func(m *discordgo.Message) {
			err := args.CmdHandler.db.DeleteReport(rep.ID)
			if err != nil {
				util.SendEmbedError(args.Session, args.Channel.ID,
					fmt.Sprintf("An error occured while deleting report from database: ```\n%s\n```", err.Error()))
				return
			}

			repRevEmb := &discordgo.MessageEmbed{
				Color:       util.ReportRevokedColor,
				Title:       "REPORT REVOCATION",
				Description: "Revoked reports are deleted from the database and no more visible in any commands.",
				Fields: []*discordgo.MessageEmbedField{
					&discordgo.MessageEmbedField{
						Name:  "Revoke Executor",
						Value: args.User.Mention(),
					},
					&discordgo.MessageEmbedField{
						Name:  "Revocation Reason",
						Value: reason,
					},
					rep.AsEmbedField(),
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
