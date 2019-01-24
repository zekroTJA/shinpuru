package commands

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"

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
	return "`report <userResolvable> [<type>] <reason>` - report a user *(if type is empty, its defaultly 0 = warn)*\n" +
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
	repID := util.ReportNodes[repType].Generate()

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
		},
		Session:        args.Session,
		UserID:         args.User.ID,
		DeleteMsgAfter: true,
		AcceptFunc: func(msg *discordgo.Message) {
			rep := &util.Report{
				ID:         repID,
				Type:       repType,
				GuildID:    args.Guild.ID,
				ExecutorID: args.User.ID,
				VictimID:   victim.User.ID,
				Msg:        repMsg,
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
