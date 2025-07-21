package server

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/nerfthisdev/backend-test-task/internal/auth"
	"github.com/nerfthisdev/backend-test-task/internal/repository"
)

type LoginHandler struct {
	users  *repository.UserRepository
	tokens *auth.JWTService
}

func NewLoginHandler(users *repository.UserRepository, tokens *auth.JWTService) *LoginHandler {
	return &LoginHandler{users: users, tokens: tokens}
}

type loginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}

// ServeHTTP authenticates the user and returns a token.
// @Summary Login user
// @Tags auth
// @Accept json
// @Produce json
// @Param data body loginRequest true "credentials"
// @Success 200 {object} loginResponse
// @Failure 401 {string} string
// @Router /login [post]
func (h *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	user, err := h.users.GetByUsername(r.Context(), req.Login)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := h.tokens.GenerateAccessToken(user.GUID)
	if err != nil {
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("X-Auth-Token", token)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(loginResponse{Token: token})
}
