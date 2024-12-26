// internal/adapter/redis/redis_cache_impl.go
package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8" // Redis uchun Go kutubxonasi
)

type RedisCacheImpl struct {
	client *redis.Client
}

// NewRedisCacheImpl - RedisCache interfeysini amalga oshiradigan structni yaratadi
func NewRedisCacheImpl(client *redis.Client) *RedisCacheImpl {
	return &RedisCacheImpl{
		client: client,
	}
}

// Set - Cachega ma'lumot qo'shish
func (r *RedisCacheImpl) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	val, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("error marshaling value: %w", err)
	}

	err = r.client.Set(ctx, key, val, expiration).Err()
	if err != nil {
		return fmt.Errorf("error setting value in cache: %w", err)
	}
	return nil
}

// Get - Cache'dan ma'lumot olish
func (r *RedisCacheImpl) Get(ctx context.Context, key string, result interface{}) error {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("error getting value from cache: %w", err)
	}

	err = json.Unmarshal([]byte(val), result)
	if err != nil {
		return fmt.Errorf("error unmarshaling value from cache: %w", err)
	}
	return nil
}

// Delete - Cache'dan ma'lumot o'chirish
func (r *RedisCacheImpl) Delete(ctx context.Context, key string) error {
	err := r.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("error deleting value from cache: %w", err)
	}
	return nil
}
