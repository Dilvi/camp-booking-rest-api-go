package app

import (
	"net/http"

	"github.com/dilvi/camp-booking-rest-api-go/internal/handler"
)

func NewRouter(authHandler *handler.AuthHandler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", handler.HealthHandler)
	mux.HandleFunc("/auth/register", authHandler.Register)

	return mux
}