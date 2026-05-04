package storage

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStorage struct {
	Pool *pgxpool.Pool
}

func NewPostgresStorage(ctx context.Context, connString string) (*PostgresStorage, error) {
	pool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("unable to ping the database: %w", err)
	}

	log.Println("Successfully connected to Postgres!")
	return &PostgresStorage{Pool: pool}, nil
}

func (s *PostgresStorage) Close() {
	s.Pool.Close()
}
