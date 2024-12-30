package user_usecase

import (
	"context"
	"time"

	"github.com/abdulazizax/yelp/config"
	"github.com/abdulazizax/yelp/internal/adapter/auth"
	"github.com/abdulazizax/yelp/internal/adapter/redis/cache"
	"github.com/abdulazizax/yelp/internal/entity"
	user_entity "github.com/abdulazizax/yelp/internal/entity/user"
	user_repo "github.com/abdulazizax/yelp/internal/repository/user"
	"github.com/abdulazizax/yelp/pkg/logger"
	"github.com/abdulazizax/yelp/pkg/security"
)

type UserUsecase interface {
	CreateUser(ctx context.Context, req *user_entity.CreateUser) error
	SignIn(ctx context.Context, signInRequest *user_entity.SignInRequest) (token string, err error)
	SendVerificationCode(ctx context.Context, email string) (time.Duration, error)
	UpdateUserPassword(ctx context.Context, req *user_entity.UpdateUserPassword) error
}

type userUsecase struct {
	repo  user_repo.UserRepository
	cache cache.RedisCache
	email security.EmailService
	log   logger.Interface
	cfg   *config.Config
}

func NewUserUsecase(repo user_repo.UserRepository, cache cache.RedisCache, cfg *config.Config, email security.EmailService, log logger.Interface) UserUsecase {
	return &userUsecase{
		repo:  repo,
		cache: cache,
		email: email,
		log:   log,
		cfg:   cfg,
	}
}

func (u *userUsecase) CreateUser(ctx context.Context, req *user_entity.CreateUser) error {
	err := u.repo.CheckUserExists(ctx, req.Email)
	if err != nil {
		u.log.Error("Error checking user existence", "error", err)
		return err
	}

	// Validate password
	err = security.ValidatePassword(req.Password)
	if err != nil {
		u.log.Error("Password validation failed", "error", err)
		return err
	}

	// Create new user
	err = u.repo.CreateUser(ctx, req)
	if err != nil {
		u.log.Error("Failed to create user", "error", err)
		return entity.ErrFailedToCreate
	}

	u.log.Info("User created successfully", "email", req.Email)
	return nil
}

func (u *userUsecase) SignIn(ctx context.Context, signInRequest *user_entity.SignInRequest) (token string, err error) {
	userInfo, err := u.repo.GetUserInfoByEmail(ctx, signInRequest.User.Email)
	if err == entity.ErrUserNotFound {
		return "", entity.ErrUserNotFound
	} else if err != nil {
		return "", err
	}

	err = security.CheckPassword(userInfo.PasswordHash, signInRequest.User.Password)
	if err != nil {
		return "", entity.ErrIncorrectPassword
	}

	err = u.repo.SignIn(ctx, userInfo.ID, signInRequest)
	if err != nil {
		return "", err
	}

	var tokenClaims = auth.TokenClaims{
		UserID: userInfo.ID,
		Email:  signInRequest.User.Email,
		Role:   userInfo.UserRole,
		Type:   userInfo.UserType,
	}

	token, err = auth.GenerateAccessToken(u.cfg.JWT.SecretKey, &tokenClaims)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (u *userUsecase) UpdateUserPassword(ctx context.Context, req *user_entity.UpdateUserPassword) error {
	return u.repo.UpdateUserPassword(ctx, req)
}

func (u *userUsecase) SendVerificationCode(ctx context.Context, email string) (time.Duration, error) {
	err := u.repo.CheckUserExists(ctx, email)
	if err != nil {
		u.log.Error("Error checking user existence", "error", err)
		return 0, err
	}

	return u.email.SendVerificationCode(ctx, email)
}
