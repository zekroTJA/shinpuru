package commands

import "github.com/bwmarrin/discordgo"

type CommandArgs struct {
	Channel    *discordgo.Channel
	User       *discordgo.User
	Guild      *discordgo.Guild
	Message    *discordgo.Message
	Args       []string
	Session    *discordgo.Session
	CmdHandler *CmdHandler
}
