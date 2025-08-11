# Everato Frontend

This directory contains the React-based frontend for the Everato event management platform.

## Overview

The Everato frontend is built as a modern Single Page Application (SPA) using React and TypeScript. It provides an interactive, responsive user interface that communicates with the Go backend API to deliver a seamless event management experience.

## Technology Stack

- **Framework**: React 18 with TypeScript
- **Build Tool**: Vite (for fast development and optimized production builds)
- **Routing**: React Router v7
- **Styling**: Tailwind CSS (utility-first CSS framework)
- **HTTP Client**: Native Fetch API

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

## Getting Started

### Prerequisites

- Node.js 18+
- pnpm (recommended) or npm

### Installation

```bash
# Install dependencies
pnpm install
```

### Development

```bash
# Start development server
pnpm dev
```

This will start the Vite development server, typically on http://localhost:3000. The development server includes:

- Hot Module Replacement for fast updates
- Error overlay for easy debugging
- API proxy configuration to the backend

### Building for Production

```bash
# Build for production
pnpm build
```

This creates optimized production files in the `dist` directory that can be served by any static file server, including the Go backend.

## Component Architecture

### Core Components

The frontend is built with a component-based architecture following React best practices:

1. **Layout Component (`Layout.tsx`)**
    - Provides consistent page structure with Navbar and Footer
    - Serves as a wrapper for all page content

2. **Navbar Component (`Navbar.tsx`)**
    - Responsive navigation header
    - Contains links to main sections and authentication actions

3. **Footer Component (`Footer.tsx`)**
    - Site-wide footer with navigation and company information
    - Social media links and copyright details

### Page Components

Each route corresponds to a page component:

1. **Home Page (`home.tsx`)**
    - Landing page with featured events
    - Makes API request to check backend health

2. **Login Page (`auth/login.tsx`)**
    - User authentication form
    - Form validation and error handling
    - API integration with backend auth endpoints

## Routing Implementation

Client-side routing is implemented using React Router v7:

```jsx
// routes.tsx
export default function AppRoutes() {
    return (
        <Routes>
            <Route path="/" element={<HomePage />} />
            <Route path="/auth/login" element={<LoginPage />} />
            <Route path="*" element={<NotFoundPage />} />
        </Routes>
    );
}
```

The router is initialized in `main.tsx`:

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

### Key API Endpoints

- **Authentication**: `/api/v1/auth/login`, `/api/v1/auth/refresh`
- **Events**: `/api/v1/events/all`, `/api/v1/events/create`
- **Health Check**: `/api/v1/health`

## Styling

Everato uses Tailwind CSS for utility-first styling:

```jsx
// Example of Tailwind classes
<div className="container mx-auto p-6">
    <h2 className="text-blue-400 text-3xl font-bold mb-4">
        Welcome to Everato!
    </h2>
    <p className="text-gray-700 mb-8">Your modern event management platform.</p>
</div>
```

## State Management

Currently, component-level state management is implemented using React hooks:

```jsx
const [email, setEmail] = useState("");
const [password, setPassword] = useState("");
const [isLoading, setIsLoading] = useState(false);
const [error, setError] = useState("");
```

## Planned Enhancements

1. **Global State Management**
    - Implementation of Context API or Redux
    - React Query for data fetching and caching

2. **Form Handling**
    - Integration of form libraries like Formik or React Hook Form
    - Validation with libraries like Yup or Zod

3. **Component Library**
    - Implementing or integrating a component library
    - Dark mode support

4. **Testing**
    - Jest and React Testing Library for unit tests
    - Cypress for end-to-end testing

## Integration with Backend

The frontend is designed to work seamlessly with the Everato Go backend. In development mode, API requests are proxied to the backend server. In production, the React build artifacts are served directly by the Go backend.

## Browser Support

The application is designed to support modern browsers:

- Chrome/Edge (latest 2 versions)
- Firefox (latest 2 versions)
- Safari (latest 2 versions)

## Contributing

When contributing to the frontend, please follow these guidelines:

1. Follow the established component structure
2. Use TypeScript for all new components
3. Maintain consistent styling with Tailwind CSS
4. Write clean, readable, and maintainable code
5. Consider accessibility in all UI components

## Resources

- [React Documentation](https://react.dev/)
- [TypeScript Documentation](https://www.typescriptlang.org/docs/)
- [Vite Documentation](https://vitejs.dev/guide/)
- [React Router Documentation](https://reactrouter.com/en/main)
- [Tailwind CSS Documentation](https://tailwindcss.com/docs)
