package repository

import (
	sq "github.com/Masterminds/squirrel"
	user_repo "github.com/abdulazizax/yelp/internal/repository/user"
	"github.com/abdulazizax/yelp/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	UserRepo() user_repo.UserRepository
}

type repository struct {
	userRepo user_repo.UserRepository
}

func NewRepository(db *pgxpool.Pool, queryBuilder sq.StatementBuilderType, log logger.Interface) Repository {
	return &repository{
		userRepo: user_repo.NewUserRepository(db, queryBuilder, log),
	}
}

func (r *repository) UserRepo() user_repo.UserRepository {
	return r.userRepo
}
