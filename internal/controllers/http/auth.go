package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/andrepulo/Calendar/internal/auth"
	"github.com/andrepulo/Calendar/internal/users"
)

const tokenHeader = "Authorization"

func (h *handlers) signin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	login, password, ok := r.BasicAuth()
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	token, err := h.auth.SignIn(ctx, login, password)
	if err != nil {
		h.handleError(ctx, w, err)
		return
	}

	w.Header().Add(tokenHeader, "Bearer "+token)
}

func (h *handlers) signup(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	login, password, ok := r.BasicAuth()
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	token, err := h.auth.SignUp(ctx, login, password)
	if err != nil {
		h.handleError(ctx, w, err)
		return
	}

	w.Header().Add(tokenHeader, "Bearer "+token)
}

func (h *handlers) verify(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, err := auth.UserIDFromContext(ctx)
	if err != nil {
		h.handleError(ctx, w, err)
		return
	}

	user, err := h.users.Get(ctx, users.UserFilter{
		ID: &userID,
	})
	if err != nil {
		h.handleError(ctx, w, err)
		return
	}

	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		h.handleError(ctx, w, err)
		return
	}
}
