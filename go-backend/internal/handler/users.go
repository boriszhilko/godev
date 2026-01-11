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

func (h *Handler) handleUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	switch r.Method {
	case http.MethodGet:
		h.listUsers(w, r)
	case http.MethodPost:
		h.createUser(w, r)
	case http.MethodOptions:
		h.handleCORS(w)
	default:
		h.writeError(w, http.StatusMethodNotAllowed, "Method not allowed", "METHOD_NOT_ALLOWED")
	}
}

func (h *Handler) listUsers(w http.ResponseWriter, r *http.Request) {
	cacheKey := cache.UsersKey()
	if cached, found := h.cache.Get(cacheKey); found {
		json.NewEncoder(w).Encode(cached)
		return
	}

	users := h.store.GetUsers()
	response := model.UsersResponse{
		Users: users,
		Count: len(users),
	}

	h.cache.Set(cacheKey, response)

	json.NewEncoder(w).Encode(response)
}

func (h *Handler) createUser(w http.ResponseWriter, r *http.Request) {
	var req model.CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid JSON format", "INVALID_JSON")
		return
	}

	// Validate name
	if !validator.NonEmpty(req.Name) {
		h.writeError(w, http.StatusBadRequest, "Name is required and cannot be empty", "INVALID_NAME")
		return
	}

	// Validate email
	if !validator.NonEmpty(req.Email) {
		h.writeError(w, http.StatusBadRequest, "Email is required and cannot be empty", "INVALID_EMAIL")
		return
	}

	if !validator.Email(req.Email) {
		h.writeError(w, http.StatusBadRequest, "Invalid email format", "INVALID_EMAIL_FORMAT")
		return
	}

	// Validate role
	if !validator.NonEmpty(req.Role) {
		h.writeError(w, http.StatusBadRequest, "Role is required and cannot be empty", "INVALID_ROLE")
		return
	}

	// Check if email already exists
	if h.store.UserExistsByEmail(req.Email) {
		h.writeError(w, http.StatusBadRequest, "Email already exists", "EMAIL_EXISTS")
		return
	}

	user := h.store.CreateUser(req.Name, req.Email, req.Role)

	h.InvalidateUserCaches()

	h.writeJSON(w, http.StatusCreated, user)
}

func (h *Handler) handleUserByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract ID from path
	path := strings.TrimPrefix(r.URL.Path, "/api/users/")
	id, err := strconv.Atoi(path)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user := h.store.GetUserByID(id)
	if user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(user)
}
