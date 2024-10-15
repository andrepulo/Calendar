package users

import (
	"context"
	"fmt"

	"go.openly.dev/pointy"
)

type repository interface {
	Create(ctx context.Context, user User) error
	Update(ctx context.Context, userID UserID, userChanges UserChanges) error
	Get(ctx context.Context, f UserFilter) (User, error)
}

type UserService struct {
	passwords *PasswordService
	repo      repository
}

func NewUserService(repo repository, passwordService *PasswordService) *UserService {
	return &UserService{
		passwords: passwordService,
		repo:      repo,
	}
}

func (s *UserService) Create(ctx context.Context, login, password string) (User, error) {
	var zero User

	hashPassword, err := s.passwords.FromString(password)
	if err != nil {
		return User{}, fmt.Errorf("prepare password: %w", err)
	}
	user, err := NewUser(login, hashPassword)
	if err != nil {
		return zero, fmt.Errorf("create new user: %w", err)
	}
	err = s.repo.Create(ctx, user)
	if err != nil {
		return zero, fmt.Errorf("persist new user: %w", err)
	}

	return user, nil
}

func (s *UserService) Update(ctx context.Context, userID UserID, password string) error {
	hashPass, err := s.passwords.FromString(password)
	if err != nil {
		return fmt.Errorf("prepare password: %w", err)
	}
	changes := UserChanges{
		Password: pointy.Pointer(hashPass),
	}
	return s.repo.Update(ctx, userID, changes)
}

func (s *UserService) Get(ctx context.Context, filter UserFilter) (User, error) {
	return s.repo.Get(ctx, filter)
}
