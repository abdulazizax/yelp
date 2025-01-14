package usecase

import (
	"github.com/abdulazizax/yelp/config"
	"github.com/abdulazizax/yelp/internal/usecase/repo"
	"github.com/abdulazizax/yelp/pkg/logger"
	"github.com/abdulazizax/yelp/pkg/postgres"
)

// UseCase -.
type UseCase struct {
	UserRepo               UserRepoI
	SessionRepo            SessionRepoI
	BusinessRepo           BusinessRepoI
	BusinessCategoryRepo   BusinessCategoryRepoI
	BusinessAttachmentRepo BusinessAttachmentRepoI
	ReviewRepo             ReviewRepoI
	ReviewAttachmentRepo   ReviewAttachmentRepoI
}

// New -.
func New(pg *postgres.Postgres, config *config.Config, logger *logger.Logger) *UseCase {
	return &UseCase{
		UserRepo:               repo.NewUserRepo(pg, config, logger),
		SessionRepo:            repo.NewSessionRepo(pg, config, logger),
		BusinessRepo:           repo.NewBusinessRepo(pg, config, logger),
		BusinessCategoryRepo:   repo.NewBusinessCategoryRepo(pg, config, logger),
		BusinessAttachmentRepo: repo.NewBusinessAttachmentRepo(pg, config, logger),
		ReviewRepo:             repo.NewReviewRepo(pg, config, logger),
		ReviewAttachmentRepo:   repo.NewReviewAttachmentRepo(pg, config, logger),
	}
}
