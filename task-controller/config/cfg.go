package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config содержит все параметры подключения.
type Config struct {
	DBUrl     string
	RedisAddr string
	RedisPass string
}

// Load загружает конфиг из переменных окружения или файла .env.
func Load() (*Config, error) {
	_ = godotenv.Load() // Не fatal — можно запускать и без .env

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dbUrl := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPass, dbHost, dbPort, dbName,
	)

	return &Config{
		DBUrl:     dbUrl,
		RedisAddr: os.Getenv("REDIS_HOST"),
		RedisPass: os.Getenv("REDIS_PASSWORD"),
	}, nil
}

// Validate проверяет обязательные поля конфига.
func (c *Config) Validate() error {
	if c.DBUrl == "" {
		return fmt.Errorf("DBUrl is empty")
	}
	if c.RedisAddr == "" {
		return fmt.Errorf("RedisAddr is empty")
	}
	return nil
}
