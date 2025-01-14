// Package usecase implements application business logic. Each logic group in own file.
package usecase

import (
	"context"

	"github.com/abdulazizax/yelp/internal/entity"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks_test.go -package=usecase_test

type (
	// UserRepo -.
	UserRepoI interface {
		Create(ctx context.Context, req entity.User) (entity.User, error)
		GetSingle(ctx context.Context, req entity.UserSingleRequest) (entity.User, error)
		GetList(ctx context.Context, req entity.GetListFilter) (entity.UserList, error)
		Update(ctx context.Context, req entity.User) (entity.User, error)
		Delete(ctx context.Context, req entity.Id) error
		UpdateField(ctx context.Context, req entity.UpdateFieldRequest) (entity.RowsEffected, error)
	}

	// SessionRepo -.
	SessionRepoI interface {
		Create(ctx context.Context, req entity.Session) (entity.Session, error)
		GetSingle(ctx context.Context, req entity.Id) (entity.Session, error)
		GetList(ctx context.Context, req entity.GetListFilter) (entity.SessionList, error)
		Update(ctx context.Context, req entity.Session) (entity.Session, error)
		Delete(ctx context.Context, req entity.Id) error
		UpdateField(ctx context.Context, req entity.UpdateFieldRequest) (entity.RowsEffected, error)
	}

	// BusinessRepo
	BusinessRepoI interface {
		Create(ctx context.Context, req entity.Business) (entity.Business, error)
		GetSingle(ctx context.Context, req entity.BusinessSingleRequest) (entity.Business, error)
		GetList(ctx context.Context, req entity.GetListFilter) (entity.BusinessList, error)
		Update(ctx context.Context, req entity.Business) (entity.Business, error)
		Delete(ctx context.Context, req entity.Id) error
		UpdateField(ctx context.Context, req entity.UpdateFieldRequest) (entity.RowsEffected, error)
	}

	// BusinessCategoryRepo
	BusinessCategoryRepoI interface {
		Create(ctx context.Context, req entity.BusinessCategory) (entity.BusinessCategory, error)
		GetSingle(ctx context.Context, req entity.BusinessCategorySingleRequest) (entity.BusinessCategory, error)
		GetList(ctx context.Context, req entity.GetListFilter) (entity.BusinessCategoryList, error)
		Update(ctx context.Context, req entity.BusinessCategory) (entity.BusinessCategory, error)
		Delete(ctx context.Context, req entity.Id) error
	}

	// BusinessAttachmentRepo
	BusinessAttachmentRepoI interface {
		Create(ctx context.Context, req entity.BusinessAttachment) (entity.BusinessAttachment, error)
		MultipleUpsert(ctx context.Context, req entity.BusinessAttachmentMultipleInsertRequest) ([]entity.BusinessAttachment, error)
		GetSingle(ctx context.Context, req entity.Id) (entity.BusinessAttachment, error)
		GetList(ctx context.Context, req entity.GetListFilter) (entity.BusinessAttachmentList, error)
		Delete(ctx context.Context, req entity.Id) error
	}

	// ReviewRepo
	ReviewRepoI interface {
		Create(ctx context.Context, req entity.Review) (entity.Review, error)
		GetSingle(ctx context.Context, req entity.Id) (entity.Review, error)
		GetList(ctx context.Context, req entity.GetListFilter) (entity.ReviewList, error)
		Update(ctx context.Context, req entity.Review) (entity.Review, error)
		Delete(ctx context.Context, req entity.Id) error
	}

	// ReviewAttachmentRepo
	ReviewAttachmentRepoI interface {
		Create(ctx context.Context, req entity.ReviewAttachment) (entity.ReviewAttachment, error)
		MultipleUpsert(ctx context.Context, req entity.ReviewAttachmentMultipleInsertRequest) ([]entity.ReviewAttachment, error)
		GetSingle(ctx context.Context, req entity.Id) (entity.ReviewAttachment, error)
		GetList(ctx context.Context, req entity.GetListFilter) (entity.ReviewAttachmentList, error)
		Delete(ctx context.Context, req entity.Id) error
	}
)
