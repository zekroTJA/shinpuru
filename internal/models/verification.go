package models

import "time"

type VerificationQueueEntry struct {
	GuildID   string    `json:"guildid,omitempty"`
	UserID    string    `json:"userid"`
	Timestamp time.Time `json:"timestamp"`
}
