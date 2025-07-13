package config

import (
	"fmt"
	"os"

	"log"

	"github.com/joho/godotenv"
)

type DbaseCfg struct {
	DBUser string
	DBPass string
	DBHost string
	DBPort string
	DBName string

	RedisAddr     string
	RedisPassword string
}

// Загрузка env
func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: .env file not found or could not be loaded: %v", err)
	}
}

// Получить переменную из env
func EnvLoad(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("env variable %s is not set", key)
	}
	return value
}

func Config() *DbaseCfg {
	return &DbaseCfg{
		DBUser:        EnvLoad("DB_USER"),
		DBPass:        EnvLoad("DB_PASSWORD"),
		DBHost:        EnvLoad("DB_HOST"),
		DBPort:        EnvLoad("DB_PORT"),
		DBName:        EnvLoad("DB_NAME"),
		RedisAddr:     EnvLoad("REDIS_ADDR"),
		RedisPassword: EnvLoad("REDIS_PASSWORD"),
	}
}

func (c *DbaseCfg) PostgresURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", c.DBUser, c.DBPass, c.DBHost, c.DBPort, c.DBName)
}
