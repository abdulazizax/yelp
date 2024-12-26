package user_usecase

import (
	"context"
	"fmt"

	"github.com/abdulazizax/yelp/internal/adapter/redis/cache"
	"github.com/abdulazizax/yelp/internal/entity"
	user_entity "github.com/abdulazizax/yelp/internal/entity/user"
	user_repo "github.com/abdulazizax/yelp/internal/repository/user"
	"github.com/abdulazizax/yelp/pkg/logger"
	"github.com/abdulazizax/yelp/pkg/security"
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
	isExists, err := u.repo.CheckUserExists(ctx, req.Email)
	if err != nil {
		return err
	}

	if isExists {
		return entity.ErrUserAlreadyExists
	}

	err = security.ValidatePassword(req.Password)
	if err != nil {
		return err
	}

	err = u.repo.CreateUser(ctx, req)
	if err != nil {
		return entity.ErrFailedToCreate
	}

	return nil
}

// func (u *userUsecase) SignIn(ctx context.Context, signIn user_entity.SignInUser, session user_entity.CreateSession) error {
// 	password_hash, err := u.repo.GetUserPassword(ctx, signIn.Email)
// 	if err == entity.ErrUserNotFound {
// 		return entity.ErrUserNotFound
// 	} else if err != nil {
// 		return err
// 	}

// 	err = security.CheckPassword(password_hash, signIn.Password)
// 	if err != nil {
// 		return err
// 	}

// 	err = u.repo.ActivateUser(ctx, signIn.Email)
// 	if err != nil {
// 		return err
// 	}

// 	err = u.repo.CreateSession(ctx, &session)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

func (u *userUsecase) SignIn(ctx context.Context, signIn user_entity.SignInUser, session user_entity.CreateSession) error {
	// Retrieve user password from the repository
	password_hash, err := u.repo.GetUserPassword(ctx, signIn.Email)
	if err == entity.ErrUserNotFound {
		return entity.ErrUserNotFound
	} else if err != nil {
		return err
	}

	// Check the provided password
	err = security.CheckPassword(password_hash, signIn.Password)
	if err != nil {
		return err
	}

	tx, err := u.repo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	err = u.repo.ActivateUser(ctx, tx, signIn.Email)
	if err != nil {
		return err
	}

	err = u.repo.CreateSession(ctx, tx, &session)
	if err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
