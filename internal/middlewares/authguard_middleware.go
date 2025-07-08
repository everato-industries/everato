// Package middlewares provides HTTP middleware components for the Everato application.
// These middlewares handle cross-cutting concerns like authentication, logging, and request handling.
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

// AuthBearerPrefix is the standard prefix for Bearer token authentication in the Authorization header
const AuthBearerPrefix = "Bearer "

// AuthMiddleware provides JWT-based authentication for both API and view routes.
// It validates user tokens and either allows the request to proceed or returns an unauthorized response.
type AuthMiddleware struct {
	Repo     *repository.Queries // Repository for database operations
	Conn     *pgx.Conn           // Database connection
	Redirect bool                // Response behavior: true redirects to login (for SSR views), false returns JSON (for APIs)
}

// NewAuthMiddleware creates a new authentication middleware instance.
//
// Parameters:
//   - repo: Repository for database operations (can be nil for public routes)
//   - conn: Database connection (can be nil for public routes)
//   - redirect: Whether to redirect to login page (true) or return JSON response (false) on auth failure
//
// Returns:
//   - A configured AuthMiddleware instance
func NewAuthMiddleware(repo *repository.Queries, conn *pgx.Conn, redirect bool) *AuthMiddleware {
	return &AuthMiddleware{
		Repo:     repo,     // For doing further operations of DB
		Conn:     conn,     // For doing further operations of DB
		Redirect: redirect, // true for SSR views, false for APIs
	}
}

// isUnauthenticatedAllowedPath checks if the given path is accessible without authentication.
// This allows certain routes like login and registration to be accessible to unauthenticated users.
//
// Parameters:
//   - path: The request path to check
//
// Returns:
//   - true if the path can be accessed without authentication, false otherwise
func isUnauthenticatedAllowedPath(path string) bool {
	publicPaths := []string{
		"/auth/login",
		"/auth/register",
	}

	for _, p := range publicPaths {
		if path == p || strings.HasSuffix(path, p) && strings.HasPrefix(path, "/api/") {
			return true
		}
	}
	return false
}

// Guard is the middleware function that intercepts HTTP requests and verifies authentication.
// It extracts and validates JWT tokens from either cookies or the Authorization header.
//
// This middleware:
// 1. Extracts the JWT token from the request
// 2. Verifies the token's validity and expiration
// 3. Extracts the user ID from the token claims
// 4. Adds the user ID to the request context for downstream handlers
// 5. Handles authentication failures based on the configured behavior (redirect or JSON)
//
// Parameters:
//   - next: The next handler in the middleware chain
//
// Returns:
//   - An http.Handler that performs authentication before calling the next handler
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
			am.unauthorized(w, r, "Token not found or malformed")
			return
		}

		claims, err := signer.Verify(token)
		if err != nil {
			if isUnauthenticatedAllowedPath(r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}
			am.unauthorized(w, r, "Invalid or expired token")
			return
		}

		uid, ok := claims["uid"]
		if !ok || uid == "" {
			if isUnauthenticatedAllowedPath(r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}
			am.unauthorized(w, r, "Token missing UID")
			return
		}

		// Add the user ID to the request context for downstream handlers
		ctx := context.WithValue(r.Context(), "uid", uid.(string))
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

// extractToken retrieves the JWT token from either cookies or the Authorization header.
// It prioritizes cookies over headers for authentication source.
//
// Parameters:
//   - r: The HTTP request to extract the token from
//
// Returns:
//   - The extracted token string, or an empty string if no token is found
func (am *AuthMiddleware) extractToken(r *http.Request) string {
	// Try cookie first
	if cookie, err := r.Cookie("jwt"); err == nil {
		return cookie.Value
	}
	// Fallback to Authorization header
	authHeader := r.Header.Get("Authorization")
	if after, ok := strings.CutPrefix(authHeader, AuthBearerPrefix); ok {
		return after
	}
	return ""
}

// unauthorized handles the response for unauthorized requests based on the configured behavior.
// For SSR views, it redirects to the login page. For APIs, it returns a JSON error response.
//
// Parameters:
//   - w: The HTTP response writer
//   - r: The HTTP request
//   - message: The error message to include in the response
func (am *AuthMiddleware) unauthorized(w http.ResponseWriter, r *http.Request, message string) {
	if am.Redirect {
		http.Redirect(w, r, "/auth/login", http.StatusFound)
	} else {
		utils.NewHttpWriter(w, r).Status(http.StatusUnauthorized).Json(utils.M{
			"message": message,
		})
	}
}
