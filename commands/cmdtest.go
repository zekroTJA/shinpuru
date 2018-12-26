package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/util"
)

type CmdTest struct {
}

func (c *CmdTest) GetInvokes() []string {
	return []string{"test"}
}

func (c *CmdTest) GetDescription() string {
	return "just for testing purposes"
}

func (c *CmdTest) GetHelp() string {
	return ""
}

func (c *CmdTest) GetGroup() string {
	return GroupEtc
}

func (c *CmdTest) GetPermission() int {
	return 999
}

func (c *CmdTest) Exec(args *CommandArgs) error {
	am := &util.AcceptMessage{
		Session: args.Session,
		Embed: &discordgo.MessageEmbed{
			Color:       util.ColorEmbedDefault,
			Description: "Test :^)",
		},
		UserID: args.User.ID,
		AcceptFunc: func(m *discordgo.Message) {
			fmt.Println("accepted")
		},
		DeclineFunc: func(m *discordgo.Message) {
			fmt.Println("declined")
		},
	}
	am.Send(args.Channel.ID)
	return nil
}
