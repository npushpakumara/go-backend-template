package user

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a given password using bcrypt with the default cost.
func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// CheckPassword verifies that a hashed password matches the given raw password.
// Returns nil if the passwords match, otherwise returns an error.
func CheckPassword(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return errors.New("incorrect password")
		}
		return err
	}
	return nil
}
