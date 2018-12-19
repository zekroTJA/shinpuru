package commands

import (
	"../util"
)

type CmdHandler struct {
	registeredCmds         map[string]Command
	registeredCmdInstances []Command
}

func NewCmdHandler() *CmdHandler {
	return &CmdHandler{
		registeredCmds:         make(map[string]Command),
		registeredCmdInstances: make([]Command, 0),
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
