.PHONY: help build run stop test clean docker-build docker-up docker-down docker-logs dev install lint

# Default target
help:
	@echo "GoDevTest - Available Commands"
	@echo "=============================="
	@echo ""
	@echo "Development:"
	@echo "  make install     - Install all dependencies"
	@echo "  make dev         - Start all services in development mode"
	@echo "  make run         - Start all services locally (no Docker)"
	@echo "  make stop        - Stop all running services"
	@echo ""
	@echo "Testing:"
	@echo "  make test        - Run all tests"
	@echo "  make test-go     - Run Go backend tests"
	@echo "  make test-node   - Run Node.js backend tests"
	@echo "  make test-react  - Run React frontend tests"
	@echo "  make coverage    - Run tests with coverage report"
	@echo ""
	@echo "Docker:"
	@echo "  make docker-build  - Build all Docker images"
	@echo "  make docker-up     - Start all services with Docker Compose"
	@echo "  make docker-down   - Stop all Docker services"
	@echo "  make docker-logs   - View logs from all services"
	@echo "  make docker-clean  - Remove all containers and images"
	@echo ""
	@echo "Code Quality:"
	@echo "  make lint        - Run linters on all projects"
	@echo "  make fmt         - Format Go code"
	@echo ""
	@echo "Utilities:"
	@echo "  make clean       - Clean build artifacts"
	@echo "  make health      - Check health of all services"

# ============== Installation ==============

install:
	@echo "Installing Go dependencies..."
	cd go-backend && go mod tidy
	@echo "Installing Node.js backend dependencies..."
	cd node-backend && npm install
	@echo "Installing React frontend dependencies..."
	cd react-frontend && npm install
	@echo "All dependencies installed!"

# ============== Development ==============

dev:
	@echo "Starting development environment with Docker..."
	docker-compose -f docker-compose.dev.yml up --build

run: run-go run-node run-react

run-go:
	@echo "Starting Go backend..."
	cd go-backend && go run ./cmd/server &

run-node:
	@echo "Starting Node.js backend..."
	cd node-backend && npm start &

run-react:
	@echo "Starting React frontend..."
	cd react-frontend && npm run dev &

stop:
	@echo "Stopping all services..."
	-pkill -f "go run"
	-pkill -f "node.*server.js"
	-pkill -f "vite"
	@echo "All services stopped."

# ============== Testing ==============

test: test-go test-node

test-go:
	@echo "Running Go backend tests..."
	cd go-backend && go test -v ./...

test-node:
	@echo "Running Node.js backend tests..."
	cd node-backend && npm test || true

test-react:
	@echo "Running React frontend tests..."
	cd react-frontend && npm test || true

coverage:
	@echo "Running Go tests with coverage..."
	cd go-backend && go test -v -coverprofile=coverage.out ./...
	cd go-backend && go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: go-backend/coverage.html"

# ============== Docker ==============

docker-build:
	@echo "Building Docker images..."
	docker-compose build

docker-up:
	@echo "Starting Docker containers..."
	docker-compose up -d
	@echo "Services started!"
	@echo "  - Go Backend:      http://localhost:8080"
	@echo "  - Node.js Backend: http://localhost:3000"
	@echo "  - React Frontend:  http://localhost:5173"

docker-down:
	@echo "Stopping Docker containers..."
	docker-compose down

docker-logs:
	docker-compose logs -f

docker-clean:
	@echo "Removing Docker containers and images..."
	docker-compose down --rmi all --volumes --remove-orphans
	@echo "Docker cleanup complete."

# ============== Code Quality ==============

lint:
	@echo "Linting Go code..."
	cd go-backend && go vet ./...
	@echo "Linting Node.js code..."
	cd node-backend && npm run lint 2>/dev/null || echo "No lint script configured"
	@echo "Linting React code..."
	cd react-frontend && npm run lint 2>/dev/null || echo "No lint script configured"

fmt:
	@echo "Formatting Go code..."
	cd go-backend && go fmt ./...

# ============== Utilities ==============

clean:
	@echo "Cleaning build artifacts..."
	rm -rf go-backend/server
	rm -rf go-backend/coverage.out go-backend/coverage.html
	rm -rf node-backend/node_modules
	rm -rf react-frontend/node_modules react-frontend/dist
	@echo "Clean complete."

health:
	@echo "Checking service health..."
	@echo "Go Backend:"
	@curl -s http://localhost:8080/health | jq . 2>/dev/null || echo "  Not running"
	@echo "Node.js Backend:"
	@curl -s http://localhost:3000/health | jq . 2>/dev/null || echo "  Not running"
	@echo "React Frontend:"
	@curl -s -o /dev/null -w "  Status: %{http_code}\n" http://localhost:5173 2>/dev/null || echo "  Not running"

# ============== API Testing ==============

test-api:
	@echo "Testing API endpoints..."
	@echo "\n=== Health Check ==="
	curl -s http://localhost:8080/health | jq .
	@echo "\n=== Get Users ==="
	curl -s http://localhost:8080/api/users | jq .
	@echo "\n=== Get Tasks ==="
	curl -s http://localhost:8080/api/tasks | jq .
	@echo "\n=== Get Stats ==="
	curl -s http://localhost:8080/api/stats | jq .

test-create-user:
	@echo "Creating test user..."
	curl -s -X POST http://localhost:8080/api/users \
		-H "Content-Type: application/json" \
		-d '{"name":"Test User","email":"test@example.com","role":"developer"}' | jq .

test-create-task:
	@echo "Creating test task..."
	curl -s -X POST http://localhost:8080/api/tasks \
		-H "Content-Type: application/json" \
		-d '{"title":"Test Task","status":"pending","userId":1}' | jq .
