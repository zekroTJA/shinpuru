package inits

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis/v8"
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekrotja/dgrs"
)

const day = 24 * time.Hour

func InitState(container di.Container) (s *dgrs.State, err error) {
	session := container.Get(static.DiDiscordSession).(*discordgo.Session)
	rd := container.Get(static.DiRedis).(*redis.Client)

	return dgrs.New(dgrs.Options{
		RedisClient:    rd,
		DiscordSession: session,
		FetchAndStore:  true,
		Lifetimes: dgrs.Lifetimes{
			Message:  30 * day,
			Member:   30 * day,
			User:     30 * day,
			Presence: 30 * day,
		},
	})
}
