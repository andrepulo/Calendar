package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/andrepulo/Calendar/internal/apperrs"
	"github.com/andrepulo/Calendar/internal/config"
	"github.com/andrepulo/Calendar/internal/users"
	"github.com/golang-jwt/jwt"
	"go.openly.dev/pointy"
)

// tokenTTL определяет время жизни JWT токена.
const tokenTTL = 1 * time.Hour

// AuthService предоставляет методы для аутентификации пользователей.
type AuthService struct {
	users     *users.UserService
	passwords *users.PasswordService
	cfg       config.Security
}

// NewAuthService создает новый экземпляр AuthService.
func NewAuthService(
	cfg config.Security,
	users *users.UserService,
	passwords *users.PasswordService,
) *AuthService {
	return &AuthService{
		users:     users,
		passwords: passwords,
		cfg:       cfg,
	}
}

// SignIn аутентифицирует пользователя и возвращает JWT токен.
func (s *AuthService) SignIn(ctx context.Context, login, password string) (string, error) {
	// Получаем пользователя по логину
	user, err := s.users.Get(ctx, users.UserFilter{
		Login: pointy.String(login),
	})
	if err != nil {
		return "", fmt.Errorf("get user: %w", err)
	}

	// Сравниваем предоставленный пароль с сохраненным хешем пароля
	if !s.passwords.Compare(users.Password, password) {
		return "", apperrs.ErrUnauthorize
	}

	// Генерируем JWT токен для аутентифицированного пользователя
	token, err := s.generateToken(user.ID)
	if err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}

	return token, nil
}

// SignUp регистрирует нового пользователя и возвращает JWT токен.
func (s *AuthService) SignUp(ctx context.Context, login, password string) (string, error) {
	// Создаем нового пользователя
	user, err := s.users.Create(ctx, login, password)
	if err != nil {
		return "", fmt.Errorf("create new user: %w", err)
	}

	// Генерируем JWT токен для нового пользователя
	token, err := s.generateToken(user.ID)
	if err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}

	return token, nil
}

// generateToken создает JWT токен для заданного идентификатора пользователя.
func (s *AuthService) generateToken(userID users.UserID) (string, error) {
	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(tokenTTL).Unix(),
		Issuer:    string(userID),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.SecretKey))
}

// Verify проверяет валидность заданного JWT токена и возвращает идентификатор пользователя.
func (s *AuthService) Verify(token string) (users.UserID, error) {
	var zero users.UserID
	token = strings.TrimSpace(token)

	t, err := jwt.ParseWithClaims(token, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.cfg.SecretKey), nil
	})
	if err != nil {
		return zero, errors.Join(apperrs.ErrUnauthorize, fmt.Errorf("parse token: %w", err))
	}

	claims := t.Claims.(*jwt.StandardClaims)
	return users.UserID(claims.Issuer), nil
}
