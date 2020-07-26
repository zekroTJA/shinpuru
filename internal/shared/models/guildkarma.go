package models

type GuildKarma struct {
	UserID  string `json:"user_id"`
	GuildID string `json:"guild_id"`
	Value   int    `json:"value"`
}
