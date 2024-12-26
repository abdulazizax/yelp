package user_db

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/abdulazizax/yelp/internal/entity"
	user_entity "github.com/abdulazizax/yelp/internal/entity/user"
	"github.com/abdulazizax/yelp/pkg/logger"
	"github.com/abdulazizax/yelp/pkg/security"
	"github.com/google/uuid"
	"github.com/jackc/pgx"
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

// BeginTx starts a new database transaction.
func (u *UserRepo) BeginTx(ctx context.Context) (pgx.Tx, error) {
	tx, err := u.db.Begin(ctx)
	if err != nil {
		u.log.Error("Error starting transaction:", err)
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	return tx, nil
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

func (u *UserRepo) CreateSessionTx(ctx context.Context, tx pgx.Tx, req *user_entity.CreateSession) error {
	id := uuid.New().String()
	query, args, err := u.queryBuilder.Insert("sessions").
		Columns("id", "user_id", "user_agent", "platform", "ip_address").
		Values(
			id,
			req.UserID,
			req.UserAgent,
			req.Platform,
			req.IPAddress,
		).ToSql()
	if err != nil {
		u.log.Error("Error generating SQL:", err)
		return err
	}

	_, err = tx.Exec(query, args...)
	if err != nil {
		u.log.Error("Error creating row:", err)
		return err
	}

	u.log.Info("Session created successfully, session_id: ", id)
	return nil
}

func (u *UserRepo) ActivateUserTx(ctx context.Context, tx pgx.Tx, user_id string) error {
	query, args, err := u.queryBuilder.Update("users").
		Set("status", entity.StatusActive).
		Where(sq.Eq{"id": user_id}).
		ToSql()
	if err != nil {
		u.log.Error("Error generating SQL:", err)
		return err
	}

	_, err = tx.Exec(query, args...)
	if err != nil {
		u.log.Error("Error updating row:", err)
		return err
	}

	u.log.Info("User activated successfully, user_id: ", user_id)
	return nil
}

func (u *UserRepo) CheckUserExists(ctx context.Context, email string) (bool, error) {
	query, args, err := u.queryBuilder.Select("id").
		From("users").
		Where(sq.Eq{"email": email}).
		ToSql()
	if err != nil {
		u.log.Error("Error generating SQL:", err)
		return false, err
	}

	var userID string
	err = u.db.QueryRow(ctx, query, args...).Scan(&userID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		u.log.Error("Error checking user existence:", err)
		return false, err
	}

	return true, nil
}

func (u *UserRepo) GetUserPassword(ctx context.Context, email string) (string, error) {
	query, args, err := u.queryBuilder.Select("password_hash").
		From("users").
		Where(sq.Eq{"email": email}).
		ToSql()
	if err != nil {
		u.log.Error("Error generating SQL:", err)
		return "", err
	}

	var passwordHash string
	err = u.db.QueryRow(ctx, query, args...).Scan(&passwordHash)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", entity.ErrUserNotFound
		}
	}

	return passwordHash, nil
}
