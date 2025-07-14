package controller

import (
	"context"
	"fmt"
	"task-controller/internal/db"
	"task-controller/internal/logger"
	queue "task-controller/internal/redis"
	"time"
)

type Controller struct {
	DB    *db.DB
	Queue *queue.Queue
}

func New(database *db.DB, q *queue.Queue) *Controller {
	return &Controller{
		DB:    database,
		Queue: q,
	}
}

// Метод для постановки новых заказов в очередь и смены их статуса на processing
func (c *Controller) EnqueueNewOrders(ctx context.Context) {
	orders, err := c.DB.GetNewOrders(ctx)
	if err != nil {
		logger.Error("Failed to get new orders: %v", err)
		return
	}

	if len(orders) == 0 {
		logger.Info("No new orders found")
		return
	}

	for _, orderID := range orders {
		task := createTaskPayload(orderID)

		// Добавляем в очередь
		err := c.Queue.AddTask(ctx, task)
		if err != nil {
			logger.Error("Failed to add task to queue: %v", err)
			continue
		}

		// Меняем статус на processing
		if err := c.DB.SetOrderProcessing(ctx, orderID); err != nil {
			logger.Error("Failed to set order %d status to processing: %v", orderID, err)
			continue
		}

		logger.Info("Added order %d to queue and set status to processing", orderID)
	}
}

func createTaskPayload(orderID int) string {
	return fmt.Sprintf("order:%d:time:%s", orderID, time.Now().Format(time.RFC3339))
}
