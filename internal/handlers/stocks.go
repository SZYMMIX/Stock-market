package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/szymmix/stock-market/internal/models"
	"github.com/szymmix/stock-market/internal/storage"
)

type StocksHandler struct {
	Storage *storage.PostgresStorage
}

func NewStocksHandler(s *storage.PostgresStorage) *StocksHandler {
	return &StocksHandler{Storage: s}
}

func (h *StocksHandler) SetStocks(w http.ResponseWriter, r *http.Request) {
	var state models.BankState

	if err := json.NewDecoder(r.Body).Decode(&state); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := h.Storage.SetBankState(r.Context(), state.Stocks); err != nil {
		http.Error(w, "Failed to update bank state", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
