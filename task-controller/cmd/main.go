package main

import (
	"context"
	"task-controller/config"
	"task-controller/internal/controller"
	"task-controller/internal/db"
	queue "task-controller/internal/redis"
)

func main() {
	ctx := context.Background()

	cfg := config.Load()
	dbConn, err := db.New(ctx, cfg.DBUrl)
	if err != nil {
		panic(err)
	}
	defer dbConn.Close(ctx)

	q := queue.New(cfg.RedisAddr, cfg.RedisPass)

	ctrl := controller.New(dbConn, q)
	ctrl.Run(ctx)
}
