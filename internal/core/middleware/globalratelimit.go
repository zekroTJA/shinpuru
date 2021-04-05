package middleware

import (
	"fmt"
	"time"

	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shireikan"
	"github.com/zekroTJA/shireikan/middleware/ratelimit"
	"github.com/zekroTJA/timedmap"
)

type GlobalRateLimitMiddleware struct {
	limit           time.Duration
	burst           int
	limiterDuration time.Duration

	limiters *timedmap.TimedMap
}

func NewGlobalRateLimitMiddleware(burst int, limit time.Duration) *GlobalRateLimitMiddleware {
	return &GlobalRateLimitMiddleware{
		limit:           limit,
		burst:           burst,
		limiterDuration: time.Duration(burst) * limit,
		limiters:        timedmap.New(15 * time.Minute),
	}
}

func (mw *GlobalRateLimitMiddleware) Handle(
	cmd shireikan.Command,
	ctx shireikan.Context,
	layer shireikan.MiddlewareLayer,
) (ok bool, err error) {
	rl := mw.getLimiter(ctx)

	ok, next := rl.Take()
	if ok {
		return
	}

	err = util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
		fmt.Sprintf("You are being rate limited. You need to wait %s until you can "+
			"use commands again.", next.Round(time.Second).String())).
		DeleteAfter(10 * time.Second).
		Error()

	return false, nil
}

func (mw *GlobalRateLimitMiddleware) GetLayer() shireikan.MiddlewareLayer {
	return shireikan.LayerBeforeCommand
}

func (mw *GlobalRateLimitMiddleware) getLimiter(ctx shireikan.Context) (rl *ratelimit.Limiter) {
	key := ctx.GetUser().ID

	var ok bool
	if rl, ok = mw.limiters.GetValue(key).(*ratelimit.Limiter); ok {
		mw.limiters.SetExpire(key, mw.limiterDuration)
		return
	}

	rl = ratelimit.NewLimiter(mw.burst, mw.limit)
	mw.limiters.Set(key, rl, mw.limiterDuration)

	return
}
