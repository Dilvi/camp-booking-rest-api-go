package app

import (
	"net/http"

	"github.com/dilvi/camp-booking-rest-api-go/internal/handler"
	"github.com/dilvi/camp-booking-rest-api-go/internal/middleware"
)

func NewRouter(authHandler *handler.AuthHandler, childHandler *handler.ChildHandler, campHandler *handler.CampHandler, favoriteHandler *handler.FavoriteHandler,bookingHandler *handler.BookingHandler, jwtSecret string) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", handler.HealthHandler)
	mux.HandleFunc("/auth/register", authHandler.Register)
	mux.HandleFunc("/auth/login", authHandler.Login)
	mux.HandleFunc("/camps", campHandler.List)
	mux.HandleFunc("/camps/", campHandler.GetByID)

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

	mux.Handle("/favorites", authMiddleware(http.HandlerFunc(favoriteHandler.List)))

	mux.Handle("/favorites/", authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			favoriteHandler.Add(w, r)
		case http.MethodDelete:
			favoriteHandler.Remove(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle("/bookings", authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		bookingHandler.Create(w, r)
	case http.MethodGet:
		bookingHandler.List(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
	})))

	handler := middleware.LoggingMiddleware(mux)
	handler = middleware.RecoveryMiddleware(handler)

	return handler
}