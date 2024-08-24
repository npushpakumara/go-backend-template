package auth

import (
	"errors"

	apiError "github.com/npushpakumara/go-backend-template/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// hashPassword hashes a given password using bcrypt with the default cost.
func hashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// checkPassword verifies that a hashed password matches the given raw password.
// Returns nil if the passwords match, otherwise returns an error.
func checkPassword(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return apiError.ErrIncorrectPassword
		}
		return err
	}
	return nil
}
