package slashcommands

import (
	_ "embed"
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/services/config"
	"github.com/zekroTJA/shinpuru/internal/services/permissions"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/embedded"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekrotja/dgrs"
	"github.com/zekrotja/ken"
)

//go:embed embed/cmd_info.md
var infoMsg string

type Info struct{}

var (
	_ ken.SlashCommand        = (*Info)(nil)
	_ permissions.PermCommand = (*Info)(nil)
	_ ken.DmCapable           = (*Info)(nil)
)

func (c *Info) Name() string {
	return "info"
}

func (c *Info) Description() string {
	return "Display some information about this bot."
}

func (c *Info) Version() string {
	return "1.0.0"
}

func (c *Info) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (c *Info) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{}
}

func (c *Info) Domain() string {
	return "sp.etc.info"
}

func (c *Info) SubDomains() []permissions.SubPermission {
	return nil
}

func (c *Info) IsDmCapable() bool {
	return true
}

func (c *Info) Run(ctx *ken.Ctx) (err error) {
	if err = ctx.Defer(); err != nil {
		return
	}

	st := ctx.Get(static.DiState).(*dgrs.State)
	self, err := st.SelfUser()
	if err != nil {
		return err
	}

	invLink := util.GetInviteLink(self.ID)

	cfg := ctx.Get(static.DiConfig).(config.Provider)
	var privacyContacts strings.Builder
	for _, c := range cfg.Config().Privacy.Contact {
		privacyContacts.WriteString(
			fmt.Sprintf("- %s: [%s](%s)\n", c.Title, c.Value, c.URL))
	}

	emb := &discordgo.MessageEmbed{
		Color: static.ColorEmbedDefault,
		Title: "Info",
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: self.AvatarURL(""),
		},
		Description: infoMsg,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Repository",
				Value: "[github.com/zekrotja/shinpuru](https://github.com/zekrotja/shinpuru)",
			},
			{
				Name: "Version",
				Value: fmt.Sprintf("This instance is running on version **%s** (commit hash `%s`)",
					embedded.AppVersion, embedded.AppCommit),
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
				Name: "Privacy",
				Value: fmt.Sprintf("[Privacy Notice](%s)\n\nContact:\n%s",
					cfg.Config().Privacy.NoticeURL, privacyContacts.String()),
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

	return ctx.FollowUpEmbed(emb).Error
}
