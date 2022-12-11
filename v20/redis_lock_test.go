package v20

import (
	"cache-go/v20/mock"
	"context"
	"errors"
	"github.com/go-redis/redis/v9"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestS6RedisLockF8Lock(p7s6t *testing.T) {
	p7s6t.Parallel()
	p7ctrl := gomock.NewController(p7s6t)
	defer p7ctrl.Finish()

	s5s6case := []struct {
		name        string
		i9redisMock func() redis.Cmdable
		key         string
		expiration  time.Duration
		i9retry     I9LockRetry
		timeout     time.Duration
		wantLock    *S6Lock
		wantErr     string
	}{
		{
			name: "lock",
			i9redisMock: func() redis.Cmdable {
				p7cmdable := mock.NewMockCmdable(p7ctrl)
				p7cmd := redis.NewCmd(context.Background(), nil)
				p7cmd.SetVal("OK")
				p7cmdable.EXPECT().Eval(gomock.Any(), luaLock, []string{"key1"}, gomock.Any()).Return(p7cmd)
				return p7cmdable
			},
			key:        "key1",
			expiration: time.Minute,
			i9retry:    &S6LockRetry{TimeInterval: time.Millisecond, NowTime: 0, MaxTime: 1},
			timeout:    30 * time.Second,
			wantLock: &S6Lock{
				key:        "key1",
				expiration: time.Minute,
			},
		},
		{
			name: "lock_fail",
			i9redisMock: func() redis.Cmdable {
				p7cmdable := mock.NewMockCmdable(p7ctrl)
				p7cmd := redis.NewCmd(context.Background(), nil)
				p7cmd.SetErr(errors.New("redis error"))
				p7cmdable.EXPECT().Eval(gomock.Any(), luaLock, []string{"key1"}, gomock.Any()).Return(p7cmd)
				return p7cmdable
			},
			key:        "key1",
			expiration: time.Minute,
			i9retry:    &S6LockRetry{TimeInterval: time.Millisecond, NowTime: 0, MaxTime: 0},
			timeout:    30 * time.Second,
			wantErr:    "redis error",
		},
		{
			name: "lock_with_retry",
			i9redisMock: func() redis.Cmdable {
				p7cmdable := mock.NewMockCmdable(p7ctrl)
				p7cmd := redis.NewCmd(context.Background(), nil)
				p7cmd.SetErr(context.DeadlineExceeded)
				p7cmdable.EXPECT().Eval(gomock.Any(), luaLock, []string{"key1"}, gomock.Any()).Times(2).Return(p7cmd)

				p7cmd2 := redis.NewCmd(context.Background(), nil)
				p7cmd2.SetVal("OK")
				p7cmdable.EXPECT().Eval(gomock.Any(), luaLock, []string{"key1"}, gomock.Any()).Return(p7cmd2)
				return p7cmdable
			},
			key:        "key1",
			expiration: time.Minute,
			i9retry:    &S6LockRetry{TimeInterval: time.Millisecond, NowTime: 0, MaxTime: 2},
			timeout:    30 * time.Second,
			wantLock: &S6Lock{
				key:        "key1",
				expiration: time.Minute,
			},
		},
		{
			name: "lock_with_retry_fail_overtime",
			i9redisMock: func() redis.Cmdable {
				p7cmdable := mock.NewMockCmdable(p7ctrl)
				p7cmd := redis.NewCmd(context.Background(), nil)
				p7cmd.SetErr(context.DeadlineExceeded)
				p7cmdable.EXPECT().Eval(gomock.Any(), luaLock, []string{"key1"}, gomock.Any()).Times(3).Return(p7cmd)
				return p7cmdable
			},
			key:        "key1",
			expiration: time.Minute,
			i9retry:    &S6LockRetry{TimeInterval: time.Millisecond, NowTime: 0, MaxTime: 2},
			timeout:    30 * time.Second,
			wantErr:    "retry end: failed with err: context deadline exceeded",
		},
		{
			name: "lock_with_retry_fail_lock_hold_by_others",
			i9redisMock: func() redis.Cmdable {
				p7cmdable := mock.NewMockCmdable(p7ctrl)
				p7cmd := redis.NewCmd(context.Background(), nil)
				p7cmdable.EXPECT().Eval(gomock.Any(), luaLock, []string{"key1"}, gomock.Any()).Times(3).Return(p7cmd)
				return p7cmdable
			},
			key:        "key1",
			expiration: time.Minute,
			i9retry:    &S6LockRetry{TimeInterval: time.Millisecond, NowTime: 0, MaxTime: 2},
			timeout:    30 * time.Second,
			wantErr:    "retry end: lock hold by others",
		},
		{
			name: "lock_with_retry_fail_timeout",
			i9redisMock: func() redis.Cmdable {
				p7cmdable := mock.NewMockCmdable(p7ctrl)
				p7cmd := redis.NewCmd(context.Background(), nil)
				// 这里测试的时候，注意控制调用次数，调用次数不够它会报错
				p7cmdable.EXPECT().Eval(gomock.Any(), luaLock, []string{"key1"}, gomock.Any()).Times(2).Return(p7cmd)
				return p7cmdable
			},
			key:        "key1",
			expiration: time.Minute,
			i9retry:    &S6LockRetry{TimeInterval: 600 * time.Millisecond, NowTime: 0, MaxTime: 2},
			timeout:    time.Second,
			wantErr:    "context deadline exceeded",
		},
	}

	for _, t4value := range s5s6case {
		p7s6t.Run(t4value.name, func(p7s6t *testing.T) {
			mockRedisCmd := t4value.i9redisMock()
			p7RedisLock := F8NewS6RedisLock(mockRedisCmd)
			ctx, cancel := context.WithTimeout(context.Background(), t4value.timeout)
			defer cancel()

			p7s6Lock, err := p7RedisLock.F8Lock(ctx, t4value.key, t4value.expiration, time.Second, t4value.i9retry)

			if t4value.wantErr != "" {
				assert.EqualError(p7s6t, err, t4value.wantErr)
				return
			} else {
				require.NoError(p7s6t, err)
			}

			assert.Equal(p7s6t, mockRedisCmd, p7s6Lock.i9RedisClient)
			assert.Equal(p7s6t, p7RedisLock.selfTag, p7s6Lock.selfTag)
			assert.Equal(p7s6t, t4value.key, p7s6Lock.key)
			assert.Equal(p7s6t, t4value.expiration, p7s6Lock.expiration)
		})
	}
}

func TestS6LockF8Unlock(p7s6t *testing.T) {
	p7s6t.Parallel()
	p7ctrl := gomock.NewController(p7s6t)
	defer p7ctrl.Finish()

	s5s6case := []struct {
		name        string
		i9redisMock func() redis.Cmdable
		wantErr     error
	}{
		{
			name: "unlock",
			i9redisMock: func() redis.Cmdable {
				p7cmdable := mock.NewMockCmdable(p7ctrl)
				p7cmd := redis.NewCmd(context.Background(), nil)
				p7cmd.SetVal(int64(1))
				p7cmdable.EXPECT().Eval(gomock.Any(), luaUnlock, []string{"key1"}, gomock.Any()).Return(p7cmd)
				return p7cmdable
			},
		},
		{
			name: "unlock_fail",
			i9redisMock: func() redis.Cmdable {
				p7cmdable := mock.NewMockCmdable(p7ctrl)
				p7cmd := redis.NewCmd(context.Background(), nil)
				p7cmd.SetErr(errors.New("redis error"))
				p7cmdable.EXPECT().Eval(gomock.Any(), luaUnlock, []string{"key1"}, gomock.Any()).Return(p7cmd)
				return p7cmdable
			},
			wantErr: errors.New("redis error"),
		},
		{
			name: "lock_not_hold",
			i9redisMock: func() redis.Cmdable {
				p7cmdable := mock.NewMockCmdable(p7ctrl)
				p7cmd := redis.NewCmd(context.Background(), nil)
				p7cmd.SetVal(int64(0))
				p7cmdable.EXPECT().Eval(gomock.Any(), luaUnlock, []string{"key1"}, gomock.Any()).Return(p7cmd)
				return p7cmdable
			},
			wantErr: ErrLockNotHold,
		},
	}

	for _, t4value := range s5s6case {
		p7s6t.Run(t4value.name, func(p7s6t *testing.T) {
			p7s6lock := f8NewS6Lock(t4value.i9redisMock(), "temp-value", "key1", time.Minute)
			err := p7s6lock.F8Unlock(context.Background())
			assert.Equal(p7s6t, t4value.wantErr, err)
		})
	}
}

func TestS6LockF8Refresh(p7s6t *testing.T) {
	p7s6t.Parallel()
	p7ctrl := gomock.NewController(p7s6t)
	defer p7ctrl.Finish()

	s5s6case := []struct {
		name       string
		s6lockMock func() *S6Lock
		wantErr    error
	}{
		{
			name: "refresh",
			s6lockMock: func() *S6Lock {
				p7cmdable := mock.NewMockCmdable(p7ctrl)
				p7cmd := redis.NewCmd(context.Background(), nil)
				p7cmd.SetVal(int64(1))
				p7cmdable.EXPECT().Eval(gomock.Any(), luaRefresh, []string{"key1"}, gomock.Any()).Return(p7cmd)
				return &S6Lock{
					i9RedisClient: p7cmdable,
					selfTag:       "temp-value",
					key:           "key1",
					expiration:    time.Minute,
				}
			},
		},
		{
			name: "refresh_fail",
			s6lockMock: func() *S6Lock {
				p7cmdable := mock.NewMockCmdable(p7ctrl)
				p7cmd := redis.NewCmd(context.Background(), nil)
				p7cmd.SetErr(redis.Nil)
				p7cmdable.EXPECT().Eval(gomock.Any(), luaRefresh, []string{"key1"}, gomock.Any()).Return(p7cmd)
				return &S6Lock{
					i9RedisClient: p7cmdable,
					selfTag:       "temp-value",
					key:           "key1",
					expiration:    time.Minute,
				}
			},
			wantErr: redis.Nil,
		},
		{
			name: "lock_not_hold",
			s6lockMock: func() *S6Lock {
				p7cmdable := mock.NewMockCmdable(p7ctrl)
				p7cmd := redis.NewCmd(context.Background(), nil)
				p7cmd.SetVal(int64(0))
				p7cmdable.EXPECT().Eval(gomock.Any(), luaRefresh, []string{"key1"}, gomock.Any()).Return(p7cmd)
				return &S6Lock{
					i9RedisClient: p7cmdable,
					selfTag:       "temp-value",
					key:           "key1",
					expiration:    time.Minute,
				}
			},
			wantErr: ErrLockNotHold,
		},
	}

	for _, t4value := range s5s6case {
		p7s6t.Run(t4value.name, func(p7s6t *testing.T) {
			err := t4value.s6lockMock().F8Refresh(context.Background())
			assert.Equal(p7s6t, t4value.wantErr, err)
		})
	}
}

func TestS6LockF8AutoRefresh(p7s6t *testing.T) {
	p7s6t.Parallel()
	p7ctrl := gomock.NewController(p7s6t)
	defer p7ctrl.Finish()

	s5s6case := []struct {
		name        string
		s6lockMock  func() *S6Lock
		unlockAfter time.Duration
		interval    time.Duration
		timeout     time.Duration
		wantErr     error
	}{
		{
			name: "auto_refresh",
			s6lockMock: func() *S6Lock {
				p7cmdable := mock.NewMockCmdable(p7ctrl)
				p7cmd := redis.NewCmd(context.Background(), nil)
				p7cmd.SetVal(int64(1))
				p7cmdable.EXPECT().Eval(gomock.Any(), luaRefresh, []string{"key1"}, gomock.Any()).AnyTimes().Return(p7cmd)

				p7cmd2 := redis.NewCmd(context.Background())
				p7cmd2.SetVal(int64(1))
				p7cmdable.EXPECT().Eval(gomock.Any(), luaUnlock, gomock.Any(), gomock.Any()).Return(p7cmd2)
				return &S6Lock{
					i9RedisClient:  p7cmdable,
					selfTag:        "temp-value",
					key:            "key1",
					expiration:     time.Minute,
					c4UnlockSignal: make(chan struct{}, 1),
				}
			},
			unlockAfter: time.Second,
			interval:    200 * time.Millisecond,
			timeout:     200 * time.Millisecond,
			wantErr:     nil,
		},
		{
			name: "auto_refresh_fail",
			s6lockMock: func() *S6Lock {
				p7cmdable := mock.NewMockCmdable(p7ctrl)
				p7cmd := redis.NewCmd(context.Background(), nil)
				p7cmd.SetErr(errors.New("redis error"))
				p7cmdable.EXPECT().Eval(gomock.Any(), luaRefresh, []string{"key1"}, gomock.Any()).AnyTimes().Return(p7cmd)

				//p7cmd2 := redis.NewCmd(context.Background())
				//p7cmd2.SetVal(int64(1))
				//p7cmdable.EXPECT().Eval(gomock.Any(), luaUnlock, gomock.Any(), gomock.Any()).Return(p7cmd2)
				return &S6Lock{
					i9RedisClient:  p7cmdable,
					selfTag:        "temp-value",
					key:            "key1",
					expiration:     time.Minute,
					c4UnlockSignal: make(chan struct{}, 1),
				}
			},
			unlockAfter: time.Second,
			interval:    200 * time.Millisecond,
			timeout:     200 * time.Millisecond,
			wantErr:     errors.New("redis error"),
		},
		{
			name: "auto_refresh_timeout",
			s6lockMock: func() *S6Lock {
				p7cmdable := mock.NewMockCmdable(p7ctrl)
				p7cmd := redis.NewCmd(context.Background(), nil)
				p7cmd.SetErr(context.DeadlineExceeded)
				p7cmdable.EXPECT().Eval(gomock.Any(), luaRefresh, []string{"key1"}, gomock.Any()).Times(2).Return(p7cmd)

				p7cmd2 := redis.NewCmd(context.Background(), nil)
				p7cmd2.SetVal(int64(1))
				p7cmdable.EXPECT().Eval(gomock.Any(), luaRefresh, []string{"key1"}, gomock.Any()).AnyTimes().Return(p7cmd2)

				p7cmd4 := redis.NewCmd(context.Background())
				p7cmd4.SetVal(int64(1))
				p7cmdable.EXPECT().Eval(gomock.Any(), luaUnlock, gomock.Any(), gomock.Any()).Return(p7cmd4)
				return &S6Lock{
					i9RedisClient:  p7cmdable,
					selfTag:        "temp-value",
					key:            "key1",
					expiration:     time.Minute,
					c4UnlockSignal: make(chan struct{}, 1),
				}
			},
			unlockAfter: time.Second,
			interval:    200 * time.Millisecond,
			timeout:     200 * time.Millisecond,
			wantErr:     nil,
		},
	}

	for _, t4value := range s5s6case {
		p7s6t.Run(t4value.name, func(p7s6t *testing.T) {
			p7s6lock := t4value.s6lockMock()
			go func() {
				time.Sleep(t4value.unlockAfter)
				err := p7s6lock.F8Unlock(context.Background())
				assert.NoError(p7s6t, err)
			}()
			err := p7s6lock.F8AutoRefresh(t4value.interval, t4value.timeout)
			assert.Equal(p7s6t, t4value.wantErr, err)
		})
	}
}
