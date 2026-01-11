// Package cache provides a thread-safe TTL-based caching layer.
package cache

import (
	"sync"
	"sync/atomic"
	"time"
)

// Entry represents a cached item with expiration.
type Entry struct {
	Data      interface{}
	ExpiresAt time.Time
}

// Cache provides thread-safe caching with TTL expiration.
type Cache struct {
	mu      sync.RWMutex
	entries map[string]Entry
	ttl     time.Duration
	hits    atomic.Int64
	misses  atomic.Int64
}

// New creates a new Cache with the specified TTL.
// It starts a background goroutine to clean up expired entries.
func New(ttl time.Duration) *Cache {
	c := &Cache{
		entries: make(map[string]Entry),
		ttl:     ttl,
	}

	go c.cleanupExpired()

	return c
}

// Get retrieves a value from the cache.
// Returns the value and true if found and not expired, nil and false otherwise.
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, exists := c.entries[key]
	if !exists {
		c.misses.Add(1)
		return nil, false
	}

	if time.Now().After(entry.ExpiresAt) {
		c.misses.Add(1)
		return nil, false
	}

	c.hits.Add(1)
	return entry.Data, true
}

// Set stores a value in the cache with the default TTL.
func (c *Cache) Set(key string, data interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries[key] = Entry{
		Data:      data,
		ExpiresAt: time.Now().Add(c.ttl),
	}
}

// Invalidate removes specified keys from the cache.
func (c *Cache) Invalidate(keys ...string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, key := range keys {
		delete(c.entries, key)
	}
}

// InvalidateAll clears all entries from the cache.
func (c *Cache) InvalidateAll() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries = make(map[string]Entry)
}

// Stats returns cache statistics.
func (c *Cache) Stats() map[string]interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()

	hits := c.hits.Load()
	misses := c.misses.Load()
	total := hits + misses

	var hitRate float64
	if total > 0 {
		hitRate = float64(hits) / float64(total) * 100
	}

	return map[string]interface{}{
		"hits":    hits,
		"misses":  misses,
		"total":   total,
		"hitRate": hitRate,
		"entries": len(c.entries),
		"ttl":     c.ttl.String(),
	}
}

func (c *Cache) cleanupExpired() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		for key, entry := range c.entries {
			if now.After(entry.ExpiresAt) {
				delete(c.entries, key)
			}
		}
		c.mu.Unlock()
	}
}

// Key generators for common cache keys.

// UsersKey returns the cache key for users list.
func UsersKey() string {
	return "users"
}

// TasksKey returns the cache key for tasks with optional filters.
func TasksKey(status, userID string) string {
	return "tasks:" + status + ":" + userID
}

// StatsKey returns the cache key for statistics.
func StatsKey() string {
	return "stats"
}
