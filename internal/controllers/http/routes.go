package handlers

import (
	"github.com/andrepulo/Calendar/internal/auth"
	"github.com/andrepulo/Calendar/internal/logger"
	"github.com/andrepulo/Calendar/internal/users"
	"github.com/go-chi/chi/middleware"
	chi "github.com/go-chi/chi/v5"
)

type (
	handlers struct {
		users *users.UserService
		auth  *auth.AuthService
		l     logger.Logger
	}
)

func NewHandlers(
	users *users.UserService,
	auth *auth.AuthService,
	log logger.Logger,
) *chi.Mux {
	r := chi.NewMux()

	h := handlers{
		users: users,
		auth:  auth,
		l:     log,
	}
	h.build(r)

	return r
}

func (h *handlers) build(r chi.Router) {
	r.Use(middleware.Recoverer)
	r.Post("/api/v1/signup", h.signup)
	r.Post("/api/v1/signin", h.signin)
	r.Route("/api/v1", func(r chi.Router) {
		r.Use(h.authMiddleware)
		r.Post("/verify", h.verify)
		r.Patch("/users", h.updateUser)
	})
}
