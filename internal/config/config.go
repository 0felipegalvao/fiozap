package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port       string
	Address    string
	AdminToken string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	LogLevel   string
	LogType    string
	WADebug    string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		Port:       getEnv("PORT", "8080"),
		Address:    getEnv("ADDRESS", "0.0.0.0"),
		AdminToken: getEnv("ADMIN_TOKEN", generateRandomToken()),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "fiozap"),
		DBPassword: getEnv("DB_PASSWORD", "fiozap123"),
		DBName:     getEnv("DB_NAME", "fiozap"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
		LogLevel:   getEnv("LOG_LEVEL", "info"),
		LogType:    getEnv("LOG_TYPE", "console"),
		WADebug:    getEnv("WA_DEBUG", ""),
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if b, err := strconv.ParseBool(value); err == nil {
			return b
		}
	}
	return defaultValue
}

func generateRandomToken() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 32)
	for i := range b {
		b[i] = charset[i%len(charset)]
	}
	return string(b)
}
