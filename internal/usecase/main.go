package usecase

import (
	"github.com/abdulazizax/yelp/config"
	"github.com/abdulazizax/yelp/internal/usecase/repo"
	"github.com/abdulazizax/yelp/pkg/logger"
	"github.com/abdulazizax/yelp/pkg/postgres"
)

// UseCase -.
type UseCase struct {
	UserRepo    UserRepoI
	SessionRepo SessionRepoI
}

// New -.
func New(pg *postgres.Postgres, config *config.Config, logger *logger.Logger) *UseCase {
	return &UseCase{
		UserRepo:    repo.NewUserRepo(pg, config, logger),
		SessionRepo: repo.NewSessionRepo(pg, config, logger),
	}
}
