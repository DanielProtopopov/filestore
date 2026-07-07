package main

import (
	"context"
	apiconfig "filestore/cmd/api/config"
	redis2 "filestore/services/redis"
	"filestore/structs"
	"fmt"
	"log"
	"os"
	"slices"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

func DeleteExpiredFolder(ctx context.Context, folder string) error {
	errRemoveDirectory := os.RemoveAll(fmt.Sprintf("%s/%s", apiconfig.Settings.StoragePath, folder))
	if errRemoveDirectory != nil {
		return errors.Wrapf(errRemoveDirectory, "Failed to remove directory %s: %s", folder, errRemoveDirectory.Error())
	}
	_, errDelete := structs.Redis.Del(ctx, folder).Result()
	if errDelete != nil {
		return errors.Wrapf(errDelete, "Failed to delete Redis value by key %s: %s", folder, errDelete.Error())
	}

	redisIPEntries, errGetEntries := structs.Redis.Keys(ctx, "ip-*").Result()
	if errGetEntries != nil && !errors.Is(errGetEntries, redis.Nil) {
		return errors.Wrapf(errGetEntries, "Failed to get a list of all IP address sets")
	}

	for _, redisIPEntry := range redisIPEntries {
		redisIPFolders, errGetRedisIPFolders := structs.Redis.LRange(ctx, redisIPEntry, 0, -1).Result()
		if errGetRedisIPFolders != nil && !errors.Is(errGetRedisIPFolders, redis.Nil) {
			return errors.Wrapf(errGetRedisIPFolders, "Failed to get a list of IP address folders for %s", redisIPEntry)
		}

		if slices.Contains(redisIPFolders, folder) {
			_, errRemoveFolder := structs.Redis.LRem(ctx, redisIPEntry, 0, folder).Result()
			if errRemoveFolder != nil && !errors.Is(errRemoveFolder, redis.Nil) {
				return errors.Wrapf(errRemoveFolder, "Failed to remove folder %s from IP address set %s", folder, redisIPEntry)
			}
		}
	}

	return nil
}

// Start listening for events in Redis via Redis Event Sink
func HandleRedisEvents(ctx context.Context) error {
	const redisEventWorkers = 128
	eventWorkers := make(chan struct{}, redisEventWorkers)

	handleRedisEvent := func(ev *redis2.RedisEventMsg) {
		if ev.KeyEvent == "expired" {
			errDeleteExpiredFolder := DeleteExpiredFolder(ctx, ev.Key)
			if errDeleteExpiredFolder != nil {
				log.Printf("Failed to delete expired folder %s: %s", ev.Key, errDeleteExpiredFolder.Error())
			}
		}
	}

	for {
		select {
		case ev := <-apiconfig.Settings.EventSink.EventChannel:
			eventWorkers <- struct{}{}
			go func(event *redis2.RedisEventMsg) {
				defer func() { <-eventWorkers }()
				handleRedisEvent(event)
			}(ev)
		}
	}
}
