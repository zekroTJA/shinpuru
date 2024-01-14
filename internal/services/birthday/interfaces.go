package birthday

import (
	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/models"
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
	GetBirthdays(guildID string) ([]models.Birthday, error)
	SetBirthday(m models.Birthday) error
	DeleteBirthday(guildID, userID string) error
	GetGuildBirthdayChan(guildID string) (string, error)
	SetGuildBirthdayChan(guildID string, chanID string) error
}

type State interface {
	Guilds() (v []*discordgo.Guild, err error)
	Channel(id string) (v *discordgo.Channel, err error)
	Member(guildID, memberID string, forceNoFetch ...bool) (v *discordgo.Member, err error)
}

type Session interface {
	ChannelMessageSendEmbed(channelID string, embed *discordgo.MessageEmbed, options ...discordgo.RequestOption) (*discordgo.Message, error)
}
