package limiter

import (
	"sync"
	"time"

	"github.com/zekroTJA/ratelimit"
	"github.com/zekroTJA/timedmap"
)

type manager struct {
	lifetime time.Duration
	tm       *timedmap.TimedMap
	pool     *sync.Pool
}

func newManager(cleanupInterval, duration time.Duration, burst int) *manager {
	return &manager{
		lifetime: time.Duration(burst) * duration,
		tm:       timedmap.New(cleanupInterval),
		pool: &sync.Pool{
			New: func() interface{} {
				return ratelimit.NewLimiter(duration, burst)
			},
		},
	}
}

func (m *manager) retrieve(key string) *ratelimit.Limiter {
	rl, ok := m.tm.GetValue(key).(*ratelimit.Limiter)
	if ok {
		m.tm.SetExpires(key, m.lifetime)
	} else {
		rl = m.pool.Get().(*ratelimit.Limiter)
		rl.Reset()
		m.tm.Set(key, rl, m.lifetime, func(v interface{}) {
			m.pool.Put(v)
		})
	}

	return rl
}
