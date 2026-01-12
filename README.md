# GoDevTest - Three-Tier Application

A full-stack demonstration project showcasing a three-tier architecture with Go backend, Node.js API gateway, and React frontend.

## Architecture

```
React Frontend (port 5173)
        ↓
Node.js Backend (port 3000) - API Gateway
        ↓
Go Backend (port 8080) - Data Source
```

## ✅ Test Status

**Last Verified:** January 11, 2026  
**All Tests:** 30/30 Passed (100%)

- Go Backend: All 21 tests passed
- Node.js Gateway: All 8 proxy tests passed
- React Frontend: Verified accessible and functional
- Cache Performance: 73.9% hit rate
- Data Persistence: Verified working

## Features

### Implemented Features

- **CRUD Operations**: Complete Create, Read, Update operations for users and tasks
- **Data Persistence**: JSON file-based storage with atomic writes
- **Caching**: TTL-based caching with automatic invalidation (5-minute TTL)
- **Request Logging**: Structured logging with method, path, status, and duration
- **Health Checks**: Basic, detailed, liveness, and readiness endpoints
- **Input Validation**: Email format, status enum, and field presence validation
- **Thread Safety**: Proper mutex usage for concurrent access
- **Error Handling**: Consistent error responses with meaningful messages
- **CORS Support**: Cross-origin resource sharing enabled

### Bonus Features

- **Authentication Middleware**: API key-based authentication (ready to enable)
- **Rate Limiting**: Per-IP rate limiting with headers (ready to enable)
- **Comprehensive Testing**: Unit and integration tests with 70%+ coverage

## Quick Start

### Using Docker (Recommended)

```bash
# Start all services
make docker-up

# View logs
make docker-logs

# Stop services
make docker-down
```

Access the application:
- React Frontend: http://localhost:5173
- Node.js Backend: http://localhost:3000
- Go Backend: http://localhost:8080

### Local Development

```bash
# Install dependencies
make install

# Start all services (in separate terminals)
make run-go     # Terminal 1
make run-node   # Terminal 2
make run-react  # Terminal 3

# Or start in development mode with Docker
make dev
```

### Stop Services

```bash
make stop
```

## API Documentation

### Go Backend Endpoints (port 8080)

#### Health & Status

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Detailed health check with uptime, version, and component status |
| GET | `/health/live` | Simple liveness probe |
| GET | `/health/ready` | Readiness probe (checks data store) |
| GET | `/api/cache/stats` | Cache statistics (hits, misses, hit rate) |

#### Users

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/users` | List all users |
| GET | `/api/users/:id` | Get user by ID |
| POST | `/api/users` | Create new user |

**POST /api/users** - Create User
```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Alice Chen",
    "email": "alice@demo.com",
    "role": "engineer"
  }'
```

Response (201):
```json
{
  "id": 4,
  "name": "Alice Chen",
  "email": "alice@demo.com",
  "role": "engineer"
}
```

#### Tasks

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/tasks` | List all tasks (supports `?status=` and `?userId=` query params) |
| GET | `/api/tasks/:id` | Get task by ID |
| POST | `/api/tasks` | Create new task |
| PUT | `/api/tasks/:id` | Update task (partial updates supported) |

**POST /api/tasks** - Create Task
```bash
curl -X POST http://localhost:8080/api/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Setup CI/CD pipeline",
    "status": "pending",
    "userId": 1
  }'
```

Valid statuses: `pending`, `in-progress`, `completed`

**PUT /api/tasks/:id** - Update Task (Partial Update)
```bash
# Update only status
curl -X PUT http://localhost:8080/api/tasks/1 \
  -H "Content-Type: application/json" \
  -d '{"status": "completed"}'

# Update multiple fields
curl -X PUT http://localhost:8080/api/tasks/1 \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Updated title",
    "status": "in-progress"
  }'
```

#### Statistics

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/stats` | Get user and task statistics |

Response:
```json
{
  "users": {
    "total": 3
  },
  "tasks": {
    "total": 5,
    "pending": 2,
    "inProgress": 2,
    "completed": 1
  }
}
```

## Testing

### Run All Tests

```bash
make test
```

### Run Go Backend Tests

```bash
make test-go
```

### Run with Coverage

```bash
make coverage
```

This generates an HTML coverage report at `go-backend/coverage.html`.

### Test API Endpoints Manually

```bash
# Test all endpoints
make test-api

# Create test user
make test-create-user

# Create test task
make test-create-task
```

## Project Structure

```
godevtest-implementation/
├── go-backend/              # Go HTTP server (data source)
│   ├── cmd/server/         # Application entry point
│   ├── internal/           # Private packages
│   │   ├── cache/          # TTL-based caching
│   │   ├── handler/        # HTTP handlers
│   │   ├── middleware/     # Logging, auth, rate limiting
│   │   ├── model/          # Domain models and DTOs
│   │   ├── store/          # Data store with persistence
│   │   └── validator/      # Input validation
│   ├── Dockerfile          # Production Docker image
│   └── README.md           # Go backend documentation
├── node-backend/            # Node.js Express API gateway
│   ├── server.js           # Express server
│   ├── Dockerfile          # Production Docker image
│   └── README.md           # Node backend documentation
├── react-frontend/          # React frontend application
│   ├── src/                # React components and services
│   ├── Dockerfile          # Production Docker image with nginx
│   ├── nginx.conf          # Nginx configuration for SPA
│   └── README.md           # React frontend documentation
├── docker-compose.yml       # Production Docker setup
├── docker-compose.dev.yml   # Development Docker setup
├── Makefile                 # Build and run commands
└── README.md                # This file
```

## Configuration

### Environment Variables

| Variable | Service | Default | Description |
|----------|---------|---------|-------------|
| PORT | Go | 8080 | Server port |
| PORT | Node | 3000 | Server port |
| GO_BACKEND_URL | Node | http://localhost:8080 | Go backend URL (use service name in Docker) |

### Data Persistence

Data is persisted to `go-backend/data/data.json`. The file is created automatically on first write and updated atomically on every mutation.

### Caching

- **TTL**: 5 minutes (configurable in `internal/cache/cache.go`)
- **Invalidation**: Automatic on POST/PUT operations
- **Statistics**: Available at `/api/cache/stats`

## Development

### Code Quality

```bash
# Format Go code
make fmt

# Run linters
make lint
```

### Docker Commands

```bash
# Build images
make docker-build

# Start services
make docker-up

# View logs
make docker-logs

# Stop services
make docker-down

# Clean up everything
make docker-clean
```

### Health Check

```bash
make health
```

## Bonus Features (Optional)

### Enable Authentication

Add the authentication middleware in `internal/handler/handler.go`:

```go
// In Start() method
validKeys := []string{"secret-key-1", "secret-key-2"}
handler := middleware.Auth(validKeys)(middleware.Logging(mux))
```

Then include the API key in requests:
```bash
curl -H "X-API-Key: secret-key-1" http://localhost:8080/api/users
```

### Enable Rate Limiting

Add the rate limiting middleware in `internal/handler/handler.go`:

```go
// In Start() method
limiter := middleware.NewRateLimiter(100, 1*time.Minute) // 100 requests per minute
handler := middleware.RateLimit(limiter)(middleware.Logging(mux))
```

## Performance

- Thread-safe operations with proper mutex usage
- Background persistence to avoid blocking requests
- TTL-based caching with automatic cleanup
- Connection pooling in Docker setup

## License

MIT

## Author

Interview Test Project - Golang Developer Assessment
