package karma

import (
	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/models"
)

type Logger interface {
	Debugf(guildID, message string, data ...interface{}) error
	Infof(guildID, message string, data ...interface{}) error
	Warnf(guildID, message string, data ...interface{}) error
	Errorf(guildID, message string, data ...interface{}) error
	Fatalf(guildID, message string, data ...interface{}) error
}

type Database interface {
	GetKarma(userID, guildID string) (int, error)
	UpdateKarma(userID, guildID string, diff int) error
	GetKarmaState(guildID string) (bool, error)
	GetKarmaPenalty(guildID string) (bool, error)
	IsKarmaBlockListed(guildID, userID string) (bool, error)
	GetKarmaRules(guildID string) ([]models.KarmaRule, error)
}

type Session interface {
	GuildMemberRoleAdd(guildID, userID, roleID string, options ...discordgo.RequestOption) (err error)
	GuildMemberRoleRemove(guildID, userID, roleID string, options ...discordgo.RequestOption) (err error)
	UserChannelCreate(recipientID string, options ...discordgo.RequestOption) (st *discordgo.Channel, err error)
	ChannelMessageSendEmbed(channelID string, embed *discordgo.MessageEmbed, options ...discordgo.RequestOption) (*discordgo.Message, error)
	GuildMemberDeleteWithReason(guildID, userID, reason string, options ...discordgo.RequestOption) (err error)
	GuildBanCreateWithReason(guildID, userID, reason string, days int, options ...discordgo.RequestOption) (err error)
}

type State interface {
	Guild(id string, hydrate ...bool) (v *discordgo.Guild, err error)
}
