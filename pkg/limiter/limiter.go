package limiter

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

const (
	xRateLimitLimit     = "X-RateLimit-Limit"
	xRateLimitRemaining = "X-RateLimit-Remaining"
	xRateLimitReset     = "X-RateLimit-Reset"
)

var defConfigInstance = Config{
	Duration:        1 * time.Second,
	Burst:           5,
	CleanupInterval: 10 * time.Minute,
	KeyGenerator: func(ctx *fiber.Ctx) string {
		return ctx.IP()
	},
	OnLimitReached: func(ctx *fiber.Ctx) error {
		return fiber.ErrTooManyRequests
	},
}

type Config struct {
	Duration        time.Duration `json:"restoration"`
	Burst           int           `json:"burst"`
	CleanupInterval time.Duration `json:"cleanupinterval"`
	KeyGenerator    func(*fiber.Ctx) string
	OnLimitReached  fiber.Handler
}

func New(config ...Config) fiber.Handler {
	cfg := defConfig(config)
	mgr := newManager(cfg.CleanupInterval, cfg.Duration, cfg.Burst)
	burstS := strconv.Itoa(cfg.Burst)

	return func(ctx *fiber.Ctx) (err error) {
		key := cfg.KeyGenerator(ctx)
		rl := mgr.retrieve(key)
		ok, res := rl.Reserve()

		ctx.Set(xRateLimitLimit, burstS)
		ctx.Set(xRateLimitRemaining, strconv.Itoa(res.Remaining))
		ctx.Set(xRateLimitReset, strconv.Itoa(int(res.Reset.Unix())))

		if ok {
			ctx.Next()
		} else {
			err = cfg.OnLimitReached(ctx)
		}
		return
	}
}

func defConfig(configs []Config) (cfg Config) {
	cfg = defConfigInstance

	if len(configs) == 0 {
		return
	}
	pcfg := configs[0]

	if pcfg.Duration > 0 {
		cfg.Duration = pcfg.Duration
	}
	if pcfg.Burst > 0 {
		cfg.Burst = pcfg.Burst
	}
	if pcfg.CleanupInterval > 0 {
		cfg.CleanupInterval = pcfg.CleanupInterval
	}
	if pcfg.KeyGenerator != nil {
		cfg.KeyGenerator = pcfg.KeyGenerator
	}
	if pcfg.OnLimitReached != nil {
		cfg.OnLimitReached = pcfg.OnLimitReached
	}

	return
}
