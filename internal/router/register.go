package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/nerfthisdev/backend-test-task/internal/domain"
	"github.com/nerfthisdev/backend-test-task/internal/repository"
)

type RegisterHandler struct {
	users *repository.UserRepository
}

func NewRegisterHandler(users *repository.UserRepository) *RegisterHandler {
	return &RegisterHandler{users: users}
}

type registerRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type registerResponse struct {
	GUID     uuid.UUID `json:"guid"`
	Username string    `json:"username"`
}

func isValidLogin(login string) bool {
	if len(login) < 3 || len(login) > 20 {
		return false
	}
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9]+$`, login)
	return matched
}

func isValidPassword(password string) bool {
	if len(password) < 8 || len(password) > 32 {
		return false
	}
	hasLetter, _ := regexp.MatchString(`[A-Za-z]`, password)
	hasDigit, _ := regexp.MatchString(`[0-9]`, password)
	return hasLetter && hasDigit
}

// ServeHTTP registers a new user.
// @Summary Register new user
// @Tags auth
// @Accept json
// @Produce json
// @Param data body registerRequest true "credentials"
// @Success 200 {object} registerResponse
// @Failure 400 {string} string
// @Router /register [post]
func (h *RegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if !isValidLogin(req.Login) || !isValidPassword(req.Password) {
		http.Error(w, "invalid credentials", http.StatusBadRequest)
		return
	}

	if _, err := h.users.GetByUsername(r.Context(), req.Login); err == nil {
		http.Error(w, "user already exists", http.StatusConflict)
		return
	} else if !errors.Is(err, pgx.ErrNoRows) {
		http.Error(w, "failed to check user", http.StatusInternalServerError)
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "failed to hash password", http.StatusInternalServerError)
		return
	}

	user := domain.User{
		GUID:     uuid.New(),
		Username: req.Login,
		Password: string(hashed),
	}

	if err := h.users.Create(r.Context(), user); err != nil {
		http.Error(w, "failed to create user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(registerResponse{GUID: user.GUID, Username: user.Username})
}
