package main

import (
	"context"
	"filestore/structs"
	"fmt"
	"log"
	"os"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

func DeleteExpiredFiles(ctx context.Context, storagePath string) error {
	directories, errReadStoragePath := os.ReadDir(storagePath)
	if errReadStoragePath != nil {
		return errors.Wrapf(errReadStoragePath, "Failed to read directory %s", storagePath)
	}

	totalDirectoryCount := 0
	for _, directory := range directories {
		if !directory.IsDir() {
			continue
		}

		keyExists, errKeyExists := structs.Redis.Exists(ctx, directory.Name()).Result()
		if errKeyExists != nil && !errors.Is(errKeyExists, redis.Nil) {
			log.Printf("[ERROR] Failed to check if %s exists: %s", directory.Name(), errKeyExists.Error())
			continue
		}

		if keyExists <= 0 {
			errRemoveDirectory := os.RemoveAll(fmt.Sprintf("%s/%s", storagePath, directory.Name()))
			if errRemoveDirectory != nil {
				log.Printf("[ERROR] Failed to remove directory %s: %s", directory.Name(), errRemoveDirectory.Error())
				continue
			}
			totalDirectoryCount++
		}
	}

	return nil
}
