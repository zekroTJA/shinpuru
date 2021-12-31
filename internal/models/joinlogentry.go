package models

import "time"

type JoinLogEntry struct {
	GuildID   string    `json:"guild_id"`
	UserID    string    `json:"user_id"`
	Tag       string    `json:"tag"`
	Created   time.Time `json:"account_created"`
	Timestamp time.Time `json:"timestamp"`
}
