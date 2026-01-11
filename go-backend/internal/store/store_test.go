package store

import (
	"sync"
	"testing"

	"go-backend/internal/model"
)

func newTestStore() *Store {
	return NewWithData(
		[]model.User{
			{ID: 1, Name: "John Doe", Email: "john@example.com", Role: "developer"},
			{ID: 2, Name: "Jane Smith", Email: "jane@example.com", Role: "designer"},
		},
		[]model.Task{
			{ID: 1, Title: "Test task 1", Status: "pending", UserID: 1},
			{ID: 2, Title: "Test task 2", Status: "in-progress", UserID: 2},
		},
	)
}

func TestStore_GetUsers(t *testing.T) {
	s := newTestStore()
	users := s.GetUsers()

	if len(users) != 2 {
		t.Errorf("expected 2 users, got %d", len(users))
	}
}

func TestStore_GetUserByID(t *testing.T) {
	s := newTestStore()

	tests := []struct {
		name     string
		id       int
		wantName string
		wantNil  bool
	}{
		{"existing user", 1, "John Doe", false},
		{"another existing user", 2, "Jane Smith", false},
		{"non-existent user", 999, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := s.GetUserByID(tt.id)

			if tt.wantNil {
				if user != nil {
					t.Errorf("expected nil, got user with ID %d", user.ID)
				}
			} else {
				if user == nil {
					t.Errorf("expected user, got nil")
					return
				}
				if user.Name != tt.wantName {
					t.Errorf("expected name %s, got %s", tt.wantName, user.Name)
				}
			}
		})
	}
}

func TestStore_CreateUser(t *testing.T) {
	s := newTestStore()

	user := s.CreateUser("Alice Cooper", "alice@example.com", "manager")

	if user.ID != 3 {
		t.Errorf("expected ID 3, got %d", user.ID)
	}
	if user.Name != "Alice Cooper" {
		t.Errorf("expected name 'Alice Cooper', got '%s'", user.Name)
	}
	if user.Email != "alice@example.com" {
		t.Errorf("expected email 'alice@example.com', got '%s'", user.Email)
	}

	// Verify user was added
	users := s.GetUsers()
	if len(users) != 3 {
		t.Errorf("expected 3 users after creation, got %d", len(users))
	}
}

func TestStore_UserExistsByEmail(t *testing.T) {
	s := newTestStore()

	tests := []struct {
		name   string
		email  string
		exists bool
	}{
		{"existing email", "john@example.com", true},
		{"another existing email", "jane@example.com", true},
		{"non-existent email", "nobody@example.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := s.UserExistsByEmail(tt.email); got != tt.exists {
				t.Errorf("UserExistsByEmail(%q) = %v, want %v", tt.email, got, tt.exists)
			}
		})
	}
}

func TestStore_GetTasks(t *testing.T) {
	s := newTestStore()

	tests := []struct {
		name      string
		status    string
		userID    string
		wantCount int
	}{
		{"all tasks", "", "", 2},
		{"pending tasks", "pending", "", 1},
		{"in-progress tasks", "in-progress", "", 1},
		{"completed tasks", "completed", "", 0},
		{"tasks for user 1", "", "1", 1},
		{"pending tasks for user 1", "pending", "1", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tasks := s.GetTasks(tt.status, tt.userID)
			if len(tasks) != tt.wantCount {
				t.Errorf("expected %d tasks, got %d", tt.wantCount, len(tasks))
			}
		})
	}
}

func TestStore_GetTaskByID(t *testing.T) {
	s := newTestStore()

	tests := []struct {
		name      string
		id        int
		wantTitle string
		wantNil   bool
	}{
		{"existing task", 1, "Test task 1", false},
		{"non-existent task", 999, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := s.GetTaskByID(tt.id)

			if tt.wantNil {
				if task != nil {
					t.Errorf("expected nil, got task with ID %d", task.ID)
				}
			} else {
				if task == nil {
					t.Errorf("expected task, got nil")
					return
				}
				if task.Title != tt.wantTitle {
					t.Errorf("expected title %s, got %s", tt.wantTitle, task.Title)
				}
			}
		})
	}
}

func TestStore_CreateTask(t *testing.T) {
	s := newTestStore()

	task := s.CreateTask("New task", "pending", 1)

	if task.ID != 3 {
		t.Errorf("expected ID 3, got %d", task.ID)
	}
	if task.Title != "New task" {
		t.Errorf("expected title 'New task', got '%s'", task.Title)
	}

	// Verify task was added
	tasks := s.GetTasks("", "")
	if len(tasks) != 3 {
		t.Errorf("expected 3 tasks after creation, got %d", len(tasks))
	}
}

func TestStore_UpdateTask(t *testing.T) {
	s := newTestStore()

	newTitle := "Updated task"
	newStatus := "completed"

	task := s.UpdateTask(1, &newTitle, &newStatus, nil)

	if task == nil {
		t.Fatal("expected task, got nil")
	}
	if task.Title != newTitle {
		t.Errorf("expected title '%s', got '%s'", newTitle, task.Title)
	}
	if task.Status != newStatus {
		t.Errorf("expected status '%s', got '%s'", newStatus, task.Status)
	}
	// UserID should be unchanged
	if task.UserID != 1 {
		t.Errorf("expected userID 1, got %d", task.UserID)
	}
}

func TestStore_UpdateTask_NotFound(t *testing.T) {
	s := newTestStore()

	newTitle := "Updated"
	task := s.UpdateTask(999, &newTitle, nil, nil)

	if task != nil {
		t.Errorf("expected nil for non-existent task, got %+v", task)
	}
}

func TestStore_GetStats(t *testing.T) {
	s := newTestStore()

	stats := s.GetStats()

	if stats.Users.Total != 2 {
		t.Errorf("expected 2 users, got %d", stats.Users.Total)
	}
	if stats.Tasks.Total != 2 {
		t.Errorf("expected 2 tasks, got %d", stats.Tasks.Total)
	}
	if stats.Tasks.Pending != 1 {
		t.Errorf("expected 1 pending task, got %d", stats.Tasks.Pending)
	}
	if stats.Tasks.InProgress != 1 {
		t.Errorf("expected 1 in-progress task, got %d", stats.Tasks.InProgress)
	}
}

func TestStore_ConcurrentAccess(t *testing.T) {
	s := newTestStore()

	var wg sync.WaitGroup
	iterations := 100

	// Concurrent reads
	for i := 0; i < iterations; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = s.GetUsers()
			_ = s.GetTasks("", "")
			_ = s.GetStats()
		}()
	}

	// Concurrent writes
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			s.CreateUser("Test User", "test@example.com", "tester")
		}(i)
	}

	wg.Wait()
}
