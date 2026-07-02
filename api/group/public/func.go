package public

import (
	"context"
	"filestore/structs"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"sync"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

// Source - https://stackoverflow.com/a/31832326
// Posted by icza, modified by community. See post 'Timeline' for change history
// Retrieved 2026-06-02, License - CC BY-SA 4.0

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func RandStringBytesMaskImpr(n int) string {
	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdxMax letters!
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

func generateFolderNames(length uint, count uint) []string {
	var folderNames []string
	for index := uint(0); index < count; index++ {
		folderNames = append(folderNames, RandStringBytesMaskImpr(int(length)))
	}

	return folderNames
}

func getIPAddressFolders(ctx context.Context, storagePath string) ([]string, error) {
	directories, errReadStoragePath := os.ReadDir(storagePath)
	if errReadStoragePath != nil {
		return []string{}, errors.Wrapf(errReadStoragePath, "Failed to read directory %s", storagePath)
	}

	var folders []string
	for _, directory := range directories {
		if !directory.IsDir() {
			continue
		}

		folderName, errKeyExists := structs.Redis.Get(ctx, fmt.Sprintf("%s-folder-ip-address", directory)).Result()
		if errKeyExists != nil && !errors.Is(errKeyExists, redis.Nil) {
			log.Printf("[ERROR] Failed to check if %s exists: %s", directory.Name(), errKeyExists.Error())
			continue
		}

		if folderName != "" {
			folders = append(folders, folderName)
		}
	}

	return folders, nil
}

// Source - https://stackoverflow.com/a/78456550
// Posted by KiddoV
// Retrieved 2026-06-23, License - CC BY-SA 4.0
func DirectorySize(path string) (int64, error) {
	var size int64
	var mu sync.Mutex

	// Function to calculate size for a given path
	var calculateSize func(string) error
	calculateSize = func(p string) error {
		fileInfo, err := os.Lstat(p)
		if err != nil {
			return err
		}

		// Skip symbolic links to avoid counting them multiple times
		if fileInfo.Mode()&os.ModeSymlink != 0 {
			return nil
		}

		if fileInfo.IsDir() {
			entries, err := os.ReadDir(p)
			if err != nil {
				return err
			}
			for _, entry := range entries {
				if err := calculateSize(filepath.Join(p, entry.Name())); err != nil {
					return err
				}
			}
		} else {
			mu.Lock()
			size += fileInfo.Size()
			mu.Unlock()
		}
		return nil
	}

	// Start calculation from the root path
	if err := calculateSize(path); err != nil {
		return 0, err
	}

	return size, nil
}
