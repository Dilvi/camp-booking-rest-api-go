package app

import (
	"database/sql"
	"net/http"

	"github.com/dilvi/camp-booking-rest-api-go/internal/config"
	"github.com/dilvi/camp-booking-rest-api-go/internal/database"
)

type App struct {
	Config config.Config
	Router *http.ServeMux
	DB *sql.DB
}

func New(cfg config.Config) (*App, error) {
	db, err := database.NewPostgres(cfg)
	if err != nil {
		return nil, err
	}

	router := NewRouter()

	return &App{
		Config: cfg,
		Router: router,
		DB: db,
	}, nil
}