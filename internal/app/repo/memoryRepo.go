package repo

import (
	"fmt"
	"log"
	"strings"
)

type MemoryRepo struct {
	storage map[string]string
}

func NewMemoryRepo() *MemoryRepo {
	return &MemoryRepo{
		storage: make(map[string]string),
	}
}

func (mr *MemoryRepo) SaveSrcUrlByHashKey(
	hashKey string,
	srcUrl string,
) error {
	if candidate, err := mr.GetSrcUrlByHashKey(hashKey); err == nil {
		if !strings.EqualFold(srcUrl, candidate) {
			return fmt.Errorf(
				"collision, hash key %v for %v already used for URL %v ",
				hashKey,
				srcUrl,
				candidate,
			)
		}

		return nil
	}

	mr.storage[hashKey] = srcUrl

	fmt.Printf(
		"URL %v successfully saved by key %v",
		srcUrl,
		hashKey,
	)

	return nil
}

func (mr *MemoryRepo) GetSrcUrlByHashKey(
	hashKey string,
) (string, error) {
	if srcUrl, ok := mr.storage[hashKey]; ok {
		log.Printf("find value: %v", srcUrl)
		return srcUrl, nil
	}

	return "", fmt.Errorf(
		"source URL not found by key: %v",
		hashKey,
	)
}
