package apperrs

import "errors"

var (
	ErrAlreadyExist       = errors.New("already exist")
	ErrNotFound           = errors.New("not found")
	ErrConditionViolation = errors.New("condition violation")
	ErrUnauthorize        = errors.New("unauthorize")
)
