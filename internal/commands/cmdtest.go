package commands

import (
	"fmt"

	"github.com/zekroTJA/shinpuru/pkg/roleutil"
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

func (c *CmdTest) GetDomainName() string {
	return "sp.test"
}

func (c *CmdTest) GetSubPermissionRules() []SubPermission {
	return nil
}

func (c *CmdTest) IsExecutableInDMChannels() bool {
	return true
}

func (c *CmdTest) Exec(args *CommandArgs) error {
	roles := args.Guild.Roles
	roleutil.SortRoles(roles, false)
	for i, r := range roles {
		fmt.Printf("%d - %s\n", i, r.Name)
	}
	return nil
}
