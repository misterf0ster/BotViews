package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"task-controller/config"
	"task-controller/internal/controller"
	"task-controller/internal/db"
	"task-controller/internal/logger"
	queue "task-controller/internal/redis"
)

func main() {
	// Загружаем .env из папки task-controller
	config.LoadEnv()

	cfg := config.Config()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Подключаемся к базе
	dbConn, err := db.New(ctx, cfg.DBaseURL())
	if err != nil {
		logger.Fatal("Failed to connect to DB: %v", err)
	}
	defer dbConn.Close(ctx)

	// Подключаемся к Redis
	q, err := queue.New(cfg.RedisAddr, cfg.RedisPass)
	if err != nil {
		logger.Fatal("Failed to connect to Redis: %v", err)
	}

	ctrl := controller.New(dbConn, q)

	go ctrl.Run(ctx)

	<-ctx.Done()
	logger.Info("Shutting down gracefully")
}
