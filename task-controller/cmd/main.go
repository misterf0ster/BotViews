package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"task-controller/config"
	db "task-controller/internal/db"
	queue "task-controller/internal/redis"
)

func main() {
	ctx := context.Background()

	// Загрузка конфигурации
	cfg := config.Load()

	// Подключение к базе
	pgDB, err := db.Connect(ctx, cfg.PostgresDSN())
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer pgDB.Close(ctx)

	// Подключение к очереди
	q, err := queue.Connect(cfg.RedisAddr, cfg.RedisPassword)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer q.Close()

	log.Println("Controller started")

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs
	log.Println("Shutting down gracefully...")
}
