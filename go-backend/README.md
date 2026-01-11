# Go Backend - Data Source Server

A RESTful HTTP server built with Go that serves as the primary data source for the three-tier application.

## Features

- **Thread-Safe Data Store**: In-memory storage with proper mutex usage
- **File Persistence**: Atomic JSON file writes for data durability
- **TTL-Based Caching**: 5-minute cache with automatic invalidation
- **Request Logging**: Structured logging with middleware
- **Health Checks**: Multiple health endpoints for different monitoring needs
- **Input Validation**: Comprehensive validation with meaningful error messages
- **Comprehensive Tests**: Unit and integration tests with 70%+ coverage

## Project Structure

This project follows the standard Go project layout:

```
go-backend/
├── cmd/
│   └── server/
│       └── main.go           # Application entry point
├── internal/
│   ├── cache/
│   │   └── cache.go          # TTL-based caching layer
│   ├── handler/
│   │   ├── handler.go        # HTTP server setup, helpers
│   │   ├── handler_test.go   # Integration tests
│   │   ├── health.go         # Health check handlers
│   │   ├── tasks.go          # Task CRUD handlers
│   │   └── users.go          # User CRUD handlers
│   ├── middleware/
│   │   ├── auth.go           # API key authentication
│   │   ├── logging.go        # Request logging
│   │   └── ratelimit.go      # Rate limiting
│   ├── model/
│   │   └── model.go          # Domain models, DTOs
│   ├── store/
│   │   ├── persistence.go    # File-based persistence
│   │   ├── store.go          # Thread-safe data store
│   │   └── store_test.go     # Unit tests
│   └── validator/
│       ├── validator.go      # Input validation
│       └── validator_test.go # Validation tests
├── Dockerfile
├── go.mod
└── README.md
```

### Package Responsibilities

| Package | Description |
|---------|-------------|
| `cmd/server` | Application entry point and DI wiring |
| `internal/cache` | TTL-based caching with automatic cleanup |
| `internal/handler` | HTTP handlers and route registration |
| `internal/middleware` | HTTP middleware (logging, auth, rate limit) |
| `internal/model` | Domain models and request/response types |
| `internal/store` | Data storage with thread-safe operations |
| `internal/validator` | Input validation helpers |

## Running the Server

### Local Development

```bash
# Install dependencies
go mod tidy

# Run server
go run ./cmd/server

# Run with custom port
PORT=8081 go run ./cmd/server
```

### Docker

```bash
# Build image
docker build -t go-backend .

# Run container
docker run -p 8080:8080 go-backend
```

## API Endpoints

### Health & Monitoring

#### GET /health
Detailed health check with component status.

```json
{
  "status": "ok",
  "message": "Go backend is running",
  "version": "1.0.0",
  "uptime": "2h30m15s",
  "checks": {
    "datastore": "ok",
    "persistence": "ok",
    "cache": "ok"
  },
  "timestamp": "2026-01-11T20:00:00Z"
}
```

#### GET /health/live
Simple liveness probe (is the server responding?).

#### GET /health/ready
Readiness probe (is the server ready to serve traffic?).

#### GET /api/cache/stats
Cache statistics.

### Users

#### GET /api/users
List all users.

#### GET /api/users/:id
Get user by ID.

#### POST /api/users
Create a new user.

Request:
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "role": "developer"
}
```

### Tasks

#### GET /api/tasks
List all tasks. Supports filtering.

Query Parameters:
- `status`: Filter by status (`pending`, `in-progress`, `completed`)
- `userId`: Filter by user ID

#### GET /api/tasks/:id
Get task by ID.

#### POST /api/tasks
Create a new task.

Request:
```json
{
  "title": "Implement feature X",
  "status": "pending",
  "userId": 1
}
```

#### PUT /api/tasks/:id
Update an existing task (partial updates supported).

Request (all fields optional):
```json
{
  "title": "Updated title",
  "status": "completed",
  "userId": 2
}
```

### Statistics

#### GET /api/stats
Get statistics about users and tasks.

## Error Handling

All errors return a consistent format:

```json
{
  "success": false,
  "error": "Error message here",
  "code": "ERROR_CODE"
}
```

## Testing

```bash
# All tests
go test ./...

# Verbose output
go test -v ./...

# With coverage
go test -cover ./...

# Generate HTML coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

## Configuration

### Environment Variables

- `PORT`: Server port (default: 8080)

## Bonus Features

### Authentication

Enable API key authentication by wrapping the handler:

```go
validKeys := []string{"secret-key-1", "secret-key-2"}
handler := middleware.Auth(validKeys)(handler)
```

### Rate Limiting

Enable per-IP rate limiting:

```go
limiter := middleware.NewRateLimiter(100, 1*time.Minute) // 100 req/min
handler := middleware.RateLimit(limiter)(handler)
```

## License

MIT
