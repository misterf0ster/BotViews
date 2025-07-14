package queue

import (
	"context"
	"fmt"
	"task-controller/internal/logger"

	"github.com/redis/go-redis/v9"
)

// Queue реализует очередь задач в Redis.
type Queue struct {
	Client *redis.Client
}

// New подключается к Redis.
func New(addr, pass string) (*Queue, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
		DB:       0,
	})

	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}
	logger.Info("Connected to Redis")
	return &Queue{Client: client}, nil
}

// Close закрывает соединение с Redis.
func (q *Queue) Close() error {
	return q.Client.Close()
}

// AddTask добавляет задачу в очередь.
func (q *Queue) AddTask(ctx context.Context, task string) error {
	return q.Client.RPush(ctx, "tasks", task).Err()
}

// GetTask извлекает задачу из очереди.
func (q *Queue) GetTask(ctx context.Context) (string, error) {
	return q.Client.LPop(ctx, "tasks").Result()
}
