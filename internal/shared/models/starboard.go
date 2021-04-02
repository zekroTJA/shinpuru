package models

import (
	"encoding/base64"
	"encoding/json"
)

type StarboardConfig struct {
	GuildID   string
	ChannelID string
	Threshold int
	EmojiID   string
}

type StarboardEntry struct {
	MessageID   string
	StarboardID string
	GuildID     string
	ChannelID   string
	AuthorID    string
	Content     string
	MediaURLs   []string
	Score       int
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
