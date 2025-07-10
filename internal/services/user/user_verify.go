// Package user provides authentication and user management services for the Everato platform.
// It handles user creation, verification, authentication, and profile management.
package user

import (
	"net/http"

	"github.com/dtg-lucifer/everato/internal/db/repository"
	"github.com/dtg-lucifer/everato/internal/utils"
	"github.com/dtg-lucifer/everato/pkg"
	"github.com/jackc/pgx/v5"
)

// VerifyUser handles the email verification process for user accounts.
// It validates the user ID from the URL, checks verification status, and updates
// the user record to mark their account as verified if needed.
//
// The function performs the following operations:
// 1. Extracts and validates the user ID from the request URL
// 2. Retrieves the user record from the database
// 3. Checks if the user is already verified
// 4. If not verified, updates the user's verification status in a transaction
// 5. Renders appropriate HTML response based on verification results
//
// Parameters:
//   - wr: Custom HTTP writer for response handling
//   - repo: Database repository for user operations
//   - conn: Database connection for transaction management
func VerifyUser(wr *utils.HttpWriter, repo *repository.Queries, conn *pgx.Conn) {
	logger := pkg.NewLogger()
	defer logger.Close()

	// Extract the user ID from the URL query parameter
	// This ID is included in the verification email link sent to the user
	user_id := wr.R.URL.Query().Get("uid")

	// Return an error if no user ID is provided in the request
	// This prevents unauthorized verification attempts without a valid ID
	if user_id == "" {
		wr.Status(http.StatusNotAcceptable).Json(
			utils.M{
				"message": "No user id provided that means this is not an legit request, ABORTING!",
			},
		)
		return
	}

	// Convert string ID to UUID format and validate its structure
	// This ensures the ID follows the proper UUID format before database lookup
	uuid, err := utils.StringToUUID(user_id)
	if err != nil {
		wr.Status(http.StatusBadRequest).Json(
			utils.M{
				"message": "Invalid user ID format",
			},
		)
		return
	}

	// Retrieve the user record from the database using the validated UUID
	// This confirms the user exists and retrieves their current verification status
	user, err := repo.GetUserByID(wr.R.Context(), uuid)
	if err != nil {
		if err == pgx.ErrNoRows {
			wr.Status(http.StatusNotFound).Json(
				utils.M{
					"message": "User not found",
				},
			)
			return
		}

		wr.Status(http.StatusInternalServerError).Json(
			utils.M{
				"message": "Internal server error",
			},
		)
		return
	}

	// Check if the user has already been verified previously
	// If so, show a message indicating the account is already verified
	if user.Verified {
		wr.Status(http.StatusOK).Html(
			"templates/mail/user_already_verified.html",
			utils.M{
				"UserId": user_id,
			},
		)
		return
	}

	// ==========================================================
	// Update user's verification status if not already verified
	// ==========================================================

	// Begin a database transaction to ensure atomicity
	// This prevents partial updates and ensures data consistency
	tx, err := conn.Begin(wr.R.Context())
	if err != nil {
		wr.Status(http.StatusInternalServerError).Json(
			utils.M{
				"message": "Looks like it's us not you :(",
			},
		)
		return
	}

	// Update the user's verified status to true in the database
	// This marks the account as verified and enables full account access
	_, err = repo.WithTx(tx).VerifyUser(wr.R.Context(), uuid)
	if err != nil {
		// Rollback the transaction if verification fails
		// This ensures no partial database changes are committed
		tx.Rollback(wr.R.Context())
		wr.Status(http.StatusInternalServerError).Json(
			utils.M{
				"message": "Failed to verify user",
			},
		)
		return
	}

	// Commit the transaction to finalize the verification
	// This makes the verification status change permanent
	if err = tx.Commit(wr.R.Context()); err != nil {
		wr.Status(http.StatusInternalServerError).Json(
			utils.M{
				"message": "Failed to verify user",
			},
		)
		return
	}

	// Log successful verification for audit and monitoring
	// This provides a record of when each account was verified
	logger.StdoutLogger.Info("User verification successful", "user_id", user_id)
	logger.FileLogger.Info("User verification successful", "user_id", user_id)

	// Render the verification success HTML template
	// This provides visual confirmation to the user that verification succeeded
	wr.Status(http.StatusAccepted).Html(
		"templates/mail/verification_success.html",
		utils.M{
			"UserId": user_id,
		},
	)
}
