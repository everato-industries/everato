// Package user provides authentication and user management services for the Everato platform.
// It handles user creation, verification, authentication, and token management.
package user

import (
	"net/http"
	"os"
	"time"

	"github.com/dtg-lucifer/everato/internal/db/repository"
	"github.com/dtg-lucifer/everato/internal/utils"
	"github.com/dtg-lucifer/everato/pkg"
	"github.com/golang-jwt/jwt/v5"
)

// RefreshUserToken handles the refresh token logic for sub-admins and users.
// It retrieves the JWT token from cookies, validates it, and issues a new token with extended expiration.
//
// The function performs the following operations:
// 1. Extracts the JWT token from the request's cookies
// 2. Verifies the token's signature and validity
// 3. Extracts the user ID from the token claims
// 4. Issues a new token with extended expiration (24 hours)
// 5. Sets the new token as a HTTP-only cookie
// 6. Returns the new token in the response for API clients
//
// Parameters:
//   - wr: Custom HTTP writer for response handling
//   - repo: Database repository for user operations (unused in this function but kept for consistency)
//   - r: HTTP request containing the token cookie
func RefreshUserToken(wr *utils.HttpWriter, repo *repository.Queries, r *http.Request) {
	// Get token from cookie instead of request body
	// This prioritizes cookie-based authentication for better security
	cookie, err := r.Cookie("jwt")
	if err != nil || cookie.Value == "" {
		wr.Status(http.StatusBadRequest).Json(utils.M{"error": "Missing JWT token in cookie"})
		return
	}

	// Retrieve JWT secret from environment variables with fallback to default
	// In production, the secret should always be set via environment variable
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "jakdjf87a8((*8___jhadja#kjaaaoitrak" // fallback default
	}
	// Initialize the token signer with the JWT secret
	signer := pkg.NewTokenSigner(jwtSecret)

	// Verify the token signature and validity
	// This ensures the token hasn't been tampered with and isn't expired
	claims, err := signer.Verify(cookie.Value)
	if err != nil {
		wr.Status(http.StatusUnauthorized).Json(utils.M{"error": "Invalid or expired refresh token"})
		return
	}

	// Check for user ID in the claims
	// The uid claim is essential for identifying which user is refreshing the token
	userID, ok := claims["uid"].(string)
	if !ok || userID == "" {
		wr.Status(http.StatusUnauthorized).Json(utils.M{"error": "Invalid token claims"})
		return
	}

	// Issue a new token with extended expiration (24 hours)
	// This refreshes the authentication session, preventing premature expiration
	expiresAt := time.Now().Add(24 * time.Hour).Unix()
	newClaims := jwt.MapClaims{
		"uid": userID,    // Preserve the user ID from the original token
		"exp": expiresAt, // Set new expiration time 24 hours from now
	}

	// Sign the new token with the JWT secret
	newToken, err := signer.Sign(newClaims)
	if err != nil {
		wr.Status(http.StatusInternalServerError).Json(utils.M{"error": "Failed to sign new token"})
		return
	}

	// Set the new token as an HTTP-only cookie for enhanced security
	// HTTP-only cookies cannot be accessed by JavaScript, protecting against XSS attacks
	http.SetCookie(wr.W, &http.Cookie{
		Name:     "jwt",                          // Standard name for JWT cookies
		Value:    newToken,                       // The newly generated token
		Expires:  time.Now().Add(24 * time.Hour), // Cookie expiration matches token expiration
		HttpOnly: true,                           // Prevents JavaScript access to the cookie
		Path:     "/",                            // Makes the cookie available across the entire site
		// Secure: true,                     // Uncomment in production to restrict to HTTPS
		// SameSite: http.SameSiteStrictMode,// Uncomment in production for stronger CSRF protection
	})

	// Return successful response with the new token
	// This allows both cookie-based and header-based authentication options
	wr.Status(http.StatusOK).Json(utils.M{
		"token": newToken,  // Include token in response for clients using Authorization header
		"exp":   expiresAt, // Include expiration timestamp for client-side expiration handling
	})
}
