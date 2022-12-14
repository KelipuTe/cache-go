package v20

import (
	"context"
	"sync"
	"time"
)

// S6CacheWithMemoryLimit 有最大内存限制的缓存
// 装饰器模式，装饰 I9Cache（S6Local）
type S6CacheWithMemoryLimit struct {
	// 缓存本体
	i9Cache I9Cache
	// 当前内存大小，单位 byte
	NowMemory int64
	// 最大内存大小，单位 byte
	MaxMemory int64
	// lru 结构，先用切片意思一下，后面有时间再换合适的结构
	s5LRU []string
	// 锁，解决并发访问 lru 结构的问题
	// 这个锁感觉有的时候和本地缓存里面那个锁重复了，后面有时间再研究
	p7s6lock *sync.Mutex
}

func F8NewS6CacheWithMemoryLimitForTest(i9cache I9Cache, maxMemory int64) *S6CacheWithMemoryLimit {
	return &S6CacheWithMemoryLimit{
		i9Cache:   i9cache,
		NowMemory: 0,
		MaxMemory: maxMemory,
		s5LRU:     make([]string, 0),
		p7s6lock:  &sync.Mutex{},
	}
}

// F8Set 这地方需要限制 value 的类型，要不然过来那种多级指针啥的，不怎么好算内存占用
func (p7this *S6CacheWithMemoryLimit) F8Set(i9ctx context.Context, key string, value string, expiration time.Duration) error {
	p7this.p7s6lock.Lock()
	defer p7this.p7s6lock.Unlock()

	// 看看容量超了没有，超了就要删掉 lru 前面的几个
	valueLen := int64(len([]byte(value)))
	for p7this.NowMemory+valueLen > p7this.MaxMemory {
		// 每次删 lru 的第一个
		t4value, err := p7this.i9Cache.F8Get(i9ctx, p7this.s5LRU[0])
		if nil != err {
			return err
		}
		t4str := t4value.(string)
		err = p7this.i9Cache.F8Delete(i9ctx, p7this.s5LRU[0])
		if nil != err {
			return err
		}
		p7this.NowMemory = p7this.NowMemory - int64(len([]byte(t4str)))
		p7this.f8DeleteFromLRU(p7this.s5LRU[0])
	}

	// 添加缓存
	t4value, err := p7this.i9Cache.F8Get(i9ctx, key)
	if nil == err {
		// 没报错，说明缓存里有。可以覆盖，也可以先删除再添加。
		err2 := p7this.i9Cache.F8Set(i9ctx, key, value, expiration)
		if nil != err2 {
			return err2
		}
		t4str := t4value.(string)
		p7this.NowMemory = p7this.NowMemory - int64(len([]byte(t4str))) + valueLen
		p7this.f8DeleteFromLRU(key)
		p7this.f8AddToLRU(key)
		return nil
	}
	// 报错了，就当缓存里没有，添加
	err = p7this.i9Cache.F8Set(i9ctx, key, value, expiration)
	if nil != err {
		return err
	}
	p7this.NowMemory = p7this.NowMemory + int64(len([]byte(value)))
	p7this.f8DeleteFromLRU(key)
	p7this.f8AddToLRU(key)
	return nil
}

func (p7this *S6CacheWithMemoryLimit) F8Delete(i9ctx context.Context, key string) error {
	p7this.p7s6lock.Lock()
	defer p7this.p7s6lock.Unlock()

	t4value, err := p7this.i9Cache.F8Get(i9ctx, key)
	if nil != err {
		// 这里还可以继续判断是 key 不存在的情况还是其他的报错
		return err
	}
	err = p7this.i9Cache.F8Delete(i9ctx, key)
	if nil != err {
		return err
	}
	t4str := t4value.(string)
	p7this.NowMemory = p7this.NowMemory - int64(len(t4str))
	p7this.f8DeleteFromLRU(key)

	return nil
}

func (p7this *S6CacheWithMemoryLimit) f8AddToLRU(key string) {
	p7this.s5LRU = append(p7this.s5LRU, key)
}

func (p7this *S6CacheWithMemoryLimit) f8DeleteFromLRU(key string) {
	deleteIndex := -1
	for t4index, t4value := range p7this.s5LRU {
		if key == t4value {
			deleteIndex = t4index
			break
		}
	}
	if 0 <= deleteIndex {
		p7this.s5LRU = append(p7this.s5LRU[:deleteIndex], p7this.s5LRU[deleteIndex+1:]...)
	}
}
