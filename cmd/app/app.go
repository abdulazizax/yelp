package app

import (
	"log/slog"

	sq "github.com/Masterminds/squirrel"
	"github.com/abdulazizax/yelp/config"
	"github.com/abdulazizax/yelp/internal/adapter/db"
	"github.com/abdulazizax/yelp/internal/adapter/http"
	handler "github.com/abdulazizax/yelp/internal/adapter/http/handlers"
	"github.com/abdulazizax/yelp/internal/adapter/redis"
	"github.com/abdulazizax/yelp/internal/adapter/redis/cache"
	"github.com/abdulazizax/yelp/internal/repository"
	"github.com/abdulazizax/yelp/internal/usecase"

	"github.com/abdulazizax/yelp/pkg/logger"
)

func Run(cfg *config.Config) error {
	// Set up logger
	log := logger.New(cfg.Logger.Level, "application.log")

	// // Initialize Postgres connection
	db, err := db.ConnectDB(cfg)
	if err != nil {
		log.Fatal("Error while connecting to MongoDB", slog.String("err", err.Error()))
		return err
	}

	queryBuilder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	// // Initialize Redis connection
	redis, err := redis.ConnectRedis(cfg)
	if err != nil {
		log.Fatal("Error while connecting to Redis", slog.String("err", err.Error()))
		return err
	}

	cache := cache.NewRedisCacheImpl(redis)

	repo := repository.NewRepository(db, queryBuilder, log)

	usecase := usecase.NewUsecase(repo, cache, log)

	// Initialize HTTP handler
	handler := handler.NewHandlers(usecase, log)

	// // Start the HTTP server
	return http.Roter(handler, cfg)
}
