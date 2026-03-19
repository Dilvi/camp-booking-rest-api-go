package app

import (
	"database/sql"
	"net/http"

	"github.com/dilvi/camp-booking-rest-api-go/internal/config"
	"github.com/dilvi/camp-booking-rest-api-go/internal/database"
	"github.com/dilvi/camp-booking-rest-api-go/internal/handler"
	"github.com/dilvi/camp-booking-rest-api-go/internal/repository/postgres"
	"github.com/dilvi/camp-booking-rest-api-go/internal/service"
)

type App struct {
	Config config.Config
	Router http.Handler
	DB     *sql.DB
}

func New(cfg config.Config) (*App, error) {
	db, err := database.NewPostgres(cfg)
	if err != nil {
		return nil, err
	}

	userRepo := postgres.NewUserRepository(db)
	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	authHandler := handler.NewAuthHandler(authService)

	profileService := service.NewProfileService(userRepo)
	profileHandler := handler.NewProfileHandler(profileService)

	childRepo := postgres.NewChildRepository(db)
	childService := service.NewChildService(childRepo)
	childHandler := handler.NewChildHandler(childService)

	campRepo := postgres.NewCampRepository(db)
	campService := service.NewCampService(campRepo)
	campHandler := handler.NewCampHandler(campService)

	favoriteRepo := postgres.NewFavoriteRepository(db)
	favoriteService := service.NewFavoriteService(favoriteRepo)
	favoriteHandler := handler.NewFavoriteHandler(favoriteService)

	bookingRepo := postgres.NewBookingRepository(db)
	bookingService := service.NewBookingService(bookingRepo, childRepo, campRepo)
	bookingHandler := handler.NewBookingHandler(bookingService)

	router := NewRouter(authHandler, profileHandler, childHandler, campHandler, favoriteHandler, bookingHandler, cfg.JWTSecret)

	return &App{
		Config: cfg,
		Router: router,
		DB:     db,
	}, nil
}