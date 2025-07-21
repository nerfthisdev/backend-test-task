package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nerfthisdev/backend-test-task/internal/domain"
)

type PostRepository struct {
	db *pgxpool.Pool
}

func NewPostRepository(db *pgxpool.Pool) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) Create(ctx context.Context, post domain.Post) (*domain.Post, error) {
	query := `INSERT INTO posts (user_guid, title, description, image_url, price)
              VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := r.db.QueryRow(ctx, query,
		post.UserGUID, post.Title, post.Description, post.ImageURL, post.Price,
	).Scan(&post.ID)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

type ListOptions struct {
	Page     int
	PerPage  int
	SortBy   string
	Order    string
	MinPrice *float64
	MaxPrice *float64
}

type Ad struct {
	ID          int64
	UserGUID    uuid.UUID
	Username    string
	Title       string
	Description string
	ImageURL    string
	Price       float64
	CreatedAt   time.Time
}

func (r *PostRepository) List(ctx context.Context, opt ListOptions) ([]Ad, error) {
	query := `SELECT p.id, p.user_guid, u.username, p.title, p.description, p.image_url, p.price, p.created_at
                FROM posts p JOIN users u ON p.user_guid = u.guid`
	params := []any{}
	idx := 1

	if opt.MinPrice != nil {
		query += fmt.Sprintf(" WHERE p.price >= $%d", idx)
		params = append(params, *opt.MinPrice)
		idx++
	}
	if opt.MaxPrice != nil {
		if len(params) == 0 {
			query += fmt.Sprintf(" WHERE p.price <= $%d", idx)
		} else {
			query += fmt.Sprintf(" AND p.price <= $%d", idx)
		}
		params = append(params, *opt.MaxPrice)
		idx++
	}

	sortCol := "created_at"
	if opt.SortBy == "price" {
		sortCol = "price"
	}
	order := "DESC"
	if strings.ToUpper(opt.Order) == "ASC" {
		order = "ASC"
	}
	query += fmt.Sprintf(" ORDER BY p.%s %s", sortCol, order)

	limit := opt.PerPage
	if limit <= 0 {
		limit = 10
	}
	if opt.Page <= 0 {
		opt.Page = 1
	}
	offset := (opt.Page - 1) * limit
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", idx, idx+1)
	params = append(params, limit, offset)

	rows, err := r.db.Query(ctx, query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ads []Ad
	for rows.Next() {
		var a Ad
		if err := rows.Scan(&a.ID, &a.UserGUID, &a.Username, &a.Title, &a.Description, &a.ImageURL, &a.Price, &a.CreatedAt); err != nil {
			return nil, err
		}
		ads = append(ads, a)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return ads, nil
}

func (r *PostRepository) Get(ctx context.Context, id int64) (*Ad, error) {
	query := `SELECT p.id, p.user_guid, u.username, p.title, p.description, p.image_url, p.price, p.created_at
               FROM posts p JOIN users u ON p.user_guid = u.guid
               WHERE p.id = $1`

	var a Ad
	err := r.db.QueryRow(ctx, query, id).Scan(&a.ID, &a.UserGUID, &a.Username, &a.Title, &a.Description, &a.ImageURL, &a.Price, &a.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &a, nil
}
