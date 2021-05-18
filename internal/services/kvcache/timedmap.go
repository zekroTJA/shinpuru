package kvcache

import (
	"time"

	"github.com/zekroTJA/timedmap"
)

type timedmapCache struct {
	tm *timedmap.TimedMap
}

func NewTimedmapCache(tickTime time.Duration) Provider {
	return &timedmapCache{
		tm: timedmap.New(tickTime),
	}
}

func (t *timedmapCache) Get(key string) interface{} {
	return t.tm.GetValue(key)
}

func (t *timedmapCache) Set(key string, v interface{}, lifetime time.Duration) {
	t.tm.Set(key, v, lifetime)
}

func (t *timedmapCache) Del(key string) {
	t.tm.Remove(key)
}
