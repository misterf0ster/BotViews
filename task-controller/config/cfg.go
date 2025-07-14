package config

import (
	"os"
	"path/filepath"

	"task-controller/internal/logger"

	"github.com/joho/godotenv"
)

type DbaseCfg struct {
	DBUser    string
	DBPass    string
	DBHost    string
	DBPort    string
	DBName    string
	RedisAddr string
	RedisPass string
}

func EnvLoad(key string) string {
	value := os.Getenv(key)
	if value == "" {
		logger.Fatal("env variable %s is not set", key)
	}
	return value
}

func LoadEnv() {
	// Загружаем .env из папки task-controller
	envPath := filepath.Join("task-controller", ".env")
	if err := godotenv.Load(envPath); err != nil {
		logger.Warn("Warning: error loading .env file from %s: %v", envPath, err)
	}
}

func Config() *DbaseCfg {
	return &DbaseCfg{
		DBUser:    EnvLoad("DB_USER"),
		DBPass:    EnvLoad("DB_PASSWORD"),
		DBHost:    EnvLoad("DB_HOST"),
		DBPort:    EnvLoad("DB_PORT"),
		DBName:    EnvLoad("DB_NAME"),
		RedisAddr: EnvLoad("REDIS_ADDR"),
		RedisPass: EnvLoad("REDIS_PASSWORD"),
	}
}

func (c *DbaseCfg) DBaseURL() string {
	return "postgres://" + c.DBUser + ":" + c.DBPass + "@" + c.DBHost + ":" + c.DBPort + "/" + c.DBName + "?sslmode=disable"
}
