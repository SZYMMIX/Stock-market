package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/szymmix/stock-market/internal/models"
)

func TestGetLogs(t *testing.T) {
	mock := &MockStorage{
		logs: []models.AuditLogEntry{
			{Type: "buy", WalletID: "w1", StockName: "Apple"},
		},
	}
	handler := &LogsHandler{Storage: mock}

	req := httptest.NewRequest("GET", "/log", nil)
	rr := httptest.NewRecorder()

	handler.GetLogs(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}
