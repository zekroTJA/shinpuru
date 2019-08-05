package commands

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

func (c *CmdTest) Exec(args *CommandArgs) error {
	// fmt.Println(args.Session.Channel("549575608074502174"))

	// return args.CmdHandler.bck.RestoreBackup(args.Guild.ID, "6499313859982409728", )
	return args.CmdHandler.bck.HardFlush(args.Guild.ID)
}
