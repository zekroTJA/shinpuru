// Package rediscmdstore provides an implementation
// of github.com/zekrotja/ken/store.CommandStore using
// a redis client to store the command cache.
package rediscmdstore

import (
	"context"
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/zekrotja/ken/store"
)

// RedisCmdStore implements CommandStore using
// a redis client instance.
type RedisCmdStore struct {
	c       *redis.Client
	keyName string
}

var _ store.CommandStore = (*RedisCmdStore)(nil)

// New creates a new instance of RedisCmdStore using the
// passed redis client instance.
//
// You can also provide a custom key name used to store
// the command cache. Defaultly, this is "cmdcache:<hostname>",
// where <hostname> is the hostname of the system/container.
// If the hostname could not be determined, it defaults to "def".
func New(client *redis.Client, keyName ...string) (s *RedisCmdStore) {
	s = &RedisCmdStore{
		c: client,
	}

	if len(keyName) > 0 {
		s.keyName = keyName[0]
	} else {
		hostname, _ := os.Hostname()
		if hostname == "" {
			hostname = "def"
		}
		s.keyName = fmt.Sprintf("cmdcache:%s", hostname)
	}

	return
}

func (s *RedisCmdStore) Store(cmds map[string]string) (err error) {
	str, err := mapToString(cmds)
	if err != nil {
		return
	}
	err = s.c.Set(context.Background(), s.keyName, str, 0).Err()
	return
}

func (s *RedisCmdStore) Load() (cmds map[string]string, err error) {
	res := s.c.Get(context.Background(), s.keyName)
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
