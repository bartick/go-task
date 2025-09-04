package model

import "time"

type ApplicationConfig struct {
	LogLevel string
}

type DatabaseConfig struct {
	DBUser string
	DBPass string
	DBHost string
	DBName string
	DBPort string
}

type ServerConfig struct {
	Address         string
	Port            string
	GracefulTimeout time.Duration
}

type Configuration struct {
	Application ApplicationConfig
	Database    DatabaseConfig
	Server      ServerConfig
}
