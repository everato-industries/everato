// Package user provides authentication and user management services for the Everato platform.
// It handles user creation, verification, authentication, and session management.
package user

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dtg-lucifer/everato/internal/db/repository"
	"github.com/dtg-lucifer/everato/internal/utils"
	"github.com/dtg-lucifer/everato/pkg"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
)

// LoginUser handles the authentication process for existing users.
// It validates credentials, issues JWT tokens, and manages authentication sessions.
//
// The function performs the following operations:
// 1. Parses and validates login credentials from the request body
// 2. Retrieves the user from the database by email
// 3. Verifies the password using bcrypt comparison
// 4. Generates a JWT token with appropriate claims and expiration
// 5. Sets the token as an HTTP cookie and returns it in the response body
//
// Parameters:
//   - wr: Custom HTTP writer for response handling
//   - repo: Database repository for user operations
//   - conn: Database connection for transaction management
func LoginUser(wr *utils.HttpWriter, repo *repository.Queries, conn *pgx.Conn) {
	logger := pkg.NewLogger()
	defer logger.Close()

	// Parse login credentials from request body
	loginDTO := &LoginUserDTO{}
	if err := wr.ParseBody(loginDTO); err != nil {
		logger.StdoutLogger.Error(
			"Error parsing login request body",
			"err", err.Error(),
			"requestId", wr.R.Header.Get("X-Request-ID"),
		)
		wr.Status(http.StatusBadRequest).Json(
			utils.M{
				"error":   "Invalid request body",
				"message": err.Error(),
			},
		)
		return
	}

	// Validate the DTO using struct validation rules
	// This ensures email format is valid and password meets minimum requirements
	if err := loginDTO.Validate(); err != nil {
		logger.StdoutLogger.Error(
			"Invalid login credentials format",
			"err", err.Error(),
			"requestId", wr.R.Header.Get("X-Request-ID"),
		)
		wr.Status(http.StatusBadRequest).Json(
			utils.M{
				"error":   "Invalid credentials format",
				"message": err.Error(),
			},
		)
		return
	}

	// Start a database transaction for atomicity
	// This ensures either all operations succeed or none do
	tx, err := conn.Begin(wr.R.Context())
	if err != nil {
		logger.StdoutLogger.Error(
			"Error starting a transaction",
			"err", err.Error(),
			"requestId", wr.R.Header.Get("X-Request-ID"),
		)
		wr.Status(http.StatusInternalServerError).Json(
			utils.M{
				"error":   "Internal server error",
				"message": "Failed to start transaction",
			},
		)
		return
	}

	// Get user by email from the database
	// Note: We use the same error message for both email not found and password mismatch
	// to avoid revealing which one was incorrect (security best practice)
	user, err := repo.WithTx(tx).GetUserByEmail(wr.R.Context(), loginDTO.Email)
	if err != nil {
		logger.StdoutLogger.Info(
			"Failed login attempt - email not found",
			"email", loginDTO.Email,
			"requestId", wr.R.Header.Get("X-Request-ID"),
		)
		wr.Status(http.StatusUnauthorized).Json(
			utils.M{
				"error":   "Authentication failed",
				"message": "Invalid email or password",
			},
		)
		return
	}

	// Verify password using bcrypt comparison
	// This securely compares the provided password against the hashed password in the database
	err = loginDTO.VerifyPassword(user.Password)
	if err != nil {
		logger.StdoutLogger.Info(
			"Failed login attempt - password mismatch",
			"email", loginDTO.Email,
			"requestId", wr.R.Header.Get("X-Request-ID"),
		)
		wr.Status(http.StatusUnauthorized).Json(
			utils.M{
				"error":   "Authentication failed",
				"message": "Invalid email or password",
			},
		)
		return
	}

	// Generate JWT Token for authentication
	// This creates a signed token containing user identity and session information
	key := utils.GetEnv("JWT_SECRET", "SUPER_SECRET_KEY")
	signer := pkg.NewTokenSigner(key)

	// Get the token duration from environment variables with fallback to 12 hours
	duration := utils.GetEnv("JWT_EXPIRATION", "12h")
	exp, err := time.ParseDuration(duration)
	if err != nil {
		logger.StdoutLogger.Warn(
			"Invalid JWT_EXPIRATION, falling back to 12h",
			"provided", duration,
			"err", err.Error(),
		)
		exp = 12 * time.Hour
	}

	// Calculate absolute expiration time from now
	exp_time := time.Now().Add(exp)

	// Create token with standard JWT claims following best practices:
	// - sub: Subject identifier (user ID)
	// - aud: Audience (application using the token)
	// - iss: Issuer (application name)
	// - iat: Issued at time
	// - exp: Expiration time
	// - uid: Custom claim with user ID for application use
	token, err := signer.Sign(jwt.MapClaims{
		"sub": fmt.Sprintf("jwt_login_user_id_%s", user.ID.String()),
		"aud": fmt.Sprintf("jwt_login_user_name_%s", user.FirstName),
		"iss": fmt.Sprintf("%s", utils.GetEnv("APP_NAME", "everato")),
		"iat": jwt.NewNumericDate(time.Now()),
		"exp": jwt.NewNumericDate(exp_time),
		"uid": user.ID.String(),
	})
	if err != nil {
		logger.StdoutLogger.Error(
			"Error signing JWT token",
			"err", err.Error(),
			"requestId", wr.R.Header.Get("X-Request-ID"),
		)
		wr.Status(http.StatusInternalServerError).Json(
			utils.M{
				"error":   "Internal server error",
				"message": "Failed to generate JWT token",
			},
		)
		return
	}

	// Commit the transaction to finalize login
	err = tx.Commit(wr.R.Context())
	if err != nil {
		logger.StdoutLogger.Error(
			"Error committing transaction",
			"err", err.Error(),
			"requestId", wr.R.Header.Get("X-Request-ID"),
		)
		wr.Status(http.StatusInternalServerError).Json(
			utils.M{
				"error":   "Internal server error",
				"message": "Failed to commit transaction",
			},
		)
		return
	}

	// Log successful login for audit and monitoring
	logger.StdoutLogger.Info(
		"User logged in successfully",
		"email", user.Email,
		"userId", user.ID,
		"requestId", wr.R.Header.Get("X-Request-ID"),
	)

	// Set JWT as an HTTP-only cookie for session management
	// HTTP-only prevents JavaScript access, providing protection against XSS attacks
	wr.SetCookie(
		utils.CookieParams{
			Name:     "jwt",
			Value:    token,
			MaxAge:   int(exp.Seconds()),
			Path:     "/",
			Secure:   false,                // TODO: Set to true in production for HTTPS environments
			SameSite: http.SameSiteLaxMode, // Provides CSRF protection while allowing redirects
			HttpOnly: true,                 // Prevents JavaScript access to the cookie
			// Domain:   wr.R.Host, // Uncomment to restrict cookie to specific domain
		},
	)

	// Return successful response with token and user information
	// This allows both cookie-based and bearer token authentication options
	wr.Status(http.StatusOK).Json(
		utils.M{
			"message": "User logged in successfully!",
			"token":   token, // Include token in response for clients that prefer Authorization header
			"user": utils.M{ // Include minimal user data needed by frontend
				"id":        user.ID,
				"email":     user.Email,
				"firstName": user.FirstName,
				"lastName":  user.LastName,
				"verified":  user.Verified,
			},
		},
	)
}
