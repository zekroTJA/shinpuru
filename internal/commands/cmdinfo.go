package commands

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shireikan"
)

type CmdInfo struct {
}

func (c *CmdInfo) GetInvokes() []string {
	return []string{"info", "information", "description", "credits", "version", "invite"}
}

func (c *CmdInfo) GetDescription() string {
	return "Display some information about this bot."
}

func (c *CmdInfo) GetHelp() string {
	return "`info`"
}

func (c *CmdInfo) GetGroup() string {
	return shireikan.GroupGeneral
}

func (c *CmdInfo) GetDomainName() string {
	return "sp.etc.info"
}

func (c *CmdInfo) GetSubPermissionRules() []shireikan.SubPermission {
	return nil
}

func (c *CmdInfo) IsExecutableInDMChannels() bool {
	return true
}

func (c *CmdInfo) Exec(ctx shireikan.Context) error {
	invLink := util.GetInviteLink(ctx.GetSession())

	emb := &discordgo.MessageEmbed{
		Color: static.ColorEmbedDefault,
		Title: "Info",
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: ctx.GetSession().State.User.AvatarURL(""),
		},
		Description: "シンプル (shinpuru), a simple *(as the name says)*, multi purpose Discord Bot written in Go, " +
			"using bwmarrin's package [discord.go](https://github.com/bwmarrin/discordgo) as API and gateway wrapper. " +
			"The focus on this bot is not to punch in as much features and commands as possible, just some commands and " +
			"features which I thought would be useful and which were the most used with my older Discord bots, like " +
			"[zekroBot 2](https://github.com/zekroTJA/zekroBot2), and more on making this bot as reliable and stable as possible.",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Repository",
				Value: "[github.com/zekrotja/shinpuru](https://github.com/zekrotja/shinpuru)",
			},
			{
				Name: "Version",
				Value: fmt.Sprintf("This instance is running on version **%s** (commit hash `%s`)",
					util.AppVersion, util.AppCommit),
			},
			{
				Name:  "Licence",
				Value: "Covered by the [MIT Licence](https://github.com/zekroTJA/shinpuru/blob/master/LICENCE).",
			},
			{
				Name: "Invite",
				Value: fmt.Sprintf("[Invite Link](%s).\n```\n%s\n```",
					invLink, invLink),
			},
			{
				Name:  "Bug Hunters",
				Value: "Much :heart: to all [**bug hunters**](https://github.com/zekroTJA/shinpuru/blob/dev/bughunters.md).",
			},
			{
				Name:  "Development state",
				Value: "You can see current tasks [here](https://github.com/zekroTJA/shinpuru/projects).",
			},
			{
				Name: "3rd Party Dependencies and Credits",
				Value: "[Here](https://github.com/zekroTJA/shinpuru/blob/master/README.md#third-party-dependencies) you can find a list of all dependencies used.\n" +
					"Avatar of [御中元 魔法少女詰め合わせ](https://www.pixiv.net/member_illust.php?mode=medium&illust_id=44692506) from [瑞希](https://www.pixiv.net/member.php?id=137253).",
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("© 2018-%s zekro Development (Ringo Hoffmann)", time.Now().Format("2006")),
		},
	}
	_, err := ctx.GetSession().ChannelMessageSendEmbed(ctx.GetChannel().ID, emb)
	return err
}
