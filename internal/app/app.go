package app

import (
	"net/http"

	"github.com/dilvi/camp-booking-rest-api-go/internal/config"
)

type App struct {
	Config config.Config
	Router *http.ServeMux
}

func New(cfg config.Config) *App {
	router := NewRouter()
	return &App{
		Config: cfg,
		Router: router,
	}
}