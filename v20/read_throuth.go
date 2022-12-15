package v20

import (
	"context"
	"log"
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
	// 锁，解决并发操作缓存的问题
	p7s6lock *sync.RWMutex
}

func F8NewS6ReadThroughCacheForTest(i9Cache I9Cache, f8LoadData func(i9ctx context.Context, key string) (any, error), loadDataExpiration time.Duration) *S6ReadThroughCache {
	return &S6ReadThroughCache{
		i9Cache:            i9Cache,
		f8LoadData:         f8LoadData,
		loadDataExpiration: loadDataExpiration,
		p7s6lock:           &sync.RWMutex{},
	}
}

func (p7this *S6ReadThroughCache) F8Get(i9ctx context.Context, key string) (any, error) {
	p7this.p7s6lock.RLock()
	value, err := p7this.i9Cache.F8Get(i9ctx, key)
	p7this.p7s6lock.RUnlock()
	if nil == err {
		// 没异常，直接返回结果
		log.Println("load from cache")
		return value, nil
	}
	// 有异常
	if err != errKeyNotFound {
		// 如果不是缓存里没捞到，那只能抛出去
		return nil, err
	} else {
		// 如果是缓存里没捞到，那就再去数据库捞
		// 这里需要加写锁，防止并发刷新缓存
		p7this.p7s6lock.Lock()
		defer p7this.p7s6lock.Unlock()
		newValue, err2 := p7this.f8LoadData(i9ctx, key)
		if nil != err2 {
			// 从数据库捞数据失败了，可能是数据库真没有，也可能是异常
			return nil, err2
		}
		err2 = p7this.i9Cache.F8Set(i9ctx, key, newValue, p7this.loadDataExpiration)
		if nil != err2 {
			// 数据捞回来了，但是设置缓存的时候失败了
			return nil, err2
		}
		log.Println("load from db")
		return newValue, nil
	}
}

// S6ReadThroughCacheUseT 使用 read through 模式的缓存，用泛型的缓存接口
type S6ReadThroughCacheUseT[T any] struct {
	i9Cache            I9CacheUseT[T]
	f8LoadData         func(i9ctx context.Context, key string) (T, error)
	loadDataExpiration time.Duration
	p7s6lock           *sync.RWMutex
}
