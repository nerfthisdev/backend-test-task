package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/joho/godotenv"
	"github.com/nerfthisdev/backend-test-task/internal/auth"
	"github.com/nerfthisdev/backend-test-task/internal/config"
	"github.com/nerfthisdev/backend-test-task/internal/database"
	"github.com/nerfthisdev/backend-test-task/internal/logger"
	"github.com/nerfthisdev/backend-test-task/internal/repository"
	server "github.com/nerfthisdev/backend-test-task/internal/router"
	"go.uber.org/zap"
)

const defaultTimeout = time.Second * 5

// @title Marketplace API
// @version 1.0
// @description REST API for marketplace service
// @BasePath /api/v1
// @host localhost:{port} in env
func main() {
	// init .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal(err.Error())
	}

	// init config
	cfg := config.InitConfig()

	// init logger
	logger := logger.GetLogger()
	logger.Info("successfully initialized logger")

	// init context
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)

	defer cancel()

	dbpool, err := database.New(ctx, cfg)
	if err != nil {
		logger.Fatal("failed to init db", zap.Error(err))
	}

	logger.Info("successfully connected to db")

	defer dbpool.Close()

	// run migrations
	err = database.RunMigrations(cfg)
	if err != nil {
		logger.Fatal("failed to run migrations", zap.Error(err))
	}

	logger.Info("successfully ran migrations")

	usersRepo := repository.NewUserRepository(dbpool)
	postsRepo := repository.NewPostRepository(dbpool)
	tokenSvc := auth.NewJWTService(cfg)

	router := server.NewRouter(usersRepo, postsRepo, tokenSvc)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	logger.Info("server starting on :" + cfg.Port)

	if err := srv.ListenAndServe(); err != nil {
		logger.Fatal("server failed", zap.Error(err))
	}
}
