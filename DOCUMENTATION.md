# Everato - Modern Event Management Platform
## Project Documentation

---

## 1. Introduction

**Everato** is a comprehensive, open-source event management platform designed to revolutionize how organizations and individuals create, manage, and host events. Built with modern technologies and a focus on simplicity and portability, Everato provides a complete solution for event creation, ticketing, user management, and analytics - all packaged into a single, self-contained binary file.

### What is Everato?

Everato is a full-stack event management system that combines the power of a Go-based backend with a modern React frontend. The platform handles everything from event creation and management to ticketing systems, payment processing, QR code generation for secure ticket validation, and real-time analytics dashboards. It offers a seamless experience for both event organizers and attendees through an intuitive, responsive interface.

### Key Capabilities

- **Event Management**: Create, update, and manage events with customizable fields, including start/end times, locations, capacity management, and detailed descriptions
- **Ticketing System**: Flexible ticket types with pricing tiers, inventory management, and automated availability tracking
- **User Management**: Comprehensive user registration, authentication with JWT tokens, profile management, and email verification
- **Admin Dashboard**: Powerful administrative interface with role-based access control (RBAC) for managing events, users, and analytics
- **Analytics & Insights**: Real-time dashboards showing event performance, attendance metrics, and revenue tracking
- **Payment Processing**: Secure payment handling with support for multiple payment providers
- **Email Notifications**: Automated email confirmations, reminders, and marketing communications using customizable HTML templates
- **QR Code Validation**: Secure ticket validation system using unique QR codes for each ticket
- **RESTful API**: Clean, well-documented API for integration with external systems

### Vision and Mission

The mission of Everato is to democratize event management by providing a powerful, enterprise-grade platform that anyone can deploy and use without vendor lock-in or expensive subscriptions. Whether you're organizing a small community meetup, a professional conference, or a large-scale festival, Everato provides all the tools you need in a single, portable package.

---

## 2. Why Everato is Different: Advantages Over Market Leaders

While platforms like Konfhub, Eventbrite, and Meetup.com dominate the event management space, Everato offers unique advantages that set it apart from these market leaders:

### 2.1 Single Binary Distribution - Ultimate Portability

**The Everato Difference**: Unlike traditional event management platforms that require complex deployment processes with separate frontend servers, backend services, database configurations, and static asset hosting, Everato distributes as a **single, self-contained binary file**.

Using Go's powerful `embed` package, Everato embeds:
- The entire compiled frontend application (React/TypeScript with Vite)
- All static assets (CSS, JavaScript, images, fonts)
- Database migration scripts
- Email templates
- The backend server and API

This means:
- **Zero External Dependencies**: No need for separate web servers like Nginx or Apache
- **Instant Deployment**: Copy one file to any server and run it
- **Version Control Simplicity**: One binary = one version of your entire application
- **Reduced Attack Surface**: Fewer moving parts means fewer security vulnerabilities
- **Cross-Platform Compatibility**: Compile once for Linux, Windows, macOS, or ARM architectures

**Market Comparison**:
- Konfhub and similar platforms require you to use their hosted infrastructure with no self-hosting option
- Self-hosted alternatives like Pretix require multiple services, complex Docker setups, and ongoing maintenance
- Everato: One binary file, one command to run

### 2.2 True Open Source Freedom

**Complete Transparency**: Everato is fully open source under a permissive license, meaning:
- **Inspect Every Line**: Review the entire codebase for security, compliance, or customization
- **No Vendor Lock-In**: Your data, your infrastructure, your control
- **Community-Driven Development**: Contribute features, fix bugs, and shape the platform's future
- **Free Forever**: No hidden costs, no per-event fees, no feature paywalls

**Market Comparison**:
- Commercial platforms charge 2-5% per ticket plus processing fees
- Proprietary codebases prevent customization and security audits
- Data ownership and privacy concerns with third-party hosting

### 2.3 Self-Hosting with Cloud Option

**Flexible Deployment Model**: Everato offers the best of both worlds:

**Self-Hosting (Current)**:
- Download the binary and run it on your own infrastructure
- Complete control over data privacy and security
- No recurring subscription costs
- Ideal for organizations with strict data compliance requirements
- Perfect for educational institutions, non-profits, and privacy-conscious organizations

**Managed Cloud Hosting (Future)**:
- For organizations without technical resources, we will offer a managed cloud service
- Affordable pricing to sustain project development and maintenance
- Professional support and guaranteed uptime
- Automatic updates and security patches
- Easy migration path from self-hosted to cloud or vice versa

This dual model ensures:
- The project remains sustainable through ethical monetization
- Users have choices based on their needs and resources
- The core platform stays free and open source
- Development continues with community and commercial support

### 2.4 Modern Architecture and Technology Stack

**Built for Performance and Scalability**:
- **Go Backend**: Lightning-fast performance, minimal memory footprint, excellent concurrency handling
- **React Frontend**: Modern, interactive single-page application with optimal user experience
- **PostgreSQL Database**: Robust, ACID-compliant data storage with excellent scalability
- **Type-Safe Queries**: Using SQLC for compile-time SQL verification, eliminating runtime SQL errors
- **Event-Driven Architecture**: Kafka integration for reliable asynchronous processing
- **Observability**: Built-in Prometheus metrics, Grafana dashboards, and structured logging

**Market Comparison**:
- Many platforms use legacy tech stacks requiring significant resources
- Everato's Go backend uses ~20-50MB RAM vs 500MB+ for Node.js or Java alternatives
- Faster response times improve user experience and reduce infrastructure costs

### 2.5 Developer-Friendly Design

**API-First Architecture**:
- Clean RESTful API design for easy integration
- Comprehensive API documentation
- JWT-based authentication for secure third-party access
- Webhook support for real-time event notifications

**Extensibility**:
- Modular code structure makes customization straightforward
- Plugin architecture (planned) for extending functionality
- Well-documented codebase with extensive comments

---

## 3. Technologies Used

Everato leverages a carefully selected technology stack that balances performance, developer experience, and maintainability.

### 3.1 Backend Technologies

#### Core Language and Framework
- **Go 1.24+**: The primary backend language, chosen for its:
  - Excellent performance and low memory footprint
  - Built-in concurrency support (goroutines and channels)
  - Fast compilation and cross-platform build capabilities
  - Strong standard library
  - Native support for embedding files into binaries

#### Web Framework and Routing
- **Gorilla Mux**: Powerful HTTP router and URL matcher for building RESTful APIs
  - Request routing with variable extraction
  - Subrouter support for organizing routes
  - Middleware chain support

#### Database Layer
- **PostgreSQL 15+**: Primary relational database
  - ACID compliance for data integrity
  - Rich feature set including JSON support, full-text search
  - Excellent performance and scalability

- **pgx/v5**: High-performance PostgreSQL driver
  - Native Go implementation
  - Connection pooling
  - Prepared statement support

- **SQLC**: SQL compiler that generates type-safe Go code from SQL
  - Compile-time SQL verification
  - Eliminates SQL injection vulnerabilities
  - Full type safety without ORM overhead

- **golang-migrate**: Database migration tool
  - Version-controlled schema changes
  - Up/down migration support
  - Embedded migration support for production

#### Authentication and Security
- **JWT (JSON Web Tokens)**: golang-jwt/jwt/v5
  - Stateless authentication
  - Token-based authorization
  - Secure claims validation

- **bcrypt**: golang.org/x/crypto
  - Secure password hashing
  - Adaptive cost factor for future-proofing

#### Message Queue and Event Processing
- **Apache Kafka**: Distributed event streaming platform
  - Asynchronous task processing
  - Event-driven architecture support
  - High throughput and fault tolerance

#### Email and Notifications
- **gomail.v2**: Email sending library
  - SMTP support
  - HTML email template rendering
  - Attachment support

#### Observability and Monitoring
- **Prometheus**: prometheus/client_golang
  - Metrics collection and exposure
  - Custom metric definition
  - Time-series data storage

- **Structured Logging**: Custom logger implementation
  - JSON-formatted logs
  - Log levels and filtering
  - Request tracing support

#### Additional Libraries
- **UUID Generation**: google/uuid for unique identifiers
- **Configuration Management**: YAML-based configuration with gopkg.in/yaml.v2
- **Environment Variables**: godotenv for development environment management
- **Validation**: go-playground/validator/v10 for request validation

### 3.2 Frontend Technologies

#### Core Framework
- **React 19.1+**: Modern JavaScript library for building user interfaces
  - Component-based architecture
  - Virtual DOM for optimal performance
  - Hooks for state management
  - Server components support

- **TypeScript 5.8+**: Strongly typed superset of JavaScript
  - Type safety across the codebase
  - Enhanced IDE support and autocomplete
  - Compile-time error detection

#### Build Tool and Development
- **Vite 7.1+**: Next-generation frontend build tool
  - Lightning-fast hot module replacement (HMR)
  - Optimized production builds
  - Native ES modules support
  - Plugin ecosystem

#### Routing
- **React Router 7.8+**: Declarative routing for React applications
  - Client-side routing for SPA
  - Nested routes support
  - Route-based code splitting
  - Navigation guards

#### Styling
- **TailwindCSS 4.1+**: Utility-first CSS framework
  - Rapid UI development
  - Consistent design system
  - Small production bundle size
  - JIT compiler for on-demand class generation

- **@tailwindcss/vite**: Vite plugin for Tailwind integration

#### HTTP Client
- **Axios 1.11+**: Promise-based HTTP client
  - Request/response interceptors
  - Automatic request/response transformation
  - CSRF protection support
  - Request cancellation

#### Icons and Assets
- **React Icons 5.5+**: Popular icon library
  - Thousands of icons from multiple icon packs
  - Tree-shaking for optimal bundle size

#### Development Tools
- **ESLint 9.33+**: JavaScript/TypeScript linter
  - Code quality enforcement
  - React-specific rules
  - TypeScript support

- **Vite Plugin React SWC**: Fast React refresh using SWC compiler

### 3.3 Infrastructure and DevOps

#### Containerization and Orchestration
- **Docker**: Application containerization
- **Docker Compose**: Multi-container orchestration for development

#### Database Management
- **Migration Scripts**: Version-controlled SQL migrations in `internal/db/migrations/`
- **Seeding**: Development data seeding for testing

#### Monitoring Stack
- **Prometheus**: Metrics collection and alerting
- **Grafana**: Metrics visualization and dashboards
- **Loki**: Log aggregation system
- **Promtail**: Log shipping agent

#### Development Tools
- **Make**: Task automation and build orchestration
- **Air**: Live reload for Go applications during development
- **pnpm**: Fast, disk-efficient package manager for Node.js

### 3.4 Database Schema and Queries

The database layer uses SQLC to generate type-safe Go code from SQL queries:

**Schema Management**:
- Migrations in `internal/db/migrations/` (up and down migrations)
- Organized, versioned schema changes (e.g., `000001_init.up.sql`)

**Query Organization**:
- SQL queries in `internal/db/queries/`
- SQLC configuration in `sqlc.yaml`
- Generated repository code in `internal/db/repository/`

**Key Database Features**:
- UUID primary keys for distributed system compatibility
- JSONB columns for flexible event metadata
- Full-text search capabilities
- Referential integrity with foreign keys
- Indexes for query optimization

---

## 4. Application Showcase with Examples

This section demonstrates Everato's capabilities through user interface examples and backend code implementations.

### 4.1 User Interface Examples

#### 4.1.1 Event Listing Page
*[Event Listing Page Screenshot]*

The event listing page showcases all available events with:
- Responsive grid layout with event cards
- Event thumbnails and titles
- Date, time, and location information
- Ticket availability status
- Quick "View Details" and "Register" actions
- Search and filter functionality
- Pagination for large event lists

#### 4.1.2 Event Details Page
*[Event Details Page Screenshot]*

A comprehensive view of individual events featuring:
- Hero image with event branding
- Detailed event description with rich text support
- Date, time, venue, and capacity information
- Multiple ticket tier options with pricing
- Real-time availability tracking
- Event organizer information
- Social sharing buttons
- Registration/booking interface

#### 4.1.3 Admin Dashboard
*[Admin Dashboard Screenshot]*

The administrative interface provides:
- Overview of all events with status indicators
- Quick statistics: total events, attendees, revenue
- Event creation and management tools
- User management and access control
- Analytics charts and graphs
- Ticket sales tracking
- Email communication center
- System configuration options

#### 4.1.4 User Registration and Authentication
*[User Login Page Screenshot]*

Clean, modern authentication interface:
- Email/password login form
- Social authentication options (planned)
- Password reset functionality
- Email verification flow
- "Remember me" functionality
- Responsive design for mobile devices

#### 4.1.5 Ticket Confirmation Page
*[Ticket Confirmation Page Screenshot]*

Post-registration confirmation featuring:
- Booking confirmation details
- QR code for ticket validation
- Event information summary
- Calendar integration buttons (.ics download)
- Print ticket option
- Email receipt notification

### 4.2 Backend Code Examples

This section demonstrates the clean, well-structured backend code that powers Everato.

#### 4.2.1 Single Binary Embedding - The Core Innovation

The following code shows how Everato embeds the entire frontend and assets into a single binary:

```go
//go:build !dev
// +build !dev

package main

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"strings"
)

// Embed the entire compiled frontend (React/Vite build output)
//go:embed www/dist
var viewsFS embed.FS

// Embed all public static assets (images, fonts, etc.)
//go:embed www/public
var publicFS embed.FS

// Embed database migration scripts for automated schema management
//go:embed internal/db/migrations/*.sql
var migrationsFS embed.FS

// ViewsFS serves the React SPA from the embedded filesystem
func ViewsFS() http.Handler {
	sub, err := fs.Sub(viewsFS, "www/dist")
	if err != nil {
		panic("Failed to create sub filesystem: " + err.Error())
	}
	fsys := http.FS(sub)
	fsHandler := http.FileServer(fsys)

	// SPA routing: serve index.html for all non-file routes
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cleanPath := strings.TrimPrefix(r.URL.Path, "/")
		if f, err := fsys.Open(cleanPath); err == nil {
			if info, _ := f.Stat(); info != nil && !info.IsDir() {
				fsHandler.ServeHTTP(w, r)
				return
			}
		}
		// Serve index.html for client-side routing
		index, err := fsys.Open("index.html")
		if err != nil {
			http.Error(w, "index.html not found", http.StatusInternalServerError)
			return
		}
		info, _ := index.Stat()
		http.ServeContent(w, r, "index.html", info.ModTime(), index)
	})
}

// MigrationsFS provides access to embedded migration scripts
func MigrationsFS() (fs.FS, error) {
	_, err := fs.Stat(migrationsFS, "internal/db/migrations")
	if err != nil && os.IsNotExist(err) {
		return nil, fmt.Errorf("migrations directory does not exist")
	}
	return migrationsFS, nil
}
```

**Why This Matters**:
- All frontend assets are compiled into the Go binary at build time
- No need for separate static file hosting
- The binary becomes truly portable - copy and run anywhere
- Version consistency - frontend and backend are always in sync

#### 4.2.2 Application Entry Point

The main entry point demonstrates clean initialization and startup:

```go
package main

import (
	"github.com/dtg-lucifer/everato/config"
	"github.com/dtg-lucifer/everato/pkg"
	"github.com/joho/godotenv"
)

func main() {
	logger := pkg.NewLogger()

	// Load environment variables (optional, falls back to config.yaml)
	if err := godotenv.Load(); err != nil {
		logger.Error("Error loading .env file, using config.yaml")
	}

	// Load application configuration
	cfg, err := config.NewConfig("config.yaml")
	if err != nil {
		logger.Error("Error loading configuration", "err", err.Error())
		panic(err)
	}

	// Run database migrations automatically
	logger.Info("Running migrations...")
	if err := MigrateDB(cfg); err != nil {
		logger.Error("Error running migrations", "err", err.Error())
		panic(err)
	}
	logger.Info("Migrations completed successfully")

	// Initialize super admin users
	if err := SuperUserInit(cfg); err != nil {
		logger.Error("Error initializing super users", "err", err.Error())
		panic(err)
	}

	// Initialize and start the HTTP server
	server := NewServer(cfg)
	if err := server.Start(); err != nil {
		logger.Error("Error starting server", "err", err.Error())
		panic(err)
	}
}
```

**Highlights**:
- Single command startup process
- Automatic database migrations
- Configuration management with fallbacks
- Structured logging from the start
- Error handling at every step

#### 4.2.3 Event Handler - Core Business Logic

This example shows the clean, well-documented event handling code:

```go
package api

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"

	"github.com/dtg-lucifer/everato/config"
	"github.com/dtg-lucifer/everato/internal/db/repository"
	"github.com/dtg-lucifer/everato/internal/handlers"
	"github.com/dtg-lucifer/everato/internal/middlewares"
	"github.com/dtg-lucifer/everato/internal/services/event"
	"github.com/dtg-lucifer/everato/internal/utils"
	"github.com/dtg-lucifer/everato/pkg"
)

// EventHandler manages all event-related HTTP endpoints
type EventHandler struct {
	Repo     *repository.Queries // Type-safe database queries
	Conn     *pgx.Conn           // Database connection
	BasePath string              // API route prefix
	Cfg      *config.Config      // Application configuration
}

// NewEventHandler initializes the event handler with database connection
func NewEventHandler(cfg *config.Config) *EventHandler {
	logger := pkg.NewLogger()
	defer logger.Close()

	// Connect to PostgreSQL
	conn, err := pgx.Connect(
		context.Background(),
		utils.GetEnv("DB_URL", "postgres://piush:root_access@localhost:5432/everato"),
	)
	if err != nil {
		logger.StdoutLogger.Error("Error connecting to database", "err", err.Error())
		return &EventHandler{Repo: nil}
	}

	// Initialize type-safe repository (generated by SQLC)
	repo := repository.New(conn)
	return &EventHandler{
		Repo:     repo,
		Conn:     conn,
		BasePath: "/events",
		Cfg:      cfg,
	}
}

// RegisterRoutes sets up all event-related routes
func (h *EventHandler) RegisterRoutes(router *mux.Router) {
	events := router.PathPrefix(h.BasePath).Subrouter()

	// Public endpoints (no authentication required)
	events.HandleFunc("/all", h.GetAllEvents).Methods(http.MethodGet)
	events.HandleFunc("/recent", h.GetRecentEvents).Methods(http.MethodGet)
	events.HandleFunc("/{slug}", h.GetBySlug).Methods(http.MethodGet)

	// Protected endpoints (admin authentication required)
	guard := middlewares.NewAdminMiddleware(h.Repo, h.Conn, false)
	protected := events.NewRoute().Subrouter()
	protected.Use(guard.Guard)

	protected.HandleFunc("/create", h.CreateEvent).Methods(http.MethodPost)
	protected.HandleFunc("/{slug}", h.UpdateEventBySlug).Methods(http.MethodPut)
	protected.HandleFunc("/{slug}", h.DeleteEventBySlug).Methods(http.MethodDelete)
	protected.HandleFunc("/{slug}/start", h.StartEvent).Methods(http.MethodPost)
	protected.HandleFunc("/{slug}/end", h.EndEvent).Methods(http.MethodPost)
}

// CreateEvent handles event creation with validation and authorization
func (h *EventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	wr := utils.NewHttpWriter(w, r)

	// Validate database connectivity
	if h.Repo == nil {
		wr.Status(http.StatusBadGateway).Json(
			utils.M{"message": "Database connection unavailable"},
		)
		return
	}

	// Delegate to service layer for business logic
	event.CreateEvent(wr, h.Repo, h.Conn)
}
```

**Design Principles Demonstrated**:
- Clean separation of concerns (handler, service, repository layers)
- Type-safe database operations (SQLC-generated code)
- Middleware-based authentication and authorization
- Comprehensive error handling
- Extensive code documentation
- RESTful API design

#### 4.2.4 Authentication Middleware

Secure JWT-based authentication implementation:

```go
package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/dtg-lucifer/everato/internal/db/repository"
	"github.com/dtg-lucifer/everato/internal/utils"
	"github.com/dtg-lucifer/everato/pkg"
	"github.com/jackc/pgx/v5"
)

const AuthBearerPrefix = "Bearer "

// AuthMiddleware provides JWT-based authentication
type AuthMiddleware struct {
	Repo     *repository.Queries
	Conn     *pgx.Conn
	Redirect bool // true for SSR, false for API
}

func NewAuthMiddleware(repo *repository.Queries, conn *pgx.Conn, redirect bool) *AuthMiddleware {
	return &AuthMiddleware{
		Repo:     repo,
		Conn:     conn,
		Redirect: redirect,
	}
}

// Guard validates JWT tokens and adds user context
func (am *AuthMiddleware) Guard(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := am.extractToken(r)

		key := utils.GetEnv("JWT_SECRET", "SUPER_SECRET_KEY")
		signer := pkg.NewTokenSigner(key)

		if token == "" {
			if isUnauthenticatedAllowedPath(r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}
			am.unauthorized(w, r, "Token not found")
			return
		}

		// Verify and extract claims from JWT
		claims, err := signer.Verify(token)
		if err != nil {
			if isUnauthenticatedAllowedPath(r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}
			am.unauthorized(w, r, "Invalid token")
			return
		}

		// Add user ID to request context for downstream handlers
		userID := claims["user_id"].(string)
		ctx := context.WithValue(r.Context(), "user_id", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
```

**Security Features**:
- JWT token validation
- Context-based user identification
- Public path whitelisting
- Secure token extraction from headers/cookies
- Stateless authentication (no server-side sessions)

#### 4.2.5 Database Queries with SQLC

Example of type-safe SQL queries generated by SQLC:

```sql
-- name: CreateEvent :one
INSERT INTO events (
    id, title, description, slug, start_time, end_time,
    location, capacity, created_by, status
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
) RETURNING *;

-- name: GetEventBySlug :one
SELECT * FROM events WHERE slug = $1 LIMIT 1;

-- name: GetAllEvents :many
SELECT * FROM events
WHERE status = 'published'
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateEvent :one
UPDATE events
SET title = $2, description = $3, start_time = $4,
    end_time = $5, location = $6, capacity = $7
WHERE id = $1
RETURNING *;
```

SQLC generates this type-safe Go code:

```go
// Generated by sqlc. DO NOT EDIT.
type Event struct {
    ID          uuid.UUID
    Title       string
    Description string
    Slug        string
    StartTime   time.Time
    EndTime     time.Time
    Location    string
    Capacity    int32
    CreatedBy   uuid.UUID
    Status      string
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

func (q *Queries) CreateEvent(ctx context.Context, arg CreateEventParams) (Event, error) {
    // Type-safe, compile-time verified SQL execution
}
```

**Benefits**:
- Compile-time SQL validation (errors caught before runtime)
- Full type safety (no `interface{}` or reflection)
- No ORM overhead - direct SQL performance
- IDE autocomplete for database operations
- Impossible to have SQL injection vulnerabilities

### 4.3 Architecture Highlights

#### Clean Separation of Concerns

```
everato/
├── main.go                 # Entry point
├── server.go              # HTTP server setup
├── embedded_fs.go         # Binary embedding
├── config/                # Configuration management
├── internal/
│   ├── handlers/          # HTTP request handlers
│   ├── services/          # Business logic layer
│   ├── middlewares/       # HTTP middlewares
│   ├── db/
│   │   ├── migrations/    # SQL schema migrations
│   │   ├── queries/       # SQL queries for SQLC
│   │   └── repository/    # Generated type-safe code
│   └── utils/             # Helper utilities
├── pkg/                   # Reusable packages
└── www/                   # Frontend application
    ├── src/               # React source code
    ├── public/            # Static assets
    └── dist/              # Production build (embedded)
```

#### Request Flow

1. **HTTP Request** → Server receives request
2. **Middleware Chain** → Request ID, CORS, Logger, Auth
3. **Router** → Gorilla Mux routes to appropriate handler
4. **Handler** → Extracts and validates request data
5. **Service Layer** → Executes business logic
6. **Repository** → Type-safe database operations (SQLC)
7. **Response** → JSON or HTML response returned

---

## 5. Conclusion

Everato represents a paradigm shift in event management platform design. By combining modern technologies with an innovative single-binary distribution model, it delivers enterprise-grade functionality without enterprise-grade complexity.

### Current State and Roadmap

While Everato is actively under development, it already demonstrates the viability and advantages of its core architectural principles. The current implementation includes:

✅ **Completed Features**:
- Single binary compilation with embedded frontend and assets
- Complete event management CRUD operations
- User authentication and authorization with JWT
- Admin dashboard with role-based access control
- Database migrations and seeding
- RESTful API design
- Type-safe database operations with SQLC
- Middleware architecture for cross-cutting concerns
- Structured logging and observability hooks

🚧 **In Progress**:
- Ticketing system with inventory management
- Payment gateway integration
- QR code generation and validation
- Email notification system
- Analytics dashboard
- Advanced search and filtering

🎯 **Planned Features**:
- Cloud-hosted managed service offering
- Plugin architecture for extensibility
- Multi-language support (i18n)
- Mobile application (React Native)
- Webhooks for third-party integrations
- Social media integration
- Calendar synchronization
- Advanced reporting and analytics

### The Open Source Promise

Everato will always remain free and open source. This commitment ensures:

- **Transparency**: Every line of code is available for inspection
- **Community Ownership**: The platform belongs to its users and contributors
- **No Lock-In**: Organizations can fork, modify, and self-host indefinitely
- **Continuous Improvement**: Community contributions drive innovation
- **Educational Value**: Learn modern Go and React architecture from real-world code

### Sustainable Funding Model

To ensure long-term viability and development, Everato will introduce a managed cloud hosting service:

- **Self-Hosting Remains Free**: The core platform will always be free to download and deploy
- **Optional Cloud Service**: Organizations can choose professionally managed hosting
- **Fair Pricing**: Affordable rates that cover infrastructure and support sustainable development
- **Same Software**: Cloud and self-hosted versions use identical code
- **Easy Migration**: Move between self-hosted and cloud seamlessly

This dual model creates a sustainable ecosystem where:
- Developers are compensated for their work
- Organizations without technical resources get professional support
- The open-source community continues to thrive
- Innovation accelerates through diverse funding sources

### Impact and Vision

Everato aims to democratize event management by removing barriers to entry. Whether you're a:

- **Small Non-Profit**: Host fundraisers and community events without expensive platform fees
- **Educational Institution**: Manage campus events and student activities with full data control
- **Tech Startup**: Organize meetups and conferences with a platform you can customize
- **Large Enterprise**: Run corporate events with complete data sovereignty and compliance
- **Government Agency**: Host public events with transparent, auditable open-source software

Everato provides the tools you need with the freedom you deserve.

### Contributing and Community

Everato thrives on community participation. Whether you're a:

- **Developer**: Contribute code, fix bugs, add features
- **Designer**: Improve UI/UX, create themes, design marketing materials
- **Technical Writer**: Enhance documentation, create tutorials
- **Tester**: Report bugs, test new features, provide feedback
- **User**: Share your use cases, request features, spread the word

Every contribution, no matter how small, moves the project forward.

### Final Thoughts

In a world where event management platforms increasingly prioritize profit over user experience and vendor lock-in over user freedom, Everato offers a refreshing alternative. It proves that modern, powerful software can be both open source and sustainable, both sophisticated and simple to deploy.

The single binary distribution model is more than a technical achievement—it's a statement about how software should be built and distributed. It respects users' time, resources, and intelligence. It acknowledges that not everyone wants or needs cloud-based solutions, and that sometimes the best answer is the simplest one: download, run, done.

As Everato continues to evolve, it will maintain these core principles: simplicity, portability, transparency, and user empowerment. The journey from a development project to a production-ready platform is ongoing, but the foundation is solid, the vision is clear, and the future is bright.

**Everato: Your events, your platform, your way.**

---

## Appendix: Quick Start Guide

### For Developers

```bash
# Clone the repository
git clone https://github.com/yourusername/everato.git
cd everato

# Set up environment
cp .env.example .env

# Start development services (PostgreSQL, Kafka)
make db

# Run migrations
make migrate-up

# Start development server with hot reload
make dev

# Frontend development (in separate terminal)
cd www
pnpm install
pnpm dev
```

### For Production Deployment

```bash
# Build the single binary
make build

# The binary includes everything
./bin/everato

# Or with custom config
./bin/everato -config /path/to/config.yaml
```

### For End Users

1. Download the Everato binary for your platform (Linux/Windows/macOS)
2. Download the configuration template
3. Set up your PostgreSQL database
4. Configure database connection in `config.yaml`
5. Run: `./everato`
6. Access the application at `http://localhost:8080`

---

**Document Version**: 1.0
**Last Updated**: January 2026
**Project Status**: Active Development
**License**: MIT
**Repository**: https://github.com/everato-industries/everato
