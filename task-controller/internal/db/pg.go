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
		return nil, fmt.Errorf("database connection error: %w", err)
	}
	if err := conn.Ping(ctx); err != nil {
		conn.Close(ctx)
		return nil, fmt.Errorf("error checking database connection: %w", err)
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

// Получить новые заказы со статусом "new"
func (db *DB) GetNewOrders(ctx context.Context) ([]int, error) {
	rows, err := db.Conn.Query(ctx, `
		SELECT order_id FROM bot.orders WHERE status = 'new'
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
	return orderIDs, rows.Err()
}

// Установить статус заказа
func (db *DB) SetOrderStatus(ctx context.Context, orderID int, status string) error {
	_, err := db.Conn.Exec(ctx, `
		UPDATE bot.orders SET status = $1 WHERE order_id = $2
	`, status, orderID)
	if err != nil {
		return fmt.Errorf("failed to update order %d status: %w", orderID, err)
	}
	return nil
}

// Получить количество выполненных циклов для заказа
func (db *DB) GetOrderCycles(ctx context.Context, orderID int) (int, error) {
	var count int
	err := db.Conn.QueryRow(ctx, `
		SELECT completed_cycles FROM bot.order_stats WHERE order_id = $1
	`, orderID).Scan(&count)
	if err != nil {
		// Если записи нет, вернем 0
		if err == pgx.ErrNoRows {
			return 0, nil
		}
		return 0, fmt.Errorf("failed to get completed_cycles for order %d: %w", orderID, err)
	}
	return count, nil
}

// Увеличить количество циклов для заказа
func (db *DB) IncrementOrderCycles(ctx context.Context, orderID int, delta int) error {
	// Если записи в order_stats нет, создадим
	_, err := db.Conn.Exec(ctx, `
		INSERT INTO bot.order_stats (order_id, completed_cycles, updated_at)
		VALUES ($1, $2, NOW())
		ON CONFLICT (order_id) DO UPDATE
		SET completed_cycles = bot.order_stats.completed_cycles + $2,
			updated_at = NOW()
	`, orderID, delta)
	return err
}
