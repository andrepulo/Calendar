package users

import (
	"database/sql/driver"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type Password []byte

func (p Password) IsSamePassword(other string) bool {
	return nil == bcrypt.CompareHashAndPassword([]byte(p), []byte(other))
}

func (pa *Password) Scan(src any) error {
	v, ok := src.(string)
	if !ok {
		return fmt.Errorf("unexpected type: %T, expect %T", src, Password{})
	}

	*pa = Password(v)
	return nil
}

// Value returns a driver Value.
// Value must not panic.
func (pa *Password) Value() (driver.Value, error) {
	if pa == nil {
		return nil, nil
	}
	return []byte(*pa), nil
}
