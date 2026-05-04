package storage

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/szymmix/stock-market/internal/models"
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

func (s *PostgresStorage) SetBankState(ctx context.Context, stocks []models.Stock) error {
	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	_, err = tx.Exec(ctx, "DELETE FROM bank_stocks")
	if err != nil {
		return err
	}

	for _, stock := range stocks {
		_, err = tx.Exec(ctx,
			"INSERT INTO bank_stocks (stock_name, quantity) VALUES ($1, $2)",
			stock.Name, stock.Quantity)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (s *PostgresStorage) GetBankStocks(ctx context.Context) ([]models.Stock, error) {
	rows, err := s.Pool.Query(ctx, "SELECT stock_name, quantity FROM bank_stocks ORDER BY stock_name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stocks []models.Stock
	for rows.Next() {
		var st models.Stock
		if err := rows.Scan(&st.Name, &st.Quantity); err != nil {
			return nil, err
		}
		stocks = append(stocks, st)
	}
	return stocks, nil
}

func (s *PostgresStorage) Close() {
	s.Pool.Close()
}
