package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUrl     string
	RedisAddr string
	RedisPass string
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env: %v", err)
	}

	dbUrl := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	return &Config{
		DBUrl:     dbUrl,
		RedisAddr: os.Getenv("REDIS_HOST"),
		RedisPass: os.Getenv("REDIS_PASSWORD"),
	}
}
