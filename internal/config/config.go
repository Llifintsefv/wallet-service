package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBConnStr string
	Port      string
}

func NewConfig() (*Config, error) {
	err := godotenv.Load("config.env")
	if err != nil {
		slog.Error("DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME, DB_SSL_MODE must be set")
		return nil, fmt.Errorf("missing required environment variables")
	}
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbSSLMode := os.Getenv("DB_SSL_MODE")

	if dbHost == "" || dbPort == "" || dbUser == "" || dbPassword == "" || dbName == "" || dbSSLMode == "" {
		slog.Error("DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME, DB_SSL_MODE must be set")
		return nil, fmt.Errorf("missing required environment variables")
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", dbUser, dbPassword, dbHost, dbPort, dbName, dbSSLMode)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}
	port = ":" + port

	return &Config{
		DBConnStr: connStr,
		Port:      port,
	}, nil
}
