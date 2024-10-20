package users

import (
	"golang.org/x/crypto/bcrypt"
)

// PasswordService предоставляет методы для хеширования и сравнения паролей.
type PasswordService struct{}

// NewPasswordService создает новый экземпляр PasswordService.
func NewPasswordService() *PasswordService {
	return &PasswordService{}
}

// FromString хеширует строку пароля.
func (s *PasswordService) FromString(password string) (string, error) {
	// Генерируем хеш из пароля с использованием bcrypt.
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err // Возвращаем ошибку, если не удалось создать хеш.
	}
	return string(hash), nil // Возвращаем хешированный пароль в виде строки.
}

// Compare проверяет, соответствует ли данный пароль хешированному паролю.
func (s *PasswordService) Compare(hash, password string) bool {
	// Сравниваем хешированный пароль с предоставленным паролем.
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
