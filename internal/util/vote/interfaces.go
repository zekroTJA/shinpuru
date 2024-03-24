package vote

import "github.com/bwmarrin/discordgo"

type Session interface {
	User(userID string, options ...discordgo.RequestOption) (st *discordgo.User, err error)
	ChannelMessageSendComplex(channelID string, data *discordgo.MessageSend, options ...discordgo.RequestOption) (st *discordgo.Message, err error)
	MessageReactionAdd(channelID, messageID, emojiID string, options ...discordgo.RequestOption) error
	ChannelMessageEditEmbed(channelID, messageID string, embed *discordgo.MessageEmbed, options ...discordgo.RequestOption) (*discordgo.Message, error)
	MessageReactionsRemoveAll(channelID, messageID string, options ...discordgo.RequestOption) error
}
