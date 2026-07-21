package account

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// hashPassword returns a bcrypt hash of the given password.
func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("hash password: %w", err)
	}
	return string(hash), nil
}

// verifyPassword compares a password against a bcrypt hash.
// Returns nil if they match, or an error otherwise.
func verifyPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
