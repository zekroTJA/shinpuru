package guildlog

import (
	"github.com/zekroTJA/shinpuru/internal/models"
	"time"
)

type TimeProvider interface {
	Now() time.Time
}

type Database interface {
	GetGuildLogDisable(guildID string) (bool, error)
	AddGuildLogEntry(entry models.GuildLogEntry) error
}
