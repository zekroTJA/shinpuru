package metrics

import (
	"context"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

type redisWatcher struct {
	sync.RWMutex
	redis redis.Cmdable
	m     map[string]float64
	timer *time.Ticker
}

func newRedisWatcher(redis redis.Cmdable) (rw *redisWatcher) {
	rw = &redisWatcher{
		m:     make(map[string]float64),
		timer: time.NewTicker(30 * time.Second),
		redis: redis,
	}

	go rw.loop()

	return
}

func (rw *redisWatcher) loop() {
	for {
		rw.collect()
		<-rw.timer.C
	}
}

func (rw *redisWatcher) collect() {
	res := rw.redis.Info(context.Background())
	if res.Err() != nil {
		logrus.WithError(res.Err()).Error("Failed collecting redis information")
		return
	}

	kcres := rw.redis.DBSize(context.Background())
	if res.Err() != nil {
		logrus.WithError(res.Err()).Error("Failed collecting redis db size")
		return
	}

	kc := float64(kcres.Val())

	rw.Lock()
	defer rw.Unlock()

	rw.m["key_count"] = kc

	for _, line := range strings.Split(res.Val(), "\n") {
		line = strings.TrimRight(line, "\r")

		if len(line) == 0 || line[0] == '#' {
			continue
		}

		kv := strings.Split(line, ":")
		if len(kv) < 2 {
			continue
		}

		key := kv[0]
		val, err := strconv.ParseFloat(kv[1], 64)
		if err != nil {
			continue
		}

		rw.m[key] = val
	}
}

func (rw *redisWatcher) Get(key string) float64 {
	rw.RLock()
	defer rw.RUnlock()

	return rw.m[key]
}
