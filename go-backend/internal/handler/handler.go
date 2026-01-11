// Package handler provides HTTP handlers for the API.
package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"go-backend/internal/cache"
	"go-backend/internal/middleware"
	"go-backend/internal/model"
	"go-backend/internal/store"
)

// Config holds handler configuration.
type Config struct {
	Version   string
	StartTime time.Time
}

// Handler contains the HTTP handlers and their dependencies.
type Handler struct {
	store  *store.Store
	cache  *cache.Cache
	config Config
}

// New creates a new Handler with the given dependencies.
func New(s *store.Store, c *cache.Cache, cfg Config) *Handler {
	return &Handler{
		store:  s,
		cache:  c,
		config: cfg,
	}
}

// RegisterRoutes sets up all routes on the given mux.
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/health", h.handleHealth)
	mux.HandleFunc("/health/live", h.handleLiveness)
	mux.HandleFunc("/health/ready", h.handleReadiness)
	mux.HandleFunc("/api/users", h.handleUsers)
	mux.HandleFunc("/api/users/", h.handleUserByID)
	mux.HandleFunc("/api/tasks", h.handleTasks)
	mux.HandleFunc("/api/tasks/", h.handleTaskByID)
	mux.HandleFunc("/api/stats", h.handleStats)
	mux.HandleFunc("/api/cache/stats", h.handleCacheStats)
}

// Start starts the HTTP server on the given port.
func (h *Handler) Start(port string) {
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	// Apply middleware chain
	// Only logging is enabled by default

	// Optional: Enable authentication (bonus feature)
	// Example usage:
	// api keys would be taken from the database
	// validKeys := []string{"secret-key-1", "secret-key-2"}
	// handler := middleware.Auth(validKeys)(middleware.Logging(mux))

	// Optional: Enable rate limiting (bonus feature)
	// Example usage:
	// limiter := middleware.NewRateLimiter(100, 1*time.Minute) // 100 req/min
	// handler := middleware.RateLimit(limiter)(middleware.Logging(mux))

	// Optional: Enable both auth and rate limiting
	// Example usage:
	// validKeys := []string{"secret-key-1"}
	// limiter := middleware.NewRateLimiter(100, 1*time.Minute)
	// handler := middleware.Auth(validKeys)(
	//     middleware.RateLimit(limiter)(
	//         middleware.Logging(mux)))

	// Current configuration: Only logging middleware
	handler := middleware.Logging(mux)

	log.Printf("Go backend server starting on http://localhost:%s", port)
	log.Printf("Serving data directly from Go backend")

	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// writeJSON writes a JSON response with the given status code.
func (h *Handler) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// writeError writes a standardized error response.
func (h *Handler) writeError(w http.ResponseWriter, status int, message, code string) {
	response := model.ErrorResponse{
		Success: false,
		Error:   message,
		Code:    code,
	}
	h.writeJSON(w, status, response)
}

// handleCORS handles preflight OPTIONS requests.
func (h *Handler) handleCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.WriteHeader(http.StatusOK)
}

// InvalidateUserCaches clears user-related caches.
func (h *Handler) InvalidateUserCaches() {
	h.cache.Invalidate(cache.UsersKey())
	h.cache.Invalidate(cache.StatsKey())
}

// InvalidateTaskCaches clears task-related caches.
func (h *Handler) InvalidateTaskCaches() {
	h.cache.InvalidateAll()
}
