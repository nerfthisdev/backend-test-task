package domain

import "github.com/google/uuid"

type User struct {
	GUID     uuid.UUID `json:"guid"`
	Username string    `json:"username"`
	Password string    `json:"password"`
}
