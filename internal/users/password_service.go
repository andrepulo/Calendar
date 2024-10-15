package users

import (
	"golang.org/x/crypto/bcrypt"
)

// PasswordService provides methods for password hashing and comparison.
type PasswordService struct{}

// NewPasswordService creates a new instance of PasswordService.
func NewPasswordService() *PasswordService {
	return &PasswordService{}
}

// FromString hashes a password string.
func (s *PasswordService) FromString(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// Compare checks if the given password matches the hashed password.
func (s *PasswordService) Compare(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
