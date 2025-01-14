package repo

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/abdulazizax/yelp/config"
	"github.com/abdulazizax/yelp/internal/entity"
	"github.com/abdulazizax/yelp/pkg/logger"
	"github.com/abdulazizax/yelp/pkg/postgres"
	"github.com/google/uuid"
)

type BusinessRepo struct {
	pg     *postgres.Postgres
	config *config.Config
	logger *logger.Logger
}

// New -.
func NewBusinessRepo(pg *postgres.Postgres, config *config.Config, logger *logger.Logger) *BusinessRepo {
	return &BusinessRepo{
		pg:     pg,
		config: config,
		logger: logger,
	}
}

func (r *BusinessRepo) Create(ctx context.Context, req entity.Business) (entity.Business, error) {
	req.ID = uuid.NewString()

	qeury, args, err := r.pg.Builder.Insert("businesses").
		Columns(`id, name, category_id, address, owner_id`).
		Values(req.ID, req.Name, req.CategoryID, req.Address, req.OwnerID).ToSql()
	if err != nil {
		return entity.Business{}, err
	}

	_, err = r.pg.Pool.Exec(ctx, qeury, args...)
	if err != nil {
		return entity.Business{}, err
	}

	return req, nil
}

func (r *BusinessRepo) GetSingle(ctx context.Context, req entity.BusinessSingleRequest) (entity.Business, error) {
	response := entity.Business{}
	var (
		createdAt, updatedAt                       time.Time
		description, contactInfo, hoursOfOperation sql.NullString
		latitude, longitude                        sql.NullFloat64
	)

	qeuryBuilder := r.pg.Builder.
		Select(`id, name, description, category_id, address, latitude, longitude, contact_info, hours_of_operation, owner_id, created_at, updated_at`).
		From("businesses")

	switch {
	case req.ID != "":
		qeuryBuilder = qeuryBuilder.Where("id = ?", req.ID)
	case req.OwnerID != "":
		qeuryBuilder = qeuryBuilder.Where("owner_id = ?", req.OwnerID)
	case req.CategoryID != "":
		qeuryBuilder = qeuryBuilder.Where("category_id = ?", req.CategoryID)
	default:
		return entity.Business{}, fmt.Errorf("GetSingle - invalid request")
	}

	qeury, args, err := qeuryBuilder.ToSql()
	if err != nil {
		return entity.Business{}, err
	}

	err = r.pg.Pool.QueryRow(ctx, qeury, args...).
		Scan(&response.ID, &response.Name, &description, &response.CategoryID, &response.Address,
			&latitude, &longitude, &contactInfo, &hoursOfOperation, &response.OwnerID, &createdAt, &updatedAt)
	if err != nil {
		return entity.Business{}, err
	}

	response.CreatedAt = createdAt.Format(time.RFC3339)
	response.UpdatedAt = updatedAt.Format(time.RFC3339)
	if latitude.Valid {
		response.Latitude = latitude.Float64
	}
	if longitude.Valid {
		response.Longitude = longitude.Float64
	}
	if contactInfo.Valid {
		// JSON to struct (ContactInfo)
		var contactInfoStruct entity.ContactInfo
		err := json.Unmarshal([]byte(contactInfo.String), &contactInfoStruct)
		if err != nil {
			return response, err
		}
		response.ContactInfo = contactInfoStruct
	}
	if hoursOfOperation.Valid {
		// JSON to struct (HoursOfOperation)
		var hoursOfOperationStruct entity.HoursOfOperation
		err := json.Unmarshal([]byte(hoursOfOperation.String), &hoursOfOperationStruct)
		if err != nil {
			return response, err
		}
		response.HoursOfOperation = hoursOfOperationStruct
	}
	if description.Valid {
		response.Description = description.String
	}

	return response, nil
}

func (r *BusinessRepo) GetList(ctx context.Context, req entity.GetListFilter) (entity.BusinessList, error) {
	var (
		response                                   = entity.BusinessList{}
		createdAt, updatedAt                       time.Time
		description, contactInfo, hoursOfOperation sql.NullString
		latitude, longitude                        sql.NullFloat64
	)

	qeuryBuilder := r.pg.Builder.
		Select(`id, name, description, category_id, address, latitude, longitude, contact_info, hours_of_operation, owner_id, created_at, updated_at`).
		From("businesses")

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
		var item entity.Business
		err = rows.Scan(&item.ID, &item.Name, &description, &item.CategoryID, &item.Address,
			&latitude, &longitude, &contactInfo, &hoursOfOperation, &item.OwnerID, &createdAt, &updatedAt)
		if err != nil {
			return response, err
		}

		item.CreatedAt = createdAt.Format(time.RFC3339)
		item.UpdatedAt = updatedAt.Format(time.RFC3339)
		if latitude.Valid {
			item.Latitude = latitude.Float64
		}
		if longitude.Valid {
			item.Longitude = longitude.Float64
		}
		if contactInfo.Valid {
			// JSON to struct (ContactInfo)
			var contactInfoStruct entity.ContactInfo
			err := json.Unmarshal([]byte(contactInfo.String), &contactInfoStruct)
			if err != nil {
				return response, err
			}
			item.ContactInfo = contactInfoStruct
		}
		if hoursOfOperation.Valid {
			// JSON to struct (HoursOfOperation)
			var hoursOfOperationStruct entity.HoursOfOperation
			err := json.Unmarshal([]byte(hoursOfOperation.String), &hoursOfOperationStruct)
			if err != nil {
				return response, err
			}
			item.HoursOfOperation = hoursOfOperationStruct
		}
		if description.Valid {
			item.Description = description.String
		}

		response.Items = append(response.Items, item)
	}

	countQuery, args, err := r.pg.Builder.Select("COUNT(1)").From("businesses").Where(where).ToSql()
	if err != nil {
		return response, err
	}

	err = r.pg.Pool.QueryRow(ctx, countQuery, args...).Scan(&response.Count)
	if err != nil {
		return response, err
	}

	return response, nil
}

func (r *BusinessRepo) Update(ctx context.Context, req entity.Business) (entity.Business, error) {
	mp := map[string]interface{}{
		"name":               req.Name,
		"description":        req.Description,
		"category_id":        req.CategoryID,
		"address":            req.Address,
		"latitude":           req.Latitude,
		"longitude":          req.Longitude,
		"contact_info":       req.ContactInfo,
		"hours_of_operation": req.HoursOfOperation,
		"owner_id":           req.OwnerID,
		"updated_at":         "now()",
	}

	qeury, args, err := r.pg.Builder.Update("businesses").SetMap(mp).Where("id = ?", req.ID).ToSql()
	if err != nil {
		return entity.Business{}, err
	}

	_, err = r.pg.Pool.Exec(ctx, qeury, args...)
	if err != nil {
		return entity.Business{}, err
	}

	return req, nil
}

func (r *BusinessRepo) Delete(ctx context.Context, req entity.Id) error {
	qeury, args, err := r.pg.Builder.Delete("businesses").Where("id = ?", req.ID).ToSql()
	if err != nil {
		return err
	}

	_, err = r.pg.Pool.Exec(ctx, qeury, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *BusinessRepo) UpdateField(ctx context.Context, req entity.UpdateFieldRequest) (entity.RowsEffected, error) {
	mp := map[string]interface{}{}
	response := entity.RowsEffected{}

	for _, item := range req.Items {
		mp[item.Column] = item.Value
	}

	qeury, args, err := r.pg.Builder.Update("businesses").SetMap(mp).Where(PrepareFilter(req.Filter)).ToSql()
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
