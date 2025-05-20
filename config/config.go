package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	Server *Server
	Redis  *Redis
}

type Server struct {
	Port string
}

type Redis struct {
	Addr     string
	Password string
	DB       int
}

func LoadConfig() *Config {
	godotenv.Load()

	var server Server
	server.Port = getEnv("SERVER_PORT", "8080")

	var redis Redis
	redis.Addr = getEnv("REDIS_ADDR", "localhost:6379")
	redis.Password = getEnv("REDIS_PASSWORD", "")
	redis.DB = 0

	return &Config{
		Server: &server,
		Redis:  &redis,
	}
}

func getEnv(key string, defaultValue string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}
	return val
}
