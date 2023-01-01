package v20

import (
	"context"
	"errors"
	"log"
	"testing"
	"time"
)

func TestS6ReadThroughCacheF8Set(p7s6t *testing.T) {
	i9cache := F8NewS6LocalForTest()
	f8LoadData := func(i9ctx context.Context, key string) (any, error) {
		switch key {
		case "key1":
			return "aaaa", nil
		case "key2":
			return "bbbb", nil
		}
		return "", errors.New("not found in db")
	}
	p7s6RTCache := F8NewS6ReadThroughCacheForTest(i9cache, f8LoadData, 3*time.Second)
	value, err := p7s6RTCache.F8Get(context.Background(), "key1")
	log.Println("key1", value, err)
	value, err = p7s6RTCache.F8Get(context.Background(), "key2")
	log.Println("key2", value, err)
	value, err = p7s6RTCache.F8Get(context.Background(), "key3")
	log.Println("key3", value, err)

	value, err = p7s6RTCache.F8Get(context.Background(), "key1")
	log.Println("key1", value, err)
	time.Sleep(5 * time.Second)
	value, err = p7s6RTCache.F8Get(context.Background(), "key1")
	log.Println("key1", value, err)
}
