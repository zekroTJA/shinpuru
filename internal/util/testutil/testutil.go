package testutil

import "github.com/bwmarrin/discordgo"

func DiscordRestError(code int) error {
	return &discordgo.RESTError{
		Message: &discordgo.APIErrorMessage{
			Code: code,
		},
	}
}
