package storage

import (
	"context"

	"github.com/szymmix/stock-market/internal/models"
)

func (s *PostgresStorage) ExecuteTrade(ctx context.Context, walletID, stockName, tradeType string) error {
	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var bankQty int
	err = tx.QueryRow(ctx, "SELECT quantity FROM bank_stocks WHERE stock_name = $1", stockName).Scan(&bankQty)
	if err != nil {
		return ErrStockNotFound
	}

	_, err = tx.Exec(ctx, "INSERT INTO wallets (id) VALUES ($1) ON CONFLICT (id) DO NOTHING", walletID)
	if err != nil {
		return err
	}

	if tradeType == "buy" {
		if bankQty < 1 {
			return ErrInsufficientStock
		}

		_, err = tx.Exec(ctx, "UPDATE bank_stocks SET quantity = quantity - 1 WHERE stock_name = $1", stockName)
		if err != nil {
			return err
		}

		_, err = tx.Exec(ctx, `
			INSERT INTO wallet_stocks (wallet_id, stock_name, quantity) 
			VALUES ($1, $2, 1) 
			ON CONFLICT (wallet_id, stock_name) 
			DO UPDATE SET quantity = wallet_stocks.quantity + 1`,
			walletID, stockName)
		if err != nil {
			return err
		}

	} else {
		var walletQty int
		err = tx.QueryRow(ctx, "SELECT quantity FROM wallet_stocks WHERE wallet_id = $1 AND stock_name = $2",
			walletID, stockName).Scan(&walletQty)

		if err != nil || walletQty < 1 {
			return ErrInsufficientStock
		}

		_, err = tx.Exec(ctx, "UPDATE bank_stocks SET quantity = quantity + 1 WHERE stock_name = $1", stockName)
		if err != nil {
			return err
		}

		_, err = tx.Exec(ctx, "UPDATE wallet_stocks SET quantity = quantity - 1 WHERE wallet_id = $1 AND stock_name = $2",
			walletID, stockName)
		if err != nil {
			return err
		}
	}

	_, err = tx.Exec(ctx, "INSERT INTO audit_logs (type, wallet_id, stock_name) VALUES ($1, $2, $3)",
		tradeType, walletID, stockName)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (s *PostgresStorage) GetWallet(ctx context.Context, walletID string) ([]models.Stock, error) {
	rows, err := s.Pool.Query(ctx, "SELECT stock_name, quantity FROM wallet_stocks WHERE wallet_id = $1", walletID)
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

func (s *PostgresStorage) GetStockQuantityInWallet(ctx context.Context, walletID, stockName string) (int, error) {
	var qty int
	err := s.Pool.QueryRow(ctx, "SELECT quantity FROM wallet_stocks WHERE wallet_id = $1 AND stock_name = $2",
		walletID, stockName).Scan(&qty)
	if err != nil {
		return 0, nil
	}
	return qty, nil
}
