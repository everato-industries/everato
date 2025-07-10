package admin

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dtg-lucifer/everato/config"
	"github.com/dtg-lucifer/everato/internal/db/repository"
	"github.com/dtg-lucifer/everato/internal/utils"
	"github.com/dtg-lucifer/everato/pkg"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func Login(wr *utils.HttpWriter, repo *repository.Queries, conn *pgx.Conn, cfg *config.Config) {
	// This function will handle the login logic for admin accounts.
	logger := pkg.NewLogger()
	defer logger.Close()

	// Get the DTO
	login_dto := NewAdminLoginDTO()
	err := wr.ParseBody(&login_dto)
	if err != nil {
		logger.StdoutLogger.Error("Error parsing login request body", "err", err.Error())
		wr.Status(http.StatusBadRequest).Json(
			utils.M{
				"message": "Invalid request body",
				"error":   err.Error(),
			},
		)
		return
	}

	// Validate the DTO using struct validation rules
	err = login_dto.Validate()
	if err != nil {
		logger.StdoutLogger.Error("Validation error in login request", "err", err.Error())
		wr.Status(http.StatusBadRequest).Json(
			utils.M{
				"message": "Validation error",
				"error":   err.Error(),
			},
		)
		return
	}

	// First line of defense, if the provided username and email is not a match from the config
	// then it is not an authorized user
	//
	//
	var match bool
	for i, su := range cfg.SuperUsers {
		if su.Email == login_dto.Email {
			match = true
			logger.StdoutLogger.Info("Super user found", "index", i, "email", su.Email)
			break
		}
	}
	if !match {
		logger.StdoutLogger.Error("Super user not found", "email", login_dto.Email)
		wr.Status(http.StatusUnauthorized).Json(
			utils.M{
				"message": "Unauthorized",
				"error":   "Super user not found",
			},
		)
		return
	}

	// *
	// Therefore start searching the database for confirmation if the details are valid or not
	// *

	// start the transaction for atomicity
	tx, err := conn.Begin(wr.R.Context())
	if err != nil {
		logger.StdoutLogger.Error("Error starting a transaction", "err", err.Error())
		wr.Status(http.StatusInternalServerError).Json(
			utils.M{
				"message": "Internal server error",
				"error":   err.Error(),
			},
		)
		return
	}

	// // Hash the password of the admin
	// err = login_dto.HashPassword()
	// if err != nil {
	// 	logger.StdoutLogger.Error("Error hashing the password", "err", err.Error())
	// 	tx.Rollback(wr.R.Context()) // Rollback the transaction
	// 	wr.Status(http.StatusInternalServerError).Json(
	// 		utils.M{
	// 			"message": "Internal server error",
	// 			"error":   err.Error(),
	// 		},
	// 	)
	// 	return
	// }

	admin_by_email, err := repo.WithTx(tx).GetAdminByEmail(wr.R.Context(), login_dto.Email)
	if err != nil {
		logger.StdoutLogger.Error("Error getting admin by email", "err", err.Error())
		tx.Rollback(wr.R.Context()) // Rollback the transaction if there is an error
		wr.Status(http.StatusInternalServerError).Json(
			utils.M{
				"message": "Internal server error",
				"error":   err.Error(),
			},
		)
		return
	}

	if ok, err := login_dto.VerifyPassword(admin_by_email.Password); err != nil && !ok {
		tx.Rollback(wr.R.Context()) // Rollback the transaction if the password is incorrect
		wr.Status(http.StatusUnauthorized).Json(
			utils.M{
				"message": "Invalid credentials",
				"error":   "Incorrect password",
			},
		)
		return
	}

	var admin_uid pgtype.UUID // Uid for the admin

	admin_uid = admin_by_email.ID // Get the admin ID from the email query result

	// Create the JWT token
	key := utils.GetEnv("JWT_SECRET", "SUPER_SECRET_KEY")
	signer := pkg.NewTokenSigner(key)

	// Get the token duration from environment variables with fallback to 12 hours
	duration := utils.GetEnv("JWT_EXPIRATION", "12h")
	exp, err := time.ParseDuration(duration)
	if err != nil {
		logger.StdoutLogger.Error(
			"Error parsing JWT expiration duration",
			"err", err.Error(),
			"requestId", wr.R.Header.Get("X-Request-ID"),
		)
		wr.Status(http.StatusInternalServerError).Json(
			utils.M{
				"error":   "Internal server error",
				"message": "Failed to parse JWT expiration duration",
			},
		)
		return
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
		"sub": fmt.Sprintf("jwt_login_user_id_%s", admin_uid.String()),
		"aud": fmt.Sprintf("jwt_login_user_name_%s", admin_by_email.Username),
		"iss": fmt.Sprintf("%s", utils.GetEnv("APP_NAME", "everato")),
		"iat": jwt.NewNumericDate(time.Now()),
		"exp": jwt.NewNumericDate(exp_time),
		"uid": admin_uid.String(),
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
		"Admin logged in successfully",
		"email", login_dto.Email,
		"admin_id", admin_uid.String(),
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
			Domain:   wr.R.URL.Hostname(),  // Uncomment to restrict cookie to specific domain
		},
	)

	// Return successful response with token and user information
	// This allows both cookie-based and bearer token authentication options
	wr.Status(http.StatusOK).Json(
		utils.M{
			"message": "User logged in successfully!",
			"token":   token, // Include token in response for clients that prefer Authorization header
			"user": utils.M{
				"id":          admin_uid.String(),
				"email":       admin_by_email.Email,
				"username":    admin_by_email.Username,
				"role":        admin_by_email.Role,
				"permissions": admin_by_email.Permissions,
			},
		},
	)
}
