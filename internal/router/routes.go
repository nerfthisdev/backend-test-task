package server

import (
	"net/http"

	_ "github.com/nerfthisdev/backend-test-task/docs"
	"github.com/nerfthisdev/backend-test-task/internal/auth"
	"github.com/nerfthisdev/backend-test-task/internal/logging"
	"github.com/nerfthisdev/backend-test-task/internal/repository"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

func NewRouter(users *repository.UserRepository, posts *repository.PostRepository, tokens *auth.JWTService, logger *zap.Logger) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/swagger/", httpSwagger.WrapHandler)
	mux.Handle("POST /api/v1/register", NewRegisterHandler(users))
	mux.Handle("POST /api/v1/login", NewLoginHandler(users, tokens))

	createAd := NewCreateAdHandler(posts)
	listAds := NewListAdsHandler(posts)
	getAd := NewGetAdHandler(posts)

	mux.Handle("POST /api/v1/ads", auth.AuthMiddleware(tokens)(createAd))
	mux.Handle("GET /api/v1/ads", listAds)
	mux.Handle("GET /api/v1/ads/{id}", getAd)

	return logging.LoggingMiddleware(logger)(mux)
}
