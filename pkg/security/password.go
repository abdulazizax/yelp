package security

import (
	"errors"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes the plain-text password.
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// CheckPassword compares the hashed password with the plain-text password.
func CheckPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// ValidatePassword checks if the password meets basic security requirements.
func ValidatePassword(password string) error {
	const minPasswordLength = 8
	const maxPasswordLength = 128 // Gmail-style upper limit

	// Check password length
	if len(password) < minPasswordLength {
		return errors.New("password must be at least 8 characters long")
	}
	if len(password) > maxPasswordLength {
		return errors.New("password must not exceed 128 characters")
	}

	// Check for the presence of different character types
	var hasUpper, hasLower, hasNumber bool
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasNumber = true
		}
	}

	// Minimal requirements:
	// - At least one lowercase letter
	// - At least one uppercase letter or one number
	if !hasLower {
		return errors.New("password must contain at least one lowercase letter")
	}
	if !hasUpper && !hasNumber {
		return errors.New("password must contain at least one uppercase letter or one number")
	}

	// The password meets all criteria
	return nil
}
