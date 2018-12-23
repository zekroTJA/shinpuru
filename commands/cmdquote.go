package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/zekroTJA/shinpuru/util"
)

type CmdQuote struct {
}

func (c *CmdQuote) GetInvokes() []string {
	return []string{"quote", "q"}
}

func (c *CmdQuote) GetDescription() string {
	return "quote a messgage from any chat"
}

func (c *CmdQuote) GetHelp() string {
	return "`quote <msgID/msgURL>`"
}

func (c *CmdQuote) GetGroup() string {
	return GroupChat
}

func (c *CmdQuote) GetPermission() int {
	return 0
}

func (c *CmdQuote) Exec(args *CommandArgs) error {
	if len(args.Args) < 1 {
		msg, err := util.SendEmbedError(args.Session, args.Channel.ID,
			"Please enter a message ID or URL which should be quoted.")
		util.DeleteMessageLater(args.Session, msg, 5*time.Second)
		return err
	}

	if strings.HasPrefix(args.Args[0], "https://discordapp.com/channels/") {
		urlSplit := strings.Split(args.Args[0], "/")
		args.Args[0] = urlSplit[len(urlSplit)-1]
	}

	msgSearchEmb := &discordgo.MessageEmbed{
		Color:       util.ColorEmbedGray,
		Description: fmt.Sprintf(":hourglass_flowing_sand:  Searching for message in channel <#%s>...", args.Channel.ID),
	}

	msgSearch, err := args.Session.ChannelMessageSendEmbed(args.Channel.ID, msgSearchEmb)
	if err != nil {
		return err
	}

	quoteMsg, err := args.Session.ChannelMessage(args.Channel.ID, args.Args[0])
	if err != nil {
		for _, ch := range args.Guild.Channels {
			if ch.Type != discordgo.ChannelTypeGuildText {
				continue
			}
			msgSearchEmb.Description = fmt.Sprintf(":hourglass_flowing_sand:  Searching for message in channel <#%s>...", ch.ID)
			args.Session.ChannelMessageEditEmbed(args.Channel.ID, msgSearch.ID, msgSearchEmb)
			quoteMsg, err = args.Session.ChannelMessage(ch.ID, args.Args[0])
			if err == nil && quoteMsg != nil {
				break
			}
		}
	}

	if err != nil || quoteMsg == nil {
		msgSearchEmb.Description = "Could not find any message with this ID. :disappointed:"
		msgSearchEmb.Color = util.ColorEmbedError
		_, err := args.Session.ChannelMessageEditEmbed(args.Channel.ID, msgSearch.ID, msgSearchEmb)
		util.DeleteMessageLater(args.Session, msgSearch, 5*time.Second)
		return err
	}

	if len(quoteMsg.Content) < 1 {
		msgSearchEmb.Description = "Found messages contect is empty. Maybe, it is an embed message itself, which can not be quoted."
		msgSearchEmb.Color = util.ColorEmbedError
		_, err := args.Session.ChannelMessageEditEmbed(args.Channel.ID, msgSearch.ID, msgSearchEmb)
		util.DeleteMessageLater(args.Session, msgSearch, 8*time.Second)
		return err
	}

	quoteMsgChannel, err := args.Session.Channel(quoteMsg.ChannelID)
	if err != nil {
		return err
	}

	msgSearchEmb = &discordgo.MessageEmbed{
		Color: util.ColorEmbedDefault,
		Author: &discordgo.MessageEmbedAuthor{
			IconURL: quoteMsg.Author.AvatarURL(""),
			Name:    quoteMsg.Author.Username + "#" + quoteMsg.Author.Discriminator,
		},
		Description: quoteMsg.Content +
			fmt.Sprintf("\n\n*[jump to message](https://discordapp.com/channels/%s/%s/%s)*", args.Guild.ID, quoteMsg.ChannelID, quoteMsg.ID),
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("#%s - quoted by: %s#%s", quoteMsgChannel.Name, args.User.Username, args.User.Discriminator),
		},
		Timestamp: string(quoteMsg.Timestamp),
	}
	args.Session.ChannelMessageEditEmbed(args.Channel.ID, msgSearch.ID, msgSearchEmb)
	return nil
}
