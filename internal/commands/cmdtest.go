package commands

import (
	"fmt"

	"github.com/zekroTJA/shinpuru/internal/core/database"
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
	db, _ := ctx.GetObject("db").(database.Database)
	fmt.Println(db.AddKarmaBlockList(ctx.GetGuild().ID, ctx.GetUser().ID))
	fmt.Println(db.IsKarmaBlockListed(ctx.GetGuild().ID, ctx.GetUser().ID))
	fmt.Println(db.RemoveKarmaBlockList(ctx.GetGuild().ID, ctx.GetUser().ID))
	fmt.Println(db.IsKarmaBlockListed(ctx.GetGuild().ID, ctx.GetUser().ID))

	return nil
}
