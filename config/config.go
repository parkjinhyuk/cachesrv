package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	ServerPort    string
	RedisAddr     string
	RedisPassword string
	RedisDB       int
}

func LoadConfig() *Config {
	godotenv.Load()
	port := getEnv("SERVER_PORT", "8080")
	redisAddr := getEnv("REDIS_ADDR", "localhost:6379")
	redisPassword := getEnv("REDIS_PASSWORD", "")
	redisDB := 0

	return &Config{
		ServerPort:    port,
		RedisAddr:     redisAddr,
		RedisPassword: redisPassword,
		RedisDB:       redisDB,
	}
}

func getEnv(key string, defaultValue string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}
	return val
}
