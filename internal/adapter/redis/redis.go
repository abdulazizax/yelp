package redis

import (
	"context"
	"fmt"
	"log"

	"github.com/abdulazizax/yelp/config"
	"github.com/go-redis/redis/v8"
)

func ConnectRedis(config *config.Config) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.Redis.Host, config.Redis.Port),
		Username: config.Redis.Username,
		Password: config.Redis.Password,
		DB:       0,
	})

	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("unable to connect to Redis: %w", err)
	}

	log.Printf("--------------------------- Connected to Redis at %s:%s ---------------------------\n", config.Redis.Host, config.Redis.Port)

	return rdb, nil
}
