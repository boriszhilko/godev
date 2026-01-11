package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"go-backend/internal/model"
)

func (h *Handler) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	checks := make(map[string]string)

	// Check data store
	users := h.store.GetUsers()
	if users != nil {
		checks["datastore"] = "ok"
	} else {
		checks["datastore"] = "error"
	}

	// Check persistence
	if err := h.store.Persist(); err != nil {
		checks["persistence"] = "warning: " + err.Error()
	} else {
		checks["persistence"] = "ok"
	}

	// Check cache
	cacheStats := h.cache.Stats()
	if cacheStats != nil {
		checks["cache"] = "ok"
	} else {
		checks["cache"] = "error"
	}

	response := model.DetailedHealthResponse{
		Status:    "ok",
		Message:   "Go backend is running",
		Version:   h.config.Version,
		Uptime:    time.Since(h.config.StartTime).String(),
		Checks:    checks,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) handleLiveness(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := model.HealthResponse{
		Status:  "ok",
		Message: "Server is alive",
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) handleReadiness(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if the data store is accessible
	users := h.store.GetUsers()
	if users == nil {
		h.writeError(w, http.StatusServiceUnavailable, "Data store not ready", "NOT_READY")
		return
	}

	response := model.HealthResponse{
		Status:  "ready",
		Message: "Server is ready to serve traffic",
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(response)
}
