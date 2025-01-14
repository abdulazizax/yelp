package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/abdulazizax/yelp/config"
	"github.com/abdulazizax/yelp/internal/entity"
	"github.com/abdulazizax/yelp/pkg/logger"
	"github.com/abdulazizax/yelp/pkg/postgres"
	"github.com/google/uuid"
)

type BusinessCategoryRepo struct {
	pg     *postgres.Postgres
	config *config.Config
	logger *logger.Logger
}

// New -.
func NewBusinessCategoryRepo(pg *postgres.Postgres, config *config.Config, logger *logger.Logger) *BusinessCategoryRepo {
	return &BusinessCategoryRepo{
		pg:     pg,
		config: config,
		logger: logger,
	}
}

func (r *BusinessCategoryRepo) Create(ctx context.Context, req entity.BusinessCategory) (entity.BusinessCategory, error) {
	req.ID = uuid.NewString()

	qeury, args, err := r.pg.Builder.Insert("business_categories").
		Columns(`id, name`).
		Values(req.ID, req.Name).ToSql()
	if err != nil {
		return entity.BusinessCategory{}, err
	}

	_, err = r.pg.Pool.Exec(ctx, qeury, args...)
	if err != nil {
		return entity.BusinessCategory{}, err
	}

	return req, nil
}

func (r *BusinessCategoryRepo) GetSingle(ctx context.Context, req entity.BusinessCategorySingleRequest) (entity.BusinessCategory, error) {
	response := entity.BusinessCategory{}
	var (
		createdAt, updatedAt time.Time
	)

	qeuryBuilder := r.pg.Builder.
		Select(`id, name, created_at, updated_at`).
		From("business_categories")

	switch {
	case req.ID != "":
		qeuryBuilder = qeuryBuilder.Where("id = ?", req.ID)
	case req.ID != "":
		qeuryBuilder = qeuryBuilder.Where("name = ?", req.Name)
	default:
		return entity.BusinessCategory{}, fmt.Errorf("GetSingle - invalid request")
	}

	qeury, args, err := qeuryBuilder.ToSql()
	if err != nil {
		return entity.BusinessCategory{}, err
	}

	err = r.pg.Pool.QueryRow(ctx, qeury, args...).
		Scan(&response.ID, &response.Name, &createdAt, &updatedAt)
	if err != nil {
		return entity.BusinessCategory{}, err
	}

	response.CreatedAt = createdAt.Format(time.RFC3339)
	response.UpdatedAt = updatedAt.Format(time.RFC3339)

	return response, nil
}

func (r *BusinessCategoryRepo) GetList(ctx context.Context, req entity.GetListFilter) (entity.BusinessCategoryList, error) {
	var (
		response             = entity.BusinessCategoryList{}
		createdAt, updatedAt time.Time
	)

	qeuryBuilder := r.pg.Builder.
		Select(`id, name, created_at, updated_at`).
		From("business_categories")

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
		var item entity.BusinessCategory
		err = rows.Scan(&item.ID, &item.Name, &createdAt, &updatedAt)
		if err != nil {
			return response, err
		}

		item.CreatedAt = createdAt.Format(time.RFC3339)
		item.UpdatedAt = updatedAt.Format(time.RFC3339)

		response.Items = append(response.Items, item)
	}

	countQuery, args, err := r.pg.Builder.Select("COUNT(1)").From("business_categories").Where(where).ToSql()
	if err != nil {
		return response, err
	}

	err = r.pg.Pool.QueryRow(ctx, countQuery, args...).Scan(&response.Count)
	if err != nil {
		return response, err
	}

	return response, nil
}

func (r *BusinessCategoryRepo) Update(ctx context.Context, req entity.BusinessCategory) (entity.BusinessCategory, error) {
	mp := map[string]interface{}{
		"name":               req.Name,
		"updated_at":         "now()",
	}

	qeury, args, err := r.pg.Builder.Update("business_categories").SetMap(mp).Where("id = ?", req.ID).ToSql()
	if err != nil {
		return entity.BusinessCategory{}, err
	}

	_, err = r.pg.Pool.Exec(ctx, qeury, args...)
	if err != nil {
		return entity.BusinessCategory{}, err
	}

	return req, nil
}

func (r *BusinessCategoryRepo) Delete(ctx context.Context, req entity.Id) error {
	qeury, args, err := r.pg.Builder.Delete("business_categories").Where("id = ?", req.ID).ToSql()
	if err != nil {
		return err
	}

	_, err = r.pg.Pool.Exec(ctx, qeury, args...)
	if err != nil {
		return err
	}

	return nil
}
