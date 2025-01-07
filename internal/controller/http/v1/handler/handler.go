package handler

import (
	"github.com/abdulazizax/yelp/config"
	"github.com/abdulazizax/yelp/internal/usecase"
	"github.com/abdulazizax/yelp/pkg/logger"
	rediscache "github.com/golanguzb70/redis-cache"
)

type Handler struct {
	Logger  *logger.Logger
	Config  *config.Config
	UseCase *usecase.UseCase
	Redis   rediscache.RedisCache
}

func NewHandler(l *logger.Logger, c *config.Config, useCase *usecase.UseCase, redis rediscache.RedisCache) *Handler {
	return &Handler{
		Logger:  l,
		Config:  c,
		UseCase: useCase,
		Redis:   redis,
	}
}
