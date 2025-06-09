package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string
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
	godotenv.Load()
	return Config{
		Port: getEnv("APP_PORT", "8080"),
		Db: DbConfig{
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", "password"),
			Address:  fmt.Sprintf("%s:%s", getEnv("DB_HOST", "127.0.0.1"), getEnv("DB_PORT", "3306")),
			Name:     getEnv("DB_NAME", "character-db"),
		},
		Auth: KeycloakConfig{
			BaseURL:      getEnv("AUTH_BASE_URL", "http://localhost:7080"),
			ClientID:     getEnv("AUTH_CLIENT_ID", "playground"),
			ClientSecret: getEnv("AUTH_CLIENT_SECRET", "bsseWxYIokXVvfgQe0PU6dkXy24Hj4no"),
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
