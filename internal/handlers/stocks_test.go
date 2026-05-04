package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/szymmix/stock-market/internal/models"
)

func TestStocksHandler_GetStocks(t *testing.T) {
	mockData := []models.Stock{{Name: "Apple", Quantity: 10}}
	mock := &MockStorage{bankStocks: mockData}
	h := NewStocksHandler(mock)

	req := httptest.NewRequest("GET", "/stocks", nil)
	rr := httptest.NewRecorder()

	h.GetStocks(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}

	var resp models.BankState
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	if err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if len(resp.Stocks) != 1 || resp.Stocks[0].Name != "Apple" {
		t.Error("returned data mismatch")
	}
}

func TestStocksHandler_SetStocks(t *testing.T) {
	tests := []struct {
		name     string
		body     string
		err      error
		expected int
	}{
		{"Valid payload", `{"stocks":[{"name":"A","quantity":1}]}`, nil, http.StatusOK},
		{"Invalid JSON", `{"stocks": "broken"}`, nil, http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &MockStorage{err: tt.err}
			h := NewStocksHandler(mock)
			req := httptest.NewRequest("POST", "/stocks", bytes.NewBufferString(tt.body))
			rr := httptest.NewRecorder()

			h.SetStocks(rr, req)

			if rr.Code != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, rr.Code)
			}
		})
	}
}
