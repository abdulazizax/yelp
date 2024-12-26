package handler

import (
	user_handler "github.com/abdulazizax/yelp/internal/adapter/http/handlers/user"
	"github.com/abdulazizax/yelp/internal/usecase"
	"github.com/abdulazizax/yelp/pkg/logger"
)

type Handlers interface {
	UserHandler() user_handler.UserHandler
}

type handlers struct {
	user user_handler.UserHandler
}

func NewHandlers(usecase usecase.Usecase, log logger.Interface) Handlers {
	return &handlers{
		user: user_handler.NewUserHandler(usecase.UserUsecase(), log),
	}
}

func (h *handlers) UserHandler() user_handler.UserHandler {
	return h.user
}
