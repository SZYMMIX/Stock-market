package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/szymmix/stock-market/internal/models"
	"github.com/szymmix/stock-market/internal/storage"
)

type WalletsHandler struct {
	Storage *storage.PostgresStorage
}

func NewWalletsHandler(s *storage.PostgresStorage) *WalletsHandler {
	return &WalletsHandler{Storage: s}
}

func (h *WalletsHandler) TradeStock(w http.ResponseWriter, r *http.Request) {
	walletID := chi.URLParam(r, "wallet_id")
	stockName := chi.URLParam(r, "stock_name")

	var body struct {
		Type string `json:"type"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err := h.Storage.ExecuteTrade(r.Context(), walletID, stockName, body.Type)

	if err != nil {
		if errors.Is(err, storage.ErrStockNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else if errors.Is(err, storage.ErrInsufficientStock) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, "Internal error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *WalletsHandler) GetWallet(w http.ResponseWriter, r *http.Request) {
	walletID := chi.URLParam(r, "wallet_id")
	stocks, err := h.Storage.GetWallet(r.Context(), walletID)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	response := models.Wallet{ID: walletID, Stocks: stocks}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}

func (h *WalletsHandler) GetStockQuantity(w http.ResponseWriter, r *http.Request) {
	walletID := chi.URLParam(r, "wallet_id")
	stockName := chi.URLParam(r, "stock_name")

	qty, err := h.Storage.GetStockQuantityInWallet(r.Context(), walletID, stockName)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	_, _ = fmt.Fprintf(w, "%d", qty)
}
