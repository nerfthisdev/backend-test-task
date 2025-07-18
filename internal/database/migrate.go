package database

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"

	"github.com/nerfthisdev/backend-test-task/internal/config"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(cfg config.Config) error {
	dburi := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBAddress,
		cfg.DBPort,
		cfg.DBName,
	)

	m, err := migrate.New(
		"file:///app/migrations",
		dburi,
	)
	if err != nil {
		return fmt.Errorf("migration init error: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration up error: %w", err)
	}

	return nil
}
