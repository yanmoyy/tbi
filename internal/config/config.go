package config

import (
	"fmt"
	"os"
	"strings"

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

type GraphQL struct {
	IndexerURLs []string
}

type Config struct {
	DB      DB
	GraphQL GraphQL
}

func Load() *Config {
	return LoadWithPath(".env")
}

func LoadWithPath(path string) *Config {
	err := godotenv.Load(path)
	if err != nil {
		panic("failed to load .env file")
	}

	// DB
	dbHost := ensureEnv("DB_HOST")
	dbPort := ensureEnv("DB_PORT")
	dbUser := ensureEnv("DB_USER")
	dbPassword := ensureEnv("DB_PASSWORD")
	dbName := ensureEnv("DB_NAME")
	dbSSLMode := ensureEnv("DB_SSL_MODE")

	// GraphQL
	urls := ensureEnv("GRAPHQL_INDEXER_URLS")
	indexerURLs := strings.Split(urls, ",")

	return &Config{
		DB: DB{
			Host:     dbHost,
			Port:     dbPort,
			User:     dbUser,
			Password: dbPassword,
			Name:     dbName,
			SSLMode:  dbSSLMode,
		},
		GraphQL: GraphQL{
			IndexerURLs: indexerURLs,
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
