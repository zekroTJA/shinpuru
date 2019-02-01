package commands

import (
	"fmt"
	"time"

	"github.com/zekroTJA/timedmap"
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
	t := timedmap.New(1 * time.Second)
	t.Set("test", "hey", 1*time.Second)
	time.AfterFunc(2*time.Second, func() {
		fmt.Println(t.GetValue("test").(string))
	})
	return nil
}
