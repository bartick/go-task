package main

import (
	"os"
	"time"

	"github.com/bartick/ringover-task/app/model"
	"github.com/joho/godotenv"
)

var config *model.Configuration

func LoadConfig() error {
	err := godotenv.Load()
	if err != nil {
		log.Error("No .env file found, reading from environment")
	}

	config = &model.Configuration{
		Application: model.ApplicationConfig{
			LogLevel: getEnv("LOG_LEVEL", "info"),
		},
		Server: model.ServerConfig{
			Address:         getEnv("SERVER_ADDRESS", "localhost"),
			Port:            getEnv("SERVER_PORT", "3000"),
			GracefulTimeout: 5 * time.Second,
		},
		Database: model.DatabaseConfig{
			DBUser: getEnv("DB_USER", "root"),
			DBPass: getEnv("DB_PASS", ""),
			DBHost: getEnv("DB_HOST", "localhost"),
			DBPort: getEnv("DB_PORT", "4000"),
			DBName: getEnv("DB_NAME", "tasking"),
		},
	}

	return nil

}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
