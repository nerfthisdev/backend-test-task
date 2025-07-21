package main

import (
	"context"
	"fmt"
	"log"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"

	"github.com/nerfthisdev/backend-test-task/internal/config"
	"github.com/nerfthisdev/backend-test-task/internal/database"
	"github.com/nerfthisdev/backend-test-task/internal/domain"
	"github.com/nerfthisdev/backend-test-task/internal/repository"
)

func main() {
	// Load environment variables from .env if present
	_ = godotenv.Load()

	cfg := config.InitConfig()
	ctx := context.Background()

	dbpool, err := database.New(ctx, cfg)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer dbpool.Close()

	if err := database.RunMigrations(cfg); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	userRepo := repository.NewUserRepository(dbpool)
	postRepo := repository.NewPostRepository(dbpool)

	gofakeit.Seed(0)

	const numUsers = 5
	const adsPerUser = 5

	var users []domain.User
	for i := 0; i < numUsers; i++ {
		pw, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
		u := domain.User{
			GUID:     uuid.New(),
			Username: gofakeit.Username(),
			Password: string(pw),
		}
		if err := userRepo.Create(ctx, u); err != nil {
			log.Fatalf("failed to create user: %v", err)
		}
		users = append(users, u)
	}

	for _, u := range users {
		for i := 0; i < adsPerUser; i++ {
			post := domain.Post{
				UserGUID:    u.GUID,
				Title:       gofakeit.Sentence(3),
				Description: gofakeit.Paragraph(1, 2, 5, " "),
				ImageURL:    fmt.Sprintf("https://picsum.photos/seed/%d/640/480", gofakeit.Number(1, 100000)),
				Price:       gofakeit.Price(10, 1000),
			}
			if _, err := postRepo.Create(ctx, post); err != nil {
				log.Fatalf("failed to create post: %v", err)
			}
		}
	}

	log.Printf("created %d users with %d ads\n", numUsers, numUsers*adsPerUser)
}
