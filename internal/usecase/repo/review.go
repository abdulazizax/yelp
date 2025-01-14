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

type ReviewRepo struct {
	pg     *postgres.Postgres
	config *config.Config
	logger *logger.Logger
}

// New -.
func NewReviewRepo(pg *postgres.Postgres, config *config.Config, logger *logger.Logger) *ReviewRepo {
	return &ReviewRepo{
		pg:     pg,
		config: config,
		logger: logger,
	}
}

func (r *ReviewRepo) Create(ctx context.Context, req entity.Review) (entity.Review, error) {
	req.ID = uuid.NewString()

	qeury, args, err := r.pg.Builder.Insert("reviews").
		Columns(`id, business_id, user_id, rating, comment`).
		Values(req.ID, req.BusinessID, req.UserID, req.Rating, req.Comment).ToSql()
	if err != nil {
		return entity.Review{}, err
	}

	_, err = r.pg.Pool.Exec(ctx, qeury, args...)
	if err != nil {
		return entity.Review{}, err
	}

	return req, nil
}

func (r *ReviewRepo) GetSingle(ctx context.Context, req entity.Id) (entity.Review, error) {
	response := entity.Review{}
	var (
		createdAt, updatedAt time.Time
		comment              sql.NullString
	)

	qeuryBuilder := r.pg.Builder.
		Select(`id, business_id, user_id, rating, comment, created_at, updated_at`).
		From("reviews")

	switch {
	case req.ID != "":
		qeuryBuilder = qeuryBuilder.Where("id = ?", req.ID)
	default:
		return entity.Review{}, fmt.Errorf("GetSingle - invalid request")
	}

	qeury, args, err := qeuryBuilder.ToSql()
	if err != nil {
		return entity.Review{}, err
	}

	err = r.pg.Pool.QueryRow(ctx, qeury, args...).
		Scan(&response.ID, &response.BusinessID, &response.UserID, &response.Rating,
			&comment, &createdAt, &updatedAt)
	if err != nil {
		return entity.Review{}, err
	}

	response.CreatedAt = createdAt.Format(time.RFC3339)
	response.UpdatedAt = updatedAt.Format(time.RFC3339)
	if comment.Valid {
		response.Comment = comment.String
	}

	return response, nil
}

func (r *ReviewRepo) GetList(ctx context.Context, req entity.GetListFilter) (entity.ReviewList, error) {
	var (
		response             = entity.ReviewList{}
		createdAt, updatedAt time.Time
		comment              sql.NullString
	)

	qeuryBuilder := r.pg.Builder.
		Select(`id, business_id, user_id, rating, comment, created_at, updated_at`).
		From("reviews")

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
		var item entity.Review
		err = rows.Scan(&item.ID, &item.BusinessID, &item.UserID, &item.Rating,
			&comment, &createdAt, &updatedAt)
		if err != nil {
			return response, err
		}

		item.CreatedAt = createdAt.Format(time.RFC3339)
		item.UpdatedAt = updatedAt.Format(time.RFC3339)
		if comment.Valid {
			item.Comment = comment.String
		}

		response.Items = append(response.Items, item)
	}

	countQuery, args, err := r.pg.Builder.Select("COUNT(1)").From("reviews").Where(where).ToSql()
	if err != nil {
		return response, err
	}

	err = r.pg.Pool.QueryRow(ctx, countQuery, args...).Scan(&response.Count)
	if err != nil {
		return response, err
	}

	return response, nil
}

func (r *ReviewRepo) Update(ctx context.Context, req entity.Review) (entity.Review, error) {
	mp := map[string]interface{}{
		"rating":     req.Rating,
		"comment":    req.Comment,
		"updated_at": "now()",
	}

	qeury, args, err := r.pg.Builder.Update("reviews").SetMap(mp).Where("id = ?", req.ID).ToSql()
	if err != nil {
		return entity.Review{}, err
	}

	_, err = r.pg.Pool.Exec(ctx, qeury, args...)
	if err != nil {
		return entity.Review{}, err
	}

	return req, nil
}

func (r *ReviewRepo) Delete(ctx context.Context, req entity.Id) error {
	qeury, args, err := r.pg.Builder.Delete("reviews").Where("id = ?", req.ID).ToSql()
	if err != nil {
		return err
	}

	_, err = r.pg.Pool.Exec(ctx, qeury, args...)
	if err != nil {
		return err
	}

	return nil
}
