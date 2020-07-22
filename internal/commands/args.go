package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/pkg/ctypes"
)

// CommandArgs wraps instances about a command
// execution like channel, user, guild, message,
// command arguments, session and the command
// handler. This is then passed to the commands
// execution handler.
type CommandArgs struct {
	Channel    *discordgo.Channel
	User       *discordgo.User
	Guild      *discordgo.Guild
	Message    *discordgo.Message
	Args       ctypes.StringArray
	Session    *discordgo.Session
	CmdHandler *CmdHandler
	IsDM       bool
}
