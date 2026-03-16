package database

import (
	"database/sql"
	"fmt"

	"github.com/dilvi/camp-booking-rest-api-go/internal/config"
	_ "github.com/lib/pq"
)

func NewPostgres(cfg config.Config) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBSSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}