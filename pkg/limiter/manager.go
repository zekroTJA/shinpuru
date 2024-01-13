package limiter

import (
	"time"

	"github.com/zekroTJA/ratelimit"
	"github.com/zekroTJA/timedmap"
	"github.com/zekrotja/safepool"
)

type manager struct {
	lifetime time.Duration
	tm       *timedmap.TimedMap
	pool     safepool.SafePool[*safepool.ResetWrapper[*ratelimit.Limiter]]
}

func newManager(cleanupInterval, duration time.Duration, burst int) *manager {
	return &manager{
		lifetime: time.Duration(burst) * duration,
		tm:       timedmap.New(cleanupInterval),
		pool: safepool.New(func() *safepool.ResetWrapper[*ratelimit.Limiter] {
			return safepool.Wrap(ratelimit.NewLimiter(duration, burst), func(v *ratelimit.Limiter) {
				v.Reset()
			})
		}),
	}
}

func (m *manager) retrieve(key string) *ratelimit.Limiter {
	rl, ok := m.tm.GetValue(key).(*ratelimit.Limiter)
	if ok {
		m.tm.SetExpires(key, m.lifetime)
	} else {
		pv := m.pool.Get()
		rl = pv.Inner
		m.tm.Set(key, rl, m.lifetime, func(v interface{}) {
			m.pool.Put(pv)
		})
	}

	return rl
}
