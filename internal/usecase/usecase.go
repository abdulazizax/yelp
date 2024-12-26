package usecase

import (
	"github.com/abdulazizax/yelp/internal/adapter/redis/cache"
	"github.com/abdulazizax/yelp/internal/repository"
	user_usecase "github.com/abdulazizax/yelp/internal/usecase/user"
	"github.com/abdulazizax/yelp/pkg/logger"
)

type Usecase interface {
	UserUsecase() user_usecase.UserUsecase
}

type usecase struct {
	user user_usecase.UserUsecase
}

func NewUsecase(repo repository.Repository, cache cache.RedisCache, log logger.Interface) Usecase {
	return &usecase{
		user: user_usecase.NewUserUsecase(repo.UserRepo(), cache, log),
	}
}

func (u *usecase) UserUsecase() user_usecase.UserUsecase {
	return u.user
}
