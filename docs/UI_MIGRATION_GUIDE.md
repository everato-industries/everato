# UI Migration Guide: Server-Side Templates to React

## Overview

This guide provides practical steps and best practices for migrating from server-side HTML templates (like Go templates or templ) to a modern React frontend. It addresses the technical challenges, process recommendations, and strategies for a smooth transition.

## Migration Approach

We recommend an incremental migration approach rather than a complete rewrite. This allows for:

- Continuous delivery of improvements
- Reduced risk of regression
- Ability to learn and adjust along the way
- Maintaining application functionality throughout the process

## Step 1: Preparation

### Assessment

1. **Inventory existing templates**
   - Document all existing templates and their dependencies
   - Identify pages with the highest user impact
   - Map out current routing and URL structure

2. **Define API contract**
   - Document data requirements for each page/component
   - Define API endpoints needed to support the React frontend
   - Ensure authentication mechanisms are compatible

3. **Set up development environment**
   - Configure React project using Vite or Create React App
   - Set up TypeScript for type safety
   - Install necessary dependencies (React Router, etc.)
   - Configure build pipeline for integration with backend

### Example API Contract Documentation

```json
{
  "endpoint": "/api/v1/events/:id",
  "method": "GET",
  "requestParams": {
    "id": "UUID of the event"
  },
  "responseFormat": {
    "id": "string (UUID)",
    "title": "string",
    "description": "string",
    "startTime": "string (ISO8601)",
    "endTime": "string (ISO8601)",
    "location": "string",
    "totalSeats": "number",
    "availableSeats": "number"
  },
  "errorResponses": [
    {
      "status": 404,
      "message": "Event not found"
    }
  ]
}
```

## Step 2: Backend Adaptation

### Convert Endpoints to Return JSON

1. **Create API handlers**
   - Implement RESTful endpoints for each data requirement
   - Ensure proper error handling and status codes
   - Add appropriate authentication middleware

2. **Modify existing handlers**
   - Create API versions of existing server-side handlers
   - Return JSON responses instead of rendering templates

### Example Go Handler Conversion

**Before (Server-Side Template):**

```go
func (h *EventHandler) GetEvent(w http.ResponseWriter, r *http.Request) {
    id := mux.Vars(r)["id"]
    event, err := h.repo.GetEventByID(r.Context(), id)
    if err != nil {
        http.Error(w, "Event not found", http.StatusNotFound)
        return
    }

    // Render server-side template
    tmpl, err := template.ParseFiles("templates/event.html")
    if err != nil {
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    tmpl.Execute(w, event)
}
```

**After (JSON API):**

```go
func (h *EventHandler) GetEventAPI(w http.ResponseWriter, r *http.Request) {
    id := mux.Vars(r)["id"]
    event, err := h.repo.GetEventByID(r.Context(), id)
    if err != nil {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusNotFound)
        json.NewEncoder(w).Encode(map[string]string{
            "message": "Event not found",
        })
        return
    }

    // Return JSON response
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(event)
}
```

## Step 3: React Implementation

### Setup React Components

1. **Create component hierarchy**
   - Identify reusable components from templates
   - Build layout components first (header, footer, etc.)
   - Implement page components that match server routes

2. **Implement routing**
   - Set up React Router with routes matching backend URLs
   - Handle authentication and protected routes
   - Implement redirects and error pages

### Example React Component

```jsx
import { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import Layout from '../components/Layout';

function EventPage() {
  const [event, setEvent] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const { id } = useParams();

  useEffect(() => {
    async function fetchEvent() {
      try {
        const response = await fetch(`/api/v1/events/${id}`);
        if (!response.ok) {
          throw new Error('Event not found');
        }
        const data = await response.json();
        setEvent(data);
      } catch (err) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    }

    fetchEvent();
  }, [id]);

  if (loading) return <Layout><div>Loading...</div></Layout>;
  if (error) return <Layout><div>Error: {error}</div></Layout>;

  return (
    <Layout>
      <div className="container mx-auto p-4">
        <h1 className="text-3xl font-bold">{event.title}</h1>
        <p className="mt-2 text-gray-600">{event.description}</p>
        <div className="mt-4">
          <p><strong>Start:</strong> {new Date(event.startTime).toLocaleString()}</p>
          <p><strong>End:</strong> {new Date(event.endTime).toLocaleString()}</p>
          <p><strong>Location:</strong> {event.location}</p>
          <p><strong>Available Seats:</strong> {event.availableSeats} / {event.totalSeats}</p>
        </div>
      </div>
    </Layout>
  );
}

export default EventPage;
```

## Step 4: Authentication Strategy

1. **Implement JWT-based authentication**
   - Create login/logout functionality
   - Store tokens securely (HTTP-only cookies preferred)
   - Add authentication headers to API requests

2. **Handle protected routes**
   - Create authentication context/provider
   - Implement route guards for protected pages
   - Handle session expiration and token refresh

### Example Authentication Context

```jsx
import { createContext, useState, useEffect } from 'react';

export const AuthContext = createContext();

export function AuthProvider({ children }) {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);

  // Check if user is already logged in on mount
  useEffect(() => {
    async function checkAuthStatus() {
      try {
        const response = await fetch('/api/v1/auth/me', {
          credentials: 'include' // Send cookies
        });
        if (response.ok) {
          const userData = await response.json();
          setUser(userData);
        }
      } catch (err) {
        console.error('Auth check failed:', err);
      } finally {
        setLoading(false);
      }
    }

    checkAuthStatus();
  }, []);

  // Login function
  const login = async (email, password) => {
    const response = await fetch('/api/v1/auth/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email, password }),
      credentials: 'include'
    });

    if (!response.ok) {
      throw new Error('Login failed');
    }

    const userData = await response.json();
    setUser(userData);
    return userData;
  };

  // Logout function
  const logout = async () => {
    await fetch('/api/v1/auth/logout', {
      method: 'POST',
      credentials: 'include'
    });
    setUser(null);
  };

  return (
    <AuthContext.Provider value={{ user, loading, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
}
```

## Step 5: Progressive Migration

### Coexistence Strategy

1. **Run both UIs simultaneously**
   - Configure routing to serve React app for migrated routes
   - Keep server-rendered templates for not-yet-migrated routes
   - Use feature flags to control which UI is served

2. **Migration path**
   - Start with simple, standalone pages
   - Progress to more complex, interactive pages
   - Leave critical workflows for last when you have more experience

### Example Backend Router Configuration

```go
func setupRoutes(r *mux.Router) {
    // API endpoints for React frontend
    api := r.PathPrefix("/api/v1").Subrouter()
    api.HandleFunc("/events", GetEventsAPI).Methods("GET")
    api.HandleFunc("/events/{id}", GetEventAPI).Methods("GET")
    // More API endpoints...

    // Serve React frontend for migrated routes
    r.HandleFunc("/events/{id}", ServeReactApp).Methods("GET")
    r.HandleFunc("/dashboard", ServeReactApp).Methods("GET")

    // Original template-based routes for non-migrated pages
    r.HandleFunc("/admin", AdminPageHandler).Methods("GET")
    // More template routes...
}

// ServeReactApp serves the React SPA
func ServeReactApp(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "www/dist/index.html")
}
```

## Step 6: Styling Strategy

1. **Choose a CSS approach**
   - Utility-first CSS (like Tailwind CSS)
   - CSS-in-JS libraries
   - Component-based styling

2. **Maintain visual consistency**
   - Extract common colors, fonts, spacing to variables/tokens
   - Create reusable UI components for common elements
   - Implement responsive designs consistently

### Example Tailwind Configuration

```js
// tailwind.config.js
module.exports = {
  content: [
    "./src/**/*.{js,jsx,ts,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        primary: {
          DEFAULT: '#3b82f6',
          dark: '#2563eb',
          light: '#60a5fa'
        },
        secondary: {
          DEFAULT: '#10b981',
          dark: '#059669',
          light: '#34d399'
        }
      },
      fontFamily: {
        sans: ['Inter', 'sans-serif'],
        heading: ['Poppins', 'sans-serif']
      },
      // Extract values from your current CSS
    }
  },
  plugins: [],
}
```

## Step 7: Testing Strategy

1. **Unit tests**
   - Test React components in isolation
   - Use React Testing Library and Jest
   - Mock API calls and external dependencies

2. **Integration tests**
   - Test API endpoints
   - Ensure correct data formatting
   - Test authentication flows

3. **End-to-end tests**
   - Test critical user journeys
   - Use Cypress or Playwright
   - Verify interactions between frontend and backend

### Example Component Test

```jsx
import { render, screen, waitFor } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import EventPage from './EventPage';

// Mock fetch
global.fetch = jest.fn();

describe('EventPage', () => {
  beforeEach(() => {
    fetch.mockClear();
  });

  it('renders event details when API call succeeds', async () => {
    const mockEvent = {
      id: '123',
      title: 'Test Event',
      description: 'A test event description',
      startTime: '2023-10-15T18:00:00Z',
      endTime: '2023-10-15T20:00:00Z',
      location: 'Test Venue',
      totalSeats: 100,
      availableSeats: 50
    };

    fetch.mockResolvedValueOnce({
      ok: true,
      json: async () => mockEvent
    });

    render(
      <BrowserRouter>
        <EventPage />
      </BrowserRouter>
    );

    // Test loading state
    expect(screen.getByText('Loading...')).toBeInTheDocument();

    // Test loaded content
    await waitFor(() => {
      expect(screen.getByText('Test Event')).toBeInTheDocument();
      expect(screen.getByText('A test event description')).toBeInTheDocument();
      expect(screen.getByText(/Test Venue/)).toBeInTheDocument();
      expect(screen.getByText(/50 \/ 100/)).toBeInTheDocument();
    });
  });

  it('shows error message when API call fails', async () => {
    fetch.mockResolvedValueOnce({
      ok: false
    });

    render(
      <BrowserRouter>
        <EventPage />
      </BrowserRouter>
    );

    await waitFor(() => {
      expect(screen.getByText(/Error:/)).toBeInTheDocument();
    });
  });
});
```

## Step 8: Deployment Strategy

1. **Build integration**
   - Configure build pipeline for React app
   - Automate deployment of React assets
   - Version frontend and backend together

2. **Performance optimization**
   - Enable code splitting and lazy loading
   - Optimize bundle size
   - Implement caching strategies

3. **Monitoring and analytics**
   - Add error tracking
   - Implement performance monitoring
   - Track user behavior analytics

### Example Production Deployment Workflow

1. Build the React application:
   ```bash
   cd www
   npm run build
   ```

2. Copy the static assets to the server's static directory:
   ```bash
   cp -r www/dist/* public/
   ```

3. Configure the Go server to serve the React app for appropriate routes and handle API calls:
   ```go
   // Static files
   r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./public"))))

   // API routes
   apiRouter := r.PathPrefix("/api/v1/").Subrouter()
   apiRouter.HandleFunc("/events", GetEventsAPI).Methods("GET")
   // More API routes...

   // Catch-all route for SPA
   r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
       http.ServeFile(w, r, "./public/index.html")
   })
   ```

## Common Challenges and Solutions

### 1. State Management

**Challenge:** Managing application state that was previously handled by the server.

**Solution:**
- Start with React's useState and useEffect hooks
- Progress to React Context for shared state
- Consider Redux or Zustand for more complex state management
- Use React Query for server state management

### 2. Authentication

**Challenge:** Maintaining authentication state between backend and frontend.

**Solution:**
- Use HTTP-only cookies for secure token storage
- Implement JWT authentication with refresh tokens
- Create protected route components in React

### 3. Form Handling

**Challenge:** Complex form validation previously handled server-side.

**Solution:**
- Use form libraries like Formik or React Hook Form
- Implement client-side validation with libraries like Yup or Zod
- Maintain server-side validation as a backup

### 4. SEO

**Challenge:** React SPAs can be more challenging for search engines.

**Solution:**
- Use React Helmet for managing meta tags
- Consider server-side rendering (SSR) or static site generation (SSG) for critical pages
- Implement proper page titles and meta descriptions

## Conclusion

Migrating from server-side templates to a React frontend is a significant undertaking, but with careful planning and incremental implementation, it can be accomplished successfully. The benefits include improved user experience, better developer productivity, and a more maintainable codebase.

Remember to:

1. Take an incremental approach
2. Maintain functionality throughout the migration
3. Test thoroughly at each step
4. Document API contracts and component interfaces
5. Prioritize user experience in the migration order

By following this guide, you can transition smoothly from server-side rendering to a modern React frontend while minimizing disruption to your users and development team.
