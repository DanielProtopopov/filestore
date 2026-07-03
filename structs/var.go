package structs

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"
)

var (
	Redis     *redis.Client
	Localizer *i18n.Localizer
	Scheduler *cron.Cron
)
