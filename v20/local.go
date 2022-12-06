package v20

import (
	"context"
	"sync"
	"time"
)

// 本地缓存
type S6Local struct {
	// 数据
	m3data map[string]*s6Unit
	// 读写锁，解决并发访问 map 的问题
	p7s6lock *sync.RWMutex
}

func F8NewS6Local() *S6Local {
	return &S6Local{
		m3data:   make(map[string]*s6Unit, 64),
		p7s6lock: &sync.RWMutex{},
	}
}

func (p7this *S6Local) F8Get(i9ctx context.Context, key string) (any, error) {
	// 判断缓存里有没有
	// 这里加读锁就行了
	p7this.p7s6lock.RLock()
	t4unit, ok := p7this.m3data[key]
	// 读完就可以解锁
	p7this.p7s6lock.RUnlock()
	if !ok {
		return nil, errKeyNotFound
	}

	// 判断缓存过期没有
	now := time.Now()
	if !t4unit.f8CheckDeadline(now) {
		// 缓存过期，删除缓存
		// 删除操作，需要加写锁
		p7this.p7s6lock.Lock()
		defer p7this.p7s6lock.Unlock()
		// 二次校验，防止别的线程抢先操作了
		t4unit, ok = p7this.m3data[key]
		if !ok {
			return nil, errKeyNotFound
		}
		if !t4unit.f8CheckDeadline(now) {
			p7this.f8Delete(i9ctx, key)
		}
		// 缓存过期可以归类为找不到
		return nil, errKeyNotFound
	}
	// 这里记得把数据拿出来
	return t4unit.value, nil
}

func (p7this *S6Local) F8Set(i9ctx context.Context, key string, value any, expiration time.Duration) error {
	// 修改操作，需要加写锁
	p7this.p7s6lock.Lock()
	defer p7this.p7s6lock.Unlock()
	// 计算过期时间
	var deadline time.Time = time.Time{}
	if 0 < expiration {
		deadline = time.Now().Add(expiration)
	}
	// 设置缓存
	p7this.m3data[key] = &s6Unit{
		value:    value,
		deadline: deadline,
	}
	return nil
}

func (p7this *S6Local) F8Delete(i9ctx context.Context, key string) error {
	// 删除操作，需要加写锁
	p7this.p7s6lock.Lock()
	defer p7this.p7s6lock.Unlock()
	p7this.f8Delete(i9ctx, key)
	return nil
}

func (p7this *S6Local) f8Delete(i9ctx context.Context, key string) {
	// 一般删除的时候，肯定外面都是拿着锁进来的，这里就不需要管并发访问 map 的问题了
	_, ok := p7this.m3data[key]
	if ok {
		delete(p7this.m3data, key)
	}
}
