package commands

import (
	"fmt"
	"net/url"

	"github.com/zekroTJA/shireikan"
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
	// const emojiID = "742680847429271654"

	// for _, e := range ctx.GetGuild().Emojis {
	// 	fmt.Printf("%+v\n", e)
	// }

	fmt.Println(url.QueryEscape(":myrunes:742680847429271654"))
	return ctx.GetSession().MessageReactionAdd(ctx.GetChannel().ID, ctx.GetMessage().ID,
		url.QueryEscape(":myrunes:742680847429271654"))
}
