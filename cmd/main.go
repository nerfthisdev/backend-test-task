package main

import (
	"context"
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/nerfthisdev/backend-test-task/internal/config"
	"github.com/nerfthisdev/backend-test-task/internal/database"
	"github.com/nerfthisdev/backend-test-task/internal/logger"
	"go.uber.org/zap"
)

const defaultTimeout = time.Second * 5

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

	dbpool, err := database.InitDB(ctx, cfg)
	if err != nil {
		logger.Fatal("failed to init db", zap.Error(err))
	}

	logger.Info("successfully connected to db")

	dbpool.Close()
}
