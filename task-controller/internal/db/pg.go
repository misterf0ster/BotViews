package db

import (
	"context"
	"fmt"
	"task-controller/internal/logger"

	"github.com/jackc/pgx/v5"
)

// DB содержит прямое соединение с Postgres (без пула).
type DB struct {
	Conn *pgx.Conn
}

// New открывает соединение с Postgres.
func New(ctx context.Context, url string) (*DB, error) {
	conn, err := pgx.Connect(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("database connection error: %w", err)
	}
	if err := conn.Ping(ctx); err != nil {
		conn.Close(ctx)
		return nil, fmt.Errorf("error checking database connection: %w", err)
	}
	logger.Info("Connected to database")
	return &DB{Conn: conn}, nil
}

// Close закрывает соединение.
func (db *DB) Close(ctx context.Context) {
	if db.Conn != nil {
		db.Conn.Close(ctx)
		logger.Info("Database connection closed")
	}
}

// GetOrderCycles возвращает число выполненных циклов заказа.
func (db *DB) GetOrderCycles(ctx context.Context, orderID int) (int, error) {
	var count int
	err := db.Conn.QueryRow(ctx, `
		SELECT completed_cycles FROM bot.order_stats WHERE order_id=$1
	`, orderID).Scan(&count)
	return count, err
}

// IncrementOrderCycles увеличивает число циклов заказа.
func (db *DB) IncrementOrderCycles(ctx context.Context, orderID int, delta int) error {
	_, err := db.Conn.Exec(ctx, `
		UPDATE bot.order_stats
		SET completed_cycles = completed_cycles + $1, updated_at=NOW()
		WHERE order_id=$2
	`, delta, orderID)
	return err
}

// SetOrderProcessing переводит заказ в статус 'processing'.
func (db *DB) SetOrderProcessing(ctx context.Context, orderID int) error {
	_, err := db.Conn.Exec(ctx, `
		UPDATE bot.orders SET status = 'processing' WHERE order_id = $1
	`, orderID)
	if err != nil {
		return fmt.Errorf("failed to update order %d status: %w", orderID, err)
	}
	return nil
}

// GetNewOrders возвращает id заказов в статусе 'new'.
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
	if rows.Err() != nil {
		return orderIDs, rows.Err()
	}

	return orderIDs, nil
}
