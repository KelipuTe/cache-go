package v20

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// #### type ####

// s6LocalInitSize 缓存初始化时的容量
var s6LocalInitSize int = 16

// S6Local 本地缓存
type S6Local struct {
	// m3data 缓存数据，这里理论上需要预估缓存容量
	m3data map[string]*s6Unit
	// p7s6lock 读写锁，解决并发访问 map 的问题
	p7s6lock *sync.RWMutex
	// f8OnEvicted 删除缓存的事件
	// 这个可以用依赖注入或者 option 设计模式来设置
	f8OnEvicted func(i9ctx context.Context, key string, p7s6unit *s6Unit) error
	// 缓存关闭信号
	c4CloseSignal chan struct{}
	// 控制缓存关闭只能一次
	p7s6CloseOnce *sync.Once
}

// #### func ####

func F8NewS6LocalForTest() *S6Local {
	return &S6Local{
		m3data:        make(map[string]*s6Unit, s6LocalInitSize),
		p7s6lock:      &sync.RWMutex{},
		c4CloseSignal: make(chan struct{}, 1),
		p7s6CloseOnce: &sync.Once{},
		f8OnEvicted: func(i9ctx context.Context, key string, p7s6unit *s6Unit) error {
			fmt.Printf("f8OnEvicted, key: %s, value %v\r\n", key, p7s6unit.value)
			return nil
		},
	}
}

// #### type func ####

func (p7this *S6Local) F8Get(i9ctx context.Context, key string) (any, error) {
	// 判断缓存里有没有。这里加读锁就行了，读完就可以解锁。
	p7this.p7s6lock.RLock()
	p7s6unit, ok := p7this.m3data[key]
	p7this.p7s6lock.RUnlock()
	if !ok {
		return nil, errKeyNotFound
	}
	// 判断缓存过期没有
	now := time.Now()
	if !p7s6unit.f8CheckDeadline(now) {
		// 缓存过期，删除缓存，需要加写锁。
		p7this.p7s6lock.Lock()
		defer p7this.p7s6lock.Unlock()
		// 二次校验，防止别的线程抢先操作删除了
		p7s6unit, ok = p7this.m3data[key]
		if !ok {
			return nil, errKeyNotFound
		}
		if !p7s6unit.f8CheckDeadline(now) {
			err := p7this.f8Delete(i9ctx, key, p7s6unit)
			if nil != err {
				return nil, err
			}
		}
		// 缓存过期可以归类为找不到
		return nil, errKeyNotFound
	}
	// 这里记得把缓存的值拿出来
	return p7s6unit.value, nil
}

func (p7this *S6Local) F8Set(i9ctx context.Context, key string, value any, expiration time.Duration) error {
	var deadline time.Time = time.Time{}

	// 修改操作，需要加写锁
	p7this.p7s6lock.Lock()
	defer p7this.p7s6lock.Unlock()
	// 计算过期时间
	if 0 < expiration {
		deadline = time.Now().Add(expiration)
	}
	// 设置缓存的时候，把传进来的 any 用 s6Unit 包装一下
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
	// 二次校验，因为 f8Delete 里面可能会重复触发 f8OnEvicted
	p7s6unit, ok := p7this.m3data[key]
	if !ok {
		return nil
	}
	return p7this.f8Delete(i9ctx, key, p7s6unit)
}

// f8Delete 删除缓存
func (p7this *S6Local) f8Delete(i9ctx context.Context, key string, p7s6unit *s6Unit) error {
	// 一般删除的时候，肯定外面都是拿着锁进来的，这里就不需要管并发访问 map 的问题了
	_, ok := p7this.m3data[key]
	if ok {
		delete(p7this.m3data, key)
	}
	if nil != p7this.f8OnEvicted {
		return p7this.f8OnEvicted(i9ctx, key, p7s6unit)
	}
	return nil
}

// F8Close 关闭缓存
func (p7this *S6Local) F8Close() error {
	// 这里可以做一些关闭缓存的时候需要处理的事情，比如把缓存刷新到数据库之类的
	// 关闭的时候要注意，防止多次关闭 chan 引发 panic
	p7this.p7s6CloseOnce.Do(func() {
		p7this.c4CloseSignal <- struct{}{}
		close(p7this.c4CloseSignal)
	})
	return nil
}

// F8AutoClean 自动清理过期缓存
func (p7this *S6Local) F8AutoClean() {
	t4CleanTicker := time.NewTicker(time.Second)
	go func() {
		for {
			select {
			case <-t4CleanTicker.C:
				// 遍历的时候要控制数量，防止自动清理执行太长时间，这里加了写锁的
				var count int = 0
				p7this.p7s6lock.Lock()
				for key, p7s6unit := range p7this.m3data {
					if !p7s6unit.f8CheckDeadline(time.Now()) {
						_ = p7this.f8Delete(context.Background(), key, p7s6unit)
					}
					count++
					if 1000 < count {
						break
					}
				}
				p7this.p7s6lock.Unlock()
			case <-p7this.c4CloseSignal:
				return
			}
		}
	}()
}
