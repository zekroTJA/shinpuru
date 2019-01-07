package commands

import (
	"fmt"

	"github.com/zekroTJA/shinpuru/core"
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

func (c *CmdTest) SetPermission(permLvl int) {}

func (c *CmdTest) Exec(args *CommandArgs) error {
	resp, _ := core.HTTPRequest("GET", util.DiscordAPIEndpoint+"/users/"+args.User.ID, map[string]string{
		"Authorization": "Bot " + args.CmdHandler.config.Discord.Token,
	}, nil)
	fmt.Println(resp.BodyAsMap())
	return nil
}
