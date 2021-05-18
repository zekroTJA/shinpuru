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
	// gl := ctx.GetObject(static.DiGuildLog).(guildlog.Logger)

	// severity, err := ctx.GetArgs().Get(0).AsInt()
	// if err != nil {
	// 	return err
	// }

	// var f func(string, string, ...interface{}) error

	// switch severity {
	// case 0:
	// 	f = gl.Debugf
	// case 1:
	// 	f = gl.Infof
	// case 2:
	// 	f = gl.Warnf
	// case 3:
	// 	f = gl.Errorf
	// case 4:
	// 	f = gl.Fatalf
	// }

	// return f(ctx.GetGuild().ID, strings.Join(ctx.GetArgs()[1:], " "))

	// db := ctx.GetObject(static.DiDatabase).(database.Database)
	// st := ctx.GetObject(static.DiObjectStorage).(storage.Storage)

	// return util.FlushAllGuildData(db, st, ctx.GetGuild().ID)

	fmt.Println(static.AdditionalPermissions)
	return nil
}
