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

type ReviewAttachmentRepo struct {
	pg     *postgres.Postgres
	config *config.Config
	logger *logger.Logger
}

// New -.
func NewReviewAttachmentRepo(pg *postgres.Postgres, config *config.Config, logger *logger.Logger) *ReviewAttachmentRepo {
	return &ReviewAttachmentRepo{
		pg:     pg,
		config: config,
		logger: logger,
	}
}

func (r *ReviewAttachmentRepo) Create(ctx context.Context, req entity.ReviewAttachment) (entity.ReviewAttachment, error) {
	req.Id = uuid.NewString()

	qeury, args, err := r.pg.Builder.Insert("review_attachments").
		Columns(`id, review_id, filepath, content_type`).
		Values(req.Id, req.ReviewId, req.FilePath, req.ContentType).ToSql()
	if err != nil {
		return entity.ReviewAttachment{}, err
	}

	_, err = r.pg.Pool.Exec(ctx, qeury, args...)
	if err != nil {
		return entity.ReviewAttachment{}, err
	}

	return req, nil
}

func (r *ReviewAttachmentRepo) MultipleUpsert(ctx context.Context, req entity.ReviewAttachmentMultipleInsertRequest) ([]entity.ReviewAttachment, error) {
	hasNewAttachment := false

	tx, err := r.pg.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	insertQuery := r.pg.Builder.Insert("review_attachments").
		Columns(`id, review_id, filepath, content_type`)

	for i, attachment := range req.Attachments {
		if attachment.Id == "" {
			hasNewAttachment = true

			attachment.Id = uuid.NewString()
			req.Attachments[i].Id = attachment.Id
			insertQuery = insertQuery.Values(attachment.Id, req.ReviewId, attachment.FilePath, attachment.ContentType)
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
			r.logger.Error("error while inserting review_attachments", err)
			return nil, err
		}
	}

	query, args, err := r.pg.Builder.Select("id").From("review_attachments").
		Where(squirrel.Eq{"review_id": req.ReviewId}).ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		r.logger.Error("error while getting ids of review_attachments", err)
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
		query, args, err := r.pg.Builder.Delete("review_attachments").Where("id = ?", e).ToSql()
		if err != nil {
			return nil, err
		}

		_, err = tx.Exec(ctx, query, args...)
		if err != nil {
			r.logger.Error("error while deleting review_attachments", err)
			return nil, err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		r.logger.Error("error while commiting review_attachments", err)
		return nil, err
	}

	attachments, err := r.GetList(ctx, entity.GetListFilter{
		Page:  1,
		Limit: 10,
		Filters: []entity.Filter{
			{
				Column: "review_id",
				Type:   "eq",
				Value:  req.ReviewId,
			},
		},
	})
	if err != nil {
		r.logger.Error("error while getting review_attachments", err)
		return nil, err
	}

	return attachments.Items, nil
}

func (r *ReviewAttachmentRepo) GetSingle(ctx context.Context, req entity.Id) (entity.ReviewAttachment, error) {
	response := entity.ReviewAttachment{}
	var (
		createdAt, updatedAt time.Time
	)

	qeuryBuilder := r.pg.Builder.
		Select(`id, review_id, filepath, content_type, created_at, updated_at`).
		From("review_attachments")

	switch {
	case req.ID != "":
		qeuryBuilder = qeuryBuilder.Where("id = ?", req.ID)
	default:
		return entity.ReviewAttachment{}, fmt.Errorf("GetSingle - invalid request")
	}

	qeury, args, err := qeuryBuilder.ToSql()
	if err != nil {
		return entity.ReviewAttachment{}, err
	}

	err = r.pg.Pool.QueryRow(ctx, qeury, args...).
		Scan(&response.Id, &response.ReviewId, &response.FilePath, &response.ContentType, &createdAt, &updatedAt)
	if err != nil {
		return entity.ReviewAttachment{}, err
	}

	response.CreatedAt = createdAt.Format(time.RFC3339)
	response.UpdatedAt = updatedAt.Format(time.RFC3339)

	return response, nil
}

func (r *ReviewAttachmentRepo) GetList(ctx context.Context, req entity.GetListFilter) (entity.ReviewAttachmentList, error) {
	var (
		response             = entity.ReviewAttachmentList{}
		createdAt, updatedAt time.Time
	)

	qeuryBuilder := r.pg.Builder.
		Select(`id, review_id, filepath, content_type, created_at, updated_at`).
		From("review_attachments")

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
		var item entity.ReviewAttachment
		err = rows.Scan(&item.Id, &item.ReviewId, &item.FilePath, &item.ContentType, &createdAt, &updatedAt)
		if err != nil {
			return response, err
		}

		item.CreatedAt = createdAt.Format(time.RFC3339)
		item.UpdatedAt = updatedAt.Format(time.RFC3339)

		response.Items = append(response.Items, item)
	}

	countQuery, args, err := r.pg.Builder.Select("COUNT(1)").From("review_attachments").Where(where).ToSql()
	if err != nil {
		return response, err
	}

	err = r.pg.Pool.QueryRow(ctx, countQuery, args...).Scan(&response.Count)
	if err != nil {
		return response, err
	}

	return response, nil
}

func (r *ReviewAttachmentRepo) Delete(ctx context.Context, req entity.Id) error {
	qeury, args, err := r.pg.Builder.Delete("review_attachments").Where("id = ?", req.ID).ToSql()
	if err != nil {
		return err
	}

	_, err = r.pg.Pool.Exec(ctx, qeury, args...)
	if err != nil {
		return err
	}

	return nil
}
