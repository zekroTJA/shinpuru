package commands

import (
	"github.com/zekroTJA/shinpuru/core"
	"github.com/zekroTJA/shinpuru/util"
)

type CmdHandler struct {
	registeredCmds         map[string]Command
	registeredCmdInstances []Command
	db                     core.Database
	config                 *core.Config
}

func NewCmdHandler(db core.Database, config *core.Config) *CmdHandler {
	return &CmdHandler{
		registeredCmds:         make(map[string]Command),
		registeredCmdInstances: make([]Command, 0),
		db:     db,
		config: config,
	}
}

func (c *CmdHandler) RegisterCommand(cmd Command) {
	c.registeredCmdInstances = append(c.registeredCmdInstances, cmd)
	for _, invoke := range cmd.GetInvokes() {
		if _, ok := c.registeredCmds[invoke]; ok {
			util.Log.Warningf("Command invoke '%s' was registered more than once!", invoke)
		}
		c.registeredCmds[invoke] = cmd
	}
}

func (c *CmdHandler) GetCommand(invoke string) (Command, bool) {
	cmd, ok := c.registeredCmds[invoke]
	return cmd, ok
}

func (c *CmdHandler) UpdateCommandPermissions(perms map[string]int) {
	for k, v := range perms {
		if cmd, ok := c.registeredCmds[k]; ok {
			cmd.SetPermission(v)
		}
	}
}
