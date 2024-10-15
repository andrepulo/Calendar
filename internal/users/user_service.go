package users

import (
	"context"
	"fmt"

	"go.openly.dev/pointy"
)

// repository определяет интерфейс для работы с репозиторием пользователей.
type repository interface {
	Create(ctx context.Context, user User) error
	Update(ctx context.Context, userID UserID, userChanges UserChanges) error
	Get(ctx context.Context, f UserFilter) (User, error)
}

// UserService предоставляет методы для работы с пользователями.
type UserService struct {
	passwords *PasswordService
	repo      repository
}

// NewUserService создает новый экземпляр UserService.
func NewUserService(repo repository, passwordService *PasswordService) *UserService {
	return &UserService{
		passwords: passwordService,
		repo:      repo,
	}
}

// Create создает нового пользователя с заданным логином и паролем.
func (s *UserService) Create(ctx context.Context, login, password string) (User, error) {
	var zero User

	// Хешируем пароль.
	hashPassword, err := s.passwords.FromString(password)
	if err != nil {
		return User{}, fmt.Errorf("prepare password: %w", err)
	}
	// Создаем нового пользователя.
	user, err := NewUser(login, hashPassword)
	if err != nil {
		return zero, fmt.Errorf("create new user: %w", err)
	}
	// Сохраняем пользователя в репозитории.
	err = s.repo.Create(ctx, user)
	if err != nil {
		return zero, fmt.Errorf("persist new user: %w", err)
	}

	return user, nil
}

// Update обновляет пароль пользователя.
func (s *UserService) Update(ctx context.Context, userID UserID, password string) error {
	// Хешируем новый пароль.
	hashPass, err := s.passwords.FromString(password)
	if err != nil {
		return fmt.Errorf("prepare password: %w", err)
	}
	// Создаем изменения для пользователя.
	changes := UserChanges{
		Password: pointy.Pointer(hashPass),
	}
	// Обновляем пользователя в репозитории.
	return s.repo.Update(ctx, userID, changes)
}

// Get возвращает пользователя по заданному фильтру.
func (s *UserService) Get(ctx context.Context, filter UserFilter) (User, error) {
	return s.repo.Get(ctx, filter)
}
