package storage

import (
	"context"

	"github.com/szymmix/stock-market/internal/models"
)

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
