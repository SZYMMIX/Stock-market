package handlers

import (
	"context"

	"github.com/szymmix/stock-market/internal/models"
)

type MockStorage struct {
	bankStocks []models.Stock
	wallet     []models.Stock
	quantity   int
	logs       []models.AuditLogEntry
	err        error
}

func (m *MockStorage) GetBankStocks(ctx context.Context) ([]models.Stock, error) {
	return m.bankStocks, m.err
}

func (m *MockStorage) SetBankState(ctx context.Context, stocks []models.Stock) error {
	return m.err
}

func (m *MockStorage) ExecuteTrade(ctx context.Context, walletID, stockName, tradeType string) error {
	return m.err
}

func (m *MockStorage) GetWallet(ctx context.Context, walletID string) ([]models.Stock, error) {
	return m.wallet, m.err
}

func (m *MockStorage) GetStockQuantityInWallet(ctx context.Context, walletID, stockName string) (int, error) {
	return m.quantity, m.err
}

func (m *MockStorage) GetAuditLogs(ctx context.Context) ([]models.AuditLogEntry, error) {
	return m.logs, m.err
}
