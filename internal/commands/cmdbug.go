package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/util"
)

type CmdBug struct {
	PermLvl int
}

func (c *CmdBug) GetInvokes() []string {
	return []string{"bug", "bugreport", "issue", "suggestion"}
}

func (c *CmdBug) GetDescription() string {
	return "Get information how to submit a bug report or feature request"
}

func (c *CmdBug) GetHelp() string {
	return "`bug`"
}

func (c *CmdBug) GetGroup() string {
	return GroupEtc
}

func (c *CmdBug) GetPermission() int {
	return c.PermLvl
}

func (c *CmdBug) SetPermission(permLvl int) {
	c.PermLvl = permLvl
}

func (c *CmdBug) Exec(args *CommandArgs) error {
	emb := &discordgo.MessageEmbed{
		Color: util.ColorEmbedDefault,
		Title: "How to report a bug or request a feature",
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name: "GitHub Issue",
				Value: "You can submit an issue by using the GitHub Issue tracker. " +
					"\n*For that, you will need a GitHub account.*:\n" +
					":link:  [**GitHub Issues**](https://github.com/zekroTJA/shinpuru/issues/new/choose)",
			},
			&discordgo.MessageEmbedField{
				Name: "Google Forms",
				Value: "Alternatively, you can also submit an issue by unsing the follwoing form: " +
					"\n*This will be transformed into an issue on GitHub later.*:\n" +
					":link:  [**Google Form**](https://docs.google.com/forms/d/e/1FAIpQLScKnY2FUDqmLVg2TjdBqSAyL-LlD55y7h5JcqsT887KwLPkIg/viewform?usp=sf_link)",
			},
			&discordgo.MessageEmbedField{
				Name:  "Bug Hunters",
				Value: "Much :heart: to all [**bug hunters**](https://github.com/zekroTJA/shinpuru/blob/dev/bughunters.md).",
			},
		},
	}
	_, err := args.Session.ChannelMessageSendEmbed(args.Channel.ID, emb)
	return err
}
