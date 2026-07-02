package config

import "github.com/shopspring/decimal"

type (
	ServerConfig struct {
		Host     string
		Port     int
		Domain   string
		Protocol string
	}

	SentryConfig struct {
		DSN             string
		TraceSampleRate decimal.Decimal
		EnableTracing   bool
	}

	EnvironmentConfig struct {
		FullName  string
		ShortName string
	}

	RedisConfig struct {
		Host     string
		Port     int
		Password string
		DB       int
		Proto    string
	}
)
