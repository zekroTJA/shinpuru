// Package limiter provides a fiber middleware
// for a bucket based request rate limiter.
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
	Next: nil,
}

// Config provides configuration values
// for the rate limiter middleware.
type Config struct {
	// Duration until a new token is added
	// to the bucket.
	//
	// Default: 1 * time.Second
	Duration time.Duration `json:"restoration"`
	// Burst is the amount of tokens which
	// can be contained in a bucket (maximum
	// amount of requests which can be done
	// simultaniously).
	//
	// Default: 5
	Burst int `json:"burst"`
	// CleanupInterval specifies the interval
	// duration for the underlying timedmap
	// holding all ratelimiter mappings.
	//
	// Default: 10 * time.Minute
	CleanupInterval time.Duration `json:"cleanupinterval"`
	// KeyGenerator is the function used to
	// get a unique, user bound key from a
	// request context.
	//
	// Default: func(ctx *fiber.Ctx) string { return ctx.IP() }
	KeyGenerator func(*fiber.Ctx) string
	// OnLimitReached is the handler function
	// executed when the rate limit was hit.
	//
	// Default: func(ctx *fiber.Ctx) error { return fiber.ErrTooManyRequests }
	OnLimitReached fiber.Handler
	// Next specifies a function which is called
	// before the middleware is executed. If the
	// function is set und returns true, the
	// middleware is skipped.
	Next func(c *fiber.Ctx) bool
}

// New initializes new rate limiter middleware
// instance and returns the middleware handler
// which can be registered to a fiber router.
//
// You can pass a custom config variables via
// the Config parameter. Values which are not
// set are taken from the default config.
func New(config ...Config) fiber.Handler {
	cfg := defConfig(config)
	mgr := newManager(cfg.CleanupInterval, cfg.Duration, cfg.Burst)
	burstS := strconv.Itoa(cfg.Burst)

	return func(ctx *fiber.Ctx) (err error) {
		if cfg.Next != nil && cfg.Next(ctx) {
			return ctx.Next()
		}

		key := cfg.KeyGenerator(ctx)
		rl := mgr.retrieve(key)
		ok, res := rl.Reserve()

		ctx.Set(xRateLimitLimit, burstS)
		ctx.Set(xRateLimitRemaining, strconv.Itoa(res.Remaining))
		ctx.Set(xRateLimitReset, strconv.Itoa(int(res.Reset.Unix())))

		if ok {
			err = ctx.Next()
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
	if pcfg.Next != nil {
		cfg.Next = pcfg.Next
	}

	return
}
