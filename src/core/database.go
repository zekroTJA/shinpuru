package core

import "errors"

var ErrDatabaseNotFound = errors.New("value not found")

type Database interface {
	Connect(credentials ...interface{}) error
	Close()

	GetGuildPrefix(guildID string) (string, error)
	SetGuildPrefix(guildID, newPrefix string) error
	GetMemberPermissionLevel(guildID string, memberID string) (int, error)
}

func IsErrDatabaseNotFound(err error) bool {
	return err == ErrDatabaseNotFound
}
