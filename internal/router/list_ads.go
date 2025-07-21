package server

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/nerfthisdev/backend-test-task/internal/auth"
	"github.com/nerfthisdev/backend-test-task/internal/repository"
)

type ListAdsHandler struct {
	posts *repository.PostRepository
}

func NewListAdsHandler(posts *repository.PostRepository) *ListAdsHandler {
	return &ListAdsHandler{posts: posts}
}

type adResponse struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	ImageURL    string  `json:"image_url"`
	Price       float64 `json:"price"`
	AuthorLogin string  `json:"author_login"`
	IsOwner     bool    `json:"is_owner,omitempty"`
}

// ServeHTTP returns a list of ads with filters.
// @Summary List ads
// @Tags ads
// @Accept json
// @Produce json
// @Param page query int false "page"
// @Param per_page query int false "per page"
// @Param sort_by query string false "sort field" Enums(price, created_at)
// @Param order query string false "order" Enums(asc, desc)
// @Param min_price query number false "min price"
// @Param max_price query number false "max price"
// @Success 200 {array} adResponse
// @Failure 400 {string} string
// @Router /ads [get]
func (h *ListAdsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	if page <= 0 {
		page = 1
	}
	perPage, _ := strconv.Atoi(q.Get("per_page"))
	if perPage <= 0 {
		perPage = 10
	}
	sortBy := strings.ToLower(q.Get("sort_by"))
	if sortBy != "price" && sortBy != "created_at" {
		sortBy = "created_at"
	}
	order := strings.ToLower(q.Get("order"))
	if order != "asc" && order != "desc" {
		if sortBy == "created_at" {
			order = "desc"
		} else {
			order = "asc"
		}
	}
	var minPricePtr, maxPricePtr *float64
	if v := q.Get("min_price"); v != "" {
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}
		minPricePtr = &f
	}
	if v := q.Get("max_price"); v != "" {
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}
		maxPricePtr = &f
	}

	opts := repository.ListOptions{
		Page:     page,
		PerPage:  perPage,
		SortBy:   sortBy,
		Order:    order,
		MinPrice: minPricePtr,
		MaxPrice: maxPricePtr,
	}

	posts, err := h.posts.List(r.Context(), opts)
	if err != nil {
		http.Error(w, "failed to list ads", http.StatusInternalServerError)
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

	resp := make([]adResponse, 0, len(posts))
	for _, p := range posts {
		item := adResponse{
			Title:       p.Title,
			Description: p.Description,
			ImageURL:    p.ImageURL,
			Price:       p.Price,
			AuthorLogin: p.Username,
		}
		if hasUser {
			item.IsOwner = p.UserGUID == current
		}
		resp = append(resp, item)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
