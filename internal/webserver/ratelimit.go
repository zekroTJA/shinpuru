package webserver

import (
	"fmt"
	"time"

	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/zekroTJA/ratelimit"
	"github.com/zekroTJA/timedmap"
)

const (
	cleanupInterval = 15 * time.Minute
	entryLifetime   = 1 * time.Hour
)

// A RateLimitManager maintains all
// rate limiters for each connection.
type RateLimitManager struct {
	limits  *timedmap.TimedMap
	handler []*rateLimitHandler
}

type rateLimitHandler struct {
	id      int
	handler routing.Handler
}

// NewRateLimitManager creates a new instance
// of RateLimitManager.
func NewRateLimitManager() *RateLimitManager {
	return &RateLimitManager{
		limits:  timedmap.New(cleanupInterval),
		handler: make([]*rateLimitHandler, 0),
	}
}

// GetHandler returns a new afsthttp-routing
// handler which manages per-route and connection-
// based rate limiting.
// Rate limit information is added as 'X-RateLimit-Limit',
// 'X-RateLimit-Remaining' and 'X-RateLimit-Reset'
// headers.
// This handler aborts the execution of following
// handlers when rate limit is exceed and throws
// a json error body in combination with a 429
// status code.
func (rlm *RateLimitManager) GetHandler(limit time.Duration, burst int) routing.Handler {
	rlh := &rateLimitHandler{
		id: len(rlm.handler),
	}

	rlh.handler = func(ctx *routing.Context) error {
		limiterID := fmt.Sprintf("%d#%s",
			rlh.id, getIPAddr(ctx))
		ok, res := rlm.GetLimiter(limiterID, limit, burst).Reserve()

		ctx.Response.Header.Set("X-RateLimit-Limit", fmt.Sprintf("%d", res.Burst))
		ctx.Response.Header.Set("X-RateLimit-Remaining", fmt.Sprintf("%d", res.Remaining))
		ctx.Response.Header.Set("X-RateLimit-Reset", fmt.Sprintf("%d", res.Reset.Unix()))

		if !ok {
			ctx.Abort()
			ctx.Response.Header.SetContentType("application/json")
			ctx.SetStatusCode(429)
			ctx.SetBodyString(
				"{\n  \"code\": 429,\n  \"message\": \"you are being rate limited\"\n}")
		}

		return nil
	}

	rlm.handler = append(rlm.handler, rlh)

	return rlh.handler
}

// GetLimiter tries to get an existent limiter
// from the limiter map. If there is no limiter
// existent for this address, a new limiter
// will be created and added to the map.
func (rlm *RateLimitManager) GetLimiter(addr string, limit time.Duration, burst int) *ratelimit.Limiter {
	var ok bool
	var limiter *ratelimit.Limiter

	if rlm.limits.Contains(addr) {
		limiter, ok = rlm.limits.GetValue(addr).(*ratelimit.Limiter)
		if !ok {
			limiter = rlm.createLimiter(addr, limit, burst)
		}
	} else {
		limiter = rlm.createLimiter(addr, limit, burst)
	}

	return limiter
}

// createLimiter creates a new limiter and
// adds it to the limiters map by the passed
// address.
func (rlm *RateLimitManager) createLimiter(addr string, limit time.Duration, burst int) *ratelimit.Limiter {
	limiter := ratelimit.NewLimiter(limit, burst)
	rlm.limits.Set(addr, limiter, entryLifetime)
	return limiter
}
