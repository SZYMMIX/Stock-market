package storage

import (
	"context"

	"github.com/szymmix/stock-market/internal/models"
)

type Storage interface {
	ExecuteTrade(ctx context.Context, walletID, stockName, tradeType string) error
	GetWallet(ctx context.Context, walletID string) ([]models.Stock, error)
	GetStockQuantityInWallet(ctx context.Context, walletID, stockName string) (int, error)
	GetBankStocks(ctx context.Context) ([]models.Stock, error)
	SetBankState(ctx context.Context, stocks []models.Stock) error
	GetAuditLogs(ctx context.Context) ([]models.AuditLogEntry, error)
}
