package app

import (
	"net/http"

	"github.com/dilvi/camp-booking-rest-api-go/internal/handler"
	"github.com/dilvi/camp-booking-rest-api-go/internal/middleware"
)

func NewRouter(authHandler *handler.AuthHandler, childHandler *handler.ChildHandler, jwtSecret string) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", handler.HealthHandler)
	mux.HandleFunc("/auth/register", authHandler.Register)
	mux.HandleFunc("/auth/login", authHandler.Login)

	authMiddleware := middleware.AuthMiddleware(jwtSecret)

	mux.Handle("/auth/me", authMiddleware(http.HandlerFunc(authHandler.Me)))
	mux.Handle("/children", authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			childHandler.Create(w, r)
		case http.MethodGet:
			childHandler.List(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})))
	mux.Handle("/children/", authMiddleware(http.HandlerFunc(childHandler.Update)))

	return mux
}