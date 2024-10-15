package auth

import (
	"context"
	"errors"

	"github.com/andrepulo/Calendar/internal/users"
)

type userKeyType struct{}

var userKey = userKeyType{}

func UserIDFromContext(ctx context.Context) (users.UserID, error) {
	userID, ok := ctx.Value(userKey).(users.UserID)
	if !ok {
		return "", errors.New("empty user id")
	}
	return userID, nil
}

func WithUserID(ctx context.Context, userID users.UserID) context.Context {
	return context.WithValue(ctx, userKey, userID)
}
