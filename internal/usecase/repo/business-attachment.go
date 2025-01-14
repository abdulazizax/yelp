package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/abdulazizax/yelp/config"
	"github.com/abdulazizax/yelp/internal/entity"
	"github.com/abdulazizax/yelp/pkg/logger"
	"github.com/abdulazizax/yelp/pkg/postgres"
	"github.com/google/uuid"
)

type BusinessAttachmentRepo struct {
	pg     *postgres.Postgres
	config *config.Config
	logger *logger.Logger
}

// New -.
func NewBusinessAttachmentRepo(pg *postgres.Postgres, config *config.Config, logger *logger.Logger) *BusinessAttachmentRepo {
	return &BusinessAttachmentRepo{
		pg:     pg,
		config: config,
		logger: logger,
	}
}

func (r *BusinessAttachmentRepo) Create(ctx context.Context, req entity.BusinessAttachment) (entity.BusinessAttachment, error) {
	req.Id = uuid.NewString()

	qeury, args, err := r.pg.Builder.Insert("business_attachments").
		Columns(`id, business_id, filepath, content_type`).
		Values(req.Id, req.BusinessId, req.FilePath, req.ContentType).ToSql()
	if err != nil {
		return entity.BusinessAttachment{}, err
	}

	_, err = r.pg.Pool.Exec(ctx, qeury, args...)
	if err != nil {
		return entity.BusinessAttachment{}, err
	}

	return req, nil
}

func (r *BusinessAttachmentRepo) MultipleUpsert(ctx context.Context, req entity.BusinessAttachmentMultipleInsertRequest) ([]entity.BusinessAttachment, error) {
	hasNewAttachment := false

	tx, err := r.pg.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	insertQuery := r.pg.Builder.Insert("business_attachments").
		Columns(`id, business_id, filepath, content_type`)

	for i, attachment := range req.Attachments {
		if attachment.Id == "" {
			hasNewAttachment = true

			attachment.Id = uuid.NewString()
			req.Attachments[i].Id = attachment.Id
			insertQuery = insertQuery.Values(attachment.Id, req.BusinessId, attachment.FilePath, attachment.ContentType)
		}
	}

	existingAttachments := make(map[string]bool)
	for _, attachment := range req.Attachments {
		if attachment.Id != "" {
			existingAttachments[attachment.Id] = true
		}
	}

	if hasNewAttachment {
		query, args, err := insertQuery.ToSql()
		if err != nil {
			return nil, err
		}

		_, err = tx.Exec(ctx, query, args...)
		if err != nil {
			r.logger.Error("error while inserting business_attachments", err)
			return nil, err
		}
	}

	query, args, err := r.pg.Builder.Select("id").From("business_attachments").
		Where(squirrel.Eq{"business_id": req.BusinessId}).ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		r.logger.Error("error while getting ids of business_attachments", err)
		return nil, err
	}
	defer rows.Close()

	deletedAttachmentIds := []string{}

	for rows.Next() {
		var id string
		err = rows.Scan(&id)
		if err != nil {
			return nil, err
		}

		if !existingAttachments[id] {
			deletedAttachmentIds = append(deletedAttachmentIds, id)
		}
	}

	for _, e := range deletedAttachmentIds {
		query, args, err := r.pg.Builder.Delete("business_attachments").Where("id = ?", e).ToSql()
		if err != nil {
			return nil, err
		}

		_, err = tx.Exec(ctx, query, args...)
		if err != nil {
			r.logger.Error("error while deleting business_attachments", err)
			return nil, err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		r.logger.Error("error while commiting business_attachments", err)
		return nil, err
	}

	attachments, err := r.GetList(ctx, entity.GetListFilter{
		Page:  1,
		Limit: 10,
		Filters: []entity.Filter{
			{
				Column: "business_id",
				Type:   "eq",
				Value:  req.BusinessId,
			},
		},
	})
	if err != nil {
		r.logger.Error("error while getting business_attachments", err)
		return nil, err
	}

	return attachments.Items, nil
}

func (r *BusinessAttachmentRepo) GetSingle(ctx context.Context, req entity.Id) (entity.BusinessAttachment, error) {
	response := entity.BusinessAttachment{}
	var (
		createdAt, updatedAt time.Time
	)

	qeuryBuilder := r.pg.Builder.
		Select(`id, business_id, filepath, content_type, created_at, updated_at`).
		From("business_attachments")

	switch {
	case req.ID != "":
		qeuryBuilder = qeuryBuilder.Where("id = ?", req.ID)
	default:
		return entity.BusinessAttachment{}, fmt.Errorf("GetSingle - invalid request")
	}

	qeury, args, err := qeuryBuilder.ToSql()
	if err != nil {
		return entity.BusinessAttachment{}, err
	}

	err = r.pg.Pool.QueryRow(ctx, qeury, args...).
		Scan(&response.Id, &response.BusinessId, &response.FilePath, &response.ContentType, &createdAt, &updatedAt)
	if err != nil {
		return entity.BusinessAttachment{}, err
	}

	response.CreatedAt = createdAt.Format(time.RFC3339)
	response.UpdatedAt = updatedAt.Format(time.RFC3339)

	return response, nil
}

func (r *BusinessAttachmentRepo) GetList(ctx context.Context, req entity.GetListFilter) (entity.BusinessAttachmentList, error) {
	var (
		response             = entity.BusinessAttachmentList{}
		createdAt, updatedAt time.Time
	)

	qeuryBuilder := r.pg.Builder.
		Select(`id, business_id, filepath, content_type, created_at, updated_at`).
		From("business_attachments")

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
		var item entity.BusinessAttachment
		err = rows.Scan(&item.Id, &item.BusinessId, &item.FilePath, &item.ContentType, &createdAt, &updatedAt)
		if err != nil {
			return response, err
		}

		item.CreatedAt = createdAt.Format(time.RFC3339)
		item.UpdatedAt = updatedAt.Format(time.RFC3339)

		response.Items = append(response.Items, item)
	}

	countQuery, args, err := r.pg.Builder.Select("COUNT(1)").From("business_attachments").Where(where).ToSql()
	if err != nil {
		return response, err
	}

	err = r.pg.Pool.QueryRow(ctx, countQuery, args...).Scan(&response.Count)
	if err != nil {
		return response, err
	}

	return response, nil
}

func (r *BusinessAttachmentRepo) Delete(ctx context.Context, req entity.Id) error {
	qeury, args, err := r.pg.Builder.Delete("business_attachments").Where("id = ?", req.ID).ToSql()
	if err != nil {
		return err
	}

	_, err = r.pg.Pool.Exec(ctx, qeury, args...)
	if err != nil {
		return err
	}

	return nil
}
