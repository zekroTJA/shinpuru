package commands

import (
	"fmt"

	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shireikan"
)

type CmdTest struct {
}

func (c *CmdTest) GetInvokes() []string {
	return []string{"test"}
}

func (c *CmdTest) GetDescription() string {
	return "Just for testing purposes."
}

func (c *CmdTest) GetHelp() string {
	return ""
}

func (c *CmdTest) GetGroup() string {
	return shireikan.GroupEtc
}

func (c *CmdTest) GetDomainName() string {
	return "sp.test"
}

func (c *CmdTest) GetSubPermissionRules() []shireikan.SubPermission {
	return nil
}

func (c *CmdTest) IsExecutableInDMChannels() bool {
	return true
}

func (c *CmdTest) Exec(ctx shireikan.Context) error {
	h := ctx.GetObject(static.DiLegacyCommandHandler).(shireikan.Handler)

	for c := range h.GetCommandMap() {
		fmt.Println(c)
	}

	return nil
}
