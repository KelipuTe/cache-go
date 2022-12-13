package v20

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestS6LocalF8Get(p7s6t *testing.T) {
	s5s6case := []struct {
		name       string
		f8GetCache func() *S6Local
		key        string
		wantValue  any
		wantErr    error
	}{
		{
			name: "key_not_exist",
			f8GetCache: func() *S6Local {
				p7s6cache := F8NewS6LocalForTest()
				return p7s6cache
			},
			key:       "key1",
			wantValue: nil,
			wantErr:   errKeyNotFound,
		},
		{
			name: "key_exist",
			f8GetCache: func() *S6Local {
				p7s6cache := F8NewS6LocalForTest()
				p7s6cache.m3data["key1"] = &s6Unit{value: "aaaa", deadline: time.Time{}}
				return p7s6cache
			},
			key:       "key1",
			wantValue: "aaaa",
			wantErr:   nil,
		},
	}

	for _, t4value := range s5s6case {
		p7s6t.Run(t4value.name, func(p7s6t2 *testing.T) {
			i9cache := t4value.f8GetCache()
			value, err := i9cache.F8Get(context.Background(), t4value.key)
			assert.Equal(p7s6t2, t4value.wantErr, err)
			if nil != err {
				return
			}
			assert.Equal(p7s6t2, t4value.wantValue, value)
		})
	}
}

func TestS6LocalF8Set(p7s6t *testing.T) {
	s5s6case := []struct {
		name       string
		f8GetCache func() *S6Local
		key        string
		value      any
		wantErr    error
	}{
		{
			name: "key_not_exist",
			f8GetCache: func() *S6Local {
				p7s6cache := F8NewS6LocalForTest()
				return p7s6cache
			},
			key:     "key1",
			value:   "aaaa",
			wantErr: nil,
		},
		{
			name: "key_exist",
			f8GetCache: func() *S6Local {
				p7s6cache := F8NewS6LocalForTest()
				p7s6cache.m3data["key1"] = &s6Unit{value: "aaaa", deadline: time.Time{}}
				return p7s6cache
			},
			key:     "key1",
			value:   "bbbb",
			wantErr: nil,
		},
	}

	for _, t4value := range s5s6case {
		p7s6t.Run(t4value.name, func(p7s6t2 *testing.T) {
			i9cache := t4value.f8GetCache()
			err := i9cache.F8Set(context.Background(), t4value.key, t4value.value, 0)
			assert.Equal(p7s6t2, t4value.wantErr, err)
			if nil != err {
				return
			}
			value, err := i9cache.F8Get(context.Background(), t4value.key)
			assert.Equal(p7s6t2, t4value.value, value)
		})
	}
}

func TestS6LocalF8Delete(p7s6t *testing.T) {
	s5s6case := []struct {
		name       string
		f8GetCache func() *S6Local
		key        string
		wantErr    error
	}{
		{
			name: "key_not_exist",
			f8GetCache: func() *S6Local {
				p7s6cache := F8NewS6LocalForTest()
				p7s6cache.m3data["key2"] = &s6Unit{value: "bbbb", deadline: time.Time{}}
				return p7s6cache
			},
			key:     "key1",
			wantErr: nil,
		},
		{
			name: "key_exist",
			f8GetCache: func() *S6Local {
				p7s6cache := F8NewS6LocalForTest()
				p7s6cache.m3data["key1"] = &s6Unit{value: "aaaa", deadline: time.Time{}}
				p7s6cache.m3data["key2"] = &s6Unit{value: "bbbb", deadline: time.Time{}}
				return p7s6cache
			},
			key:     "key1",
			wantErr: nil,
		},
	}

	for _, t4value := range s5s6case {
		p7s6t.Run(t4value.name, func(p7s6t2 *testing.T) {
			i9cache := t4value.f8GetCache()
			err := i9cache.F8Delete(context.Background(), t4value.key)
			assert.Equal(p7s6t2, t4value.wantErr, err)
			if nil != err {
				return
			}
		})
	}
}
