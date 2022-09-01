package discordutil

import (
	"github.com/bwmarrin/discordgo"
)

// WrapHandler takes a handler function which uses ISession for
// the Discord session instance and T as event playload and returns
// a valid handler function using *discordgo.Session for the session
// and T as event payload.
//
// This can be used for handler functions which are unit tested
// and therefore need to use ISession to pass in mocked session
// instances.
func WrapHandler[T any](f func(s ISession, e T)) func(s *discordgo.Session, h T) {
	return func(s *discordgo.Session, h T) {
		f(s, h)
	}
}
