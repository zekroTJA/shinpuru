package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/util"
)

type CmdKick struct {
	PermLvl int
}

func (c *CmdKick) GetInvokes() []string {
	return []string{"kick", "userkick"}
}

func (c *CmdKick) GetDescription() string {
	return "kick users with creating a report entry"
}

func (c *CmdKick) GetHelp() string {
	return "`kick <UserResolvable> <Reason>`"
}

func (c *CmdKick) GetGroup() string {
	return GroupModeration
}

func (c *CmdKick) GetPermission() int {
	return c.PermLvl
}

func (c *CmdKick) SetPermission(permLvl int) {
	c.PermLvl = permLvl
}

func (c *CmdKick) Exec(args *CommandArgs) error {
	if len(args.Args) < 2 {
		msg, err := util.SendEmbedError(args.Session, args.Channel.ID,
			"Invalid command arguments. Please use `help kick` to see how to use this command.")
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

	repMsg := strings.Join(args.Args[1:], " ")
	var repType int
	for i, v := range util.ReportTypes {
		if v == "KICK" {
			repType = i
		}
	}
	repID := util.ReportNodes[repType].Generate()

	acceptMsg := util.AcceptMessage{
		Embed: &discordgo.MessageEmbed{
			Color:       util.ReportColors[repType],
			Title:       "Kick Check",
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
				return
			}
			args.Session.ChannelMessageSendEmbed(args.Channel.ID, rep.AsEmbed())
			if modlogChan, err := args.CmdHandler.db.GetGuildModLog(args.Guild.ID); err == nil {
				args.Session.ChannelMessageSendEmbed(modlogChan, rep.AsEmbed())
			}
			dmChan, err := args.Session.UserChannelCreate(victim.User.ID)
			if err == nil {
				args.Session.ChannelMessageSendEmbed(dmChan.ID, rep.AsEmbed())
			}
			err = args.Session.GuildMemberDeleteWithReason(args.Guild.ID, victim.User.ID, repMsg)
			if err != nil {
				util.SendEmbedError(args.Session, args.Channel.ID,
					"Failed kicking member: ```\n"+err.Error()+"\n```")
				return
			}
		},
	}

	_, err = acceptMsg.Send(args.Channel.ID)

	return err
}
