package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/szymmix/stock-market/internal/models"
	"github.com/szymmix/stock-market/internal/storage"
)

func TestWalletsHandler_TradeStock(t *testing.T) {
	tests := []struct {
		name           string
		tradeType      string
		mockError      error
		expectedStatus int
	}{
		{
			name:           "Success buy",
			tradeType:      "buy",
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Stock not found - 404",
			tradeType:      "buy",
			mockError:      storage.ErrStockNotFound,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Insufficient stock - 400",
			tradeType:      "buy",
			mockError:      storage.ErrInsufficientStock,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &MockStorage{err: tt.mockError}
			handler := NewWalletsHandler(mock)

			body := fmt.Sprintf(`{"type": "%s"}`, tt.tradeType)
			req := httptest.NewRequest("POST", "/wallets/1/stocks/Apple", bytes.NewBufferString(body))

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("wallet_id", "1")
			rctx.URLParams.Add("stock_name", "Apple")
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			rr := httptest.NewRecorder()

			handler.TradeStock(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %v, got %v", tt.expectedStatus, rr.Code)
			}
		})
	}
}

func TestWalletsHandler_GetWallet(t *testing.T) {
	mockStocks := []models.Stock{{Name: "Apple", Quantity: 5}}
	mock := &MockStorage{wallet: mockStocks}
	h := NewWalletsHandler(mock)

	req := httptest.NewRequest("GET", "/wallets/user1", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("wallet_id", "user1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()
	h.GetWallet(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}

	var resp models.Wallet
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	if err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if resp.ID != "user1" || len(resp.Stocks) != 1 {
		t.Error("wallet response mismatch")
	}
}

func TestWalletsHandler_GetStockQuantity(t *testing.T) {
	mock := &MockStorage{quantity: 42}
	h := NewWalletsHandler(mock)

	req := httptest.NewRequest("GET", "/wallets/user1/stocks/Apple", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("wallet_id", "user1")
	rctx.URLParams.Add("stock_name", "Apple")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()
	h.GetStockQuantity(rr, req)

	if rr.Body.String() != "42" {
		t.Errorf("expected 42, got %s", rr.Body.String())
	}
}
