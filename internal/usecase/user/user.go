package user_usecase

import (
	"context"

	"github.com/abdulazizax/yelp/internal/adapter/redis/cache"
	user_entity "github.com/abdulazizax/yelp/internal/entity/user"
	user_repo "github.com/abdulazizax/yelp/internal/repository/user"
	"github.com/abdulazizax/yelp/pkg/logger"
)

type UserUsecase interface {
	CreateUser(ctx context.Context, req *user_entity.CreateUser) error
}

type userUsecase struct {
	repo  user_repo.UserRepository
	cache cache.RedisCache
	log   logger.Interface
}

func NewUserUsecase(repo user_repo.UserRepository, cache cache.RedisCache, log logger.Interface) UserUsecase {
	return &userUsecase{
		repo:  repo,
		cache: cache,
		log:   log,
	}
}

func (u *userUsecase) CreateUser(ctx context.Context, req *user_entity.CreateUser) error {
	return u.repo.CreateUser(ctx, req)
}
