package user

import (
	"net/http"

	"github.com/dtg-lucifer/everato/internal/db/repository"
	"github.com/dtg-lucifer/everato/internal/utils"
	"github.com/dtg-lucifer/everato/pkg"
	"github.com/jackc/pgx/v5"
)

// VerifyUser handles the user verification process.
//
// Parameters:
//   - wr - pointer to a custom HttpWriter implementation
//   - repo - pointer to the repository
//   - conn - pointer to the actual postgres connection
func VerifyUser(wr *utils.HttpWriter, repo *repository.Queries, conn *pgx.Conn) {
	logger := pkg.NewLogger()
	defer logger.Close()

	// Get the user id from the URL query
	user_id := wr.R.URL.Query().Get("uid")

	// Return an error in case no id is provided
	if user_id == "" {
		wr.Status(http.StatusNotAcceptable).Json(
			utils.M{
				"message": "No user id provided that means this is not an legit request, ABORTING!",
			},
		)
		return
	}

	// Check if the id is valid
	uuid, err := utils.StringToUUID(user_id)
	if err != nil {
		wr.Status(http.StatusBadRequest).Json(
			utils.M{
				"message": "Invalid user ID format",
			},
		)
		return
	}

	// Retrieve the user with the given ID
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

	// Check if the user is already verified
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
	// If no then update user's verification status
	// ==========================================================

	tx, err := conn.Begin(wr.R.Context()) // Start the transaction for ATOMICITY
	if err != nil {
		wr.Status(http.StatusInternalServerError).Json(
			utils.M{
				"message": "Looks like it's us not you :(",
			},
		)
		return
	}

	_, err = repo.WithTx(tx).VerifyUser(wr.R.Context(), uuid) // Verify the user in the database
	if err != nil {
		tx.Rollback(wr.R.Context()) // Rollback the transaction in case of error
		wr.Status(http.StatusInternalServerError).Json(
			utils.M{
				"message": "Failed to verify user",
			},
		)
		return
	}

	if err = tx.Commit(wr.R.Context()); err != nil { // Commit the transaction
		wr.Status(http.StatusInternalServerError).Json(
			utils.M{
				"message": "Failed to verify user",
			},
		)
		return
	}

	logger.StdoutLogger.Info("Verifying user", "user_id", user_id)
	logger.FileLogger.Info("Verifying user", "user_id", user_id)

	// At the end render the success message
	wr.Status(http.StatusAccepted).Html(
		"templates/mail/verification_success.html",
		utils.M{
			"UserId": user_id,
		},
	)
}
