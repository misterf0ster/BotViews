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
	MaxBotsPerOrder      = 20
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
	ticker := time.NewTicker(5 * time.Second)
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

// Основная логика обработки заказов
func (c *Controller) ProcessOrders(ctx context.Context) {
	orders, err := c.DB.GetNewOrders(ctx)
	if err != nil {
		logger.Error("Failed to get new orders: %v", err)
		return
	}

	// Добавим уже обрабатываемые заказы (status='processing'), чтобы динамически распределять задачи
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

// Получить список заказов в статусе 'processing'
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

	// Рандомное количество ботов для работы с заказом (5-20)
	botsCount := rand.Intn(MaxBotsPerOrder-MinBotsPerOrder+1) + MinBotsPerOrder

	// Количество циклов, которое каждый бот должен выполнить за час
	// Так как 1 цикл ~40 сек, за час бот делает примерно 90 циклов
	// Распределяем: общее количество циклов на заказ в час делим на кол-во ботов
	// Для простоты положим каждый бот делает 90 циклов в час
	cyclesPerBot := 90

	// Создаём задачи в очереди для каждого бота
	for i := 0; i < botsCount; i++ {
		taskPayload := createTaskPayload(orderID, cyclesPerBot)
		if err := c.Queue.AddTask(ctx, taskPayload); err != nil {
			logger.Error("Failed to add task to queue: %v", err)
		} else {
			logger.Info(fmt.Sprintf("Added task for order %d: %s", orderID, taskPayload))
		}
	}

	// Обновим статус, если заказ был новый (для новых)
	if err := c.DB.SetOrderStatus(ctx, orderID, "processing"); err != nil {
		logger.Error("Failed to update order status to processing: %v", err)
	}

	return nil
}

func createTaskPayload(orderID int, cycles int) string {
	return fmt.Sprintf("order:%d:cycles:%d:time:%s", orderID, cycles, time.Now().Format(time.RFC3339))
}

// Функция для обработки отчёта от бота (пример)
// Бот после выполнения цикла отправляет сюда результат для инкремента completed_cycles
func (c *Controller) ReportCycleDone(ctx context.Context, orderID int, cyclesDone int) error {
	err := c.DB.IncrementOrderCycles(ctx, orderID, cyclesDone)
	if err != nil {
		return fmt.Errorf("failed to increment cycles for order %d: %w", orderID, err)
	}
	logger.Info(fmt.Sprintf("Incremented %d cycles for order %d", cyclesDone, orderID))
	return nil
}
