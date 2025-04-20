package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	PublicHost           string
	Port                 string
	MONGOUser            string
	MONGOPass            string
	MONGODB              string
	MONGOAddress         string
	AccessDuration       int64
	RefreshTokenDuration int64
	CSRFTokenLength      int64
	JWTSecret            string
	UserServiceURL       string
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()

	return Config{
		PublicHost:           getEnv("PUBLIC_HOST", "127.0.0.1"),
		Port:                 getEnv("PORT", ":8080"),
		MONGOUser:            getEnv("MONGO_USER", "admin"),
		MONGOPass:            getEnv("MONGO_PASS", "password"),
		MONGODB:              getEnv("MONGO_DB", "AuthServer"),
		MONGOAddress:         fmt.Sprintf("%s:%s", getEnv("MONGO_HOST", "localhost"), getEnv("MONGO_PORT", "27017")),
		AccessDuration:       getEnvAsInt("ACCESS_DURATION", 10),
		RefreshTokenDuration: getEnvAsInt("REFRESH_EXPIRATION", 10),
		CSRFTokenLength:      getEnvAsInt("CSRF_TOKEN_LENGTH", 128),
		JWTSecret:            getEnv("JWT_SECRET", "secret"),
		UserServiceURL:       getEnv("USER_SERVICE_URL", "http://localhost:8080"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fallback
		}
		return i
	}
	return fallback
}
