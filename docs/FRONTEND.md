# Everato Frontend Implementation

## Overview

The Everato platform has been enhanced with a modern React frontend that provides an improved user experience through client-side rendering and interactivity. This document outlines the implementation details, component structure, and integration with the backend API.

## Technology Stack

- **Framework**: React 18 with TypeScript
- **Build Tool**: Vite (for fast development and optimized production builds)
- **Routing**: React Router v7
- **Styling**: Tailwind CSS (utility-first CSS framework)
- **HTTP Client**: Native Fetch API
- **Deployment**: Static files served alongside Go backend

## Project Structure

```
everato/www/
├── public/             # Static assets served as-is
├── src/                # Source code
│   ├── assets/         # Static assets that will be processed
│   ├── components/     # Reusable UI components
│   │   ├── Footer.tsx  # Site footer component
│   │   ├── Layout.tsx  # Layout wrapper component
│   │   └── Navbar.tsx  # Navigation bar component
│   ├── pages/          # Page components
│   │   ├── auth/       # Authentication-related pages
│   │   │   └── login.tsx # Login page component
│   │   └── home.tsx    # Home page component
│   ├── app.tsx         # Main application component
│   ├── index.css       # Global styles (includes Tailwind)
│   ├── main.tsx        # Application entry point
│   ├── routes.tsx      # Route definitions
│   └── vite-env.d.ts   # Vite type definitions
```

## Component Architecture

### Core Components

The frontend is built with a component-based architecture following React best practices:

1. **Layout Component (`Layout.tsx`)**
   - Provides consistent page structure
   - Includes Navbar and Footer components
   - Wraps main content

2. **Navbar Component (`Navbar.tsx`)**
   - Responsive navigation header
   - Contains links to main sections
   - Displays login/signup actions

3. **Footer Component (`Footer.tsx`)**
   - Site-wide footer with navigation links
   - Social media links
   - Copyright information
   - Organized in responsive grid layout

### Page Components

Each route corresponds to a page component:

1. **Home Page (`home.tsx`)**
   - Landing page for the application
   - Displays featured events
   - Makes API request to check backend health

2. **Login Page (`auth/login.tsx`)**
   - Authentication form
   - Form validation
   - API integration with backend auth endpoints
   - Error handling and loading states

## Routing Implementation

Client-side routing is implemented using React Router v7:

```jsx
// routes.tsx
export default function AppRoutes() {
    return (
        <Routes>
            <Route path="/" element={<HomePage />} />
            <Route path="/auth/login" element={<LoginPage />} />

            {/* 404 Not Found route */}
            <Route
                path="*"
                element={
                    <div className="text-center text-red-500">
                        404 Not Found
                    </div>
                }
            />
        </Routes>
    );
}
```

The router is initialized in `main.tsx` with a basename to support deployment in a subdirectory:

```jsx
// main.tsx
createRoot(document.getElementById("root")!).render(
    <StrictMode>
        <BrowserRouter basename="/www">
            <App />
        </BrowserRouter>
    </StrictMode>,
);
```

## API Integration

### Backend Communication

The frontend communicates with the Go backend API using the Fetch API:

```jsx
// Example API call from LoginPage.tsx
const response = await fetch("/api/v1/auth/login", {
    method: "POST",
    headers: {
        "Content-Type": "application/json",
    },
    body: JSON.stringify({ email, password }),
});

if (response.ok) {
    // Handle successful login
    window.location.href = "/";
} else {
    // Handle error
    const data = await response.json();
    setError(data.message || "Login failed. Please try again.");
}
```

### Authentication Flow

1. User enters credentials in the login form
2. Frontend sends credentials to `/api/v1/auth/login`
3. Backend validates credentials and returns JWT token
4. Token is stored in both HTTP-only cookie and local state
5. User is redirected to home page or dashboard

### Health Check

The home page implements a health check to verify API connectivity:

```jsx
useEffect(() => {
    async function fetch_health() {
        try {
            const response = await fetch("/api/v1/health");
            if (!response.ok) {
                console.error("Health check failed");
            }
            const data = await response.json();
            console.log("Health check from API: ", data);
        } catch (error) {
            console.error("Error fetching health:", error);
        }
    }
    fetch_health();
}, []);
```

## Styling Approach

Everato uses Tailwind CSS for utility-first styling:

1. **Global Styles**
   - Basic imports in `index.css`
   - Tailwind's utility classes for responsive design

2. **Component-Level Styling**
   - Inline Tailwind classes for component styling
   - Responsive design using Tailwind's breakpoint utilities
   - Example: `className="container mx-auto p-6"`

## State Management

Currently, state management is handled at the component level using React hooks:

```jsx
const [email, setEmail] = useState("");
const [password, setPassword] = useState("");
const [isLoading, setIsLoading] = useState(false);
const [error, setError] = useState("");
```

As the application grows, more sophisticated state management solutions like Context API or Redux may be implemented.

## Performance Considerations

1. **Code Splitting**
   - Routes are natural boundaries for code splitting
   - Implemented via React Router's lazy loading

2. **Optimized Builds**
   - Vite's production build optimizes for performance
   - Minification and tree-shaking to reduce bundle size

3. **Responsive Design**
   - Mobile-first approach using Tailwind's responsive utilities
   - Performance optimization for various device sizes

## Deployment

The React frontend is built into static assets that are served by the Go backend:

1. Build process creates optimized production assets:
   ```bash
   cd everato/www
   pnpm build
   ```

2. Assets are output to the `dist` directory

3. Go server configured to serve these static files alongside API endpoints

4. All API routes have the prefix `/api/v1/` to distinguish them from frontend routes

## Future Enhancements

1. **State Management**
   - Implement Context API or Redux for global state
   - Add React Query for data fetching and caching

2. **Form Handling**
   - Add form libraries like Formik or React Hook Form
   - Implement validation libraries like Yup or Zod

3. **UI Enhancements**
   - Add a component library like Radix UI or Headless UI
   - Implement dark mode support
   - Add animations and transitions

4. **Testing**
   - Add Jest and React Testing Library for unit tests
   - Implement end-to-end testing with Cypress

## Conclusion

The React frontend implementation for Everato provides a modern, component-based user interface that communicates with the Go backend via RESTful API endpoints. This approach delivers a better user experience through client-side rendering and improved interactivity while maintaining a clear separation of concerns between frontend and backend logic.
