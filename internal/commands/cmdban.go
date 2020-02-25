package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/shared"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/acceptmsg"
	"github.com/zekroTJA/shinpuru/internal/util/imgstore"
	"github.com/zekroTJA/shinpuru/internal/util/snowflakenodes"
	"github.com/zekroTJA/shinpuru/internal/util/static"
)

type CmdBan struct {
}

func (c *CmdBan) GetInvokes() []string {
	return []string{"ban", "userban"}
}

func (c *CmdBan) GetDescription() string {
	return "ban users with creating a report entry"
}

func (c *CmdBan) GetHelp() string {
	return "`ban <UserResolvable> <Reason>`"
}

func (c *CmdBan) GetGroup() string {
	return GroupModeration
}

func (c *CmdBan) GetDomainName() string {
	return "sp.guild.mod.ban"
}

func (c *CmdBan) GetSubPermissionRules() []SubPermission {
	return nil
}

func (c *CmdBan) Exec(args *CommandArgs) error {
	if len(args.Args) < 2 {
		msg, err := util.SendEmbedError(args.Session, args.Channel.ID,
			"Invalid command arguments. Please use `help ban` to see how to use this command.")
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

	if victim.User.ID == args.User.ID {
		msg, err := util.SendEmbedError(args.Session, args.Channel.ID,
			"You can not ban yourself...")
		util.DeleteMessageLater(args.Session, msg, 6*time.Second)
		return err
	}

	authorMemb, err := args.Session.GuildMember(args.Guild.ID, args.User.ID)
	if err != nil {
		return err
	}

	if util.RolePosDiff(victim, authorMemb, args.Guild) >= 0 {
		msg, err := util.SendEmbedError(args.Session, args.Channel.ID,
			"You can only ban members with lower permissions than yours.")
		util.DeleteMessageLater(args.Session, msg, 8*time.Second)
		return err
	}

	repMsg := strings.Join(args.Args[1:], " ")
	var repType int
	for i, v := range static.ReportTypes {
		if v == "BAN" {
			repType = i
		}
	}
	repID := snowflakenodes.NodesReport[repType].Generate()

	var attachment string
	repMsg, attachment = imgstore.ExtractFromMessage(repMsg, args.Message.Attachments)
	if attachment != "" {
		img, err := imgstore.DownloadFromURL(attachment)
		if err == nil && img != nil {
			if err = args.CmdHandler.db.SaveImageData(img); err != nil {
				return err
			}
			attachment = img.ID.String()
		}
	}

	acceptMsg := acceptmsg.AcceptMessage{
		Embed: &discordgo.MessageEmbed{
			Color:       static.ReportColors[repType],
			Title:       "Ban Check",
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
					Value: static.ReportTypes[repType],
				},
				&discordgo.MessageEmbedField{
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
			rep, err := shared.PushBan(
				args.Session,
				args.CmdHandler.db,
				args.CmdHandler.config.WebServer.PublicAddr,
				args.Guild.ID,
				args.User.ID,
				victim.User.ID,
				repMsg,
				attachment)

			if err != nil {
				util.SendEmbedError(args.Session, args.Channel.ID,
					"Failed banning member: ```\n"+err.Error()+"\n```")
			} else {
				args.Session.ChannelMessageSendEmbed(args.Channel.ID, rep.AsEmbed(args.CmdHandler.config.WebServer.PublicAddr))
			}
		},
	}

	_, err = acceptMsg.Send(args.Channel.ID)

	return err
}
