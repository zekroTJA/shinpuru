package inits

import (
	"reflect"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis/v8"
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/services/config"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/timeutil"
	"github.com/zekrotja/dgrs"
	"github.com/zekrotja/rogu/log"
)

func getLifetimes(cfg config.Provider) (dgrs.Lifetimes, bool, error) {
	lifetimes := cfg.Config().Cache.Lifetimes

	var target dgrs.Lifetimes

	vlt := reflect.ValueOf(lifetimes)
	vtg := reflect.ValueOf(&target)

	set := false

	for i := 0; i < vlt.NumField(); i++ {
		ds := vlt.Field(i).String()
		if ds == "" {
			continue
		}

		d, err := timeutil.ParseDuration(ds)
		if err != nil {
			return dgrs.Lifetimes{}, false, err
		}

		if d == 0 {
			continue
		}

		vtg.Elem().FieldByName(vlt.Type().Field(i).Name).Set(reflect.ValueOf(d))
		set = true
	}

	return target, set, nil
}

func InitState(container di.Container) (s *dgrs.State, err error) {
	session := container.Get(static.DiDiscordSession).(*discordgo.Session)
	rd := container.Get(static.DiRedis).(*redis.Client)
	cfg := container.Get(static.DiConfig).(config.Provider)

	lf, set, err := getLifetimes(cfg)
	if err != nil {
		return nil, err
	}

	if !set {
		lf.General = 7 * 24 * time.Hour
		log.Tagged("State").Warn().
			Field("d", lf.General).
			Msg("No cache lifetimes have been set; applying default duration for all fields")
	}

	// When a value for `General` is set, all 0 value durations
	// will be set to the vaue of `General`. So it is effectively
	// the default caching duration, if not further specified.
	lf.OverrrideZero = true

	return dgrs.New(dgrs.Options{
		RedisClient:    rd,
		DiscordSession: session,
		FetchAndStore:  true,
		Lifetimes:      lf,
	})
}
