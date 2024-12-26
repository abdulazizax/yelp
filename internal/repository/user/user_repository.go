package user_repo

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"github.com/abdulazizax/yelp/internal/adapter/db/repos"
	user_entity "github.com/abdulazizax/yelp/internal/entity/user"
	"github.com/abdulazizax/yelp/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	CreateUser(ctx context.Context, req *user_entity.CreateUser) error
}

func NewUserRepository(db *pgxpool.Pool, queryBuilder sq.StatementBuilderType, log logger.Interface) UserRepository {
	return repos.NewUserRepo(db, queryBuilder, log)
}
