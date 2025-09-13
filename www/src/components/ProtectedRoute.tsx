import React from "react";
import { Navigate, useLocation } from "react-router-dom";
import { isAuthenticated } from "../lib/auth";

interface ProtectedRouteProps {
    children: React.ReactNode;
    redirectTo?: string;
}

/**
 * ProtectedRoute component that requires admin authentication
 * Redirects to admin login page if user is not authenticated
 */
export default function ProtectedRoute({
    children,
    redirectTo = "/admin",
}: ProtectedRouteProps) {
    const location = useLocation();

    if (!isAuthenticated()) {
        // Redirect to login page with the current location as state
        // so we can redirect back after successful login
        return (
            <Navigate
                to={redirectTo}
                state={{ from: location }}
                replace
            />
        );
    }

    return <>{children}</>;
}
