package inits

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis/v8"
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekrotja/dgrs"
)

func InitState(container di.Container) (s *dgrs.State, err error) {
	session := container.Get(static.DiDiscordSession).(*discordgo.Session)
	rd := container.Get(static.DiRedis).(redis.Cmdable)

	return dgrs.New(dgrs.Options{
		RedisClient:    rd,
		DiscordSession: session,
		FetchAndStore:  true,
		Lifetimes: dgrs.Lifetimes{
			Message: 14 * 24 * time.Hour, // 14 Days
		},
	})
}
