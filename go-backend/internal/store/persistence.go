package store

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"go-backend/internal/model"
)

const dataFilePath = "data/data.json"

// PersistentData represents the data structure stored in the JSON file.
type PersistentData struct {
	Users []model.User `json:"users"`
	Tasks []model.Task `json:"tasks"`
}

// LoadData loads data from the JSON file.
// Returns empty data if the file doesn't exist.
func LoadData() (*PersistentData, error) {
	if _, err := os.Stat(dataFilePath); os.IsNotExist(err) {
		return &PersistentData{
			Users: []model.User{},
			Tasks: []model.Task{},
		}, nil
	}

	data, err := os.ReadFile(dataFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read data file: %w", err)
	}

	var persistentData PersistentData
	if err := json.Unmarshal(data, &persistentData); err != nil {
		return nil, fmt.Errorf("failed to parse data file: %w", err)
	}

	return &persistentData, nil
}

// SaveData saves data to the JSON file atomically.
func SaveData(data *PersistentData) error {
	dir := filepath.Dir(dataFilePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	// Write atomically: temp file then rename
	tempFile := dataFilePath + ".tmp"
	if err := os.WriteFile(tempFile, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write data file: %w", err)
	}

	if err := os.Rename(tempFile, dataFilePath); err != nil {
		os.Remove(tempFile)
		return fmt.Errorf("failed to rename data file: %w", err)
	}

	return nil
}

// Initialize loads data from file or uses defaults and returns a Store.
func Initialize() *Store {
	persistentData, err := LoadData()
	if err != nil {
		log.Printf("Warning: Failed to load data from file: %v. Using default data.", err)
		return defaultStore()
	}

	// If loaded data is empty, use defaults
	if len(persistentData.Users) == 0 && len(persistentData.Tasks) == 0 {
		return defaultStore()
	}

	return NewWithData(persistentData.Users, persistentData.Tasks)
}

// defaultStore returns a Store with sample data.
func defaultStore() *Store {
	return NewWithData(
		[]model.User{
			{ID: 1, Name: "John Doe", Email: "john@example.com", Role: "developer"},
			{ID: 2, Name: "Jane Smith", Email: "jane@example.com", Role: "designer"},
			{ID: 3, Name: "Bob Johnson", Email: "bob@example.com", Role: "manager"},
		},
		[]model.Task{
			{ID: 1, Title: "Implement authentication", Status: "pending", UserID: 1},
			{ID: 2, Title: "Design user interface", Status: "in-progress", UserID: 2},
			{ID: 3, Title: "Review code changes", Status: "completed", UserID: 3},
		},
	)
}

// Persist saves the current state of the Store to file.
func (s *Store) Persist() error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data := &PersistentData{
		Users: s.users,
		Tasks: s.tasks,
	}

	return SaveData(data)
}
