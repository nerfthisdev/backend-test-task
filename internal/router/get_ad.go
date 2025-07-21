package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/nerfthisdev/backend-test-task/internal/auth"
	"github.com/nerfthisdev/backend-test-task/internal/repository"
)

// ServeHTTP returns ad by id.
// @Summary Get ad
// @Tags ads
// @Accept json
// @Produce json
// @Param id path int true "Ad ID"
// @Success 200 {object} singleAdResponse
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Router /ads/{id} [get]
type GetAdHandler struct {
	posts *repository.PostRepository
}

func NewGetAdHandler(posts *repository.PostRepository) *GetAdHandler {
	return &GetAdHandler{posts: posts}
}

type singleAdResponse struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	ImageURL    string  `json:"image_url"`
	Price       float64 `json:"price"`
	AuthorLogin string  `json:"author_login"`
	IsOwner     bool    `json:"is_owner,omitempty"`
}

func (h *GetAdHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	post, err := h.posts.Get(r.Context(), id)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	var current uuid.UUID
	var hasUser bool
	if idStr, ok := auth.UserIDFromContext(r.Context()); ok {
		if uid, err := uuid.Parse(idStr); err == nil {
			current = uid
			hasUser = true
		}
	}

	resp := singleAdResponse{
		Title:       post.Title,
		Description: post.Description,
		ImageURL:    post.ImageURL,
		Price:       post.Price,
		AuthorLogin: post.Username,
	}
	if hasUser {
		resp.IsOwner = post.UserGUID == current
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
