package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
)

func Set[T any](r *RedisMiddleware, key string, val T) error {
	return r.client.Set(context.Background(), key, val, 0).Err()
}

func Get[T any](
	r *RedisMiddleware,
	key string,
	fallback func() (T, error),
) (val T, err error) {
	err = r.client.Get(context.Background(), key).Scan(&val)
	if err == redis.Nil {
		if val, err = fallback(); err != nil {
			return
		}
		err = Set(r, key, val)
	}
	return
}
