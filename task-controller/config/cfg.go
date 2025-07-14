package config

import (
	"net/url"
	"os"

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
	envPath := ".env"
	if err := godotenv.Load(envPath); err != nil {
		logger.Warn("Warning: error loading .env file from %s: %v", envPath, err)
	}
}

func Config() *DbaseCfg {
	user := url.QueryEscape(EnvLoad("DB_USER"))
	pass := url.QueryEscape(EnvLoad("DB_PASSWORD"))
	host := url.QueryEscape(EnvLoad("DB_HOST"))
	port := url.QueryEscape(EnvLoad("DB_PORT"))
	dbname := url.QueryEscape(EnvLoad("DB_NAME"))

	return &DbaseCfg{
		DBUser:    user,
		DBPass:    pass,
		DBHost:    host,
		DBPort:    port,
		DBName:    dbname,
		RedisAddr: EnvLoad("REDIS_ADDR"),
		RedisPass: EnvLoad("REDIS_PASSWORD"),
	}
}

func (c *DbaseCfg) DBaseURL() string {
	return "postgres://" + c.DBUser + ":" + c.DBPass + "@" + c.DBHost + ":" + c.DBPort + "/" + c.DBName + "?sslmode=disable"
}
