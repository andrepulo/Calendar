package users

import (
	"fmt"

	"github.com/andrepulo/Calendar/internal/apperrs"
	"github.com/google/uuid"
)

type UserID = string

type User struct {
	ID       UserID   `json:"id"`
	Login    string   `json:"login"`
	Password Password `json:"-"`
}

func NewUser(login string, password Password) (User, error) {
	var zero User
	if len(login) == 0 {
		return zero, fmt.Errorf("login is empty: %w", apperrs.ErrConditionViolation)
	}
	if len(password) == 0 {
		return zero, fmt.Errorf("password is empty: %w", apperrs.ErrConditionViolation)
	}

	return User{
		ID:       uuid.New().String(),
		Login:    login,
		Password: password,
	}, nil
}

type UserChanges struct {
	Password *Password
}

type UserFilter struct {
	ID    *UserID
	Login *string
}
