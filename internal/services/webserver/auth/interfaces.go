package auth

import (
	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/util"
	"time"
)

type TimeProvider interface {
	Now() time.Time
}

type Session interface {
	util.MessageSession

	User(userID string, options ...discordgo.RequestOption) (st *discordgo.User, err error)
}

type Database interface {
	GetUserByRefreshToken(token string) (string, time.Time, error)
	SetUserRefreshToken(userID, token string, expires time.Time) error
	RevokeUserRefreshToken(userID string) error
}
