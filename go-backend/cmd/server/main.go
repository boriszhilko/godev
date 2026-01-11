// Package main is the entry point for the Go backend server.
package main

import (
	"os"
	"time"

	"go-backend/internal/cache"
	"go-backend/internal/handler"
	"go-backend/internal/store"
)

const (
	defaultPort = "8080"
	version     = "1.0.0"
)

func main() {
	startTime := time.Now()

	// Initialize data store from persistence
	dataStore := store.Initialize()

	// Initialize cache with 5 minute TTL
	appCache := cache.New(5 * time.Minute)

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// Create handler with dependencies
	h := handler.New(dataStore, appCache, handler.Config{
		Version:   version,
		StartTime: startTime,
	})

	// Start the server
	h.Start(port)
}
