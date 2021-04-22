package backupmodels

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

// Entry wraps the database entry of a backup.
type Entry struct {
	GuildID   string    `json:"guild_id"`
	Timestamp time.Time `json:"timestamp"`
	FileID    string    `json:"file_id"`
}

// Object wraps a backup structure with an unique
// ID (snowflake), timestamp of creation, the guild
// properties and each channels, roles and members
// properties.
type Object struct {
	ID        string     `json:"id"`
	Timestamp time.Time  `json:"timestamp"`
	Guild     *Guild     `json:"guild"`
	Channels  []*Channel `json:"channels"`
	Roles     []*Role    `json:"roles"`
	Members   []*Member  `json:"members"`
}

// Guild contains general properties of the guild.
type Guild struct {
	Name                        string `json:"name"`
	AfkChannelID                string `json:"afk_channel_id"`
	AfkTimeout                  int    `json:"afk_timeout"`
	VerificationLevel           int    `json:"verification_level"`
	DefaultMessageNotifications int    `json:"default_message_notifications"`
}

// Channel contains general properties of the channel.
type Channel struct {
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

// Role contains general properties of the role.
type Role struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Mentionable bool   `json:"mentionable"`
	Hoist       bool   `json:"hoist"`
	Color       int    `json:"color"`
	Position    int    `json:"position"`
	Permissions int64  `json:"permissions"`
}

// Member contains general properties of the member.
type Member struct {
	ID    string   `json:"id"`
	Nick  string   `json:"nick"`
	Deaf  bool     `json:"deaf"`
	Mute  bool     `json:"mute"`
	Roles []string `json:"roles"`
}
