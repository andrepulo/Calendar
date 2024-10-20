package application

import (
	"github.com/andrepulo/Calendar/internal/config"
	"github.com/andrepulo/Calendar/internal/databases"
	"github.com/andrepulo/Calendar/internal/users"
	"github.com/andrepulo/Calendar/internal/users/repository"
)

type usersDomain struct {
	passwords *users.PasswordService
	users     *users.UserService
}

func buildUsersDomain(cfg *config.Config, db *databases.DB) usersDomain {
	repo := repository.NewRepository(db)
	passService := users.NewPasswordService(&cfg.Security)
	userService := users.NewUserService(repo, passService)
	return usersDomain{
		passwords: passService,
		users:     userService,
	}
}
