# Everato - Modern Event Management Platform

**Everato** is a comprehensive event management platform designed as a monolithic, server-side rendered application. Built with modern Go technology, Everato provides a complete solution for event creation, management, ticketing, and analytics in a single, efficient binary.

## Overview

Everato combines all functionality into a cohesive platform that handles everything from event creation to analytics, ticketing systems, payment processing, and administration through a unified, server-side rendered interface.

![Everato Platform Overview](public/static/arch_00.png)

## Key Features

- **Single Binary Deployment**: The entire application runs from a single Go binary with an external configuration file, making deployment and scaling simple.
- **Server-Side Rendering**: Fast, SEO-friendly pages with reduced client-side JavaScript requirements.
- **Event Management**: Create, update, and manage events with customizable fields.
- **Ticketing System**: Flexible ticket types, pricing tiers, and inventory management.
- **User Management**: Comprehensive user registration, authentication, and profile management.
- **Analytics Dashboard**: Real-time insights into event performance, attendance, and revenue.
- **Payment Processing**: Secure payment handling with multiple provider options.
- **Email Notifications**: Automated confirmations, reminders, and marketing communications.
- **QR Code Generation**: Secure ticket validation through unique QR codes.

## Architecture Highlights

- **Monolithic Design**: All components are integrated into a single application, eliminating microservice complexity.
- **SSR Performance**: Server-side rendering delivers faster initial page loads and improved SEO.
- **Event Bus**: Internal event processing using Kafka for reliable asynchronous operations.
- **Database Integration**: Direct PostgreSQL connectivity with migration tooling.
- **Comprehensive Logging**: Structured logging for monitoring and debugging.
- **Observability Stack**: Prometheus, Grafana, Loki, and Promtail for complete system visibility.

## Tech Stack

- **Backend**: Go with modern frameworks and libraries
- **Database**: PostgreSQL with pgx driver
- **ORM/Query**: SQLC for type-safe SQL
- **Messaging**: Kafka & Zookeeper
- **Frontend**: Server-side rendered templates with minimal JavaScript
- **UI Framework**: TailwindCSS for styling
- **Templating**: Templ for type-safe HTML templates
- **Authentication**: JWT-based authentication
- **Development**: Docker for local development environment, Air for hot reloading

## Getting Started

### Prerequisites

- Go 1.24+
- PostgreSQL 15+
- Docker and Docker Compose (for development environment)
- Make (for running development commands)
- Node.js and npm/pnpm (for TailwindCSS compilation)

### Quick Start

1. Clone the repository:

    ```
    git clone https://github.com/yourusername/everato.git
    cd everato
    ```

2. Set up environment variables:

    ```
    cp .env.example .env
    # Edit .env file with your configuration
    ```

# Start the database:

    ```
    make db
    ```

4. Start the logging stack (optional):

    ```
    make logs
    ```

    This starts Prometheus, Grafana, Loki, and Promtail for monitoring and observability.
    - Grafana: http://localhost:3000 (default login: admin/admin)
    - Prometheus: http://localhost:9090
    - Loki: http://localhost:3100

5. Run migrations:

    ```
    make migrate-up
    ```

6. Run the application:

    ```
    # For development with hot reload
    make dev

    # For production build
    make build
    ./bin/everato
    ```

## Development

### Building

```
make build
```

### Testing

```
make test
```

### Database Management

```
# Create a new migration
make migrate-new

# Apply migrations
make migrate-up

# Roll back a migration
make migrate-down

# Seed the database with sample data
make seed
```

## Project Structure

Everato follows a well-organized directory structure that separates concerns and promotes maintainability:

```
everato/
├── assets/                # Project assets like architecture diagrams
├── components/            # UI components for templ rendering
├── config/                # Configuration management
├── docker/                # Docker-related files for development
│   ├── init/              # Database initialization scripts
│   ├── Dockerfile         # Production Docker image definition
│   ├── prometheus.yml     # Prometheus configuration
│   └── promtail-config.yaml # Promtail configuration
├── internal/              # Private application code
│   ├── db/                # Database-related code
│   │   ├── migrations/    # SQL migration files
│   │   ├── queries/       # SQLC query definitions
│   │   ├── repository/    # Generated database access code (auto-generated)
│   │   └── seed/          # Database seeding utilities
│   ├── handlers/          # HTTP request handlers
│   │   ├── handler_interface.go  # Common interface for all handlers
│   │   └── v1/            # API version 1 handlers
│   │       ├── api/       # REST API endpoints
│   │       └── views/     # Server-side rendered views
│   ├── middlewares/       # HTTP middleware components
│   │   ├── authguard_middleware.go  # Authentication middleware
│   │   ├── cors_middleware.go       # CORS handling
│   │   ├── logger_middleware.go     # Request logging
│   │   ├── requestid_middleware.go  # Request ID generation
│   │   └── timeout_middleware.go    # Request timeout handling
│   ├── services/          # Business logic
│   │   ├── event/         # Event-related services
│   │   ├── mailer/        # Email notification services
│   │   └── user/          # User management services
│   └── utils/             # Utility functions
│       ├── handler_utils.go  # Handler utilities
│       ├── http_utils.go     # HTTP utilities
│       └── utils.go          # General utilities
├── pages/                 # Page templates (templ)
├── pkg/                   # Shared public libraries
│   ├── jwt.go             # JWT handling
│   ├── logger.go          # Application logger
│   └── template.go        # Template utilities
├── public/                # Static assets (served directly)
│   ├── css/               # Compiled CSS files
│   └── js/                # JavaScript files
├── scripts/               # Utility scripts for database management
├── styles/                # Source CSS files (TailwindCSS)
└── templates/             # HTML templates
    ├── mail/              # Email templates
    └── views/             # View templates
```

### Key Files

- `main.go` - Application entry point
- `server.go` - HTTP server implementation
- `migrations_dev.go`/`migrations_prod.go` - Database migration handling
- `embedded_fs.go`/`live_fs.go` - File system handling (dev vs prod)
- `.env` - Environment variables (not committed to git)
- `config.yaml` - Application configuration
- `Makefile` - Development and build commands
- `.air.toml` - Configuration for hot reloading

### File Responsibilities

#### Configuration Files

- `config/config.go` - Loads and manages application configuration
- `.env` & `.env.example` - Environment variables for local development
- `config.yaml` - External configuration for deployment settings

#### Database

- `internal/db/migrations/*.sql` - Database schema definitions
- `internal/db/queries/*.sql` - SQLC query definitions
- `internal/db/repository/*.go` - Generated database access code
- `sqlc.yaml` - SQLC code generation configuration

#### API Implementation

- `internal/handlers/v1/api/*.go` - REST API handlers
- `internal/middlewares/*.go` - HTTP middleware components
- `internal/services/*.go` - Business logic implementation

#### View Rendering

- `components/*.templ` - Reusable UI components
- `pages/*.templ` - Page templates
- `templates/views/*.html` - Traditional HTML templates

#### Frontend Assets

- `styles/root.css` - Source CSS using TailwindCSS
- `public/css/styles.css` - Compiled CSS
- `public/js/script.js` - Client-side JavaScript

### Development Workflow

1. **Local Setup**
    - Copy `.env.example` to `.env` and configure values
    - Install dependencies with `make install`
        - This installs: SQLC, golangci-lint, Air, TailwindCSS, golang-migrate, and templ
    - Start the database with `make db`
    - Start the logging stack with `make logs` (optional)
    - Apply migrations with `make migrate-up`

2. **Development Cycle**
    - Run the application with hot reloading: `make dev`
    - For UI changes, use `make watch` to automatically rebuild templates and CSS
    - Create database migrations: `make migrate-new`
    - Seed test data: `make seed`
    - Format code: `make fmt`
    - Lint code: `make lint`

3. **Code Organization Principles**
    - Business logic lives in `internal/services/`
    - HTTP handlers only handle request/response, delegating to services
    - Database queries are defined in SQL and generated with SQLC
    - UI components use the templ templating language for type-safe HTML
    - Configuration is loaded from both environment variables and config files
    - DTOs (Data Transfer Objects) handle data validation and transformation
    - Structured logging for traceability and monitoring
    - Middleware-based HTTP request processing

4. **Build Modes**
    - **Development Mode**: Uses file system directly for live reloading
        - Build with `-tags=dev` flag
        - Reads templates and migrations from disk
        - Enables hot reloading via Air
    - **Production Mode**: Embeds static assets into binary
        - Build without tags for a self-contained binary
        - Templates and migrations embedded in executable
        - Single binary deployment with config file

5. **API Structure**
    - REST API endpoints under `/api/v1/`
    - Authentication via JWT (stored in HTTP-only cookies)
    - Structured error responses with request IDs for traceability
    - Routes grouped by domain (auth, events, users, etc.)
    - Handler interface pattern for consistent route registration
    - Response formatting with custom HttpWriter utility
    - OpenAPI/Swagger documentation (planned)
    - Prometheus metrics for performance monitoring

6. **Monitoring & Observability**
    - Prometheus metrics collection for application performance
    - Grafana dashboards for visualization (http://localhost:3000)
    - Loki for log aggregation
    - Promtail for log shipping
    - Request IDs for cross-service traceability
    - Structured JSON logging with different log levels
    - Logging stack can be started with `make logs`
    - All components are configured in docker-compose.yaml

7. **Testing**
    - Unit tests: `make test`
    - Integration tests use real database connections
    - Logging and metrics for performance analysis

## Deployment

Everato can be deployed as a single binary with an accompanying configuration file:

1. Build for production:

    ```
    make build
    ```

2. Copy the binary and configuration file to your server:

    ```
    scp bin/everato config.yaml user@your-server:/path/to/deployment/
    ```

3. Run the application:
    ```
    ./everato -config config.yaml
    ```

## Contributing

Contributions are welcome! Please read our [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct and the process for submitting pull requests.

## License

This project is copyright (c) 2025 Piush Bose. All rights reserved.

This source code is made publicly visible for portfolio and educational purposes only. You are not permitted to copy, distribute, modify, or use this code or any part of it without explicit written permission from the author - see the [LICENSE](LICENSE) file for details.
