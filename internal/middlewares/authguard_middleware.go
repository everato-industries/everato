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

const (
	AuthBearerPrefix = "Bearer "
)

type AuthGuardMiddleware struct {
	Repo *repository.Queries
	Conn *pgx.Conn
}

func NewAuthGuardMiddleware(repo *repository.Queries, conn *pgx.Conn) *AuthGuardMiddleware {
	return &AuthGuardMiddleware{
		Repo: repo,
		Conn: conn,
	}
}

func (a *AuthGuardMiddleware) AuthGuard(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wr := utils.NewHttpWriter(w, r)

		// Check for the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			wr.Status(http.StatusUnauthorized).Json(
				utils.M{
					"message": "Authorization header is missing",
				},
			)
			return
		}

		// Validate the JWT token (this is a placeholder, implement your own validation logic)
		key := utils.GetEnv("JWT_SECRET", "SUPER_SECRET_KEY")
		signer := pkg.NewTokenSigner(key)

		if len(authHeader) <= len(AuthBearerPrefix) {
			wr.Status(http.StatusUnauthorized).Json(
				utils.M{
					"message": "Authorization header is malformed",
				},
			)
			return
		}

		if !strings.HasPrefix(authHeader, AuthBearerPrefix) {
			wr.Status(http.StatusUnauthorized).Json(
				utils.M{
					"message": "Authorization header must start with 'Bearer '",
				},
			)
			return
		}

		token := authHeader[len(AuthBearerPrefix):] // Get the token as the substring without the bearer prefix

		claims, err := signer.Verify(token)
		if err != nil {
			wr.Status(http.StatusUnauthorized).Json(
				utils.M{
					"message": "Invalid or expired token",
				},
			)
			return
		}

		if uid, ok := claims["uid"]; !ok || uid == "" {
			wr.Status(http.StatusUnauthorized).Json(
				utils.M{
					"message": "Token does not contain a valid user ID",
				},
			)
			return
		}

		// Check if the UID is valid or not

		// Set the user ID in the request context for further use
		new_ctx := context.WithValue(r.Context(), "uid", claims["uid"].(string))

		// Create a new request with the updated context
		new_request, err := http.NewRequestWithContext(
			new_ctx,
			r.Method,
			r.URL.String(),
			r.Body,
		)
		if err != nil {
			wr.Status(http.StatusUnauthorized).Json(
				utils.M{
					"message": "Failed to create a new request with context",
				},
			)
			return
		}

		// If the token is valid, proceed to the next handler
		h.ServeHTTP(w, new_request) // Set the request to be the new updated one
	})
}
