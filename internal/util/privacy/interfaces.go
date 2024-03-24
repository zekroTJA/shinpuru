package privacy

import (
	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/models"
	"github.com/zekroTJA/shinpuru/internal/services/backup/backupmodels"
)

type Session interface {
	User(userID string, options ...discordgo.RequestOption) (st *discordgo.User, err error)
	ChannelMessageSendComplex(channelID string, data *discordgo.MessageSend, options ...discordgo.RequestOption) (st *discordgo.Message, err error)
	MessageReactionAdd(channelID, messageID, emojiID string, options ...discordgo.RequestOption) error
	ChannelMessageEditEmbed(channelID, messageID string, embed *discordgo.MessageEmbed, options ...discordgo.RequestOption) (*discordgo.Message, error)
	MessageReactionsRemoveAll(channelID, messageID string, options ...discordgo.RequestOption) error
}

type Database interface {
	GetBackups(guildID string) ([]backupmodels.Entry, error)
	GetReportsGuildCount(guildID string) (int, error)
	GetReportsGuild(guildID string, offset, limit int) ([]models.Report, error)
	FlushGuildData(guildID string) error
	FlushUserData(userID string) (res map[string]int, err error)
}

type Storage interface {
	DeleteObject(bucketName, objectName string) error
}

type State interface {
	UserGuilds(id string) (res []string, err error)
	RemoveGuild(id string, dehydrate ...bool) error
	RemoveMember(guildID, memberID string) (err error)
	RemoveUser(id string) (err error)
}
