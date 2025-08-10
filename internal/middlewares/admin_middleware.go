// Package middlewares provides HTTP middleware components for the Everato application.
// AdminMiddleware specifically handles authentication for admin-only routes.
package middlewares

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/dtg-lucifer/everato/internal/db/repository"
	"github.com/dtg-lucifer/everato/internal/utils"
	"github.com/dtg-lucifer/everato/pkg"
	"github.com/jackc/pgx/v5"
)

// AdminMiddleware provides JWT-based authentication specifically for admin routes.
// It validates admin tokens and either allows the request to proceed or returns an unauthorized response.
type AdminMiddleware struct {
	Repo     *repository.Queries // Repository for database operations
	Conn     *pgx.Conn           // Database connection
	Redirect bool                // Response behavior: true redirects to login (for SSR views), false returns JSON (for APIs)
}

// NewAdminMiddleware creates a new admin authentication middleware instance.
//
// Parameters:
//   - repo: Repository for database operations
//   - conn: Database connection
//   - redirect: Whether to redirect to login page (true) or return JSON response (false) on auth failure
//
// Returns:
//   - A configured AdminMiddleware instance
func NewAdminMiddleware(repo *repository.Queries, conn *pgx.Conn, redirect bool) *AdminMiddleware {
	return &AdminMiddleware{
		Repo:     repo,     // For doing further operations of DB
		Conn:     conn,     // For doing further operations of DB
		Redirect: redirect, // true for SSR views, false for APIs
	}
}

// isAdminPublicPath checks if the given path is accessible without admin authentication.
// This allows only the admin login route to be accessible without authentication.
//
// Parameters:
//   - path: The request path to check
//
// Returns:
//   - true if the path can be accessed without admin authentication, false otherwise
func isAdminPublicPath(path string) bool {
	// Only the login route is public for admins
	if strings.HasSuffix(path, "/login") {
		pathParts := strings.Split(path, "/")
		for i, part := range pathParts {
			if part == "admin" && i < len(pathParts)-1 && pathParts[i+1] == "login" {
				return true
			}
		}
	}
	return false
}

// Guard is the middleware function that intercepts HTTP requests and verifies admin authentication.
// It extracts and validates JWT tokens, ensuring the user has admin privileges.
//
// This middleware:
// 1. Extracts the JWT token from the request
// 2. Verifies the token's validity and expiration
// 3. Extracts the user ID and verifies admin status
// 4. Adds the admin ID to the request context for downstream handlers
// 5. Handles authentication failures based on the configured behavior (redirect or JSON)
//
// Parameters:
//   - next: The next handler in the middleware chain
//
// Returns:
//   - An http.Handler that performs authentication before calling the next handler
func (am *AdminMiddleware) Guard(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip authentication for public admin paths (login)
		if isAdminPublicPath(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		token := am.extractToken(r)

		key := utils.GetEnv("JWT_SECRET", "SUPER_SECRET_KEY")
		signer := pkg.NewTokenSigner(key)

		if token == "" {
			am.unauthorized(w, r, "Admin authentication required")
			return
		}

		claims, err := signer.Verify(token)
		if err != nil {
			am.unauthorized(w, r, "Invalid or expired admin token")
			return
		}

		// Extract UID from claims
		uid, ok := claims["uid"]
		if !ok || uid == "" {
			am.unauthorized(w, r, "Token missing UID")
			return
		}

		// // Extract role from claims (if it exists)
		// role, roleOk := claims["role"]

		// // Verify that the user is an admin
		// // First check if role claim exists and has admin value
		// if roleOk && role == "admin" {
		// 	// Add the admin ID to the request context for downstream handlers
		// 	ctx := context.WithValue(r.Context(), "admin_id", uid.(string))
		// 	ctx = context.WithValue(ctx, "is_admin", true)
		// 	r = r.WithContext(ctx)
		// 	next.ServeHTTP(w, r)
		// 	return
		// }

		// If no role claim or role is not admin, verify against database
		if am.Repo != nil && am.Conn != nil {
			isAdmin, err := am.verifyAdminStatus(r.Context(), uid.(string))
			if err != nil {
				logger := pkg.NewLogger()
				defer logger.Close()
				logger.StdoutLogger.Error("Error verifying admin status", "error", err.Error())
				am.unauthorized(w, r, "Error verifying admin status")
				return
			}

			if isAdmin {
				// Add the admin ID to the request context for downstream handlers
				ctx := context.WithValue(r.Context(), "admin_id", uid.(string))
				new_ctx := context.WithValue(ctx, "is_admin", true)
				new_r := r.WithContext(new_ctx)
				next.ServeHTTP(w, new_r)
				return
			}
		}

		// Not an admin
		am.unauthorized(w, r, "Admin privileges required")
	})
}

// verifyAdminStatus checks if the user with the given ID has admin privileges.
//
// Parameters:
//   - ctx: The context for the database operation
//   - uid: The user ID to check
//
// Returns:
//   - bool: true if the user is an admin, false otherwise
//   - error: Any error that occurred during verification
func (am *AdminMiddleware) verifyAdminStatus(ctx context.Context, uid string) (bool, error) {
	// TODO: Implement actual admin verification logic against the database
	logger := pkg.NewLogger()
	defer logger.Close()

	if am.Repo == nil {
		logger.StdoutLogger.Error("Database repository is nil, cannot verify admin status")
		return false, nil
	}

	s_uid, err := utils.StringToUUID(uid)
	if err != nil {
		logger.StdoutLogger.Error("Error converting UID to UUID", "uid", uid, "err", err.Error())
		return false, err
	}

	// Get admin by ID without using a transaction for this simple read operation
	admin, err := am.Repo.GetAdminById(ctx, s_uid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.StdoutLogger.Info("Admin not found", "uid", uid)
			return false, nil // Not an admin
		}
		logger.StdoutLogger.Error("Error retrieving admin by ID", "uid", uid, "err", err.Error())
		return false, err // Database error
	}

	// If we got here, the admin exists
	logger.StdoutLogger.Info("Admin verified successfully", "uid", uid, "role", admin.Role)
	return true, nil
}

// extractToken retrieves the JWT token from either cookies or the Authorization header.
// It prioritizes cookies over headers for authentication source.
//
// Parameters:
//   - r: The HTTP request to extract the token from
//
// Returns:
//   - The extracted token string, or an empty string if no token is found
func (am *AdminMiddleware) extractToken(r *http.Request) string {
	// Try admin-specific cookie first
	if cookie, err := r.Cookie("admin_jwt"); err == nil {
		return cookie.Value
	}

	// Try regular JWT cookie next
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
func (am *AdminMiddleware) unauthorized(w http.ResponseWriter, r *http.Request, message string) {
	if am.Redirect {
		http.Redirect(w, r, "/admin/login", http.StatusFound)
	} else {
		utils.NewHttpWriter(w, r).Status(http.StatusUnauthorized).Json(utils.M{
			"message": message,
			"status":  "unauthorized",
		})
	}
}
