package v20

import (
	"context"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestS6CacheWithMemoryLimitF8Set(p7s6t *testing.T) {
	s5s6case := []struct {
		name          string
		cache         func() *S6CacheWithMemoryLimit
		key           string
		value         string
		wantNowMemory int64
		wantS5LRU     []string
		wantErr       error
	}{
		{
			name: "key_not_exist",
			cache: func() *S6CacheWithMemoryLimit {
				return &S6CacheWithMemoryLimit{
					i9Cache:   F8NewS6LocalForTest(),
					NowMemory: 0,
					MaxMemory: 16,
					s5LRU:     make([]string, 0),
					p7s6lock:  &sync.Mutex{},
				}
			},
			key:           "key1",
			value:         "aaaa",
			wantNowMemory: 4,
			wantS5LRU:     []string{"key1"},
			wantErr:       nil,
		},
		{
			name: "key_exist",
			cache: func() *S6CacheWithMemoryLimit {
				return &S6CacheWithMemoryLimit{
					i9Cache: &S6Local{
						m3data: map[string]*s6Unit{
							"key1": &s6Unit{value: "aaaa", deadline: time.Time{}},
						},
						p7s6lock: &sync.RWMutex{},
					},
					NowMemory: 4,
					MaxMemory: 16,
					s5LRU:     []string{"key1"},
					p7s6lock:  &sync.Mutex{},
				}
			},
			key:           "key1",
			value:         "bbbbbb",
			wantNowMemory: 6,
			wantS5LRU:     []string{"key1"},
			wantErr:       nil,
		},
		{
			name: "add_key",
			cache: func() *S6CacheWithMemoryLimit {
				return &S6CacheWithMemoryLimit{
					i9Cache: &S6Local{
						m3data: map[string]*s6Unit{
							"key1": &s6Unit{value: "aaaa", deadline: time.Time{}},
						},
						p7s6lock: &sync.RWMutex{},
					},
					NowMemory: 4,
					MaxMemory: 16,
					s5LRU:     []string{"key1"},
					p7s6lock:  &sync.Mutex{},
				}
			},
			key:           "key2",
			value:         "bbbb",
			wantNowMemory: 8,
			wantS5LRU:     []string{"key1", "key2"},
			wantErr:       nil,
		},
		{
			name: "delete_key",
			cache: func() *S6CacheWithMemoryLimit {
				return &S6CacheWithMemoryLimit{
					i9Cache: &S6Local{
						m3data: map[string]*s6Unit{
							"key1": &s6Unit{value: "aaaa", deadline: time.Time{}},
							"key2": &s6Unit{value: "bbbb", deadline: time.Time{}},
							"key3": &s6Unit{value: "cccc", deadline: time.Time{}},
							"key4": &s6Unit{value: "dddd", deadline: time.Time{}},
						},
						p7s6lock: &sync.RWMutex{},
					},
					NowMemory: 16,
					MaxMemory: 16,
					s5LRU:     []string{"key1", "key2", "key3", "key4"},
					p7s6lock:  &sync.Mutex{},
				}
			},
			key:           "key5",
			value:         "eeee",
			wantNowMemory: 16,
			wantS5LRU:     []string{"key2", "key3", "key4", "key5"},
			wantErr:       nil,
		},
		{
			name: "delete_key_twice",
			cache: func() *S6CacheWithMemoryLimit {
				return &S6CacheWithMemoryLimit{
					i9Cache: &S6Local{
						m3data: map[string]*s6Unit{
							"key1": &s6Unit{value: "aaaa", deadline: time.Time{}},
							"key2": &s6Unit{value: "bbbb", deadline: time.Time{}},
							"key3": &s6Unit{value: "cccc", deadline: time.Time{}},
							"key4": &s6Unit{value: "dddd", deadline: time.Time{}},
						},
						p7s6lock: &sync.RWMutex{},
					},
					NowMemory: 16,
					MaxMemory: 16,
					s5LRU:     []string{"key1", "key2", "key3", "key4"},
					p7s6lock:  &sync.Mutex{},
				}
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
			cache := t4value.cache()
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
