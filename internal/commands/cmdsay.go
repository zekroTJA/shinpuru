package commands

import (
	"encoding/json"
	"flag"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/zekrotja/discordgo"

	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shireikan"
)

var embedColors = map[string]int{
	"red":    0xf44336,
	"pink":   0xE91E63,
	"violet": 0x9C27B0,
	"blue":   0x2196F3,
	"cyan":   0x00BCD4,
	"green":  0x8BC34A,
	"yellow": 0xFFEB3B,
	"orange": 0xffc107,
	"white":  0xF5F5F5,
	"black":  0x263238,
}

var (
	msgLinkRx = regexp.MustCompile(`https?:\/\/(?:canary\.)?discord(?:app)?\.com\/channels\/\d+\/(\d+)\/(\d+)`)
)

type CmdSay struct {
}

func (c *CmdSay) GetInvokes() []string {
	return []string{"say", "msg"}
}

func (c *CmdSay) GetDescription() string {
	return "Send an embedded message with the bot."
}

func (c *CmdSay) GetSubPermissionRules() []shireikan.SubPermission {
	return nil
}

func (c *CmdSay) GetHelp() string {
	colors := make([]string, len(embedColors))
	i := 0
	for k := range embedColors {
		colors[i] = k
		i++
	}
	return "`say [flags] <message>`\n\n**Flags:** \n```\n" +
		`-c string
      color (default "orange")
-e string
      Message Link or [ChannelID/]MessageID of the message to be edited
-f string
      footer
-raw
      parses following content as raw embed from json (see https://discord.com/developers/docs/resources/channel#embed-object)
-t string
      title
` + "\n```\n**Colors:**\n" + strings.Join(colors, ", ")
}

func (c *CmdSay) GetGroup() string {
	return shireikan.GroupChat
}

func (c *CmdSay) GetDomainName() string {
	return "sp.chat.say"
}

func (c *CmdSay) IsExecutableInDMChannels() bool {
	return true
}

func (c *CmdSay) Exec(ctx shireikan.Context) (err error) {
	f := flag.NewFlagSet("sayflags", flag.ContinueOnError)
	fcolor := f.String("c", "orange", "color")
	ftitle := f.String("t", "", "title")
	ffooter := f.String("f", "", "footer")
	fraw := f.Bool("raw", false, "parses following content as raw embed from json (see https://discord.com/developers/docs/resources/channel#embed-object)")
	fedit := f.String("e", "", "Message Link or [ChannelID/]MessageID of the message to be edited")

	if err = f.Parse(ctx.GetArgs()); err != nil {
		return err
	}

	var editMsg *discordgo.Message

	if *fedit != "" {
		msgID, chanID := getMsgID(*fedit, ctx.GetChannel().ID)
		editMsg, err = ctx.GetSession().ChannelMessage(chanID, msgID)
		if err != nil {
			return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
				fmt.Sprintf("The message to be edited could not be found:\n```\n%s\n```", err.Error())).
				DeleteAfter(10 * time.Second).Error()
		}
		if editMsg.Author.ID != ctx.GetSession().State.User.ID {
			return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
				"You can only edit messages which were created by shinpuru.").
				DeleteAfter(8 * time.Second).Error()
		}
	}

	authorField := &discordgo.MessageEmbedAuthor{
		IconURL: ctx.GetUser().AvatarURL(""),
		Name:    ctx.GetUser().Username,
	}

	var emb *discordgo.MessageEmbed
	if *fraw {
		offset := strings.IndexRune(ctx.GetMessage().Content, '{')
		if offset < 0 || offset >= len(ctx.GetMessage().Content) {
			return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
				"Wrong JSON format. The JSON object must start with `{`."+
					"If you need help building an embed with raw json, take a look here:\nhttps://discord.com/developers/docs/resources/channel#embed-object").
				DeleteAfter(20 * time.Second).Error()
		}
		content := ctx.GetMessage().Content[offset:]
		err := json.Unmarshal([]byte(content), &emb)
		if err != nil {
			return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
				fmt.Sprintf("Failed parsing message embed from input: ```\n%s\n```", err.Error())+
					"If you need help building an embed with raw json, take a look here:\nhttps://discord.com/developers/docs/resources/channel#embed-object").
				DeleteAfter(20 * time.Second).Error()
		}
		emb.Author = authorField
	} else {
		content := strings.Join(f.Args(), " ")
		if len(content) < 1 {
			return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
				"Please enter something you want to say :wink:").
				DeleteAfter(8 * time.Second).Error()
		}
		embColor, ok := embedColors[strings.ToLower(*fcolor)]
		if !ok {
			return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
				fmt.Sprintf("Sorry, but I don't know the color `%s`. Please enter `help say` to get a list of valid colors.", *fcolor)).
				DeleteAfter(8 * time.Second).Error()
		}

		emb = &discordgo.MessageEmbed{
			Title:       *ftitle,
			Color:       embColor,
			Author:      authorField,
			Description: content,
		}

		if *ffooter != "" {
			emb.Footer = &discordgo.MessageEmbedFooter{
				Text: *ffooter,
			}
		}
	}

	if editMsg != nil {
		_, err = ctx.GetSession().ChannelMessageEditEmbed(editMsg.ChannelID, editMsg.ID, emb)
	} else {
		_, err = ctx.GetSession().ChannelMessageSendEmbed(ctx.GetChannel().ID, emb)
	}

	return
}

func getMsgID(v, altChanID string) (msgID, chanID string) {
	res := msgLinkRx.FindAllStringSubmatch(v, -1)
	if res != nil && len(res) >= 1 && len(res[0]) >= 3 {
		chanID = res[0][1]
		msgID = res[0][2]
		return
	}

	i := strings.Index(v, "/")
	if i >= 0 && i < len(v)-1 {
		chanID = v[:i]
		msgID = v[i+1:]
		return
	}

	msgID = v
	chanID = altChanID

	return
}
