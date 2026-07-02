package redis

import (
	"github.com/redis/go-redis/v9"
)

type (
	Service struct {
		Redis     *redis.Client
		EventSink *RedisEventSink
		Keys      map[string]chan *RedisEventMsg
	}
)
