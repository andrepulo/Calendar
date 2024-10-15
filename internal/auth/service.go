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

const tokenTTL = 1 * time.Hour

type AuthService struct {
	users     *users.UserService
	passwords *users.PasswordService
	cfg       config.Security
}

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

func (s *AuthService) SignIn(ctx context.Context, login, password string) (string, error) {
	user, err := s.users.Get(ctx, users.UserFilter{
		Login: pointy.String(login),
	})
	if err != nil {
		return "", fmt.Errorf("get user: %w", err)
	}

	if !s.passwords.Compare(users.Password, password) {
		return "", apperrs.ErrUnauthorize
	}

	token, err := s.generateToken(user.ID)
	if err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}

	return token, nil
}

func (s *AuthService) SignUp(ctx context.Context, login, password string) (string, error) {
	user, err := s.users.Create(ctx, login, password)
	if err != nil {
		return "", fmt.Errorf("create new user: %w", err)
	}

	token, err := s.generateToken(user.ID)
	if err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}

	return token, nil
}

func (s *AuthService) generateToken(userID users.UserID) (string, error) {
	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(tokenTTL).Unix(),
		Issuer:    string(userID),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.SecretKey))
}

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
