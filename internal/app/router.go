package app

import (
	"net/http"
	"strings"

	"github.com/dilvi/camp-booking-rest-api-go/internal/handler"
	"github.com/dilvi/camp-booking-rest-api-go/internal/middleware"
)

func NewRouter(authHandler *handler.AuthHandler, profileHandler *handler.ProfileHandler, childHandler *handler.ChildHandler, campHandler *handler.CampHandler, favoriteHandler *handler.FavoriteHandler, bookingHandler *handler.BookingHandler, bikeRouteHandler *handler.BikeRouteHandler, jwtSecret string) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", handler.HealthHandler)
	mux.HandleFunc("/auth/register", authHandler.Register)
	mux.HandleFunc("/auth/login", authHandler.Login)
	mux.HandleFunc("/camps", campHandler.List)
	mux.HandleFunc("/camps/", campHandler.GetByID)
	mux.HandleFunc("/bike-routes", bikeRouteHandler.List)

	authMiddleware := middleware.AuthMiddleware(jwtSecret)

	mux.Handle("/bike-routes/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/ai-analysis") {
			authMiddleware(http.HandlerFunc(bikeRouteHandler.AIAnalysis)).ServeHTTP(w, r)
			return
		}

		bikeRouteHandler.GetByID(w, r)
	}))

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
	mux.Handle("/children/", authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			childHandler.Update(w, r)
		case http.MethodDelete:
			childHandler.Delete(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})))

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
	mux.Handle("/bookings/", authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && strings.HasSuffix(r.URL.Path, "/bike-routes"):
			bookingHandler.ListBikeRoutes(w, r)
		case r.Method == http.MethodPost && strings.Contains(r.URL.Path, "/bike-routes/") && strings.HasSuffix(r.URL.Path, "/ai-analysis"):
			bookingHandler.BikeRouteAIAnalysis(w, r)
		case r.Method == http.MethodGet:
			bookingHandler.GetByID(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle("/profile", authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			profileHandler.Get(w, r)
		case http.MethodPut:
			profileHandler.Update(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle("/profile/password", authMiddleware(http.HandlerFunc(profileHandler.UpdatePassword)))

	handler := middleware.LoggingMiddleware(mux)
	handler = middleware.RecoveryMiddleware(handler)

	return handler
}
