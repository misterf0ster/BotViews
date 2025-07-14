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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		logger.Info("Received shutdown signal")
		cancel()
	}()

	cfg, err := config.Load()
	if err != nil {
		logger.Error("Config load error: %v", err)
		os.Exit(1)
	}
	if err := cfg.Validate(); err != nil {
		logger.Error("Config validation error: %v", err)
		os.Exit(1)
	}

	dbConn, err := db.New(ctx, cfg.DBUrl)
	if err != nil {
		logger.Error("DB connection error: %v", err)
		os.Exit(1)
	}
	defer dbConn.Close(ctx)

	q, err := queue.New(cfg.RedisAddr, cfg.RedisPass)
	if err != nil {
		logger.Error("Redis connection error: %v", err)
		os.Exit(1)
	}
	defer q.Close()

	ctrl := controller.New(dbConn, q)
	ctrl.Run(ctx)
}
