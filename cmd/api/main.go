package main

import (
	"context"
	"filestore/api"
	apiconfig "filestore/cmd/api/config"
	"filestore/internal/translations"
	"filestore/structs"
	"log"
	"runtime/debug"
	"time"

	"github.com/joho/godotenv"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func main() {
	ctx := context.Background()
	errLoadEnv := godotenv.Load(".env")
	if errLoadEnv != nil {
		log.Printf("[ERROR] Error loading .env file: %s", errLoadEnv.Error())
	}

	apiconfig.Init(ctx)
	log.Println("Configuration was successfully loaded")

	go func() {
		t := time.Tick(time.Minute)
		for {
			<-t
			debug.FreeOSMemory()
		}
	}()

	structs.Localizer = i18n.NewLocalizer(apiconfig.Settings.Bundle, translations.DefaultLanguage)

	// Process Redis events (to replace cron job)
	go HandleRedisEvents(ctx)

	api.Serve()
}
