package users

import (
	"database/sql/driver"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// Password представляет тип для хранения хешированного пароля.
type Password []byte

// IsSamePassword сравнивает хешированный пароль с предоставленным паролем.
func (p Password) IsSamePassword(other string) bool {
	// Сравниваем хешированный пароль с предоставленным паролем.
	return nil == bcrypt.CompareHashAndPassword([]byte(p), []byte(other))
}

// Scan реализует интерфейс sql.Scanner для типа Password.
func (pa *Password) Scan(src any) error {
	// Проверяем, что источник является строкой.
	v, ok := src.(string)
	if !ok {
		return fmt.Errorf("unexpected type: %T, expect %T", src, Password{})
	}

	// Присваиваем значение хешированного пароля.
	*pa = Password(v)
	return nil
}

// Value возвращает значение для хранения в базе данных.
// Value не должен вызывать панику.
func (pa *Password) Value() (driver.Value, error) {
	// Если пароль пустой, возвращаем nil.
	if pa == nil {
		return nil, nil
	}
	// Возвращ��ем хешированный пароль.
	return []byte(*pa), nil
}
