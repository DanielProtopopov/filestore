package public

import (
	"context"
	"crypto/rand"
	"filestore/structs"
	"fmt"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"slices"
	"sync"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

func getToken(length int) string {
	token := ""
	codeAlphabet := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

	for i := 0; i < length; i++ {
		token += string(codeAlphabet[cryptoRandSecure(int64(len(codeAlphabet)))])
	}
	return token
}

func cryptoRandSecure(max int64) int64 {
	nBig, err := rand.Int(rand.Reader, big.NewInt(max))
	if err != nil {
		log.Println(err)
	}
	return nBig.Int64()
}

func generateFolderNames(length uint, count uint) []string {
	var folderNames []string
	for index := uint(0); index < count; index++ {
		folderNames = append(folderNames, getToken(int(length)))
	}

	return folderNames
}

func getIPAddressFolders(ctx context.Context, storagePath string, ipAddress string) ([]string, error) {
	directories, errReadStoragePath := os.ReadDir(storagePath)
	if errReadStoragePath != nil {
		return []string{}, errors.Wrapf(errReadStoragePath, "Failed to read directory %s", storagePath)
	}

	var folders []string
	for _, directory := range directories {
		if !directory.IsDir() {
			continue
		}

		foldersByIPAddress, errKeyExists := structs.Redis.LRange(ctx, fmt.Sprintf("ip-%s", ipAddress), 0, -1).Result()
		if errKeyExists != nil && !errors.Is(errKeyExists, redis.Nil) {
			log.Printf("[ERROR] Failed to check if %s exists: %s", directory.Name(), errKeyExists.Error())
			continue
		}

		// Directory gets deleted when expired, so matching key to folder means it's not yet expired
		if slices.Contains(foldersByIPAddress, directory.Name()) {
			folders = append(folders, directory.Name())
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
