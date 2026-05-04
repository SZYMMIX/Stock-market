package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/szymmix/stock-market/internal/models"
	"github.com/szymmix/stock-market/internal/storage"
)

type LogsHandler struct {
	Storage *storage.PostgresStorage
}

func (h *LogsHandler) GetLogs(w http.ResponseWriter, r *http.Request) {
	entries, err := h.Storage.GetAuditLogs(r.Context())
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(models.AuditLogResponse{Log: entries})
}
