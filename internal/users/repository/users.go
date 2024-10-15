package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/andrepulo/Calendar/internal/users"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/andrepulo/Calendar/internal/apperrs"
	"github.com/andrepulo/Calendar/internal/databases"
	"github.com/georgysavva/scany/v2/dbscan"
	"github.com/rgurov/pgerrors"
)

// Repository представляет репозиторий для работы с пользователями.
type Repository struct {
	db *databases.DB
}

// NewRepository создает новый экземпляр репозитория.
// Вся логика работы с БД и любым другим хранилищем инкапсулирована в слое репозитория.
// Методы репозитория принимают доменную модель и возвращают доменную модель.
func NewRepository(db *databases.DB) *Repository {
	return &Repository{db}
}

const createQuery = `
INSERT INTO users (
 id, login, password
) VALUES (
 $1, $2, $3
)
`

// Create добавляет нового пользователя в базу данных.
func (r *Repository) Create(ctx context.Context, user users.User) error {
	_, err := r.db.ExecContext(ctx, createQuery,
		user.ID,
		user.Login,
		user.Password,
	)
	if err != nil {
		// Проверяем, если ошибка связана с уникальностью логина.
		if strings.Contains(err.Error(), pgerrors.UniqueViolation) {
			return fmt.Errorf("user with login %s: %w", user.Login, apperrs.ErrAlreadyExist)
		}
		return fmt.Errorf("user with login: %w", err)
	}

	return nil
}

// Update обновляет данные пользователя в б��зе данных.
func (r *Repository) Update(ctx context.Context, userID users.UserID, changes users.UserChanges) error {
	qb := squirrel.Update("users").Where(squirrel.Eq{"id": userID})
	qb = qb.SetMap(changesBuilder(changes).ToMap())

	query, args, err := qb.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return fmt.Errorf("build query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("exec query: %w", err)
	}

	return nil
}

// Get возвращает пользователя из базы данных по за��анному фильтру.
func (r *Repository) Get(ctx context.Context, filter users.UserFilter) (_ users.User, err error) {
	var zero users.User

	qb := squirrel.Select("id", "login", "password").From("users")
	qb = userFilterBuilder(filter).apply(qb)
	query, args, err := qb.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return zero, fmt.Errorf("build conditions: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return zero, fmt.Errorf("exec query: %w", err)
	}
	defer func() {
		err = errors.Join(err, rows.Close())
	}()

	var u users.User
	err = dbscan.ScanOne(&u, rows)
	if err != nil {
		return zero, fmt.Errorf("scan rows: %w", err)
	}

	return u, nil
}

// changesBuilder помогает строить карту изменений для обновления пользователя.
type changesBuilder users.UserChanges

// ToMap преобразует изменения пользователя в карту.
func (b changesBuilder) ToMap() map[string]any {
	fields := make(map[string]any)
	if b.Password != nil {
		fields["password"] = *b.Password
	}

	return fields
}

// userFilterBuilder помогает строить условия фильтрации для запроса пользователя.
type userFilterBuilder users.UserFilter

// apply применяет фильтры к запросу.
func (f userFilterBuilder) apply(qb squirrel.SelectBuilder) squirrel.SelectBuilder {
	if f.ID != nil {
		qb = qb.Where(squirrel.Eq{"id": *f.ID})
	}
	if f.Login != nil {
		qb = qb.Where(squirrel.Eq{"login": *f.Login})
	}

	return qb
}
