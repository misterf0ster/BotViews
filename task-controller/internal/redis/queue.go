package queue

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type Queue struct {
	Client *redis.Client
}

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
	return &Queue{Client: client}, nil
}

func (q *Queue) Close() error {
	return q.Client.Close()
}

// Добавить задачу в очередь для конкретного заказа
func (q *Queue) AddTaskForOrder(ctx context.Context, orderID int, task string) error {
	key := fmt.Sprintf("tasks:order:%d", orderID)
	return q.Client.RPush(ctx, key, task).Err()
}

// Получить задачу для конкретного заказа
func (q *Queue) GetTaskForOrder(ctx context.Context, orderID int) (string, error) {
	key := fmt.Sprintf("tasks:order:%d", orderID)
	return q.Client.LPop(ctx, key).Result()
}

// Получить количество задач для заказа
func (q *Queue) CountTasksForOrder(ctx context.Context, orderID int) (int64, error) {
	key := fmt.Sprintf("tasks:order:%d", orderID)
	return q.Client.LLen(ctx, key).Result()
}
