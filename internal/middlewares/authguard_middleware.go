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

type AuthMiddleware struct {
	Repo     *repository.Queries
	Conn     *pgx.Conn
	Redirect bool // true for SSR views, false for APIs
}

func NewAuthMiddleware(repo *repository.Queries, conn *pgx.Conn, redirect bool) *AuthMiddleware {
	return &AuthMiddleware{
		Repo:     repo,     // For doing further operations of DB
		Conn:     conn,     // For doing further operations of DB
		Redirect: redirect, // true for SSR views, false for APIs
	}
}

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

		ctx := context.WithValue(r.Context(), "uid", uid.(string))
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

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

func (am *AuthMiddleware) unauthorized(w http.ResponseWriter, r *http.Request, message string) {
	if am.Redirect {
		http.Redirect(w, r, "/auth/login", http.StatusFound)
	} else {
		utils.NewHttpWriter(w, r).Status(http.StatusUnauthorized).Json(utils.M{
			"message": message,
		})
	}
}
