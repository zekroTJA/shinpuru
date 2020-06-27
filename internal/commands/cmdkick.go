package commands

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/shared"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/imgstore"
	"github.com/zekroTJA/shinpuru/internal/util/snowflakenodes"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/acceptmsg"
)

type CmdKick struct {
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

func (c *CmdKick) GetDomainName() string {
	return "sp.guild.mod.kick"
}

func (c *CmdKick) GetSubPermissionRules() []SubPermission {
	return nil
}

func (c *CmdKick) Exec(args *CommandArgs) error {
	if len(args.Args) < 2 {
		return util.SendEmbedError(args.Session, args.Channel.ID,
			"Invalid command arguments. Please use `help kick` to see how to use this command.").
			DeleteAfter(8 * time.Second).Error()
	}
	victim, err := util.FetchMember(args.Session, args.Guild.ID, args.Args[0])
	if err != nil || victim == nil {
		return util.SendEmbedError(args.Session, args.Channel.ID,
			"Sorry, could not find any member :cry:").
			DeleteAfter(8 * time.Second).Error()
	}

	if victim.User.ID == args.User.ID {
		return util.SendEmbedError(args.Session, args.Channel.ID,
			"You can not kick yourself...").
			DeleteAfter(8 * time.Second).Error()
	}

	authorMemb, err := args.Session.GuildMember(args.Guild.ID, args.User.ID)
	if err != nil {
		return err
	}

	if util.RolePosDiff(victim, authorMemb, args.Guild) >= 0 {
		return util.SendEmbedError(args.Session, args.Channel.ID,
			"You can only kick members with lower permissions than yours.").
			DeleteAfter(8 * time.Second).Error()
	}

	repMsg := strings.Join(args.Args[1:], " ")
	var repType int
	for i, v := range static.ReportTypes {
		if v == "KICK" {
			repType = i
		}
	}
	repID := snowflakenodes.NodesReport[repType].Generate()

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
			Title:       "Kick Check",
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
			rep, err := shared.PushKick(
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
					"Failed kicking member: ```\n"+err.Error()+"\n```")
			} else {
				args.Session.ChannelMessageSendEmbed(args.Channel.ID, rep.AsEmbed(args.CmdHandler.config.WebServer.PublicAddr))
			}
		},
	}

	_, err = acceptMsg.Send(args.Channel.ID)

	return err
}
