// internal/adapter/redis/redis_cache.go
package cache

import (
	"context"
	"time"
)

// RedisCache - Redis caching uchun interfeys
type RedisCache interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string, result interface{}) error
	Delete(ctx context.Context, key string) error
}
