package user_repo

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	user_db "github.com/abdulazizax/yelp/internal/adapter/db/repos/user"
	user_entity "github.com/abdulazizax/yelp/internal/entity/user"
	"github.com/abdulazizax/yelp/pkg/logger"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	CreateUser(ctx context.Context, req *user_entity.CreateUser) error
	CreateSession(ctx context.Context, tx pgx.Tx, req *user_entity.CreateSession) error
	ActivateUser(ctx context.Context, tx pgx.Tx, user_id string) error
	CheckUserExists(ctx context.Context, email string) (bool, error)
	GetUserPassword(ctx context.Context, email string) (string, error)
	BeginTx(ctx context.Context) (pgx.Tx, error)
}

func NewUserRepository(db *pgxpool.Pool, queryBuilder sq.StatementBuilderType, log logger.Interface) UserRepository {
	return user_db.NewUserRepo(db, queryBuilder, log)
}
