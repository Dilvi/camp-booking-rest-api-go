package app

import (
	"net/http"

	"github.com/dilvi/camp-booking-rest-api-go/internal/handler"
	"github.com/dilvi/camp-booking-rest-api-go/internal/middleware"
)

func NewRouter(authHandler *handler.AuthHandler, jwtSecret string) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", handler.HealthHandler)
	mux.HandleFunc("/auth/register", authHandler.Register)
	mux.HandleFunc("/auth/login", authHandler.Login)

	authMiddleware := middleware.AuthMiddleware(jwtSecret)

	mux.Handle("/auth/me", authMiddleware(http.HandlerFunc(authHandler.Me)))

	return mux
}