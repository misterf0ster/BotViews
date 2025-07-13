package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
)

type DB struct {
	Psql *pgx.Conn
}

func Open(url string) (*DB, error) {
	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(context.Background()); err != nil {
		conn.Close(context.Background())
		return nil, err
	}

	log.Println("Successfully connected to PostgreSQL")
	return &DB{Psql: conn}, nil
}

func (db *DB) Close() {
	if db.Psql != nil {
		db.Psql.Close(context.Background())
		log.Println("Closed PostgreSQL connection")
	}
}
