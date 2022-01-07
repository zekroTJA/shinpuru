package verification

type Provider interface {
	GetEnabled(guildID string) (ok bool, err error)
	SetEnabled(guildID string, enabled bool) (err error)

	IsVerified(userID string) (ok bool, err error)
	EnqueueVerification(guildID, userID string) (err error)
	Verify(userID string) (err error)
	KickRoutine()
}
