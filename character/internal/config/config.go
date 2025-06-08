package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string
	Db   DbConfig
}

type DbConfig struct {
	User     string
	Password string
	Address  string
	Name     string
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
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
