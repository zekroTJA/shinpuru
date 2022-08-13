package testutil

import (
	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/mock"
)

func DiscordRestError(code int) error {
	return &discordgo.RESTError{
		Message: &discordgo.APIErrorMessage{
			Code: code,
		},
	}
}

func Nil[T any]() any {
	return mock.MatchedBy(func(v *T) bool {
		return v == nil
	})
}
