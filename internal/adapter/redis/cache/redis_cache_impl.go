// internal/adapter/redis/redis_cache_.go
package cache

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/abdulazizax/yelp/pkg/logger"
	"github.com/go-redis/redis/v8"
)

type redisCache struct {
	redisDb *redis.Client
	log     logger.Interface
}

// NewRedisCache - RedisCache interfeysini amalga oshiradigan structni yaratadi
func NewRedisCache(redisDb *redis.Client, log logger.Interface) RedisCache {
	return &redisCache{
		redisDb: redisDb,
		log:     log,
	}
}

func (r *redisCache) StoreEmailAndCode(ctx context.Context, email string, code int, duration time.Duration) error {
	codeKey := "verification_code:" + email
	err := r.redisDb.Set(ctx, codeKey, code, time.Minute*duration).Err()
	if err != nil {
		r.log.Error("Error while storing verification code", slog.String("error", err.Error()))
		return err
	}
	return nil
}

func (r *redisCache) GetCodeByEmail(ctx context.Context, email string) (int, error) {
	codeKey := "verification_code:" + email
	codeStr, err := r.redisDb.Get(ctx, codeKey).Result()
	if err == redis.Nil {
		return 0, nil
	} else if err != nil {
		r.log.Error("Error while getting verification code", slog.String("error", err.Error()))
		return 0, err
	}

	var code int
	_, err = fmt.Sscanf(codeStr, "%d", &code)
	if err != nil {
		r.log.Error("Error while parsing verification code", slog.String("error", err.Error()))
		return 0, err
	}

	return code, nil
}
