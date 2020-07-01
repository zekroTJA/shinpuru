package commands

import "github.com/bwmarrin/discordgo"

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
	Args       []string
	Session    *discordgo.Session
	CmdHandler *CmdHandler
}
