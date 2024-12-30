package user_repo

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	user_db "github.com/abdulazizax/yelp/internal/adapter/db/repos/user"
	user_entity "github.com/abdulazizax/yelp/internal/entity/user"
	"github.com/abdulazizax/yelp/pkg/logger"
	"github.com/abdulazizax/yelp/pkg/security"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	CreateUser(ctx context.Context, req *user_entity.CreateUser) error
	SignIn(ctx context.Context, user_id string, signInRequest *user_entity.SignInRequest) error
	CheckUserExists(ctx context.Context, email string) error
	GetUserInfoByEmail(ctx context.Context, email string) (*user_entity.UserInfo, error)
	UpdateUserPassword(ctx context.Context, req *user_entity.UpdateUserPassword) error
}

func NewUserRepository(db *pgxpool.Pool, queryBuilder sq.StatementBuilderType, email security.EmailService, log logger.Interface) UserRepository {
	return user_db.NewUserRepo(db, queryBuilder, email, log)
}
