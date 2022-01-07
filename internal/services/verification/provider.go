package verification

type Provider interface {
	IsVerified(userID string) (ok bool, err error)
	EnqueueVerification(guildID, userID string) (err error)
	Verify(userID string) (err error)
}
