package apiconfig

import (
	"context"
	"filestore/config"
	"filestore/internal/translations"
	redis2 "filestore/services/redis"
	"filestore/structs"
	"fmt"
	"log"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/getsentry/sentry-go"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/redis/go-redis/v9"
	"github.com/shopspring/decimal"
	"golang.org/x/text/language"
)

type Config struct {
	Server                  config.ServerConfig      // API server configuration
	Sentry                  config.SentryConfig      // Sentry configuration
	Environment             config.EnvironmentConfig // Server environment configuration
	Bundle                  *i18n.Bundle             // I18n bundle instance (localization)
	I18n                    *i18n.Localizer          // I18n configuration (i18n)
	Redis                   config.RedisConfig       // Redis configuration
	EventSink               *redis2.RedisEventSink   // Redis event sink
	StoragePath             string                   // Path to storing temporary files
	MaximumStorageSize      int64                    // Maximum storage size (in bytes)
	MaximumExpirationPeriod uint                     // Maximum file expiration time (in seconds)
	MaximumUploadSize       int64                    // Maximum file upload size (in megabytes)
	MaximumFilesPerIP       uint                     // Maximum files per IP address
	GoogleTag               string                   // Google Tag Manager key
	Debug                   bool                     // Debugging flag
}

var Settings *Config

func Init(ctx context.Context) {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	for _, lang := range translations.Languages {
		_, errLoadLang := bundle.LoadMessageFile(fmt.Sprintf("%s/translate.%s.toml", "data/i18n", lang))
		if errLoadLang != nil {
			log.Panicf("Error loading translations for language %s: %s", lang, errLoadLang.Error())
			return
		}
	}

	if !config.GetEnvAsBool("DEBUG", false) {
		log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	} else {
		log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	}

	Settings = &Config{
		Bundle: bundle,
		I18n:   i18n.NewLocalizer(bundle, translations.Languages...),
		Debug:  config.GetEnvAsBool("DEBUG", true),
		Server: config.ServerConfig{
			Host:     config.GetEnv("SERVER_HOST", "0.0.0.0"),
			Port:     config.GetEnvAsInt("SERVER_PORT", 80),
			Domain:   config.GetEnv("SERVER_DOMAIN", "localhost"),
			Protocol: config.GetEnv("SERVER_PROTOCOL", "http"),
		},
		Sentry: config.SentryConfig{
			DSN:             config.GetEnv("SENTRY_DSN", ""),
			TraceSampleRate: config.GetEnvAsDecimal("SENTRY_TRACE_RATE", decimal.NewFromFloat(1.0)),
			EnableTracing:   config.GetEnvAsBool("SENTRY_ENABLE_TRACING", true),
		},
		Environment: config.EnvironmentConfig{
			FullName:  config.GetEnv("ENVIRONMENT_FULL", "development"),
			ShortName: config.GetEnv("ENVIRONMENT_SHORT", "dev"),
		},
		Redis: config.RedisConfig{
			Host:     config.GetEnv("REDIS_HOST", "localhost"),
			Port:     config.GetEnvAsInt("REDIS_PORT", 6379),
			Password: config.GetEnv("REDIS_PASSWORD", ""),
			DB:       config.GetEnvAsInt("REDIS_DB", 0),
			Proto:    config.GetEnv("REDIS_PROTOCOL", "tcp"),
		},
		StoragePath:             config.GetEnv("STORAGE_PATH", ""),
		MaximumStorageSize:      int64(config.GetEnvAsInt("MAXIMUM_STORAGE_SIZE", 250*1024*1024*1024)),
		MaximumUploadSize:       int64(config.GetEnvAsInt("MAXIMUM_UPLOAD_SIZE", 50*1024*1024)),
		MaximumExpirationPeriod: config.GetEnvAsUInt("MAXIMUM_EXPIRATION_PERIOD", 4*60*60),
		MaximumFilesPerIP:       config.GetEnvAsUInt("MAXIMUM_FILES_PER_IP", 10),
		GoogleTag:               config.GetEnv("GOOGLE_TAG", ""),
	}

	client := redis.NewClient(&redis.Options{
		Network:  Settings.Redis.Proto,
		Addr:     fmt.Sprintf("%s:%d", Settings.Redis.Host, Settings.Redis.Port),
		Password: Settings.Redis.Password, DB: Settings.Redis.DB, MaxRetries: 3,
	})
	if errPing := client.Ping(context.Background()).Err(); errPing != nil {
		log.Panicf("Error pinging Redis: %s", errPing)
	} else {
		log.Println("Pinged Redis successfully")
	}
	structs.Redis = client
	Settings.EventSink = redis2.NewRedisEventSink(client, 0, "*")

	if Settings.Sentry.DSN != "" {
		errInitSentry := InitSentry(Settings)
		if errInitSentry != nil {
			log.Panicf("Error connecting to Sentry: %s", errInitSentry.Error())
		} else {
			log.Printf("Successfully connected to Sentry on address %s", Settings.Sentry.DSN)
		}
		sentry.ConfigureScope(func(scope *sentry.Scope) {
			scope.SetTag("Module", "Server")
		})
	} else {
		log.Println("Sentry DSN is not set, no logging there will occur")
	}
}

func InitSentry(config *Config) error {
	traceRate, _ := config.Sentry.TraceSampleRate.Float64()
	errInitSentry := sentry.Init(sentry.ClientOptions{
		Dsn: config.Sentry.DSN, Debug: config.Debug, EnableTracing: config.Sentry.EnableTracing,
		TracesSampleRate: traceRate})
	if errInitSentry != nil {
		return errInitSentry
	}

	// Flush buffered events before the program terminates.
	defer func() {
		errRecover := recover()
		if errRecover != nil {
			sentry.CurrentHub().Recover(errRecover)
			sentry.Flush(time.Second * 5)
		}
	}()

	return nil
}
