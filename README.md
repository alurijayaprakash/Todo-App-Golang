
# Todo App - Go REST API

A production-ready Go REST API for managing todo items with PostgreSQL persistence, in-memory caching, and conditional Swagger UI based on environment.

## Features

✅ **Clean Architecture** - Domain-driven design with clear separation of concerns  
✅ **PostgreSQL Support** - Persistent data storage with migrations  
✅ **In-Memory Repository** - Swap between PostgreSQL and memory storage  
✅ **Caching Layer** - Built-in memory cache for improved performance  
✅ **Environment-Based Swagger** - Swagger UI controlled via `APP_ENV` variable  
✅ **Docker Support** - Multi-stage Alpine-based build for minimal images  
✅ **Kubernetes Ready** - Complete K8s deployment manifests included  
✅ **Structured Logging** - Clear bootstrap and request logging  

---

## Project Structure

```
.
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go            # Configuration loader
│   ├── handler/
│   │   ├── domain/
│   │   │   ├── cache.go         # Cache interface
│   │   │   └── todo.go          # Todo domain models
│   │   └── http/
│   │       ├── router.go        # HTTP router setup
│   │       └── todo.go          # HTTP handlers
│   ├── repository/
│   │   ├── memory/
│   │   │   ├── cache.go         # Memory cache implementation
│   │   │   └── memory.go        # In-memory repository
│   │   └── postgres/
│   │       └── postgres.go      # PostgreSQL repository
│   └── service/
│       └── todo.go              # Business logic service
├── docs/                         # Swagger generated docs
├── migrations/                   # Database migration scripts
├── deployments/                  # K8s manifests
├── Dockerfile                    # Production-optimized build
├── docker-compose.yml            # Local development stack
├── .env                          # Environment configuration
└── go.mod                        # Go module definition
```

---

## Prerequisites

- **Go 1.24.1** or higher
- **PostgreSQL 16** (optional - can use in-memory storage)
- **Docker & Docker Compose** (optional - for containerized deployment)
- **kubectl** (optional - for Kubernetes deployment)

---

## Quick Start

### 1. Local Development (In-Memory Storage)

```bash
# Clone repository
git clone <repo-url>
cd Todo-App-Golang

# Install dependencies
go mod download

# Run with in-memory storage (development)
APP_ENV=development DB_TYPE=memory go run ./cmd/server/main.go
```

Server starts at: `http://localhost:8080`  
Swagger UI: `http://localhost:8080/swagger/index.html`

### 2. Local Development with PostgreSQL

Update `.env`:
```env
DB_TYPE=postgres
DB_HOST=localhost
DB_PORT=5432
DB_NAME=todo_db
DB_USER=postgres
DB_PASSWORD=securepassword
```

Start PostgreSQL:
```bash
# Using Docker
docker run -d \
  --name postgres \
  -e POSTGRES_PASSWORD=securepassword \
  -e POSTGRES_DB=todo_db \
  -p 5432:5432 \
  postgres:16-alpine

# Run migrations
psql -h localhost -U postgres -d todo_db -f migrations/001_initial_schema.sql
```

Run server:
```bash
go run ./cmd/server/main.go
```

### 3. Docker Compose (Complete Stack)

```bash
# Start all services (PostgreSQL + API)
docker-compose up -d

# View logs
docker-compose logs -f todo-api

# Access services
# API: http://localhost:8080
# Swagger: http://localhost:8080/swagger/index.html (development)
# Database: postgres://postgres:securepassword@localhost:5432/todo_db
```

---

## Configuration

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `APP_ENV` | `development` | Environment mode (`development` or `production`) |
| `SERVER_PORT` | `8080` | HTTP server port |
| `DB_TYPE` | `postgres` | Repository type (`memory` or `postgres`) |
| `DB_HOST` | `localhost` | PostgreSQL host |
| `DB_PORT` | `5432` | PostgreSQL port |
| `DB_NAME` | `todo_db` | Database name |
| `DB_USER` | `postgres` | Database user |
| `DB_PASSWORD` | - | Database password |
| `DB_SSLMODE` | `disable` | SSL mode for PostgreSQL |

### Swagger Availability

- **Development** (`APP_ENV=development`): Swagger UI accessible at `/swagger/index.html`
- **Production** (`APP_ENV=production`): Swagger UI disabled (404 responses)

---

## API Endpoints

### Create Todo

```bash
POST /api/v1/todos
Content-Type: application/json

{
  "title": "Buy groceries",
  "description": "Milk, eggs, bread"
}
```

**Response** (201 Created):
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "title": "Buy groceries",
  "description": "Milk, eggs, bread",
  "completed": false,
  "created_at": "2026-06-10T10:30:00Z"
}
```

### Get Todo by ID

```bash
GET /api/v1/todos/{id}
```

**Response** (200 OK):
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "title": "Buy groceries",
  "description": "Milk, eggs, bread",
  "completed": false,
  "created_at": "2026-06-10T10:30:00Z"
}
```

---

## Docker Deployment

### Build Image

```bash
docker build -t todo-api:latest .
```

**Image Size**: ~15MB  
**Base**: Alpine 3.19  
**Build Strategy**: Multi-stage (build on larger image, deploy on minimal)

### Run Container

**Development** (Swagger enabled):
```bash
docker run -d \
  -e APP_ENV=development \
  -e DB_TYPE=postgres \
  -e DB_HOST=<db-host> \
  -e DB_PORT=5432 \
  -e DB_NAME=todo_db \
  -e DB_USER=postgres \
  -e DB_PASSWORD=<password> \
  -p 8080:8080 \
  todo-api:latest
```

**Production** (Swagger disabled):
```bash
docker run -d \
  -e APP_ENV=production \
  -e DB_TYPE=postgres \
  -e DB_HOST=<db-host> \
  -e DB_PORT=5432 \
  -e DB_NAME=todo_db \
  -e DB_USER=postgres \
  -e DB_PASSWORD=<password> \
  -p 8080:8080 \
  todo-api:latest
```

---

## Kubernetes Deployment

Deploy to Kubernetes cluster:

```bash
# Create namespace and deploy all resources
kubectl apply -f deployments/k8s-deployment.yaml

# Verify deployment
kubectl get pods -n todo-app

# Port forward to access API
kubectl port-forward -n todo-app svc/todo-api-service 8080:8080

# Access API at http://localhost:8080
```

**Note**: Swagger UI is disabled in Kubernetes (APP_ENV=production).

---

## Development

### Build Locally

```bash
go build -o server ./cmd/server/main.go
./server
```

### Format Code

```bash
gofmt -w ./...
```

### Run Tests

```bash
go test ./...
```

### Generate API Documentation

Swagger docs are auto-generated from code comments. To regenerate:

```bash
swag init -g cmd/server/main.go
```

---

## Architecture

### Layered Design

```
HTTP Layer (handlers)
    ↓
Service Layer (business logic)
    ↓
Repository Layer (data access)
    ├── PostgreSQL
    └── In-Memory
```

### Dependency Injection

Dependencies are injected at startup in `main.go`:
- Repository implementation selected based on `DB_TYPE`
- Service receives repository and cache
- Handlers receive service

---

## File Dependencies

| File | Depends On |
|------|-----------|
| `main.go` | `config`, `service`, `http` |
| `config.go` | Environment & flags |
| `service.go` | Repository, Cache |
| `http/router.go` | Service, `AppEnv` |
| `http/todo.go` | Service |
| `repository/memory.go` | Todo domain |
| `repository/postgres.go` | Todo domain, database/sql |

---

## Troubleshooting

### Swagger Not Showing?
- Check `APP_ENV` is set to `development`
- Verify server is running: `curl http://localhost:8080/api/v1/todos`

### Database Connection Failed?
- Verify PostgreSQL is running
- Check `DB_HOST`, `DB_PORT`, credentials in `.env`
- Ensure migrations have been applied

### Docker Build Fails?
- Verify Docker daemon is running
- Check Dockerfile syntax: `docker build -t test .`

---

## Security Considerations

- ✅ Non-root execution in Docker
- ✅ Environment-based Swagger disable for production
- ✅ Minimal attack surface (Alpine 3.19)
- ✅ Prepared statements for SQL injection prevention
- ✅ Input validation at handler layer

---

## License

MIT License - See LICENSE file for details

---

## Support

For issues, questions, or contributions, please open an issue on GitHub.

**Project Repository**: [Todo-App-Golang](https://github.com/yourusername/Todo-App-Golang)