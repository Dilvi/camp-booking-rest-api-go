package app

import (
	"net/http"

	"github.com/dilvi/camp-booking-rest-api-go/internal/handler"
)

func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", handler.HealthHandler)
	return mux
}