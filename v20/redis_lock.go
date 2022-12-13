package v20

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v9"
	"sync"
	"time"
)

// luaLock 加锁脚本

//go:embed lua/lock.lua
var luaLock string

//go:embed lua/unlock.lua
var luaUnlock string

//go:embed lua/refresh.lua
var luaRefresh string

var ErrLockHoldByOthers = errors.New("lock hold by others")
var ErrLockNotHold = errors.New("lock not hold")

// #### type ####

// redis 锁
type S6RedisLock struct {
	i9RedisClient redis.Cmdable
	// 用于标记自己，去 redis 加锁的时候，对应 key 的 value
	selfTag string
}

// 锁
// 锁里面带着的数据是为了持有时间到了之后续约用的
type S6Lock struct {
	i9RedisClient redis.Cmdable
	selfTag       string
	key           string
	expiration    time.Duration

	// 控制解锁只能一次
	s6UnlockOnce sync.Once
	// 解锁信号
	c4UnlockSignal chan struct{}
}

// 加锁重试接口
type I9LockRetry interface {
	// 有没有下次重试，下次重试的时间间隔
	F8NextTry() (bool, time.Duration)
}

type S6LockRetry struct {
	// 重试的时间间隔
	TimeInterval time.Duration
	// 第几次尝试
	NowTime int
	// 最大重试几次
	MaxTime int
}

// #### func ####

func F8NewS6RedisLock(i9RedisClient redis.Cmdable) *S6RedisLock {
	return &S6RedisLock{
		i9RedisClient: i9RedisClient,
		selfTag:       "temp-value",
	}
}

func f8NewS6Lock(i9RedisClient redis.Cmdable, selfTag string, key string, expiration time.Duration) *S6Lock {
	return &S6Lock{
		i9RedisClient:  i9RedisClient,
		selfTag:        selfTag,
		key:            key,
		expiration:     expiration,
		s6UnlockOnce:   sync.Once{},
		c4UnlockSignal: make(chan struct{}, 1),
	}
}

// #### type func ####

// 带重试次数的加锁
// i9ctx 控制整体加锁的超时时间
// key 锁的名字
// expiration 锁的持有时间
// timeout 每次加锁的超时时间
// i9retry 重试
func (p7this *S6RedisLock) F8Lock(i9ctx context.Context, key string, expiration time.Duration, timeout time.Duration, i9retry I9LockRetry) (*S6Lock, error) {
	// 用于控制重试的间隔
	var p7RetryTimer *time.Timer
	defer func() {
		if nil != p7RetryTimer {
			p7RetryTimer.Stop()
		}
	}()

	for {
		// 拿一个超时的 i9ctx
		t4ctx, t4f8cancel := context.WithTimeout(i9ctx, timeout)
		// 执行 redis 加锁 lua 脚本，成功的话这里应该返回 OK
		result, err := p7this.i9RedisClient.Eval(t4ctx, luaLock, []string{key}, p7this.selfTag, expiration.Seconds()).Result()
		t4f8cancel()
		if nil != err && !errors.Is(err, context.DeadlineExceeded) {
			// 如果不是超时，那大概率是寄了，没必要重试了
			// 但是这里会有小概率是网络波动造成的 redis 失联，不过这里可以丢给上层去处理
			return nil, err
		}

		// redis 返回 OK 表示拿到锁了
		if "OK" == result {
			return f8NewS6Lock(p7this.i9RedisClient, p7this.selfTag, key, expiration), nil
		}

		// 如果走到这里了，表示超时了或者加锁失败了，需要进入重试流程
		needNext, interval := i9retry.F8NextTry()
		if !needNext {
			// 不需要重试
			if nil != err {
				err = fmt.Errorf("failed with err: %w", err)
			} else {
				err = ErrLockHoldByOthers
			}
			return nil, fmt.Errorf("retry end: %w", err)
		}
		// 重试
		if nil == p7RetryTimer {
			p7RetryTimer = time.NewTimer(interval)
		} else {
			p7RetryTimer.Reset(interval)
		}
		select {
		case <-p7RetryTimer.C:
			// 等重试
		case <-i9ctx.Done():
			// 等整体超时
			return nil, i9ctx.Err()
		}
	}
}

// 单次尝试加锁
func (p7this *S6RedisLock) F8TryLock(i9ctx context.Context, key string, expiration time.Duration) (*S6Lock, error) {
	ok, err := p7this.i9RedisClient.SetNX(i9ctx, key, p7this.selfTag, expiration).Result()
	if nil != err {
		// 访问 redis 直接报错了
		return nil, err
	}
	if !ok {
		// 访问 redis 没问题，但是加不上锁，说明已经锁被占有
		return nil, ErrLockHoldByOthers
	}
	// 走到这里说明加到锁了
	return f8NewS6Lock(p7this.i9RedisClient, p7this.selfTag, key, expiration), nil
}

// 解锁
func (p7this *S6Lock) F8Unlock(i9ctx context.Context) error {
	// 执行 redis 解锁 lua 脚本，成功的话这里应该返回 1
	result, err := p7this.i9RedisClient.Eval(i9ctx, luaUnlock, []string{p7this.key}, p7this.selfTag).Int64()
	defer func() {
		// 发出解锁信号
		p7this.s6UnlockOnce.Do(func() {
			p7this.c4UnlockSignal <- struct{}{}
			close(p7this.c4UnlockSignal)
		})
	}()
	// redis 直接报错
	if err != nil {
		if err == redis.Nil {
			return ErrLockNotHold
		}
		return err
	}
	// redis 返回错误
	if 1 != result {
		return ErrLockNotHold
	}
	return nil
}

// 刷新锁的时间
func (p7this *S6Lock) F8Refresh(i9ctx context.Context) error {
	// 执行 redis 解锁 lua 脚本，成功的话这里应该返回 1
	result, err := p7this.i9RedisClient.Eval(i9ctx, luaRefresh, []string{p7this.key}, p7this.selfTag, p7this.expiration.Seconds()).Int64()
	// redis 直接报错
	if err != nil {
		return err
	}
	// redis 返回错误
	if result != 1 {
		return ErrLockNotHold
	}
	return nil
}

// interval 自动刷新锁的时间的时间间隔，这个间隔可以短一点，需要给重试留时间
// timeout 每次自动刷新的超时时间
func (p7this *S6Lock) F8AutoRefresh(interval time.Duration, timeout time.Duration) error {
	// 用于控制自动刷新的间隔
	t4RefreshTicker := time.NewTicker(interval)
	// 自动刷新超时时的重试信号
	c4retry := make(chan struct{}, 1)
	for {
		select {
		case <-t4RefreshTicker.C:
			fmt.Println("F8AutoRefresh.t4RefreshTicker.C")
			// 正常情况走这里
			t4ctx, t4f8cancel := context.WithTimeout(context.Background(), timeout)
			err := p7this.F8Refresh(t4ctx)
			t4f8cancel()

			if err == context.DeadlineExceeded {
				// 超时可以重试
				// 注意这里不能直接 c4retry <- struct{}{}。如果时间间隔和超时同时存在，
				// 这个时候 select 走了时间间隔，但是又超时了，那就会在这里阻塞住，让重试先跑。
				select {
				case c4retry <- struct{}{}:
				default:
				}
				continue
			}
			// 其他错误就不重试了
			if nil != err {
				return err
			}
		case <-c4retry:
			fmt.Println("F8AutoRefresh.c4retry")
			// 重试情况走这里
			// 逻辑和正常情况是一样的，cv 过来都可以
			t4ctx, t4f8cancel := context.WithTimeout(context.Background(), timeout)
			err := p7this.F8Refresh(t4ctx)
			t4f8cancel()

			if err == context.DeadlineExceeded {
				// 重试如果超时了，其实不管问题也不大
				select {
				// 这和上面那里是一样的，如果时间间隔和超时同时存在，
				// 这个时候 select 走了重试间隔，但是又超时了，那就会在这里阻塞住，让时间间隔先跑。
				case c4retry <- struct{}{}:
				default:
				}
				continue
			}
			if nil != err {
				return err
			}
		case <-p7this.c4UnlockSignal:
			fmt.Println("F8AutoRefresh.c4UnlockSignal")
			// 接到解锁信号，退出自动刷新
			return nil
		}
	}
}

func (p7this *S6LockRetry) F8NextTry() (bool, time.Duration) {
	p7this.NowTime++
	return p7this.NowTime <= p7this.MaxTime, p7this.TimeInterval
}
