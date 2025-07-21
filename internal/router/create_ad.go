package server

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/google/uuid"

	"github.com/nerfthisdev/backend-test-task/internal/auth"
	"github.com/nerfthisdev/backend-test-task/internal/domain"
	"github.com/nerfthisdev/backend-test-task/internal/repository"
)

type CreateAdHandler struct {
	posts *repository.PostRepository
}

func NewCreateAdHandler(posts *repository.PostRepository) *CreateAdHandler {
	return &CreateAdHandler{posts: posts}
}

type createAdRequest struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	ImageURL    string  `json:"image_url"`
	Price       float64 `json:"price"`
}

// ServeHTTP creates a new ad.
// @Summary Create ad
// @Tags ads
// @Accept json
// @Produce json
// @Param data body createAdRequest true "ad info"
// @Success 200 {object} domain.Post
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Router /ads [post]
func (h *CreateAdHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req createAdRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if len(req.Title) == 0 || len(req.Title) > 100 {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	if len(req.Description) == 0 || len(req.Description) > 500 {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	if req.Price < 0 {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	if _, err := url.ParseRequestURI(req.ImageURL); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	guid, err := uuid.Parse(userID)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	post := domain.Post{
		UserGUID:    guid,
		Title:       req.Title,
		Description: req.Description,
		ImageURL:    req.ImageURL,
		Price:       req.Price,
	}

	created, err := h.posts.Create(r.Context(), post)
	if err != nil {
		http.Error(w, "failed to create ad", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(created)
}
