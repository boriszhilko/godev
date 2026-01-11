package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"go-backend/internal/cache"
	"go-backend/internal/model"
	"go-backend/internal/store"
)

func newTestHandler() *Handler {
	s := store.NewWithData(
		[]model.User{
			{ID: 1, Name: "John Doe", Email: "john@example.com", Role: "developer"},
			{ID: 2, Name: "Jane Smith", Email: "jane@example.com", Role: "designer"},
		},
		[]model.Task{
			{ID: 1, Title: "Test task 1", Status: "pending", UserID: 1},
			{ID: 2, Title: "Test task 2", Status: "in-progress", UserID: 2},
		},
	)
	c := cache.New(5 * time.Minute)
	cfg := Config{Version: "test", StartTime: time.Now()}
	return New(s, c, cfg)
}

func TestHandler_HandleHealth(t *testing.T) {
	h := newTestHandler()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	h.handleHealth(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	var response model.DetailedHealthResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Status != "ok" {
		t.Errorf("expected status 'ok', got '%s'", response.Status)
	}
}

func TestHandler_HandleUsers_GET(t *testing.T) {
	h := newTestHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/users", nil)
	rr := httptest.NewRecorder()

	h.handleUsers(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	var response model.UsersResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Count != 2 {
		t.Errorf("expected count 2, got %d", response.Count)
	}
	if len(response.Users) != 2 {
		t.Errorf("expected 2 users, got %d", len(response.Users))
	}
}

func TestHandler_HandleUsers_POST_Valid(t *testing.T) {
	h := newTestHandler()

	body := `{"name":"Test User","email":"test@example.com","role":"developer"}`
	req := httptest.NewRequest(http.MethodPost, "/api/users", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	h.createUser(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", rr.Code)
	}

	var user model.User
	if err := json.NewDecoder(rr.Body).Decode(&user); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if user.Name != "Test User" {
		t.Errorf("expected name 'Test User', got '%s'", user.Name)
	}
	if user.Email != "test@example.com" {
		t.Errorf("expected email 'test@example.com', got '%s'", user.Email)
	}
	if user.ID == 0 {
		t.Error("expected non-zero ID")
	}
}

func TestHandler_HandleUsers_POST_InvalidEmail(t *testing.T) {
	h := newTestHandler()

	body := `{"name":"Test User","email":"invalid-email","role":"developer"}`
	req := httptest.NewRequest(http.MethodPost, "/api/users", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	h.createUser(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rr.Code)
	}

	var response model.ErrorResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Code != "INVALID_EMAIL_FORMAT" {
		t.Errorf("expected code 'INVALID_EMAIL_FORMAT', got '%s'", response.Code)
	}
}

func TestHandler_HandleTasks_GET(t *testing.T) {
	h := newTestHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/tasks", nil)
	rr := httptest.NewRecorder()

	h.handleTasks(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	var response model.TasksResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Count != 2 {
		t.Errorf("expected count 2, got %d", response.Count)
	}
}

func TestHandler_HandleTasks_POST_Valid(t *testing.T) {
	h := newTestHandler()

	body := `{"title":"New Task","status":"pending","userId":1}`
	req := httptest.NewRequest(http.MethodPost, "/api/tasks", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	h.createTask(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", rr.Code)
	}

	var task model.Task
	if err := json.NewDecoder(rr.Body).Decode(&task); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if task.Title != "New Task" {
		t.Errorf("expected title 'New Task', got '%s'", task.Title)
	}
}

func TestHandler_HandleTasks_POST_InvalidStatus(t *testing.T) {
	h := newTestHandler()

	body := `{"title":"New Task","status":"invalid","userId":1}`
	req := httptest.NewRequest(http.MethodPost, "/api/tasks", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	h.createTask(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rr.Code)
	}

	var response model.ErrorResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Code != "INVALID_STATUS" {
		t.Errorf("expected code 'INVALID_STATUS', got '%s'", response.Code)
	}
}

func TestHandler_HandleTasks_POST_InvalidUserID(t *testing.T) {
	h := newTestHandler()

	body := `{"title":"New Task","status":"pending","userId":999}`
	req := httptest.NewRequest(http.MethodPost, "/api/tasks", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	h.createTask(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rr.Code)
	}

	var response model.ErrorResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Code != "INVALID_USER_ID" {
		t.Errorf("expected code 'INVALID_USER_ID', got '%s'", response.Code)
	}
}

func TestHandler_HandleTaskByID_GET(t *testing.T) {
	h := newTestHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/tasks/1", nil)
	rr := httptest.NewRecorder()

	h.handleTaskByID(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	var task model.Task
	if err := json.NewDecoder(rr.Body).Decode(&task); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if task.ID != 1 {
		t.Errorf("expected ID 1, got %d", task.ID)
	}
}

func TestHandler_HandleTaskByID_GET_NotFound(t *testing.T) {
	h := newTestHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/tasks/999", nil)
	rr := httptest.NewRecorder()

	h.handleTaskByID(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", rr.Code)
	}
}

func TestHandler_HandleTaskByID_PUT(t *testing.T) {
	h := newTestHandler()

	body := `{"title":"Updated Task","status":"completed"}`
	req := httptest.NewRequest(http.MethodPut, "/api/tasks/1", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	h.handleTaskByID(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	var task model.Task
	if err := json.NewDecoder(rr.Body).Decode(&task); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if task.Title != "Updated Task" {
		t.Errorf("expected title 'Updated Task', got '%s'", task.Title)
	}
	if task.Status != "completed" {
		t.Errorf("expected status 'completed', got '%s'", task.Status)
	}
}

func TestHandler_HandleStats(t *testing.T) {
	h := newTestHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/stats", nil)
	rr := httptest.NewRecorder()

	h.handleStats(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	var stats model.StatsResponse
	if err := json.NewDecoder(rr.Body).Decode(&stats); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if stats.Users.Total != 2 {
		t.Errorf("expected 2 users, got %d", stats.Users.Total)
	}
	if stats.Tasks.Total != 2 {
		t.Errorf("expected 2 tasks, got %d", stats.Tasks.Total)
	}
}
