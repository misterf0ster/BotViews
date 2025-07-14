package queue

import (
	"context"

	"task-controller/internal/logger"

	"github.com/redis/go-redis/v9"
)

type Queue struct {
	Client *redis.Client
}

func New(addr, pass string) *Queue {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
		DB:       0,
	})
	logger.Info("Connected to Redis")
	return &Queue{Client: client}
}

// Добавить задачу
func (q *Queue) AddTask(ctx context.Context, task string) error {
	return q.Client.RPush(ctx, "tasks", task).Err()
}

// Получить задачу
func (q *Queue) GetTask(ctx context.Context) (string, error) {
	return q.Client.LPop(ctx, "tasks").Result()
}
