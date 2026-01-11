package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"go-backend/internal/cache"
	"go-backend/internal/model"
	"go-backend/internal/validator"
)

func (h *Handler) handleTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	switch r.Method {
	case http.MethodGet:
		h.listTasks(w, r)
	case http.MethodPost:
		h.createTask(w, r)
	case http.MethodOptions:
		h.handleCORS(w)
	default:
		h.writeError(w, http.StatusMethodNotAllowed, "Method not allowed", "METHOD_NOT_ALLOWED")
	}
}

func (h *Handler) listTasks(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	userID := r.URL.Query().Get("userId")

	cacheKey := cache.TasksKey(status, userID)
	if cached, found := h.cache.Get(cacheKey); found {
		json.NewEncoder(w).Encode(cached)
		return
	}

	tasks := h.store.GetTasks(status, userID)
	response := model.TasksResponse{
		Tasks: tasks,
		Count: len(tasks),
	}

	h.cache.Set(cacheKey, response)

	json.NewEncoder(w).Encode(response)
}

func (h *Handler) createTask(w http.ResponseWriter, r *http.Request) {
	var req model.CreateTaskRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid JSON format", "INVALID_JSON")
		return
	}

	// Validate title
	if !validator.NonEmpty(req.Title) {
		h.writeError(w, http.StatusBadRequest, "Title is required and cannot be empty", "INVALID_TITLE")
		return
	}

	// Validate status
	if !validator.Status(req.Status) {
		h.writeError(w, http.StatusBadRequest, "Invalid status. Must be one of: pending, in-progress, completed", "INVALID_STATUS")
		return
	}

	// Validate userId exists
	if h.store.GetUserByID(req.UserID) == nil {
		h.writeError(w, http.StatusBadRequest, "User ID does not exist", "INVALID_USER_ID")
		return
	}

	task := h.store.CreateTask(req.Title, req.Status, req.UserID)

	h.InvalidateTaskCaches()

	h.writeJSON(w, http.StatusCreated, task)
}

func (h *Handler) handleTaskByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Extract ID from path
	path := strings.TrimPrefix(r.URL.Path, "/api/tasks/")
	if path == "" {
		h.writeError(w, http.StatusBadRequest, "Task ID is required", "MISSING_ID")
		return
	}

	id, err := strconv.Atoi(path)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid task ID", "INVALID_ID")
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getTaskByID(w, r, id)
	case http.MethodPut:
		h.updateTask(w, r, id)
	case http.MethodOptions:
		h.handleCORS(w)
	default:
		h.writeError(w, http.StatusMethodNotAllowed, "Method not allowed", "METHOD_NOT_ALLOWED")
	}
}

func (h *Handler) getTaskByID(w http.ResponseWriter, r *http.Request, id int) {
	task := h.store.GetTaskByID(id)
	if task == nil {
		h.writeError(w, http.StatusNotFound, "Task not found", "TASK_NOT_FOUND")
		return
	}

	h.writeJSON(w, http.StatusOK, task)
}

func (h *Handler) updateTask(w http.ResponseWriter, r *http.Request, id int) {
	// Check if task exists first
	if h.store.GetTaskByID(id) == nil {
		h.writeError(w, http.StatusNotFound, "Task not found", "TASK_NOT_FOUND")
		return
	}

	var req model.UpdateTaskRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid JSON format", "INVALID_JSON")
		return
	}

	// Validate status if provided
	if req.Status != nil && !validator.Status(*req.Status) {
		h.writeError(w, http.StatusBadRequest, "Invalid status. Must be one of: pending, in-progress, completed", "INVALID_STATUS")
		return
	}

	// Validate userId if provided
	if req.UserID != nil && h.store.GetUserByID(*req.UserID) == nil {
		h.writeError(w, http.StatusBadRequest, "User ID does not exist", "INVALID_USER_ID")
		return
	}

	// Validate title if provided
	if req.Title != nil && !validator.NonEmpty(*req.Title) {
		h.writeError(w, http.StatusBadRequest, "Title cannot be empty", "INVALID_TITLE")
		return
	}

	updatedTask := h.store.UpdateTask(id, req.Title, req.Status, req.UserID)

	h.InvalidateTaskCaches()

	h.writeJSON(w, http.StatusOK, updatedTask)
}

func (h *Handler) handleStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cacheKey := cache.StatsKey()
	if cached, found := h.cache.Get(cacheKey); found {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		json.NewEncoder(w).Encode(cached)
		return
	}

	stats := h.store.GetStats()

	h.cache.Set(cacheKey, stats)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(stats)
}

func (h *Handler) handleCacheStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeError(w, http.StatusMethodNotAllowed, "Method not allowed", "METHOD_NOT_ALLOWED")
		return
	}

	stats := h.cache.Stats()
	h.writeJSON(w, http.StatusOK, stats)
}
