package queue

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type Queue struct {
	client *redis.Client
}

var ctx = context.Background()

func New(addr, pass string) *Queue {
	return &Queue{
		client: redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: pass,
			DB:       0,
		}),
	}
}

func (q *Queue) PopTask(queues ...string) (string, error) {
	result, err := q.client.BLPop(ctx, 0, queues...).Result()
	if err != nil || len(result) < 2 {
		return "", err
	}
	return result[1], nil
}
