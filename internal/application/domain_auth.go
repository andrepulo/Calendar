package application

import (
	"github.com/andrepulo/Calendar/internal/auth"
	"github.com/andrepulo/Calendar/internal/config"
)

type authDomain struct {
	auth *auth.AuthService
}

func buildAuthDomain(cfg *config.Config, ud usersDomain) authDomain {
	authService := auth.NewAuthService(cfg.Security, ud.users, ud.passwords)
	return authDomain{
		auth: authService,
	}
}
