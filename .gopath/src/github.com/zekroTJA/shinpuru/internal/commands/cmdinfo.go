package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/util"
)

type CmdInfo struct {
	PermLvl int
}

func (c *CmdInfo) GetInvokes() []string {
	return []string{"info", "information", "description", "credits", "version", "invite"}
}

func (c *CmdInfo) GetDescription() string {
	return "display some information about this bot"
}

func (c *CmdInfo) GetHelp() string {
	return "`info`"
}

func (c *CmdInfo) GetGroup() string {
	return GroupGeneral
}

func (c *CmdInfo) GetPermission() int {
	return c.PermLvl
}

func (c *CmdInfo) SetPermission(permLvl int) {
	c.PermLvl = permLvl
}

func (c *CmdInfo) Exec(args *CommandArgs) error {
	invLink := fmt.Sprintf("https://discordapp.com/api/oauth2/authorize?client_id=%s&scope=bot&permissions=%d",
		args.Session.State.User.ID, util.InvitePermission)
	emb := &discordgo.MessageEmbed{
		Color: util.ColorEmbedDefault,
		Title: "Info",
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: args.Session.State.User.AvatarURL(""),
		},
		Description: "シンプル (shinpuru), a simple *(as the name says)*, multi purpose Discord Bot written in Go, " +
			"using bwmarrin's package [discord.go](https://github.com/bwmarrin/discordgo) as API and gateway wrapper. " +
			"The focus on this bot is not to punch in as much features and commands as possible, just some commands and " +
			"features which I thought would be useful and which were the most used with my older Discord bots, like " +
			"[zekroBot 2](https://github.com/zekroTJA/zekroBot2), and more on making this bot as reliable and stable as possible.",
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name:  "Repository",
				Value: "[github.com/zekrotja/shinpuru](https://github.com/zekrotja/shinpuru)",
			},
			&discordgo.MessageEmbedField{
				Name: "Version",
				Value: fmt.Sprintf("This instance is running on version **%s** (commit hash `%s`)",
					util.AppVersion, util.AppCommit),
			},
			&discordgo.MessageEmbedField{
				Name:  "Licence",
				Value: "Covered by [MIT Licence](https://github.com/zekroTJA/shinpuru/blob/master/LICENCE).",
			},
			&discordgo.MessageEmbedField{
				Name: "Invite",
				Value: fmt.Sprintf("[Invite Link](%s).\n```\n%s\n```",
					invLink, invLink),
			},
			&discordgo.MessageEmbedField{
				Name:  "Bug Hunters",
				Value: "Much :heart: to all [**bug hunters**](https://github.com/zekroTJA/shinpuru/blob/dev/bughunters.md).",
			},
			&discordgo.MessageEmbedField{
				Name:  "Development state",
				Value: "You can see current tasks [here](https://github.com/zekroTJA/shinpuru/projects).",
			},
			&discordgo.MessageEmbedField{
				Name: "3rd Party Dependencies and Credits",
				Value: "- [bwmarrin/discordgo](https://github.com/bwmarrin/discordgo)\n" +
					"- [go-yaml/yaml](https://github.com/go-yaml/yaml)\n" +
					"- [go-sql-driver/mysql](https://github.com/Go-SQL-Driver/MySQL/)\n" +
					"- [op/go-logging](https://github.com/op/go-logging)\n\n" +
					"Avatar of [御中元 魔法少女詰め合わせ](https://www.pixiv.net/member_illust.php?mode=medium&illust_id=44692506) from [瑞希](https://www.pixiv.net/member.php?id=137253).",
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "© 2018 zekro Development (Ringo Hoffmann)",
		},
	}
	_, err := args.Session.ChannelMessageSendEmbed(args.Channel.ID, emb)
	return err
}
