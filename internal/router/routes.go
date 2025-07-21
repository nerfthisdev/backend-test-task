package server

import (
	"net/http"

	_ "github.com/nerfthisdev/backend-test-task/docs"
	"github.com/nerfthisdev/backend-test-task/internal/auth"
	"github.com/nerfthisdev/backend-test-task/internal/repository"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(users *repository.UserRepository, posts *repository.PostRepository, tokens *auth.JWTService) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/swagger/", httpSwagger.WrapHandler)
	mux.Handle("POST /register", NewRegisterHandler(users))
	mux.Handle("POST /login", NewLoginHandler(users, tokens))

	createAd := NewCreateAdHandler(posts)
	listAds := NewListAdsHandler(posts)
	getAd := NewGetAdHandler(posts)

	mux.Handle("POST /ads", auth.AuthMiddleware(tokens)(createAd))
	mux.Handle("GET /ads", listAds)
	mux.Handle("GET /ads/{id}", getAd)

	return mux
}
