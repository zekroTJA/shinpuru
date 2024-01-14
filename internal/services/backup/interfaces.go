package backup

import (
	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/services/backup/backupmodels"
	"io"
)

type Logger interface {
	Debugf(guildID, message string, data ...interface{}) error
	Infof(guildID, message string, data ...interface{}) error
	Warnf(guildID, message string, data ...interface{}) error
	Errorf(guildID, message string, data ...interface{}) error
	Fatalf(guildID, message string, data ...interface{}) error
}

type Database interface {
	AddBackup(guildID, fileID string) error
	DeleteBackup(guildID, fileID string) error
	GetBackups(guildID string) ([]backupmodels.Entry, error)
	GetGuilds() ([]string, error)
}

type Storage interface {
	PutObject(bucketName, objectName string, reader io.Reader, objectSize int64, mimeType string) error
	GetObject(bucketName, objectName string) (io.ReadCloser, int64, error)
	DeleteObject(bucketName, objectName string) error
}

type Session interface {
	GuildEdit(guildID string, g *discordgo.GuildParams, options ...discordgo.RequestOption) (st *discordgo.Guild, err error)
	GuildRoles(guildID string, options ...discordgo.RequestOption) (st []*discordgo.Role, err error)
	GuildRoleCreate(guildID string, data *discordgo.RoleParams, options ...discordgo.RequestOption) (st *discordgo.Role, err error)
	GuildRoleReorder(guildID string, roles []*discordgo.Role, options ...discordgo.RequestOption) (st []*discordgo.Role, err error)
	Channel(channelID string, options ...discordgo.RequestOption) (st *discordgo.Channel, err error)
	GuildChannelCreateComplex(guildID string, data discordgo.GuildChannelCreateData, options ...discordgo.RequestOption) (st *discordgo.Channel, err error)
	ChannelEditComplex(channelID string, data *discordgo.ChannelEdit, options ...discordgo.RequestOption) (st *discordgo.Channel, err error)
	GuildMember(guildID, userID string, options ...discordgo.RequestOption) (st *discordgo.Member, err error)
	GuildMemberEdit(guildID, userID string, data *discordgo.GuildMemberParams, options ...discordgo.RequestOption) (st *discordgo.Member, err error)
	GuildMemberNickname(guildID, userID, nickname string, options ...discordgo.RequestOption) (err error)
	GuildRoleDelete(guildID, roleID string, options ...discordgo.RequestOption) (err error)
	ChannelDelete(channelID string, options ...discordgo.RequestOption) (st *discordgo.Channel, err error)
}

type State interface {
	Guild(id string, hydrate ...bool) (v *discordgo.Guild, err error)
	Channels(guildID string, forceFetch ...bool) (v []*discordgo.Channel, err error)
	Members(guildID string, forceFetch ...bool) (v []*discordgo.Member, err error)
}
