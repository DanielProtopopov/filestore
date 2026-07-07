package main

import (
	"context"
	"filestore/structs"
	"fmt"
	"io/fs"
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
	var directoriesRemoved []string
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
			directoriesRemoved = append(directoriesRemoved, directory.Name())
			totalDirectoryCount++
		}
	}

	redisIPEntries, errGetEntries := structs.Redis.Keys(ctx, "ip-*").Result()
	if errGetEntries != nil && !errors.Is(errGetEntries, redis.Nil) {
		return errors.Wrapf(errGetEntries, "Failed to get a list of all IP address sets")
	}

	for _, redisIPEntry := range redisIPEntries {
		redisIPFolders, errGetRedisIPFolders := structs.Redis.LRange(ctx, redisIPEntry, 0, -1).Result()
		if errGetRedisIPFolders != nil && !errors.Is(errGetRedisIPFolders, redis.Nil) {
			return errors.Wrapf(errGetRedisIPFolders, "Failed to get a list of all IP address folders for %s", redisIPEntry)
		}
		for _, redisIPFolder := range redisIPFolders {
			_, errDirectoryExists := os.Stat(fmt.Sprintf("%s/%s", storagePath, redisIPFolder))
			if errors.Is(errDirectoryExists, fs.ErrNotExist) {
				errRemoveRedisIPFolder := structs.Redis.LRem(ctx, redisIPEntry, 0, redisIPFolder).Err()
				if errRemoveRedisIPFolder != nil {
					log.Printf("Failed to remove entry %s from Redis IP collection %s: %s", redisIPFolder, redisIPEntry, errRemoveRedisIPFolder.Error())
					continue
				}
				ipFoldersCount, errGetCount := structs.Redis.LLen(ctx, redisIPEntry).Result()
				if errGetCount != nil && !errors.Is(errGetCount, redis.Nil) {
					log.Printf("Failed to get Redis IP folders collection count of %s: %s", redisIPEntry, errGetCount.Error())
					continue
				}
				if ipFoldersCount == 0 {
					_, errDeleteByKey := structs.Redis.Del(ctx, redisIPEntry).Result()
					if errDeleteByKey != nil && !errors.Is(errDeleteByKey, redis.Nil) {
						log.Printf("Failed to delete Redis IP collection by key %s: %s", redisIPEntry, errDeleteByKey.Error())
						continue
					}
				}
				_, errDelete := structs.Redis.Del(ctx, redisIPFolder).Result()
				if errDelete != nil {
					log.Printf("Failed to delete Redis value by key %s: %s", redisIPFolder, errDelete.Error())
					continue
				}
			}
		}
	}

	return nil
}
