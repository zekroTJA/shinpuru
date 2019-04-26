package commands

import (
	"encoding/json"
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/zekroTJA/shinpuru/internal/util"
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

type CmdSay struct {
	PermLvl int
}

func (c *CmdSay) GetInvokes() []string {
	return []string{"say", "msg"}
}

func (c *CmdSay) GetDescription() string {
	return "send an embeded message with the bot"
}

func (c *CmdSay) GetHelp() string {
	colors := make([]string, len(embedColors))
	i := 0
	for k := range embedColors {
		colors[i] = k
		i++
	}
	return "`say [flags] <message>`\n\n**Flags:** \n```\n" +
		"-c string\n" +
		"	color (default \"orange\")\n" +
		"-f string\n" +
		"	footer\n" +
		"-raw string\n" +
		"	raw embed from json (see https://discordapp.com/developers/docs/resources/channel#embed-object)\n" +
		"-t string\n" +
		"	title\n```\n" +
		"**Colors:**\n" + strings.Join(colors, ", ")
}

func (c *CmdSay) GetGroup() string {
	return GroupChat
}

func (c *CmdSay) GetPermission() int {
	return c.PermLvl
}

func (c *CmdSay) SetPermission(permLvl int) {
	c.PermLvl = permLvl
}

func (c *CmdSay) Exec(args *CommandArgs) error {
	f := flag.NewFlagSet("sayflags", flag.ContinueOnError)
	fcolor := f.String("c", "orange", "color")
	ftitle := f.String("t", "", "title")
	ffooter := f.String("f", "", "footer")
	fraw := f.Bool("raw", false, "parses following content as raw embed from json (see https://discordapp.com/developers/docs/resources/channel#embed-object)")
	f.Parse(args.Args)

	authorField := &discordgo.MessageEmbedAuthor{
		IconURL: args.User.AvatarURL(""),
		Name:    args.User.Username,
	}

	var emb *discordgo.MessageEmbed
	if *fraw {
		offset := strings.IndexRune(args.Message.Content, '{')
		if offset < 0 || offset >= len(args.Message.Content) {
			msg, err := util.SendEmbedError(args.Session, args.Channel.ID,
				"Wrong JSON format. The JSON object must start with `{`."+
					"If you need help building an embed with raw json, take a look here:\nhttps://discordapp.com/developers/docs/resources/channel#embed-object")
			util.DeleteMessageLater(args.Session, msg, 30*time.Second)
			return err
		}
		content := args.Message.Content[offset:]
		err := json.Unmarshal([]byte(content), &emb)
		if err != nil {
			msg, err := util.SendEmbedError(args.Session, args.Channel.ID,
				fmt.Sprintf("Failed parsing message embed from input: ```\n%s\n```", err.Error())+
					"If you need help building an embed with raw json, take a look here:\nhttps://discordapp.com/developers/docs/resources/channel#embed-object")
			util.DeleteMessageLater(args.Session, msg, 30*time.Second)
			return err
		}
		emb.Author = authorField
	} else {
		content := strings.Join(f.Args(), " ")
		if len(content) < 1 {
			msg, err := util.SendEmbedError(args.Session, args.Channel.ID,
				"Please enter something you want to say :wink:")
			util.DeleteMessageLater(args.Session, msg, 6*time.Second)
			return err
		}
		embColor, ok := embedColors[strings.ToLower(*fcolor)]
		if !ok {
			msg, err := util.SendEmbedError(args.Session, args.Channel.ID,
				fmt.Sprintf("Sorry, but I don't know the color `%s`. Please enter `help say` to get a list of valid colors.", *fcolor))
			util.DeleteMessageLater(args.Session, msg, 10*time.Second)
			return err
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

	_, err := args.Session.ChannelMessageSendEmbed(args.Channel.ID, emb)

	return err
}
