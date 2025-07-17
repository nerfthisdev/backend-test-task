package repository

import "github.com/jackc/pgx/v5/pgxpool"

type TokenRepository struct {
	db *pgxpool.Pool
}
