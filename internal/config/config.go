package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort string
	DBHost string
	DBPort string
	DBUser string
	DBPassword string
	DBName string
	DBSSLMode string
	JWTSecret string
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
	DBhost := os.Getenv("DB_HOST")
	DBport := os.Getenv("DB_PORT")
	DBuser := os.Getenv("DB_USER")
	DBpassword := os.Getenv("DB_PASSWORD")
	DBname := os.Getenv("DB_NAME")
	DBsslMode := os.Getenv("DB_SSLMODE")
	JWTSecret := os.Getenv("JWT_SECRET")
	return Config{AppPort: port, DBHost: DBhost, DBPort: DBport, DBUser: DBuser, DBPassword: DBpassword, DBName: DBname, DBSSLMode: DBsslMode, JWTSecret: JWTSecret}
}