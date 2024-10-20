package apperrs

import "errors"

// Определение различных ошибок, которые могут возникнуть в приложении.
var (
	// ErrAlreadyExist указывает на то, что объект уже существует.
	ErrAlreadyExist = errors.New("already exist")

	// ErrNotFound указывает на то, что объект не найден.
	ErrNotFound = errors.New("not found")

	// ErrConditionViolation указывает на нарушение условия.
	ErrConditionViolation = errors.New("condition violation")

	// ErrUnauthorize указывает на отсутствие авторизации.
	ErrUnauthorize = errors.New("unauthorize")
)
