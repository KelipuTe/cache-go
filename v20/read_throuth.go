package v20

import (
	"context"
	"sync"
	"time"
)

// S6ReadThroughCache 使用 read through 模式的缓存
// read through 的逻辑是先从缓存里捞，缓存里没有再去数据库捞。去数据库捞可以是同步的、半异步的、全异步的。
// 同步的：业务调用 miss->去数据库捞->设置缓存->返回
// 半异步的：业务调用 miss->去数据库捞->这里不设置缓存就返回了，异步执行设置缓存
// 全异步的：业务调用 miss->这里就直接返回错误了，异步执行去数据库捞->设置缓存
type S6ReadThroughCache struct {
	// i9Cache 缓存本体
	i9Cache I9Cache
	// f8LoadData 缓存 miss 的时候，去数据库加载数据
	// 这里如果一个方法搞不定所有的场景，就可以考虑设计成接口
	f8LoadData func(i9ctx context.Context, key string) (any, error)
	// loadDataExpiration 缓存 miss 的时候，设置缓存时设置的缓存超时时间
	loadDataExpiration time.Duration
	// 锁，解决并发调用 f8LoadData 的问题
	p7s6lock *sync.Mutex
}
