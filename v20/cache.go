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
