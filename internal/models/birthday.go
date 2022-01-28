package models

import "time"

type Birthday struct {
	GuildID  string    `json:"guildid"`
	UserID   string    `json:"userid"`
	Date     time.Time `json:"date"`
	ShowYear bool      `json:"showyear"`
}
