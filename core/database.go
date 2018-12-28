package core

import (
	"errors"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/util"
)

var ErrDatabaseNotFound = errors.New("value not found")

type Database interface {
	Connect(credentials ...interface{}) error
	Close()

	GetGuildPrefix(guildID string) (string, error)
	SetGuildPrefix(guildID, newPrefix string) error

	GetGuildAutoRole(guildID string) (string, error)
	SetGuildAutoRole(guildID, autoRoleID string) error

	GetGuildPermissions(guildID string) (map[string]int, error)
	SetGuildRolePermission(guildID, roleID string, permLvL int) error

	AddReport(rep *util.Report) error
	GetReportsGuild(guildID string) ([]*util.Report, error)

	GetMemberPermissionLevel(s *discordgo.Session, guildID string, memberID string) (int, error)

	GetSetting(setting string) (string, error)
	SetSetting(setting, value string) error
}

func IsErrDatabaseNotFound(err error) bool {
	return err == ErrDatabaseNotFound
}
