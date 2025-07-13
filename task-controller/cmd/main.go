package main

import (
	"context"
	"log"
)

func main() {
	cfg := config.LoadConfig()

	rdb := redis.NewClient(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)
	if err := redis.Ping(context.Background(), rdb); err != nil {
		log.Fatalf("Redis connect error: %v", err)
	}
	log.Println("Connected to Redis")

	ctrl := controller.NewController(rdb, cfg.TaskQueue)
	ctrl.Run()
}
