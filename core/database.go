package core

import (
	"errors"

	"github.com/bwmarrin/discordgo"
)

var ErrDatabaseNotFound = errors.New("value not found")

type Database interface {
	Connect(credentials ...interface{}) error
	Close()

	GetGuildPrefix(guildID string) (string, error)
	SetGuildPrefix(guildID, newPrefix string) error

	GetGuildPermissions(guildID string) (map[string]int, error)
	SetGuildRolePermission(guildID, roleID string, permLvL int) error

	GetMemberPermissionLevel(s *discordgo.Session, guildID string, memberID string) (int, error)
}

func IsErrDatabaseNotFound(err error) bool {
	return err == ErrDatabaseNotFound
}
