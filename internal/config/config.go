package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type DB struct {
	Host     string
	User     string
	Password string
	Name     string
	Port     string
	SSLMode  string
}

type Config struct {
	DB *DB
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		panic("failed to load .env file")
	}

	dbHost := ensureEnv("DB_HOST")
	dbPort := ensureEnv("DB_PORT")
	dbUser := ensureEnv("DB_USER")
	dbPassword := ensureEnv("DB_PASSWORD")
	dbName := ensureEnv("DB_NAME")
	dbSSLMode := ensureEnv("DB_SSL_MODE")

	return &Config{
		DB: &DB{
			Host:     dbHost,
			Port:     dbPort,
			User:     dbUser,
			Password: dbPassword,
			Name:     dbName,
			SSLMode:  dbSSLMode,
		},
	}
}

func ensureEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("ensureEnv: %s is not set", key))
	}
	return value
}
