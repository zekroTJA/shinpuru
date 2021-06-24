package commands

import (
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shireikan"
	"github.com/zekrotja/discordgo"
)

type CmdBug struct {
}

func (c *CmdBug) GetInvokes() []string {
	return []string{"bug", "bugreport", "issue", "suggestion"}
}

func (c *CmdBug) GetDescription() string {
	return "Get information how to submit a bug report or feature request."
}

func (c *CmdBug) GetHelp() string {
	return "`bug`"
}

func (c *CmdBug) GetGroup() string {
	return shireikan.GroupEtc
}

func (c *CmdBug) GetDomainName() string {
	return "sp.etc.bug"
}

func (c *CmdBug) GetSubPermissionRules() []shireikan.SubPermission {
	return nil
}

func (c *CmdBug) IsExecutableInDMChannels() bool {
	return true
}

func (c *CmdBug) Exec(ctx shireikan.Context) error {
	emb := &discordgo.MessageEmbed{
		Color: static.ColorEmbedDefault,
		Title: "How to report a bug or request a feature",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name: "GitHub Issue",
				Value: "You can submit an issue by using the GitHub Issue tracker. " +
					"\n*For that, you will need a GitHub account.*:\n" +
					":link:  [**GitHub Issues**](https://github.com/zekroTJA/shinpuru/issues/new/choose)",
			},
			{
				Name: "Google Forms",
				Value: "Alternatively, you can also submit an issue by unsing the follwoing form: " +
					"\n*This will be transformed into an issue on GitHub later.*:\n" +
					":link:  [**Google Form**](https://docs.google.com/forms/d/e/1FAIpQLScKnY2FUDqmLVg2TjdBqSAyL-LlD55y7h5JcqsT887KwLPkIg/viewform?usp=sf_link)",
			},
			{
				Name:  "Bug Hunters",
				Value: "Much :heart: to all [**bug hunters**](https://github.com/zekroTJA/shinpuru/blob/dev/bughunters.md).",
			},
		},
	}
	_, err := ctx.GetSession().ChannelMessageSendEmbed(ctx.GetChannel().ID, emb)
	return err
}
