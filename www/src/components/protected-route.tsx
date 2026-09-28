import React from "react";
import { Navigate, useLocation } from "react-router-dom";
import { useAuth } from "../hooks/useAuth";

interface ProtectedRouteProps {
    children: React.ReactNode;
    /** Where to redirect unauthenticated users. Defaults to the user login page. */
    redirectTo?: string;
}

/**
 * ProtectedRoute — guards routes that require a logged-in (non-admin) user.
 * Saves the original location so login can redirect back after success.
 */
export default function ProtectedRoute({
    children,
    redirectTo = "/auth/login",
}: ProtectedRouteProps) {
    const { isAuthenticated } = useAuth();
    const location = useLocation();

    if (!isAuthenticated) {
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
