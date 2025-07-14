package db

import (
	"context"
	"fmt"

	"task-controller/internal/logger"

	"github.com/jackc/pgx/v5"
)

type DB struct {
	Conn *pgx.Conn
}

func New(ctx context.Context, url string) (*DB, error) {
	conn, err := pgx.Connect(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("db connect error: %w", err)
	}
	logger.Info("Connected to database")
	return &DB{Conn: conn}, nil
}

func (db *DB) Close(ctx context.Context) {
	if db.Conn != nil {
		db.Conn.Close(ctx)
		logger.Info("Database connection closed")
	}
}

// Получение количество выполненных циклов
func (db *DB) GetOrderCycles(ctx context.Context, orderID int) (int, error) {
	var count int
	err := db.Conn.QueryRow(ctx, `
		SELECT completed_cycles FROM bot.order_stats WHERE order_id=$1
	`, orderID).Scan(&count)
	return count, err
}

// Увеличение количество циклов
func (db *DB) IncrementOrderCycles(ctx context.Context, orderID int, delta int) error {
	_, err := db.Conn.Exec(ctx, `
		UPDATE bot.order_stats
		SET completed_cycles = completed_cycles + $1, updated_at=NOW()
		WHERE order_id=$2
	`, delta, orderID)
	return err
}

// Установка статуса заказа в 'processing'
func (db *DB) SetOrderProcessing(ctx context.Context, orderID int) error {
	_, err := db.Conn.Exec(ctx, `
		UPDATE bot.orders SET status = 'processing' WHERE order_id = $1
	`, orderID)
	if err != nil {
		return fmt.Errorf("failed to update order %d status: %w", orderID, err)
	}
	return nil
}

// Получение заказов со статусом 'new'
func (db *DB) GetNewOrders(ctx context.Context) ([]int, error) {
	rows, err := db.Conn.Query(ctx, `
		SELECT order_id 
		FROM bot.orders 
		WHERE status = 'new'
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to get new orders: %w", err)
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

	return orderIDs, nil
}
