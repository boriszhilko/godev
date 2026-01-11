// Package store provides thread-safe data storage operations.
package store

import (
	"log"
	"strconv"
	"sync"

	"go-backend/internal/model"
)

// Store holds all application data with thread-safe access.
type Store struct {
	mu    sync.RWMutex
	users []model.User
	tasks []model.Task
}

// New creates a new empty Store.
func New() *Store {
	return &Store{
		users: []model.User{},
		tasks: []model.Task{},
	}
}

// NewWithData creates a Store with initial data.
func NewWithData(users []model.User, tasks []model.Task) *Store {
	return &Store{
		users: users,
		tasks: tasks,
	}
}

// GetUsers returns all users.
func (s *Store) GetUsers() []model.User {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.users
}

// GetUserByID returns a user by ID or nil if not found.
func (s *Store) GetUserByID(id int) *model.User {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for i := range s.users {
		if s.users[i].ID == id {
			return &s.users[i]
		}
	}
	return nil
}

// UserExistsByEmail checks if a user with the given email exists.
func (s *Store) UserExistsByEmail(email string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, user := range s.users {
		if user.Email == email {
			return true
		}
	}
	return false
}

// CreateUser adds a new user and returns it with a generated ID.
func (s *Store) CreateUser(name, email, role string) model.User {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Generate new ID by finding max ID + 1
	maxID := 0
	for _, user := range s.users {
		if user.ID > maxID {
			maxID = user.ID
		}
	}

	newUser := model.User{
		ID:    maxID + 1,
		Name:  name,
		Email: email,
		Role:  role,
	}

	s.users = append(s.users, newUser)

	// Persist data asynchronously
	go s.persistAsync()

	return newUser
}

// GetTasks returns tasks, optionally filtered by status and/or userID.
func (s *Store) GetTasks(status, userID string) []model.Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var filtered []model.Task
	for _, task := range s.tasks {
		matchStatus := status == "" || task.Status == status

		matchUserID := true
		if userID != "" {
			if id, err := strconv.Atoi(userID); err == nil {
				matchUserID = task.UserID == id
			} else {
				matchUserID = false
			}
		}

		if matchStatus && matchUserID {
			filtered = append(filtered, task)
		}
	}
	return filtered
}

// GetTaskByID returns a task by ID or nil if not found.
func (s *Store) GetTaskByID(id int) *model.Task {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for i := range s.tasks {
		if s.tasks[i].ID == id {
			return &s.tasks[i]
		}
	}
	return nil
}

// CreateTask adds a new task and returns it with a generated ID.
func (s *Store) CreateTask(title, status string, userID int) model.Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Generate new ID by finding max ID + 1
	maxID := 0
	for _, task := range s.tasks {
		if task.ID > maxID {
			maxID = task.ID
		}
	}

	newTask := model.Task{
		ID:     maxID + 1,
		Title:  title,
		Status: status,
		UserID: userID,
	}

	s.tasks = append(s.tasks, newTask)

	// Persist data asynchronously
	go s.persistAsync()

	return newTask
}

// UpdateTask updates a task and returns the updated task or nil if not found.
// Only non-nil fields are updated.
func (s *Store) UpdateTask(id int, title, status *string, userID *int) *model.Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.tasks {
		if s.tasks[i].ID == id {
			if title != nil {
				s.tasks[i].Title = *title
			}
			if status != nil {
				s.tasks[i].Status = *status
			}
			if userID != nil {
				s.tasks[i].UserID = *userID
			}

			// Persist data asynchronously
			go s.persistAsync()

			return &s.tasks[i]
		}
	}
	return nil
}

// GetStats returns statistics about users and tasks.
func (s *Store) GetStats() model.StatsResponse {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var stats model.StatsResponse
	stats.Users.Total = len(s.users)
	stats.Tasks.Total = len(s.tasks)

	for _, task := range s.tasks {
		switch task.Status {
		case "pending":
			stats.Tasks.Pending++
		case "in-progress":
			stats.Tasks.InProgress++
		case "completed":
			stats.Tasks.Completed++
		}
	}

	return stats
}

// persistAsync persists data asynchronously.
func (s *Store) persistAsync() {
	if err := s.Persist(); err != nil {
		log.Printf("Warning: Failed to persist data: %v", err)
	}
}
