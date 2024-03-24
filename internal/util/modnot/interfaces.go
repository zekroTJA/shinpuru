package modnot

import "github.com/bwmarrin/discordgo"

type Session interface {
	ChannelMessageSendEmbed(channelID string, embed *discordgo.MessageEmbed, options ...discordgo.RequestOption) (*discordgo.Message, error)
}

type Database interface {
	GetGuildModNot(guildID string) (string, error)
}
