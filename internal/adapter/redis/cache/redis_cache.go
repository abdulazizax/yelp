// internal/adapter/redis/redis_cache.go
package cache

import (
	"context"
	"time"
)

// RedisCache - Redis caching uchun interfeys
type RedisCache interface {
	StoreEmailAndCode(ctx context.Context, email string, code int, duration time.Duration) error
	GetCodeByEmail(ctx context.Context, email string) (int, error)
}
