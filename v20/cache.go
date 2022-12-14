package v20

import (
	"context"
	"time"
)

// I9Cache 缓存接口
type I9Cache interface {
	// F8Get 从缓存中获取
	F8Get(i9ctx context.Context, key string) (any, error)
	// F8Set 设置缓存
	F8Set(i9ctx context.Context, key string, value any, expiration time.Duration) error
	// F8Delete 删除缓存
	F8Delete(i9ctx context.Context, key string) error
}

// I9CacheUseT 用泛型的缓存接口
// 这个没有实现，就是放在这里作为一个提示，缓存这玩意可以用泛型来设计
type I9CacheUseT[T any] interface {
	// F8Get 从缓存中获取
	F8Get(i9ctx context.Context, key string) (T, error)
	// F8Set 设置缓存
	F8Set(i9ctx context.Context, key string, value T, expiration time.Duration) error
	// F8Delete 删除缓存
	F8Delete(i9ctx context.Context, key string) error
}
