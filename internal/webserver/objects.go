package webserver

import (
	"github.com/bwmarrin/discordgo"
)

type User struct {
	*discordgo.User

	AvatarURL string `json:"avatar_url"`
}
