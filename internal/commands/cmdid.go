package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/util"
)

type CmdId struct {
}

func (c *CmdId) GetInvokes() []string {
	return []string{"id", "ids"}
}

func (c *CmdId) GetDescription() string {
	return "Get the discord ID(s) by resolvable"
}

func (c *CmdId) GetHelp() string {
	return "`id (<resolvable>)`"
}

func (c *CmdId) GetGroup() string {
	return GroupEtc
}

func (c *CmdId) GetDomainName() string {
	return "sp.etc.id"
}

func (c *CmdId) Exec(args *CommandArgs) error {
	var user *discordgo.User
	var role *discordgo.Role
	var textChannel *discordgo.Channel
	var voiceChannel *discordgo.Channel

	if len(args.Args) < 1 {
		user = args.User
	} else {
		joinedArgs := strings.Join(args.Args, " ")
		if u, err := util.FetchMember(args.Session, args.Guild.ID, joinedArgs); err == nil {
			user = u.User
		}
		if r, err := util.FetchRole(args.Session, args.Guild.ID, joinedArgs); err == nil {
			role = r
		}
		if tc, err := util.FetchChannel(args.Session, args.Guild.ID, joinedArgs, func(c *discordgo.Channel) bool {
			return c.Type == discordgo.ChannelTypeGuildText
		}); err == nil {
			textChannel = tc
		}
		if vc, err := util.FetchChannel(args.Session, args.Guild.ID, joinedArgs, func(c *discordgo.Channel) bool {
			return c.Type == discordgo.ChannelTypeGuildVoice
		}); err == nil {
			voiceChannel = vc
		}
	}

	if user == nil && role == nil && textChannel == nil && voiceChannel == nil {
		msg, err := util.SendEmbedError(args.Session, args.Channel.ID,
			"Could not fetch any member, role or channel by this resolvable.")
		util.DeleteMessageLater(args.Session, msg, 8*time.Second)
		return err
	}

	emb := &discordgo.MessageEmbed{
		Color:  util.ColorEmbedDefault,
		Fields: make([]*discordgo.MessageEmbedField, 0),
	}

	if user != nil {
		emb.Fields = append(emb.Fields, &discordgo.MessageEmbedField{
			Name:  "Member",
			Value: fmt.Sprintf("<@%s> (%s#%s)\n```\n%s\n```", user.ID, user.Username, user.Discriminator, user.ID),
		})
	}
	if role != nil {
		emb.Fields = append(emb.Fields, &discordgo.MessageEmbedField{
			Name:  "Role",
			Value: fmt.Sprintf("<@&%s> (%s)\n```\n%s\n```", role.ID, role.Name, role.ID),
		})
	}
	if textChannel != nil {
		emb.Fields = append(emb.Fields, &discordgo.MessageEmbedField{
			Name:  "Text Channel",
			Value: fmt.Sprintf("<#%s> (%s)\n```\n%s\n```", textChannel.ID, textChannel.Name, textChannel.ID),
		})
	}
	if voiceChannel != nil {
		emb.Fields = append(emb.Fields, &discordgo.MessageEmbedField{
			Name:  "Voice Channel",
			Value: fmt.Sprintf("%s\n```\n%s\n```", voiceChannel.Name, voiceChannel.ID),
		})
	}
	emb.Fields = append(emb.Fields, &discordgo.MessageEmbedField{
		Name:  "Guild",
		Value: fmt.Sprintf("%s\n```\n%s\n```", args.Guild.Name, args.Guild.ID),
	})

	_, err := args.Session.ChannelMessageSendEmbed(args.Channel.ID, emb)
	return err
}
