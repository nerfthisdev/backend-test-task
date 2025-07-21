package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/nerfthisdev/backend-test-task/internal/config"
)

type JWTService struct {
	secret    string
	accessTTL time.Duration
}

func NewJWTService(cfg config.Config) *JWTService {
	return &JWTService{
		secret:    cfg.JWTSecret,
		accessTTL: cfg.AccessTTL,
	}
}

func (s *JWTService) GenerateAccessToken(guid uuid.UUID) (string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.RegisteredClaims{
		Subject:   guid.String(),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.accessTTL)),
	})

	accessTokenString, err := accessToken.SignedString([]byte(s.secret))
	if err != nil {
		return "", err
	}

	return accessTokenString, nil
}

func (s *JWTService) ValidateAccessToken(token string) (map[string]any, error) {
	parsed, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", t.Header["alg"])
		}
		return []byte(s.secret), nil
	})
	if err != nil {
		return nil, err
	}

	if !parsed.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := parsed.Claims.(jwt.MapClaims)

	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}

	return claims, nil
}
