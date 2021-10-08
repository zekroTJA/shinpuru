package slashcommands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/services/permissions"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekrotja/ken"
)

type Exec struct{}

var (
	_ ken.Command             = (*Exec)(nil)
	_ permissions.PermCommand = (*Exec)(nil)
)

func (c *Exec) Name() string {
	return "bug"
}

func (c *Exec) Description() string {
	return "Get information how to submit a bug report or feature request."
}

func (c *Exec) Version() string {
	return "1.0.0"
}

func (c *Exec) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (c *Exec) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{}
}

func (c *Exec) Domain() string {
	return "sp.etc.bug"
}

func (c *Exec) SubDomains() []permissions.SubPermission {
	return nil
}

func (c *Exec) Run(ctx *ken.Ctx) (err error) {
	if err = ctx.Defer(); err != nil {
		return
	}

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
				Value: "Alternatively, you can also submit an issue by using the following form: " +
					"\n*This will be transformed into an issue on GitHub later.*:\n" +
					":link:  [**Google Form**](https://docs.google.com/forms/d/e/1FAIpQLScKnY2FUDqmLVg2TjdBqSAyL-LlD55y7h5JcqsT887KwLPkIg/viewform?usp=sf_link)",
			},
			{
				Name:  "Bug Hunters",
				Value: "Much :heart: to all [**Exec hunters**](https://github.com/zekroTJA/shinpuru/blob/dev/bughunters.md).",
			},
		},
	}
	err = ctx.FollowUpEmbed(emb).Error

	return
}
