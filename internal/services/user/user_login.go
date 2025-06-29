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

	// Validate the DTO
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

	// Start a transaction
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

	// Get user by email
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

	// Verify password
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

	// Generate JWT Token
	key := utils.GetEnv("JWT_SECRET", "SUPER_SECRET_KEY")
	signer := pkg.NewTokenSigner(key)

	duration := utils.GetEnv("JWT_EXPIRATION", "12h") // Get the duration from the ENV variable
	exp, err := time.ParseDuration(duration)          // Parse the duration string into an real time.Time format
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

	token, err := signer.Sign(jwt.MapClaims{
		"sub": fmt.Sprintf("jwt_login_user_id_%s", user.ID.String()),  // Subject of the token, usually the user ID
		"aud": fmt.Sprintf("jwt_login_user_name_%s", user.FirstName),  // Audience of the token
		"iss": fmt.Sprintf("%s", utils.GetEnv("APP_NAME", "everato")), // Issuer of the token
		"iat": jwt.NewNumericDate(time.Now()),                         // Set the issued at time for the token
		"exp": jwt.NewNumericDate(time.Now().Add(exp)),                // Set the expiration time for the token (i.e by default it is 12h)
		"uid": user.ID.String(),                                       // User ID of the logged in user
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

	err = tx.Commit(wr.R.Context()) // Commit the transaction
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

	logger.StdoutLogger.Info(
		"User logged in successfully",
		"email", user.Email,
		"userId", user.ID,
		"requestId", wr.R.Header.Get("X-Request-ID"),
	)

	wr.Status(http.StatusOK).Json(
		utils.M{
			"message": "User logged in successfully!",
			"token":   token,
			"user": utils.M{
				"id":        user.ID,
				"email":     user.Email,
				"firstName": user.FirstName,
				"lastName":  user.LastName,
				"verified":  user.Verified,
			},
		},
	)
}
