package entity

import "fmt"

// Error type
var (
	ErrUserAlreadyExists = fmt.Errorf("user with this email already exists")
	ErrInvalidUserData   = fmt.Errorf("invalid data")
	ErrFailedToCreate    = fmt.Errorf("failed to create")
	ErrUserNotFound      = fmt.Errorf("user not found")
	ErrIncorrectPassword = fmt.Errorf("incorrect password")
	ErrInvalidPassword   = fmt.Errorf("invalid password")
	ErrUserNotExists     = fmt.Errorf("user does not exist")
)

// Structs for API response
type (
	Error struct {
		Message string `json:"message"`
	}

	Info struct {
		Message string `json:"message"`
	}
)
