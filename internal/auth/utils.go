package auth

import (
	"context"
	"errors"

	"github.com/andrepulo/Calendar/internal/users"
)

// userKeyType определяет тип ключа для хранения иденти��икатора пользователя в контексте.
type userKeyType struct{}

// userKey является ключом для хранения идентификатора пользователя в контексте.
var userKey = userKeyType{}

// UserIDFromContext извлекает идентификатор пользователя из контекста.
func UserIDFromContext(ctx context.Context) (users.UserID, error) {
	userID, ok := ctx.Value(userKey).(users.UserID)
	if !ok {
		return "", errors.New("empty user id") // Возвращает ошибку, если идентификатор пользователя пустой.
	}
	return userID, nil
}

// WithUserID добавляет идентификатор пользователя в контекст.
func WithUserID(ctx context.Context, userID users.UserID) context.Context {
	return context.WithValue(ctx, userKey, userID)
}
