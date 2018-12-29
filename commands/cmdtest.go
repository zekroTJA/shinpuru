package commands

import (
	"time"

	"github.com/bwmarrin/discordgo"
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
	args.Session.ChannelMessageSendEmbed(args.Channel.ID, &discordgo.MessageEmbed{
		Description: "HURENSOHN",
		Timestamp:   time.Now().Format(time.RFC3339),
	})
	return nil
}
