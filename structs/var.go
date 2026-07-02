package structs

import (
	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"
)

var (
	Redis     *redis.Client
	Scheduler *cron.Cron
)
