package cmdstore

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/zekrotja/ken/store"
)

type RedisCmdStore struct {
	c *redis.Client
}

func NewRedisCmdStore(client *redis.Client) *RedisCmdStore {
	return &RedisCmdStore{client}
}

var _ store.CommandStore = (*RedisCmdStore)(nil)

func (s *RedisCmdStore) Store(cmds map[string]string) (err error) {
	str, err := mapToString(cmds)
	if err != nil {
		return
	}
	err = s.c.Set(context.Background(), keyName, str, 0).Err()
	return
}

func (s *RedisCmdStore) Load() (cmds map[string]string, err error) {
	res := s.c.Get(context.Background(), keyName)
	if res.Err() != nil {
		if res.Err() == redis.Nil {
			cmds = make(map[string]string)
		} else {
			err = res.Err()
		}
		return
	}
	cmds, err = stringToMap(res.Val())
	return
}
