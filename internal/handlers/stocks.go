package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/szymmix/stock-market/internal/models"
	"github.com/szymmix/stock-market/internal/storage"
)

type StocksHandler struct {
	Storage storage.Storage
}

func NewStocksHandler(s storage.Storage) *StocksHandler {
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

func (h *StocksHandler) GetStocks(w http.ResponseWriter, r *http.Request) {
	stocks, err := h.Storage.GetBankStocks(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch bank state", http.StatusInternalServerError)
		return
	}

	response := models.BankState{
		Stocks: stocks,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}
