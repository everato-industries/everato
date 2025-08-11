# Everato Frontend

This directory contains the React-based frontend for the Everato event management platform.

## Structure

- `src/` - Source code for the React application
  - `components/` - Reusable UI components
  - `pages/` - Page components that correspond to routes
  - `routes.tsx` - React Router route definitions
  - `app.tsx` - Main application component
  - `main.tsx` - Application entry point

## Components

### Layout Components

- `Layout.tsx` - Main layout wrapper with navbar and footer
- `Navbar.tsx` - Top navigation bar
- `Footer.tsx` - Site footer with links and copyright

### Page Components

- `HomePage.tsx` - Landing page for the application
- `LoginPage.tsx` - User authentication page

## Development

### Running the Development Server

```bash
cd everato/www
pnpm install
pnpm dev
```

### Building for Production

```bash
pnpm build
```

This will create optimized assets in the `dist` directory.

## Integration with Backend

The React frontend communicates with the Go backend via API requests to endpoints defined in `everato/server.go`. Authentication is handled through the `/api/v1/auth/login` endpoint.

## Routing

The application uses React Router v7 for client-side routing. Routes are defined in `routes.tsx` and the router is initialized in `main.tsx`.

## Styling

The application uses Tailwind CSS for styling, which is set up to work with Vite through the `@tailwindcss/vite` plugin.
