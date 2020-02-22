package backupmodels

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

type BackupEntry struct {
	GuildID   string
	Timestamp time.Time
	FileID    string
}

type BackupObject struct {
	ID       string           `json:"id"`
	Guild    *BackupGuild     `json:"guild"`
	Channels []*BackupChannel `json:"channels"`
	Roles    []*BackupRole    `json:"roles"`
	Members  []*BackupMember  `json:"members"`
}

type BackupGuild struct {
	Name                        string `json:"name"`
	AfkChannelID                string `json:"afk_channel_id"`
	AfkTimeout                  int    `json:"afk_timeout"`
	VerificationLevel           int    `json:"verification_level"`
	DefaultMessageNotifications int    `json:"default_message_notifications"`
}

type BackupChannel struct {
	ID                   string                           `json:"id"`
	Name                 string                           `json:"name"`
	Topic                string                           `json:"topic"`
	Type                 int                              `json:"type"`
	NSFW                 bool                             `json:"nsfw"`
	Position             int                              `json:"position"`
	Bitrate              int                              `json:"bitrate"`
	UserLimit            int                              `json:"user_limit"`
	ParentID             string                           `json:"parent_id"`
	PermissionOverwrites []*discordgo.PermissionOverwrite `json:"permission_overwrites"`
}

type BackupRole struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Mentionable bool   `json:"mentionable"`
	Hoist       bool   `json:"hoist"`
	Color       int    `json:"color"`
	Position    int    `json:"position"`
	Permissions int    `json:"permissions"`
}

type BackupMember struct {
	ID    string   `json:"id"`
	Nick  string   `json:"nick"`
	Deaf  bool     `json:"deaf"`
	Mute  bool     `json:"mute"`
	Roles []string `json:"roles"`
}
