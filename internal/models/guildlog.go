package models

import (
	"time"

	"github.com/bwmarrin/snowflake"
)

type GuildLogSeverity int

const (
	GLAll GuildLogSeverity = iota - 1
	GLDebug
	GLInfo
	GLWarn
	GLError
	GLFatal
)

type GuildLogEntry struct {
	ID        snowflake.ID     `json:"id"`
	GuildID   string           `json:"guildid"`
	Module    string           `json:"module"`
	Message   string           `json:"message"`
	Severity  GuildLogSeverity `json:"severity"`
	Timestamp time.Time        `json:"timestamp"`
}
