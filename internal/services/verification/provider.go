package verification

import "github.com/bwmarrin/discordgo"

type Provider interface {
	GetEnabled(guildID string) (ok bool, err error)
	SetEnabled(guildID string, enabled bool) (err error)

	IsVerified(userID string) (ok bool, err error)
	EnqueueVerification(member discordgo.Member) (err error)
	Verify(userID string) (err error)
	KickRoutine()
}
