# Agent Guidelines for Everato

This document provides essential information for AI coding agents working on the Everato event management platform. Everato is a full-stack application with a Go backend API and React TypeScript frontend.

## Project Overview

**Stack**: Go 1.24+ backend, PostgreSQL 15+, React + TypeScript frontend (Vite), TailwindCSS, SQLC for type-safe queries, JWT auth

**Structure**:
- `internal/` - Go backend (handlers, services, db, middlewares, utils)
- `www/src/` - React frontend (components, pages, lib, hooks)
- `pkg/` - Public shared libraries (logger, jwt, templates)
- `internal/db/queries/` - SQL queries (SQLC generates type-safe Go code)
- `internal/db/migrations/` - Database migrations

## Build, Test, and Lint Commands

### Backend (Go)
```bash
make build              # Production build (embeds static assets)
make dev                # Development with hot reload (Air)
make test               # Run all tests
go test ./internal/handlers/v1/api/... -v              # Test specific package
go test ./internal/services/event/... -run TestName -v # Run single test
go test ./... -v -race -cover -coverprofile=coverage.out # Coverage
make fmt                # Format code
make lint               # Run golangci-lint
```

### Frontend (React + TypeScript)
```bash
cd www
pnpm install           # Install dependencies
pnpm dev               # Development server
pnpm build             # Production build
pnpm lint              # ESLint
tsc -b                 # Type check
```

### Database
```bash
make db                # Start PostgreSQL (Docker)
make migrate-up        # Apply migrations
make migrate-down      # Rollback last migration
make migrate-new       # Create new migration (prompts for name)
make sqlc-gen          # Generate Go code from SQL queries
make seed              # Seed database with test data
```

## Go Backend Code Style

### Import Organization
Standard library → blank line → third-party → blank line → local imports. Alphabetical within groups.

```go
import (
    "context"
    "net/http"

    "github.com/gorilla/mux"
    "github.com/jackc/pgx/v5"

    "github.com/dtg-lucifer/everato/internal/db/repository"
    "github.com/dtg-lucifer/everato/pkg"
)
```

### Naming Conventions
- **Packages**: lowercase, single word (e.g., `event`, `user`, `mailer`)
- **Files**: snake_case (e.g., `event_create.go`, `dashboard_handler.go`)
- **Types/Structs**: PascalCase (e.g., `DashboardHandler`, `CreateEventDTO`)
- **Functions**: PascalCase (exported), camelCase (private)
- **Variables**: camelCase (e.g., `eventDTO`, `userCount`)
- **Constants**: PascalCase or UPPER_SNAKE_CASE

### Handler Pattern
All handlers implement `Handler` interface with `RegisterRoutes(router *mux.Router)`. Structure:

```go
type DashboardHandler struct {
    Repo     *repository.Queries
    Conn     *pgx.Conn
    BasePath string
    Config   *config.Config
}

func NewDashboardHandler(cfg *config.Config) *DashboardHandler {
    // Initialize logger, connect to DB, return instance
}

func (h *DashboardHandler) RegisterRoutes(router *mux.Router) {
    r := router.PathPrefix(h.BasePath).Subrouter()
    r.HandleFunc("/stats", h.Stats).Methods(http.MethodGet)
}
```

### Error Handling
Always check errors. Use `pgx.ErrNoRows` for missing records. Log with context. Return appropriate HTTP status codes.

```go
event, err := repo.GetEventBySlug(ctx, slug)
if err != nil {
    if err == pgx.ErrNoRows {
        wr.Status(404).Json(utils.M{"error": "Event not found"})
        return
    }
    logger.Error("Failed to get event", "error", err)
    wr.Status(500).Json(utils.M{"error": "Database error"})
    return
}
```

### Database Queries (SQLC)
1. Write SQL in `internal/db/queries/*.sql` with SQLC annotations
2. Run `make sqlc-gen` to generate type-safe Go code in `internal/db/repository/`
3. Use uppercase SQL keywords, proper indentation

```sql
-- name: GetEventBySlug :one
SELECT * FROM events
WHERE slug = $1 AND deleted_at IS NULL
LIMIT 1;
```

### Logging
Use `pkg.NewLogger()`, always `defer logger.Close()`. Structured logging with key-value pairs.

```go
logger := pkg.NewLogger()
defer logger.Close()
logger.Info("Processing event", "slug", slug, "user_id", userId)
logger.Error("Failed operation", "error", err)
```

### Documentation
- Package-level comments for all packages
- GoDoc format for exported functions/types (comment starts with name)
- Document complex functions with parameter descriptions

## TypeScript/React Frontend Code Style

### Import Organization
React → third-party → local (components, utils, types). Alphabetical within groups.

```typescript
import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { type Event, eventAPI } from "../lib/api";
import Layout from "../components/layout";
```

### Component Structure
Functional components with hooks. Interfaces at top. PascalCase names. Export default at bottom.

```typescript
interface DashboardProps {
    userId?: string;
}

export default function DashboardPage({ userId }: DashboardProps) {
    const [loading, setLoading] = useState(true);
    
    useEffect(() => {
        fetchData();
    }, []);
    
    return <Layout>{/* Component JSX */}</Layout>;
}
```

### TypeScript Types
- `interface` for API responses and object shapes
- `type` for unions and simple aliases
- Export shared types, place in `lib/` directory

### API Integration
Use centralized client from `lib/api.ts`. Handle errors with try-catch. Show loading states. Use async/await.

```typescript
const fetchEvents = async () => {
    try {
        setLoading(true);
        const response = await eventAPI.getAllEvents();
        setEvents(response.data.data.events);
    } catch (error) {
        console.error("Error fetching events:", error);
    } finally {
        setLoading(false);
    }
};
```

### Styling
TailwindCSS utility classes. Tab width: 4 spaces (`.prettierrc.json`). Mobile-first responsive design.

## Development Guidelines

### Commit Messages
Follow conventional commits: `type(scope): Short description`

Types: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`

### Branch Naming
- `feature/short-description` - New features
- `bugfix/issue-number-short-description` - Bug fixes
- `refactor/component-name` - Code refactoring
- `docs/what-changed` - Documentation
- `test/what-tested` - Tests

### Testing
- Write table-driven tests for Go
- Test success and error cases
- Mock external dependencies
- Use descriptive test names
- Run tests before submitting PRs

### Configuration
- Environment variables in `.env` (never commit)
- Application config in `config.yaml`
- Use `utils.GetEnv()` with defaults

### Middleware Order
1. CORS (API endpoints)
2. Request ID generation
3. Logger
4. Auth guard (protected routes)
5. Timeout

### File Organization
- One handler per file in `internal/handlers/v1/api/`
- Group services by domain in `internal/services/{domain}/`
- DTOs close to handlers
- Handlers route requests, services contain business logic
- Database queries only in `internal/db/queries/` SQL files
