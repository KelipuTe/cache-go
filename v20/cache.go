package v20

import (
	"context"
	"time"
)

type I9Cache interface {
	F8Get(i9ctx context.Context, key string) (any, error)
	F8Set(i9ctx context.Context, key string, value any, expiration time.Duration) error
	F8Delete(i9ctx context.Context, key string) error
}
