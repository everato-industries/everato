# Everato Platform Architecture

## System Overview

Everato is a modern event management platform built with a decoupled architecture consisting of a Go backend API and a React frontend. This document outlines the high-level architecture, component interactions, and technical decisions that shape the platform.

## Architectural Diagram

```mermaid
graph TD
    subgraph "Client Side"
        Browser[Web Browser]
        ReactApp[React SPA]
    end

    subgraph "Server Side"
        APIGateway[API Gateway]

        subgraph "Backend Services"
            AuthService[Auth Service]
            EventService[Event Service]
            AdminService[Admin Service]
            UserService[User Service]
        end

        Database[(PostgreSQL)]
    end

    Browser -->|HTTP/HTTPS| ReactApp
    ReactApp -->|API Requests| APIGateway
    APIGateway -->|Route Requests| AuthService
    APIGateway -->|Route Requests| EventService
    APIGateway -->|Route Requests| AdminService
    APIGateway -->|Route Requests| UserService

    AuthService -->|Query/Update| Database
    EventService -->|Query/Update| Database
    AdminService -->|Query/Update| Database
    UserService -->|Query/Update| Database
```

## Technology Stack

### Frontend
- **Framework**: React with TypeScript
- **Build Tool**: Vite
- **Routing**: React Router v7
- **Styling**: Tailwind CSS
- **API Communication**: Fetch API

### Backend
- **Language**: Go
- **Web Framework**: Custom with Gorilla Mux
- **Database**: PostgreSQL
- **ORM**: sqlc for type-safe SQL
- **Authentication**: JWT-based auth

## Component Breakdown

### Frontend Components

The React application follows a component-based architecture with the following structure:

```mermaid
graph TD
    App[App Component] --> Router[Router]
    Router --> Layout[Layout]
    Layout --> Navbar[Navbar]
    Layout --> MainContent[Main Content]
    Layout --> Footer[Footer]

    MainContent --> HomePage[Home Page]
    MainContent --> LoginPage[Login Page]
    MainContent --> EventsPage[Events Page]
    MainContent --> AdminDashboard[Admin Dashboard]
```

#### Key Frontend Components:

1. **Layout Components**
   - `Layout.tsx`: Main layout wrapper that provides consistent structure
   - `Navbar.tsx`: Navigation component with routing links
   - `Footer.tsx`: Footer with site information and links

2. **Page Components**
   - `HomePage.tsx`: Landing page for the application
   - `LoginPage.tsx`: Authentication page for users
   - Event-related pages for listing, creating, and managing events

3. **Route Configuration**
   - Client-side routing handled by React Router
   - Path definitions in `routes.tsx`

### Backend Components

The Go backend is organized into several logical services, each handling specific domain functionality:

```mermaid
graph TD
    subgraph "API Layer"
        Handlers[HTTP Handlers]
        Middleware[Middleware Chain]
    end

    subgraph "Service Layer"
        AuthService[Auth Service]
        EventService[Event Service]
        AdminService[Admin Service]
        UserService[User Service]
    end

    subgraph "Data Access Layer"
        Repository[SQL Repository]
        Models[Data Models]
    end

    subgraph "Infrastructure"
        Database[(PostgreSQL)]
        Migrations[DB Migrations]
    end

    Handlers --> Middleware
    Middleware --> ServiceLayer[Service Layer]

    ServiceLayer --> AuthService
    ServiceLayer --> EventService
    ServiceLayer --> AdminService
    ServiceLayer --> UserService

    AuthService --> Repository
    EventService --> Repository
    AdminService --> Repository
    UserService --> Repository

    Repository --> Models
    Repository --> Database
    Migrations --> Database
```

#### Key Backend Components:

1. **HTTP Handlers**
   - REST API endpoints organized by resource
   - JSON request/response handling
   - Error handling and status codes

2. **Middleware**
   - Authentication and authorization
   - Request logging
   - CORS handling
   - Request timeout management

3. **Services**
   - Business logic implementation
   - Transaction management
   - Domain-specific validation

4. **Data Access Layer**
   - Type-safe SQL queries using sqlc
   - Repository pattern for database operations
   - Data models representing database entities

## Database Schema

The PostgreSQL database schema includes the following key tables:

```mermaid
erDiagram
    USERS ||--o{ BOOKINGS : makes
    USERS ||--o{ TICKETS : owns
    SUPER_USERS ||--o{ EVENTS : manages
    EVENTS ||--o{ TICKET_TYPES : offers
    EVENTS ||--o{ BOOKINGS : has
    BOOKINGS ||--o{ TICKETS : contains
    TICKETS }|--|| TICKET_TYPES : is_of
    EVENTS ||--o{ COUPONS : provides
    BOOKINGS }o--o{ COUPONS : uses
    TICKETS ||--o{ PAYMENTS : has

    USERS {
        uuid id PK
        string first_name
        string last_name
        string email
        string password
        boolean verified
        timestamp created_at
        timestamp updated_at
    }

    SUPER_USERS {
        uuid id PK
        string email
        string username
        string name
        string password
        enum role
        enum[] permissions
        boolean verified
        timestamp created_at
        timestamp updated_at
    }

    EVENTS {
        uuid id PK
        string title
        string description
        string slug
        string banner
        string icon
        uuid admin_id FK
        timestamp start_time
        timestamp end_time
        string location
        int total_seats
        int available_seats
        timestamp created_at
        timestamp updated_at
    }

    TICKET_TYPES {
        uuid id PK
        string name
        uuid event_id FK
        float price
        int available_tickets
    }

    BOOKINGS {
        uuid id PK
        uuid event_id FK
        uuid user_id FK
        uuid coupon_id FK
        enum status
        timestamp created_at
        timestamp updated_at
    }

    TICKETS {
        uuid id PK
        float price
        enum status
        uuid event_id FK
        uuid user_id FK
        uuid ticket_type FK
        uuid booking_id FK
        string qr_code
    }
```

## Authentication Flow

```mermaid
sequenceDiagram
    actor User
    participant ReactApp as React App
    participant APIGateway as API Gateway
    participant AuthService as Auth Service
    participant Database as Database

    User->>ReactApp: Enter credentials
    ReactApp->>APIGateway: POST /api/v1/auth/login
    APIGateway->>AuthService: Forward request
    AuthService->>Database: Verify credentials
    Database-->>AuthService: User record

    alt Invalid Credentials
        AuthService-->>APIGateway: 401 Unauthorized
        APIGateway-->>ReactApp: 401 Unauthorized
        ReactApp-->>User: Display error message
    else Valid Credentials
        AuthService->>AuthService: Generate JWT
        AuthService-->>APIGateway: 200 OK with JWT
        APIGateway-->>ReactApp: 200 OK with JWT & Set-Cookie
        ReactApp->>ReactApp: Store token
        ReactApp-->>User: Redirect to dashboard
    end
```

## API Integration

The React frontend communicates with the Go backend through RESTful API endpoints. Key integration points include:

1. **Authentication**
   - Login: `POST /api/v1/auth/login`
   - Refresh Token: `POST /api/v1/auth/refresh`
   - Email Verification: `GET /api/v1/auth/verify-email?uid={user_id}`

2. **Events**
   - List Events: `GET /api/v1/events/all`
   - Create Event: `POST /api/v1/events/create`
   - Update Event: `PUT /api/v1/events/update`

3. **Admin Operations**
   - Admin Login: `POST /api/v1/admin/login`
   - Create Admin: `POST /api/v1/admin/create`
   - Get All Admins: `GET /api/v1/admin/all`

## Deployment Architecture

```mermaid
graph TD
    subgraph "Client"
        Browser[Web Browser]
    end

    subgraph "Web Server"
        StaticFiles[Static Files]
        APIServer[Go API Server]
    end

    subgraph "Database Server"
        PostgreSQL[(PostgreSQL)]
    end

    Browser -->|HTTP/HTTPS| StaticFiles
    Browser -->|API Requests| APIServer
    APIServer -->|SQL Queries| PostgreSQL
```

The deployment strategy involves:
1. Building the React application into static files
2. Serving these static files alongside the API server
3. Configuring the server to route API requests to the Go backend
4. Setting up appropriate database connections and security

## Development Workflow

The development workflow involves:

1. **Frontend Development**
   - Running the React dev server with hot reloading
   - Implementing components and pages based on designs
   - Connecting to API endpoints for data fetching and mutations

2. **Backend Development**
   - Implementing API endpoints and business logic
   - Running database migrations for schema changes
   - Testing API endpoints with tools like Postman or curl

3. **Integration**
   - Ensuring API contracts are maintained
   - Testing end-to-end flows across frontend and backend
   - Verifying authentication and authorization works correctly

## Future Considerations

1. **Scalability**
   - Consider breaking down monolithic backend into microservices
   - Implement caching strategies for frequently accessed data
   - Add load balancing for horizontal scaling

2. **Frontend Enhancements**
   - Add state management with Redux or Context API
   - Implement code splitting for improved load times
   - Add comprehensive error handling and recovery

3. **Security Enhancements**
   - Implement rate limiting for API endpoints
   - Add CSRF protection for sensitive operations
   - Regular security audits and penetration testing

This architectural overview provides a high-level understanding of how the Everato platform is structured, with particular emphasis on the integration between the React frontend and Go backend components.
