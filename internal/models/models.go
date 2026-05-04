package models

import "time"

type Stock struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}

type Wallet struct {
	ID     string  `json:"id"`
	Stocks []Stock `json:"stocks"`
}

type AuditLogEntry struct {
	Type      string    `json:"type"`
	WalletID  string    `json:"wallet_id"`
	StockName string    `json:"stock_name"`
	CreatedAt time.Time `json:"-"`
}

type BankState struct {
	Stocks []Stock `json:"stocks"`
}

type AuditLogResponse struct {
	Log []AuditLogEntry `json:"log"`
}
