package handlers

import (
	"encoding/json"
	"github.com/andrepulo/Calendar/internal/auth"
	"github.com/andrepulo/Calendar/internal/users"
	"net/http"
)

// UserHandler обрабатывает HTTP-запросы, связанные с пользователями.
type UserHandler struct {
	authService *auth.AuthService
	userService *users.UserService
}

type updateUserPayload struct {
	Password string `json:"password"`
}

// NewUserHandler создает новый экземпляр UserHandler.
func NewUserHandler(authService *auth.AuthService, userService *users.UserService) *UserHandler {
	return &UserHandler{
		authService: authService,
		userService: userService,
	}
}

// SignUp обрабатывает запрос на регистрацию нового пользователя.
func (h *UserHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	// Декодируем тело запроса в структуру req
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest) // Возвращаем ошибку, если запрос некорректен
		return
	}

	// Регистрируем нового пользователя и получаем JWT токен
	token, err := h.authService.SignUp(r.Context(), req.Login, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) // Возвращаем ошибку, если не удалось зарегистрировать пользователя
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"token": token}) // Возвращаем JWT токен в ответе
}

// SignIn обрабатывает запрос на аутентификацию пользователя.
func (h *UserHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	// Декодируем тело запроса в структуру req
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest) // Возвращаем ошибку, если запрос некорректен
		return
	}

	// Аутентифицируем пользователя и получаем JWT токен
	token, err := h.authService.SignIn(r.Context(), req.Login, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized) // Возвращаем ошибку, если не удалось аутентифицировать пользователя
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token}) // Возвращаем JWT токен в ответе
}

func (h *handlers) updateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, err := auth.UserIDFromContext(r.Context())
	if err != nil {
		h.handleError(ctx, w, err)
		return
	}
	var changes updateUserPayload
	err = json.NewDecoder(r.Body).Decode(&changes)
	if err != nil {
		h.handleError(ctx, w, err)
		return
	}
	err = r.Body.Close()
	if err != nil {
		h.handleError(ctx, w, err)
		return
	}

	err = h.users.Update(ctx, userID, changes.Password)
	if err != nil {
		h.handleError(ctx, w, err)
		return
	}
}
