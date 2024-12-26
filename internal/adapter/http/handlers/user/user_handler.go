package user_handler

import (
	"net/http"

	"github.com/abdulazizax/yelp/internal/entity"
	user_entity "github.com/abdulazizax/yelp/internal/entity/user"
	user_usecase "github.com/abdulazizax/yelp/internal/usecase/user"
	"github.com/abdulazizax/yelp/pkg/logger"
	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	SignUp(c *gin.Context)
}

type userHandler struct {
	usecase user_usecase.UserUsecase
	log     logger.Interface
}

func NewUserHandler(usecase user_usecase.UserUsecase, log logger.Interface) UserHandler {
	return &userHandler{
		usecase: usecase,
		log:     log,
	}
}

// SignUpUser creates a new user
// @Summary Create a new user
// @Description Create a new user and return the created user's ID
// @Tags users
// @Accept  json
// @Produce  json
// @Param user body user_entity.CreateUser true "user information"
// @Success 201 {object} entity.Info "message"
// @Failure 400 {object} entity.Error "Bad Request"
// @Failure 500 {object} entity.Error "Internal Server Error"
// @Router /sign-up [post]
func (u *userHandler) SignUp(c *gin.Context) {
	var user user_entity.CreateUser
	if err := c.ShouldBindJSON(&user); err != nil {
		u.log.Error("failed to bind JSON", "error", err)
		c.JSON(http.StatusBadRequest, entity.Error{Message: "Invalid user data"})
		return
	}

	err := u.usecase.CreateUser(c, &user)
	if err != nil {
		u.log.Error("failed to create user", "error", err)
		c.JSON(http.StatusInternalServerError, entity.Error{Message: "Failed to create user"})
		return
	}

	u.log.Info("User created successfully")
	c.JSON(http.StatusCreated, entity.Info{Message: "User created successfully"})
}
