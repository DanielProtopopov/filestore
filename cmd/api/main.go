package main

import (
	"context"
	"filestore/api"
	apiconfig "filestore/cmd/api/config"
	"filestore/internal/translations"
	"filestore/structs"
	"log"
	"os"
	"runtime/debug"
	"time"

	"github.com/joho/godotenv"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/robfig/cron/v3"
)

func main() {
	ctx := context.Background()
	errLoadEnv := godotenv.Load(".env")
	if errLoadEnv != nil {
		log.Printf("[ERROR] Error loading .env file: %s", errLoadEnv.Error())
	}

	apiconfig.Init(ctx)
	log.Println("Configuration was successfully loaded")

	structs.Scheduler = cron.New(cron.WithChain(cron.SkipIfStillRunning(cron.VerbosePrintfLogger(log.New(os.Stdout, "cron: ", log.LstdFlags)))))

	var standardParser = cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.DowOptional | cron.Descriptor)
	storageCheckSched, errParseStorageDuration := standardParser.Parse("*/30 * * * * *")
	if errParseStorageDuration != nil {
		log.Panicf("Error parsing scheduler for storage check: %s", errParseStorageDuration.Error())
	}

	structs.Scheduler.Schedule(storageCheckSched, cron.FuncJob(func() {
		errDeleteExpiredFiles := DeleteExpiredFiles(ctx, apiconfig.Settings.StoragePath)
		if errDeleteExpiredFiles != nil {
			log.Printf("[ERROR] Failed to check files due for deletion: %s", errDeleteExpiredFiles.Error())
		}
	}))

	go func() {
		t := time.Tick(time.Minute)
		for {
			<-t
			debug.FreeOSMemory()
		}
	}()

	structs.Localizer = i18n.NewLocalizer(apiconfig.Settings.Bundle, translations.DefaultLanguage)

	if len(structs.Scheduler.Entries()) != 0 {
		structs.Scheduler.Start()
	}

	errDeleteExpiredFiles := DeleteExpiredFiles(ctx, apiconfig.Settings.StoragePath)
	if errDeleteExpiredFiles != nil {
		log.Printf("[ERROR] Failed to check files due for deletion: %s", errDeleteExpiredFiles.Error())
	}
	api.Serve()
}
