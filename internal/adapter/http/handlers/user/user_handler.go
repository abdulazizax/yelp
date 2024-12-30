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
	SignIn(c *gin.Context)
	SendVerificationCode(c *gin.Context)
	UpdateUserPassword(c *gin.Context)
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

// SignUpUser godoc
// @Summary      Create a new user account
// @Description  Registers a new user with the provided details and returns a confirmation message
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user     body      user_entity.CreateUser     true  "User registration details"
// @Success      201      {object}  entity.Info               "User successfully registered"
// @Failure      400      {object}  entity.Error              "Invalid request payload"
// @Failure      500      {object}  entity.Error              "Internal server error"
// @Router       /sign-up [post]
func (u *userHandler) SignUp(c *gin.Context) {
	var user user_entity.CreateUser
	// Binding the request body to user struct
	if err := c.ShouldBindJSON(&user); err != nil {
		u.log.Error("failed to bind JSON", "error", err)
		c.JSON(http.StatusBadRequest, entity.Error{Message: "Invalid user data"})
		return
	}

	// Calling usecase to create the user
	err := u.usecase.CreateUser(c, &user)
	if err != nil {
		if err == entity.ErrUserAlreadyExists {
			c.JSON(http.StatusConflict, entity.Error{Message: err.Error()})
			return
		}
		// For other errors, return internal server error
		u.log.Error("failed to create user", "error", err)
		c.JSON(http.StatusInternalServerError, entity.Error{Message: err.Error()})
		return
	}

	// Successfully created user
	u.log.Info("User created successfully")
	c.JSON(http.StatusCreated, entity.Info{Message: "User created successfully"})
}

// SignIn godoc
// @Summary      User Sign In
// @Description  Authenticates a user and returns a JWT token upon successful login
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      user_entity.SignInRequest  true  "Sign In Request"
// @Success      200      {object}  gin.H               	  "JWT Token"
// @Failure      400      {object}  entity.Error              "Invalid user data"
// @Failure      401      {object}  entity.Error              "Incorrect password"
// @Failure      404      {object}  entity.Error              "User not found"
// @Failure      500      {object}  entity.Error              "Internal server error"
// @Router       /sign-in [post]
func (u *userHandler) SignIn(c *gin.Context) {
	// Binding the request body to user struct
	var req user_entity.SignInRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		u.log.Error("failed to bind JSON", "error", err)
		c.JSON(http.StatusBadRequest, entity.Error{Message: "Invalid user data"})
		return
	}

	token, err := u.usecase.SignIn(c, &req)
	if err != nil {
		// Check if the error is related to user not found
		if err == entity.ErrUserNotFound {
			c.JSON(http.StatusNotFound, entity.Error{Message: err.Error()})
			return
		}
		// Check if the error is related to incorrect password
		if err == entity.ErrIncorrectPassword {
			c.JSON(http.StatusUnauthorized, entity.Error{Message: err.Error()})
			return
		}
		// For other errors, return internal server error
		u.log.Error("failed to sign in user", "error", err)
		c.JSON(http.StatusInternalServerError, entity.Error{Message: "Failed to sign in user"})
		return
	}

	// Successfully signed in user
	u.log.Info("User signed in successfully")
	c.JSON(http.StatusOK, gin.H{"token": token})
}

// SendVerificationCode godoc
// @Summary      Send Verification Code
// @Description  Sends a verification code to the user's email for verification purposes
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      user_entity.SendVerificationCodeRequest  true  "Send Verification Code Request"
// @Success      200      {object}  user_entity.SendVerificationCodeResponse "Verification code sent successfully"
// @Failure      400      {object}  entity.Error                             "Invalid user data"
// @Failure      404      {object}  entity.Error              				 "User not found"
// @Failure      500      {object}  entity.Error                             "Failed to send verification code"
// @Router       /send-verification-code [post]
func (u *userHandler) SendVerificationCode(c *gin.Context) {
	// Binding the request body to user struct
	var req user_entity.SendVerificationCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		u.log.Error("failed to bind JSON", "error", err)
		c.JSON(http.StatusBadRequest, entity.Error{Message: "Invalid user data"})
		return
	}

	// Calling usecase to create the user
	duration, err := u.usecase.SendVerificationCode(c, req.Email)
	if err != nil {
		// Check if the error is related to user not found
		if err == entity.ErrUserNotFound {
			c.JSON(http.StatusNotFound, entity.Error{Message: err.Error()})
			return
		}

		u.log.Error("failed to send verification code", "error", err)
		c.JSON(http.StatusInternalServerError, entity.Error{Message: "Failed to send verification code"})
		return
	}

	// Successfully sent verification code
	u.log.Info("Verification code sent successfully", "duration", duration)
	c.JSON(http.StatusOK, user_entity.SendVerificationCodeResponse{
		Message:  "Verification code sent successfully",
		Duration: duration.String(),
	})
}

// UpdateUserPassword godoc
// @Summary      Update User Password
// @Description  Updates a user's password after validating the request
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      user_entity.UpdateUserPassword  true  "Update User Password Request"
// @Success      200      {object}  entity.Info                    "Password updated successfully"
// @Failure      400      {object}  entity.Error                   "Invalid user data"
// @Failure      500      {object}  entity.Error                   "Failed to update user password"
// @Router       /update-password [post]
func (u *userHandler) UpdateUserPassword(c *gin.Context) {
	// Binding the request body to user struct
	var req user_entity.UpdateUserPassword
	if err := c.ShouldBindJSON(&req); err != nil {
		u.log.Error("failed to bind JSON", "error", err)
		c.JSON(http.StatusBadRequest, entity.Error{Message: "Invalid user data"})
		return
	}

	// Calling usecase to update the user password
	err := u.usecase.UpdateUserPassword(c, &req)
	if err != nil {
		// For other errors, return internal server error
		u.log.Error("failed to update user password", "error", err)
		c.JSON(http.StatusInternalServerError, entity.Error{Message: "Failed to update user password"})
		return
	}

	// Successfully updated user password
	u.log.Info("User password updated successfully")
	c.JSON(http.StatusOK, entity.Info{Message: "User password updated successfully"})
}
