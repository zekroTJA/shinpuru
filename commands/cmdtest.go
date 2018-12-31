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
	filter := func(m *discordgo.Message) bool {
		return m.Author.ID == args.User.ID && m.Content == "a"
	}
	options := &util.MessageCollectorOptions{
		MaxMatches:         3,
		DeleteMatchesAfter: true,
	}
	mc, err := util.NewMessageCollector(args.Session, args.Channel.ID, filter, options)
	if err != nil {
		return err
	}

	mc.OnColelcted(func(msg *discordgo.Message, c *util.MessageCollector) {
		fmt.Println("Collected: ", msg.Content)
	})
	mc.OnMatched(func(msg *discordgo.Message, c *util.MessageCollector) {
		fmt.Println("Matched: ", msg.Content)
	})
	mc.OnClosed(func(reason string, c *util.MessageCollector) {
		fmt.Println(reason, len(c.CollectedMessages), len(c.CollectedMatches))
	})

	return nil
}
