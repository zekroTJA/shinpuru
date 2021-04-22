package models

import (
	"encoding/base64"
	"encoding/json"
)

type StarboardSortBy int

const (
	StarboardSortByLatest StarboardSortBy = iota
	StarboardSortByMostRated
)

type StarboardConfig struct {
	GuildID   string
	ChannelID string
	Threshold int
	EmojiID   string
	KarmaGain int
}

type StarboardEntry struct {
	MessageID   string   `json:"message_id"`
	StarboardID string   `json:"starboard_id"`
	GuildID     string   `json:"guild_id"`
	ChannelID   string   `json:"channel_id"`
	AuthorID    string   `json:"author_id"`
	Content     string   `json:"content"`
	MediaURLs   []string `json:"media_urls"`
	Score       int      `json:"score"`
	Deleted     bool     `json:"-"`
}

func (e *StarboardEntry) MediaURLsEncoded() string {
	res, _ := json.Marshal(e.MediaURLs)
	return base64.StdEncoding.EncodeToString(res)
}

func (e *StarboardEntry) SetMediaURLs(encoded string) (err error) {
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return
	}
	e.MediaURLs = make([]string, 0)
	err = json.Unmarshal(data, &e.MediaURLs)
	return
}
