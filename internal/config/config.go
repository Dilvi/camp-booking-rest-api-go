package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort string
}

func Load() Config {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found")
	}
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}
	return Config{AppPort: port}
}