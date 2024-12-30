package user_db

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"strconv"

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
	email        security.EmailService
	queryBuilder sq.StatementBuilderType
	log          logger.Interface
}

func NewUserRepo(db *pgxpool.Pool, queryBuilder sq.StatementBuilderType, email security.EmailService, log logger.Interface) *UserRepo {
	return &UserRepo{
		db:           db,
		email:        email,
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

func (u *UserRepo) SignIn(ctx context.Context, user_id string, signInRequest *user_entity.SignInRequest) error {
	// Begin a transaction
	tx, err := u.db.Begin(ctx)
	if err != nil {
		u.log.Error("Error starting transaction:", err)
		return err
	}
	defer func() {
		// If an error occurs during the transaction, rollback the transaction
		if err != nil {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				u.log.Error("Error rolling back transaction:", rollbackErr)
			}
		}
	}()

	// update 'users' table
	query, args, err := u.queryBuilder.Update("users").
		Set("status", entity.StatusActive).
		Where(sq.Eq{"email": signInRequest.User.Email}).
		ToSql()
	if err != nil {
		u.log.Error("Error generating SQL:", err)
		return err
	}

	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		u.log.Error("Error updating row:", err)
		return err
	}

	// update `sessions` table
	id := uuid.New().String()
	query, args, err = u.queryBuilder.Insert("sessions").
		Columns("id", "user_id", "user_agent", "platform", "ip_address").
		Values(
			id,
			user_id,
			signInRequest.Session.UserAgent,
			signInRequest.Session.Platform,
			signInRequest.Session.IPAddress,
		).ToSql()
	if err != nil {
		u.log.Error("Error generating SQL:", err)
		return err
	}

	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		u.log.Error("Error creating row:", err)
		return err
	}

	// Commit the transaction
	if commitErr := tx.Commit(ctx); commitErr != nil {
		u.log.Error("Error committing transaction:", commitErr)
		return commitErr
	}

	return nil
}

func (u *UserRepo) CheckUserExists(ctx context.Context, email string) error {
	query, args, err := u.queryBuilder.Select("id").
		From("users").
		Where(sq.Eq{"email": email}).
		ToSql()
	if err != nil {
		u.log.Error("Error generating SQL query", "error", err)
		return fmt.Errorf("failed to generate SQL query: %w", err)
	}

	var userID string
	err = u.db.QueryRow(ctx, query, args...).Scan(&userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) || err.Error() == "no rows in result set" {
			log.Println("User not exists")
			u.log.Info("User not exists", "email", email)
			return nil
		}
		u.log.Error("Error checking user existence", "error", err)
		return fmt.Errorf("failed to check user existence: %w", err)
	}

	u.log.Info("User exists", "email", email)
	return entity.ErrUserAlreadyExists
}

func (u *UserRepo) GetUserInfoByEmail(ctx context.Context, email string) (*user_entity.UserInfo, error) {
	query, args, err := u.queryBuilder.Select("id", "user_type", "user_role", "password_hash").
		From("users").
		Where(sq.Eq{"email": email}).
		ToSql()
	if err != nil {
		u.log.Error("Error generating SQL", "error", err)
		return nil, err
	}

	var userInfo user_entity.UserInfo

	err = u.db.QueryRow(ctx, query, args...).
		Scan(&userInfo.ID, &userInfo.UserType, &userInfo.UserRole, &userInfo.PasswordHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			u.log.Info("User not exists", "email", email)
			return nil, entity.ErrUserNotFound
		}
		u.log.Error("Error fetching user info", "error", err)
		return nil, fmt.Errorf("failed to fetch user info: %w", err)
	}

	return &userInfo, nil
}

func (u *UserRepo) UpdateUserPassword(ctx context.Context, req *user_entity.UpdateUserPassword) error {
	intVerificationCode, err := strconv.Atoi(req.VerificationCode)
	if err != nil {
		u.log.Error("Error converting verification code to integer", err)
		return err
	}

	err = u.email.VerifyEmail(context.Background(), req.Email, intVerificationCode)
	if err != nil {
		u.log.Error("Failed to verify email", slog.String("error", err.Error()))
		return fmt.Errorf("failed to verify email: %w", err)
	}

	password_hash, err := security.HashPassword(req.NewPassword)
	if err != nil {
		u.log.Error("Failed to hash password", slog.String("error", err.Error()))
		return fmt.Errorf("failed to hash password: %w", err)
	}

	query, args, err := u.queryBuilder.Update("users").
		Set("password_hash", password_hash).
		Where(sq.Eq{"email": req.Email}).
		ToSql()
	if err != nil {
		u.log.Error("Failed to build query", slog.String("error", err.Error()))
		return fmt.Errorf("failed to build query: %w", err)
	}

	_, err = u.db.Exec(ctx, query, args...)
	if err != nil {
		u.log.Error("Failed to execute query", slog.String("error", err.Error()))
		return fmt.Errorf("failed to execute query: %w", err)
	}

	return nil
}
