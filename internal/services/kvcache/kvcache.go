package kvcache

import "time"

type Provider interface {
	Get(key string) interface{}
	Set(key string, v interface{}, lifetime time.Duration)
	Del(key string)
}
