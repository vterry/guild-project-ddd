package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Addr string
	Db   DbConfig
	Auth KeycloakConfig
}

type DbConfig struct {
	User     string
	Password string
	Address  string
	Name     string
}

type KeycloakConfig struct {
	BaseURL      string
	ClientID     string
	ClientSecret string
	Realm        string
}

var Envs = initConfig()

func initConfig() Config {

	if err := godotenv.Load(); err != nil {
		fmt.Printf("Warning: .env file not found: %v\n", err)
	}

	return Config{
		Addr: getEnv("APP_ADDR", ":8080"),
		Db: DbConfig{
			User:     getEnv("DB_USER", "character"),
			Password: getEnv("DB_PASSWORD", "characterPW"),
			Address:  fmt.Sprintf("%s:%s", getEnv("DB_HOST", "127.0.0.1"), getEnv("DB_PORT", "3306")),
			Name:     getEnv("DB_NAME", "character-db"),
		},
		Auth: KeycloakConfig{
			BaseURL:      getEnv("AUTH_BASE_URL", "http://localhost:7080"),
			ClientID:     getEnv("AUTH_CLIENT_ID", "playground"),
			ClientSecret: getEnv("AUTH_CLIENT_SECRET", "PAfdvjPnUDFyTqm5fBuqjHxiAJCGQLVu"),
			Realm:        getEnv("AUTH_REALM", "playground"),
		},
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
