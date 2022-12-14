package v20

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestS6CacheWithMemoryLimitF8Set(p7s6t *testing.T) {
	s5s6case := []struct {
		name          string
		f8GetCache    func() *S6CacheWithMemoryLimit
		key           string
		value         string
		wantNowMemory int64
		wantS5LRU     []string
		wantErr       error
	}{
		{
			name: "key_not_exist",
			f8GetCache: func() *S6CacheWithMemoryLimit {
				i9cache := F8NewS6LocalForTest()
				p7s6cache := F8NewS6CacheWithMemoryLimitForTest(i9cache, 16)
				return p7s6cache
			},
			key:           "key1",
			value:         "aaaa",
			wantNowMemory: 4,
			wantS5LRU:     []string{"key1"},
			wantErr:       nil,
		},
		{
			name: "key_exist",
			f8GetCache: func() *S6CacheWithMemoryLimit {
				i9cache := F8NewS6LocalForTest()
				i9cache.m3data["key1"] = &s6Unit{value: "aaaa", deadline: time.Time{}}
				p7s6cache := F8NewS6CacheWithMemoryLimitForTest(i9cache, 16)
				p7s6cache.NowMemory = 4
				p7s6cache.s5LRU = []string{"key1"}
				return p7s6cache
			},
			key:           "key1",
			value:         "bbbbbb",
			wantNowMemory: 6,
			wantS5LRU:     []string{"key1"},
			wantErr:       nil,
		},
		{
			name: "add_key",
			f8GetCache: func() *S6CacheWithMemoryLimit {
				i9cache := F8NewS6LocalForTest()
				i9cache.m3data["key1"] = &s6Unit{value: "aaaa", deadline: time.Time{}}
				p7s6cache := F8NewS6CacheWithMemoryLimitForTest(i9cache, 16)
				p7s6cache.NowMemory = 4
				p7s6cache.s5LRU = []string{"key1"}
				return p7s6cache
			},
			key:           "key2",
			value:         "bbbb",
			wantNowMemory: 8,
			wantS5LRU:     []string{"key1", "key2"},
			wantErr:       nil,
		},
		{
			name: "add_delete_key",
			f8GetCache: func() *S6CacheWithMemoryLimit {
				i9cache := F8NewS6LocalForTest()
				i9cache.m3data["key1"] = &s6Unit{value: "aaaa", deadline: time.Time{}}
				i9cache.m3data["key2"] = &s6Unit{value: "bbbb", deadline: time.Time{}}
				i9cache.m3data["key3"] = &s6Unit{value: "cccc", deadline: time.Time{}}
				i9cache.m3data["key4"] = &s6Unit{value: "dddd", deadline: time.Time{}}
				p7s6cache := F8NewS6CacheWithMemoryLimitForTest(i9cache, 16)
				p7s6cache.NowMemory = 16
				p7s6cache.s5LRU = []string{"key1", "key2", "key3", "key4"}
				return p7s6cache
			},
			key:           "key5",
			value:         "eeee",
			wantNowMemory: 16,
			wantS5LRU:     []string{"key2", "key3", "key4", "key5"},
			wantErr:       nil,
		},
		{
			name: "add_delete_2key",
			f8GetCache: func() *S6CacheWithMemoryLimit {
				i9cache := F8NewS6LocalForTest()
				i9cache.m3data["key1"] = &s6Unit{value: "aaaa", deadline: time.Time{}}
				i9cache.m3data["key2"] = &s6Unit{value: "bbbb", deadline: time.Time{}}
				i9cache.m3data["key3"] = &s6Unit{value: "cccc", deadline: time.Time{}}
				i9cache.m3data["key4"] = &s6Unit{value: "dddd", deadline: time.Time{}}
				p7s6cache := F8NewS6CacheWithMemoryLimitForTest(i9cache, 16)
				p7s6cache.NowMemory = 16
				p7s6cache.s5LRU = []string{"key1", "key2", "key3", "key4"}
				return p7s6cache
			},
			key:           "key5",
			value:         "eeeeeeee",
			wantNowMemory: 16,
			wantS5LRU:     []string{"key3", "key4", "key5"},
			wantErr:       nil,
		},
	}

	for _, t4value := range s5s6case {
		p7s6t.Run(t4value.name, func(p7s6t2 *testing.T) {
			cache := t4value.f8GetCache()
			err := cache.F8Set(context.Background(), t4value.key, t4value.value, time.Minute)
			assert.Equal(p7s6t2, t4value.wantErr, err)
			if nil != err {
				return
			}
			assert.Equal(p7s6t2, t4value.wantNowMemory, cache.NowMemory)
			assert.Equal(p7s6t2, t4value.wantS5LRU, cache.s5LRU)
		})
	}
}

func TestS6CacheWithMemoryLimitF8Delete(p7s6t *testing.T) {
	s5s6case := []struct {
		name          string
		f8GetCache    func() *S6CacheWithMemoryLimit
		key           string
		wantNowMemory int64
		wantS5LRU     []string
		wantErr       error
	}{
		{
			name: "key_not_exist",
			f8GetCache: func() *S6CacheWithMemoryLimit {
				i9cache := F8NewS6LocalForTest()
				p7s6cache := F8NewS6CacheWithMemoryLimitForTest(i9cache, 16)
				return p7s6cache
			},
			key:           "key1",
			wantNowMemory: 0,
			wantS5LRU:     []string{},
			wantErr:       errKeyNotFound,
		},
		{
			name: "key_exist",
			f8GetCache: func() *S6CacheWithMemoryLimit {
				i9cache := F8NewS6LocalForTest()
				i9cache.m3data["key1"] = &s6Unit{value: "aaaa", deadline: time.Time{}}
				p7s6cache := F8NewS6CacheWithMemoryLimitForTest(i9cache, 16)
				p7s6cache.NowMemory = 4
				p7s6cache.s5LRU = []string{"key1"}
				return p7s6cache
			},
			key:           "key1",
			wantNowMemory: 0,
			wantS5LRU:     []string{},
			wantErr:       nil,
		},
		{
			name: "delete_key",
			f8GetCache: func() *S6CacheWithMemoryLimit {
				i9cache := F8NewS6LocalForTest()
				i9cache.m3data["key1"] = &s6Unit{value: "aaaa", deadline: time.Time{}}
				i9cache.m3data["key2"] = &s6Unit{value: "bbbb", deadline: time.Time{}}
				p7s6cache := F8NewS6CacheWithMemoryLimitForTest(i9cache, 16)
				p7s6cache.NowMemory = 8
				p7s6cache.s5LRU = []string{"key1", "key2"}
				return p7s6cache
			},
			key:           "key1",
			wantNowMemory: 4,
			wantS5LRU:     []string{"key2"},
			wantErr:       nil,
		},
	}

	for _, t4value := range s5s6case {
		p7s6t.Run(t4value.name, func(p7s6t2 *testing.T) {
			p7s6cache := t4value.f8GetCache()
			err := p7s6cache.F8Delete(context.Background(), t4value.key)
			assert.Equal(p7s6t2, t4value.wantErr, err)
			if nil != err {
				return
			}
			assert.Equal(p7s6t2, t4value.wantNowMemory, p7s6cache.NowMemory)
			assert.Equal(p7s6t2, t4value.wantS5LRU, p7s6cache.s5LRU)
		})
	}
}
