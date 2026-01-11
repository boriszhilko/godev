// Package model defines the domain models and API request/response types.
package model

// User represents a user in the system.
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

// Task represents a task assigned to a user.
type Task struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
	UserID int    `json:"userId"`
}

// UsersResponse is the response format for listing users.
type UsersResponse struct {
	Users []User `json:"users"`
	Count int    `json:"count"`
}

// TasksResponse is the response format for listing tasks.
type TasksResponse struct {
	Tasks []Task `json:"tasks"`
	Count int    `json:"count"`
}

// StatsResponse provides statistics about users and tasks.
type StatsResponse struct {
	Users struct {
		Total int `json:"total"`
	} `json:"users"`
	Tasks struct {
		Total      int `json:"total"`
		Pending    int `json:"pending"`
		InProgress int `json:"inProgress"`
		Completed  int `json:"completed"`
	} `json:"tasks"`
}

// HealthResponse is a simple health check response.
type HealthResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// DetailedHealthResponse provides detailed health status with checks.
type DetailedHealthResponse struct {
	Status    string            `json:"status"`
	Message   string            `json:"message"`
	Version   string            `json:"version"`
	Uptime    string            `json:"uptime"`
	Checks    map[string]string `json:"checks"`
	Timestamp string            `json:"timestamp"`
}

// ErrorResponse is the standard error response format.
type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Code    string `json:"code,omitempty"`
}

// CreateUserRequest is the request body for creating a user.
type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

// CreateTaskRequest is the request body for creating a task.
type CreateTaskRequest struct {
	Title  string `json:"title"`
	Status string `json:"status"`
	UserID int    `json:"userId"`
}

// UpdateTaskRequest is the request body for updating a task.
// Pointer types allow distinguishing between "not set" and "set to zero value".
type UpdateTaskRequest struct {
	Title  *string `json:"title,omitempty"`
	Status *string `json:"status,omitempty"`
	UserID *int    `json:"userId,omitempty"`
}
