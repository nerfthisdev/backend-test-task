package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nerfthisdev/backend-test-task/internal/domain"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func (r *UserRepository) Create(ctx context.Context, user domain.User) error {
	query := `INSERT INTO users (guid, username, password)
	VALUES ($1, $2, $3)`

	_, err := r.db.Exec(ctx, query, user.GUID, user.Username, user.Password)

	return err
}

func (r *UserRepository) GetByGuid(ctx context.Context, guid uuid.UUID) (*domain.User, error) {
	query := `SELECT guid, username, password FROM users WHERE guid = $1`

	row := r.db.QueryRow(ctx, query, guid)

	var user domain.User

	err := row.Scan(&user.GUID, &user.Username, &user.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	query := `SELECT guid, username, password FROM users WHERE username = $1`

	row := r.db.QueryRow(ctx, query, username)

	var user domain.User

	err := row.Scan(&user.GUID, &user.Username, &user.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) Delete(ctx context.Context, user domain.User) error {
	query := `DELETE FROM users WHERE guid = $1`

	_, err := r.db.Exec(ctx, query, user.GUID)

	return err
}
