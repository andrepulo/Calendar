package databases

import (
	"database/sql"
	"fmt"

	"github.com/andrepulo/Calendar/internal/config"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/stdlib"
)

type (
	// closeFn определяет тип функции для закрытия соединения с базой данных.
	closeFn func() error
	// DB представляет тип базы данных.
	DB = sql.DB
)

// NewDB создает новое соединение с базой данных на основе конфигурации.
func NewDB(cfg *config.DB) (*sql.DB, closeFn, error) {
	// Разбираем URI для подключения к базе данных.
	connCfg, err := pgx.ParseURI(cfg.URI)
	if err != nil {
		return nil, nil, fmt.Errorf("parse URI: %w", err)
	}
	// Открываем соединение с базой данных.
	db := stdlib.OpenDB(connCfg)
	// Проверяем соединение с базой данных.
	err = db.Ping()
	if err != nil {
		return nil, nil, fmt.Errorf("check connection: %w", err)
	}

	// Возвращаем объект базы данных и функцию для закрытия соединения.
	return db, db.Close, nil
}
