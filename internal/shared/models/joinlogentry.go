package models

import "time"

type JoinLogEntry struct {
	GuildID   string
	UserID    string
	Tag       string
	Timestamp time.Time
}
