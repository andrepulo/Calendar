package users

import (
	"fmt"

	"github.com/andrepulo/Calendar/internal/apperrs"
	"github.com/google/uuid"
)

// UserID представляет тип идентификатора пользователя.
type UserID = string

// User представляет структуру пользователя.
type User struct {
	ID       UserID   `json:"id"`    // Идентификатор пользователя
	Login    string   `json:"login"` // Логин пользователя
	Password Password `json:"-"`     // Хешированный пароль пользователя
}

// NewUser создает нового пользователя с заданным логином и паролем.
func NewUser(login string, password Password) (User, error) {
	var zero User
	// Проверяем, что логин не пустой.
	if len(login) == 0 {
		return zero, fmt.Errorf("login is empty: %w", apperrs.ErrConditionViolation)
	}
	// Проверяем, что пароль не пустой.
	if len(password) == 0 {
		return zero, fmt.Errorf("password is empty: %w", apperrs.ErrConditionViolation)
	}

	// Возвращаем нового пользователя с уникальным идентификатором.
	return User{
		ID:       uuid.New().String(),
		Login:    login,
		Password: password,
	}, nil
}

// UserChanges представляет изменения, которые можно внести в пользователя.
type UserChanges struct {
	Password *Password // Новое значение пароля
}

// UserFilter используется для фильтрации пользователей.
type UserFilter struct {
	ID    *UserID // Фильтр по идентификатору пользователя
	Login *string // Фильтр по логину пользователя
}
