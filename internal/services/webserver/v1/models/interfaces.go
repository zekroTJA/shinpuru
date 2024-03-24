package models

import "github.com/zekroTJA/shinpuru/internal/services/backup/backupmodels"

type Database interface {
	GetKarma(userID, guildID string) (int, error)
	GetKarmaSum(userID string) (int, error)
	GetGuildBackup(guildID string) (bool, error)
	GetBackups(guildID string) ([]backupmodels.Entry, error)
	GetGuildInviteBlock(guildID string) (string, error)
}
