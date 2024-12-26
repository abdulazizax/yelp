package repos

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/abdulazizax/yelp/internal/entity"
	user_entity "github.com/abdulazizax/yelp/internal/entity/user"
	"github.com/abdulazizax/yelp/pkg/logger"
	"github.com/abdulazizax/yelp/pkg/security"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	db           *pgxpool.Pool
	queryBuilder sq.StatementBuilderType
	log          logger.Interface
}

func NewUserRepo(db *pgxpool.Pool, queryBuilder sq.StatementBuilderType, log logger.Interface) *UserRepo {
	return &UserRepo{
		db:           db,
		queryBuilder: queryBuilder,
		log:          log,
	}
}

func (u *UserRepo) CreateUser(ctx context.Context, req *user_entity.CreateUser) error {
	id := uuid.New().String()
	hashedPassword, err := security.HashPassword(req.Password)
	if err != nil {
		u.log.Error("Error hashing password", err)
		return err
	}
	query, args, err := u.queryBuilder.Insert("users").
		Columns("id", "user_type", "user_role", "name", "email", "password_hash", "gender", "status").
		Values(
			id,
			entity.TypeUser,
			entity.RoleUser,
			req.Name,
			req.Email,
			hashedPassword,
			req.Gender,
			entity.StatusInverify,
		).ToSql()
	if err != nil {
		u.log.Error("Error generating SQL:", err)
		return err
	}

	_, err = u.db.Exec(ctx, query, args...)
	if err != nil {
		u.log.Error("Error creating row:", err)
		return err
	}

	u.log.Info("User created successfully, user_id: ", id)
	return nil
}
