package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config groups environment-driven configuration for the application.
type Config struct {
	AppPort   string
	JWTSecret string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       string
}

// Load returns the application configuration, loading .env if present.
func Load() Config {
	_ = godotenv.Load()

	cfg := Config{
		AppPort:       getEnv("APP_PORT", "8080"),
		JWTSecret:     getEnv("JWT_SECRET", "dev_secret_change_me"),
		DBHost:        getEnv("DB_HOST", ""),
		DBPort:        getEnv("DB_PORT", "5432"),
		DBUser:        getEnv("DB_USER", ""),
		DBPassword:    getEnv("DB_PASSWORD", ""),
		DBName:        getEnv("DB_NAME", ""),
		DBSSLMode:     getEnv("DB_SSLMODE", "disable"),
		RedisHost:     getEnv("REDIS_HOST", "localhost"),
		RedisPort:     getEnv("REDIS_PORT", "6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       getEnv("REDIS_DB", "0"),
	}

	if cfg.JWTSecret == "dev_secret_change_me" {
		log.Println("warning: using default JWT secret; set JWT_SECRET in environment for production")
	}

	return cfg
}

func getEnv(key, def string) string {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	return val
}
