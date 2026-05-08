package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort                    string
	DBHost                     string
	DBPort                     string
	DBUser                     string
	DBPassword                 string
	DBName                     string
	DBSSLMode                  string
	JWTSecret                  string
	GigaChatAuthKey            string
	GigaChatScope              string
	GigaChatModel              string
	GigaChatOAuthURL           string
	GigaChatCompletionsURL     string
	GigaChatInsecureSkipVerify bool
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

	gigaChatScope := os.Getenv("GIGACHAT_SCOPE")
	if gigaChatScope == "" {
		gigaChatScope = "GIGACHAT_API_PERS"
	}
	gigaChatModel := os.Getenv("GIGACHAT_MODEL")
	if gigaChatModel == "" {
		gigaChatModel = "GigaChat-2"
	}
	gigaChatOAuthURL := os.Getenv("GIGACHAT_OAUTH_URL")
	if gigaChatOAuthURL == "" {
		gigaChatOAuthURL = "https://ngw.devices.sberbank.ru:9443/api/v2/oauth"
	}
	gigaChatCompletionsURL := os.Getenv("GIGACHAT_COMPLETIONS_URL")
	if gigaChatCompletionsURL == "" {
		gigaChatCompletionsURL = "https://gigachat.devices.sberbank.ru/api/v1/chat/completions"
	}
	gigaChatInsecureSkipVerify, _ := strconv.ParseBool(os.Getenv("GIGACHAT_INSECURE_SKIP_VERIFY"))

	return Config{
		AppPort:                    port,
		DBHost:                     DBhost,
		DBPort:                     DBport,
		DBUser:                     DBuser,
		DBPassword:                 DBpassword,
		DBName:                     DBname,
		DBSSLMode:                  DBsslMode,
		JWTSecret:                  JWTSecret,
		GigaChatAuthKey:            os.Getenv("GIGACHAT_AUTH_KEY"),
		GigaChatScope:              gigaChatScope,
		GigaChatModel:              gigaChatModel,
		GigaChatOAuthURL:           gigaChatOAuthURL,
		GigaChatCompletionsURL:     gigaChatCompletionsURL,
		GigaChatInsecureSkipVerify: gigaChatInsecureSkipVerify,
	}
}
