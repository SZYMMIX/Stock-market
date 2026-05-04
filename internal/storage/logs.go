package storage

import (
	"context"

	"github.com/szymmix/stock-market/internal/models"
)

func (s *PostgresStorage) GetAuditLogs(ctx context.Context) ([]models.AuditLogEntry, error) {
	rows, err := s.Pool.Query(ctx, "SELECT type, wallet_id, stock_name FROM audit_logs ORDER BY id ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []models.AuditLogEntry
	for rows.Next() {
		var e models.AuditLogEntry
		if err := rows.Scan(&e.Type, &e.WalletID, &e.StockName); err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, nil
}
