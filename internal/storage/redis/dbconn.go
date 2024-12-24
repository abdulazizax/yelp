package redis

import (
	"fmt"
	"log"
	"strconv"

	"context"

	"github.com/abdulazizax/yelp/pkg/config"
	rediscache "github.com/golanguzb70/redis-cache"
)

func ConnectRedis(config *config.Config) (*rediscache.RedisCache, error) {
	// Convert Redis port from string to integer
	port, err := strconv.Atoi(config.Redis.Port)
	if err != nil {
		return nil, fmt.Errorf("invalid Redis port: %w", err)
	}

	cfg := &rediscache.Config{
		RedisHost:     config.Redis.Host,
		RedisPort:     port,
		RedisUsername: config.Redis.Username,
		RedisPassword: config.Redis.Password,
	}

	cache, err := rediscache.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to Redis: %w", err)
	}

	ctx := context.Background()
	if err := cache.Ping(ctx); err != nil {
		return nil, fmt.Errorf("unable to ping Redis: %w", err)
	}

	log.Printf("--------------------------- Connected to Redis at %s:%d ---------------------------\n", config.Redis.Host, port)
	
	return &cache, nil
}
