package verification

import (
	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/models"
	"github.com/zekroTJA/shinpuru/internal/util"
	"time"
)

type Logger interface {
	Debugf(guildID, message string, data ...interface{}) error
	Infof(guildID, message string, data ...interface{}) error
	Warnf(guildID, message string, data ...interface{}) error
	Errorf(guildID, message string, data ...interface{}) error
	Fatalf(guildID, message string, data ...interface{}) error
}

type TimeProvider interface {
	Now() time.Time
}

type Database interface {
	GetGuildVerificationRequired(guildID string) (bool, error)
	SetGuildVerificationRequired(guildID string, enable bool) error
	GetVerificationQueue(guildID, userID string) ([]models.VerificationQueueEntry, error)
	FlushVerificationQueue(guildID string) error
	AddVerificationQueue(e models.VerificationQueueEntry) error
	RemoveVerificationQueue(guildID, userID string) (bool, error)
	GetUserVerified(userID string) (bool, error)
	SetUserVerified(userID string, enabled bool) error
	GetGuildJoinMsg(guildID string) (string, string, error)
}

type Session interface {
	util.MessageSession

	GuildMemberTimeout(guildID string, userID string, until *time.Time, options ...discordgo.RequestOption) (err error)
	GuildMemberDelete(guildID, userID string, options ...discordgo.RequestOption) (err error)
	UserChannelCreate(recipientID string, options ...discordgo.RequestOption) (st *discordgo.Channel, err error)
	ChannelMessageSendComplex(channelID string, data *discordgo.MessageSend, options ...discordgo.RequestOption) (st *discordgo.Message, err error)
}
