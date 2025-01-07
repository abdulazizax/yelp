package repo

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/abdulazizax/yelp/config"
	"github.com/abdulazizax/yelp/internal/entity"
	"github.com/abdulazizax/yelp/pkg/logger"
	"github.com/abdulazizax/yelp/pkg/postgres"
	"github.com/google/uuid"
)

type UserRepo struct {
	pg     *postgres.Postgres
	config *config.Config
	logger *logger.Logger
}

// New -.
func NewUserRepo(pg *postgres.Postgres, config *config.Config, logger *logger.Logger) *UserRepo {
	return &UserRepo{
		pg:     pg,
		config: config,
		logger: logger,
	}
}

func (r *UserRepo) Create(ctx context.Context, req entity.User) (entity.User, error) {
	req.ID = uuid.NewString()

	qeury, args, err := r.pg.Builder.Insert("users").
		Columns(`id, user_type, user_role, full_name, username, email, password, gender, status`).
		Values(req.ID, req.UserType, req.UserRole, req.FullName, req.Username, req.Email, req.Password, req.Gender, req.Status).ToSql()
	if err != nil {
		return entity.User{}, err
	}

	_, err = r.pg.Pool.Exec(ctx, qeury, args...)
	if err != nil {
		return entity.User{}, err
	}

	return req, nil
}

func (r *UserRepo) GetSingle(ctx context.Context, req entity.UserSingleRequest) (entity.User, error) {
	response := entity.User{}
	var (
		createdAt, updatedAt time.Time
		bio, profile_picture sql.NullString
	)

	qeuryBuilder := r.pg.Builder.
		Select(`id, user_type, user_role, full_name, username, email, password, bio, gender, profile_picture, status, created_at, updated_at`).
		From("users")

	switch {
	case req.ID != "":
		qeuryBuilder = qeuryBuilder.Where("id = ?", req.ID)
	case req.Email != "":
		qeuryBuilder = qeuryBuilder.Where("email = ?", req.Email)
	case req.UserName != "":
		qeuryBuilder = qeuryBuilder.Where("username = ?", req.UserName)
	default:
		return entity.User{}, fmt.Errorf("GetSingle - invalid request")
	}

	qeury, args, err := qeuryBuilder.ToSql()
	if err != nil {
		return entity.User{}, err
	}

	err = r.pg.Pool.QueryRow(ctx, qeury, args...).
		Scan(&response.ID, &response.UserType, &response.UserRole, &response.FullName, &response.Username,
			&response.Email, &response.Password, &bio, &response.Gender, &profile_picture, &response.Status, &createdAt, &updatedAt)
	if err != nil {
		return entity.User{}, err
	}

	response.CreatedAt = createdAt.Format(time.RFC3339)
	response.UpdatedAt = updatedAt.Format(time.RFC3339)
	if bio.Valid {
		response.Bio = bio.String
	}
	if profile_picture.Valid {
		response.ProfilePicture = profile_picture.String
	}

	return response, nil
}

func (r *UserRepo) GetList(ctx context.Context, req entity.GetListFilter) (entity.UserList, error) {
	var (
		response             = entity.UserList{}
		createdAt, updatedAt time.Time
		bio, profile_picture sql.NullString
	)

	qeuryBuilder := r.pg.Builder.
		Select(`id, user_type, user_role, full_name, username, email, password, bio, gender, profile_picture, status, created_at, updated_at`).
		From("users")

	qeuryBuilder, where := PrepareGetListQuery(qeuryBuilder, req)

	qeury, args, err := qeuryBuilder.ToSql()
	if err != nil {
		return response, err
	}

	rows, err := r.pg.Pool.Query(ctx, qeury, args...)
	if err != nil {
		return response, err
	}
	defer rows.Close()

	for rows.Next() {
		var item entity.User
		err = rows.Scan(&item.ID, &item.UserType, &item.UserRole, &item.FullName, &item.Username,
			&item.Email, &item.Password, &bio, &item.Gender, &profile_picture, &item.Status, &createdAt, &updatedAt)
		if err != nil {
			return response, err
		}

		item.CreatedAt = createdAt.Format(time.RFC3339)
		item.UpdatedAt = updatedAt.Format(time.RFC3339)
		if bio.Valid {
			item.Bio = bio.String
		}
		if profile_picture.Valid {
			item.ProfilePicture = profile_picture.String
		}

		response.Items = append(response.Items, item)
	}

	countQuery, args, err := r.pg.Builder.Select("COUNT(1)").From("users").Where(where).ToSql()
	if err != nil {
		return response, err
	}

	err = r.pg.Pool.QueryRow(ctx, countQuery, args...).Scan(&response.Count)
	if err != nil {
		return response, err
	}

	return response, nil
}

func (r *UserRepo) Update(ctx context.Context, req entity.User) (entity.User, error) {
	mp := map[string]interface{}{
		"full_name":       req.FullName,
		"username":        req.Username,
		"status":          req.Status,
		"email":           req.Email,
		"bio":             req.Bio,
		"gender":          req.Gender,
		"user_role":       req.UserRole,
		"profile_picture": req.ProfilePicture,
		"updated_at":      "now()",
	}

	if req.Password != "" {
		mp["password"] = req.Password
	}

	qeury, args, err := r.pg.Builder.Update("users").SetMap(mp).Where("id = ?", req.ID).ToSql()
	if err != nil {
		return entity.User{}, err
	}

	_, err = r.pg.Pool.Exec(ctx, qeury, args...)
	if err != nil {
		return entity.User{}, err
	}

	return req, nil
}

func (r *UserRepo) Delete(ctx context.Context, req entity.Id) error {
	qeury, args, err := r.pg.Builder.Delete("users").Where("id = ?", req.ID).ToSql()
	if err != nil {
		return err
	}

	_, err = r.pg.Pool.Exec(ctx, qeury, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepo) UpdateField(ctx context.Context, req entity.UpdateFieldRequest) (entity.RowsEffected, error) {
	mp := map[string]interface{}{}
	response := entity.RowsEffected{}

	for _, item := range req.Items {
		mp[item.Column] = item.Value
	}

	qeury, args, err := r.pg.Builder.Update("users").SetMap(mp).Where(PrepareFilter(req.Filter)).ToSql()
	if err != nil {
		return response, err
	}

	n, err := r.pg.Pool.Exec(ctx, qeury, args...)
	if err != nil {
		return response, err
	}

	response.RowsEffected = int(n.RowsAffected())

	return response, nil
}
