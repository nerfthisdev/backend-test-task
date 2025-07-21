package domain

import "github.com/google/uuid"

type User struct {
	GUID     uuid.UUID `json:"guid"`
	Username string    `json:"username"`
	Password string    `json:"password"`
}

type Post struct {
	ID          int64     `json:"id"`
	UserGUID    uuid.UUID `json:"user_guid"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ImageURL    string    `json:"image_url"`
	Price       float64   `json:"price"`
}
