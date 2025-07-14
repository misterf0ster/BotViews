package controller

import (
	"context"
	"fmt"
	"math/rand"
	"task-controller/internal/db"
	"task-controller/internal/logger"
	queue "task-controller/internal/redis"
	"time"
)

const (
	TargetCyclesPerOrder = 12500
	MinBotsPerOrder      = 5
	MaxBotsPerOrder      = 10
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

func (c *Controller) Run(ctx context.Context) {
	ticker := time.NewTicker(45 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Info("Controller stopped")
			return
		case <-ticker.C:
			c.ProcessOrders(ctx)
		}
	}
}

func (c *Controller) ProcessOrders(ctx context.Context) {
	orders, err := c.DB.GetNewOrders(ctx)
	if err != nil {
		logger.Error("Failed to get new orders: %v", err)
		return
	}

	processingOrders, err := c.getProcessingOrders(ctx)
	if err != nil {
		logger.Error("Failed to get processing orders: %v", err)
		return
	}

	allOrders := append(orders, processingOrders...)

	if len(allOrders) == 0 {
		logger.Info("No orders to process")
		return
	}

	for _, orderID := range allOrders {
		if err := c.processSingleOrder(ctx, orderID); err != nil {
			logger.Error("Error processing order %d: %v", orderID, err)
		}
	}
}

func (c *Controller) getProcessingOrders(ctx context.Context) ([]int, error) {
	rows, err := c.DB.Conn.Query(ctx, `
		SELECT order_id FROM bot.orders WHERE status = 'processing'
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to get processing orders: %w", err)
	}
	defer rows.Close()

	var orderIDs []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			logger.Error("Failed to scan order_id: %v", err)
			continue
		}
		orderIDs = append(orderIDs, id)
	}
	return orderIDs, rows.Err()
}

func (c *Controller) processSingleOrder(ctx context.Context, orderID int) error {
	currentCycles, err := c.DB.GetOrderCycles(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to get cycles: %w", err)
	}

	if currentCycles >= TargetCyclesPerOrder {
		// Завершаем заказ
		if err := c.DB.SetOrderStatus(ctx, orderID, "complete"); err != nil {
			return fmt.Errorf("failed to set order complete: %w", err)
		}
		logger.Info(fmt.Sprintf("Order %d completed with %d cycles", orderID, currentCycles))
		return nil
	}

	// Проверяем, есть ли уже задачи для этого заказа
	queuedTasks, err := c.Queue.CountTasksForOrder(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to count tasks for order %d: %w", orderID, err)
	}
	if queuedTasks > 0 {
		logger.Info(fmt.Sprintf("Order %d already has %d tasks in queue, skipping", orderID, queuedTasks))
		return nil
	}

	botsCount := rand.Intn(MaxBotsPerOrder-MinBotsPerOrder+1) + MinBotsPerOrder
	cyclesPerBot := 90

	for i := 0; i < botsCount; i++ {
		taskPayload := createTaskPayload(orderID, cyclesPerBot)
		if err := c.Queue.AddTaskForOrder(ctx, orderID, taskPayload); err != nil {
			logger.Error("Failed to add task to queue: %v", err)
		} else {
			logger.Info(fmt.Sprintf("Added task for order %d: %s", orderID, taskPayload))
		}
	}

	// Обновим статус заказа на processing, если был новый
	if err := c.DB.SetOrderStatus(ctx, orderID, "processing"); err != nil {
		logger.Error("Failed to update order status to processing: %v", err)
	}

	return nil
}

func createTaskPayload(orderID int, cycles int) string {
	return fmt.Sprintf("order:%d:cycles:%d:time:%s", orderID, cycles, time.Now().Format(time.RFC3339))
}

// Для отчёта от бота о завершении циклов
func (c *Controller) ReportCycleDone(ctx context.Context, orderID int, cyclesDone int) error {
	err := c.DB.IncrementOrderCycles(ctx, orderID, cyclesDone)
	if err != nil {
		return fmt.Errorf("failed to increment cycles for order %d: %w", orderID, err)
	}
	logger.Info(fmt.Sprintf("Incremented %d cycles for order %d", cyclesDone, orderID))
	return nil
}
