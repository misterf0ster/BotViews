package redis

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

type Queue struct {
	Client *redis.Client
}

func Connect(addr, password string) (*Queue, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})

	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	log.Println("Successfully connected to Redis")
	return &Queue{Client: rdb}, nil
}

func (q *Queue) Close() error {
	return q.Client.Close()
}
