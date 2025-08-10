package admin

import (
	"net/http"

	"github.com/dtg-lucifer/everato/internal/db/repository"
	"github.com/dtg-lucifer/everato/internal/utils"
	"github.com/dtg-lucifer/everato/pkg"
	"github.com/jackc/pgx/v5"
)

// GetAllRoles retrieves all available admin roles from the database.
// This function queries the database for all enum values of the SUPER_USER_ROLE type.
//
// Parameters:
//   - wr: HTTP response writer wrapper
//   - repo: Database repository for querying role information
//   - conn: Database connection for transaction handling
//   - cfg: Application configuration
func GetAllRoles(wr *utils.HttpWriter, repo *repository.Queries, conn *pgx.Conn) {
	logger := pkg.NewLogger()
	defer logger.Close()

	// Start a transaction for consistency
	tx, err := conn.Begin(wr.R.Context())
	if err != nil {
		logger.StdoutLogger.Error("Error starting transaction for GetAllRoles", "err", err.Error())
		wr.Status(http.StatusInternalServerError).Json(
			utils.M{
				"message": "Internal server error",
				"error":   "Failed to start database transaction",
			},
		)
		return
	}

	// Query the database for all admin roles
	roles, err := repo.WithTx(tx).GetAdminUserRoles(wr.R.Context())
	if err != nil {
		tx.Rollback(wr.R.Context()) // Rollback the transaction
		logger.StdoutLogger.Error("Failed to retrieve admin roles", "err", err.Error())
		wr.Status(http.StatusInternalServerError).Json(
			utils.M{
				"message": "Failed to retrieve admin roles",
				"error":   err.Error(),
			},
		)
		return
	}

	// Commit the transaction
	if err = tx.Commit(wr.R.Context()); err != nil {
		logger.StdoutLogger.Error("Failed to commit transaction", "err", err.Error())
		wr.Status(http.StatusInternalServerError).Json(
			utils.M{
				"message": "Internal server error",
				"error":   "Failed to complete database operation",
			},
		)
		return
	}

	// Transform the role enum values into a structured format
	rolesList := make([]utils.M, len(roles))
	for i, role := range roles {
		// Create a description for each role
		var description string
		switch role {
		case "SUPER_ADMIN":
			description = "Full system access with all permissions"
		case "ADMIN":
			description = "Administrative access to manage events and users"
		case "EDITOR":
			description = "Limited access to edit content only"
		default:
			description = "Role with custom permissions"
		}

		rolesList[i] = utils.M{
			"value":       role,
			"description": description,
		}
	}

	// Return the roles as a JSON response
	wr.Status(http.StatusOK).Json(
		utils.M{
			"message": "Admin roles retrieved successfully",
			"count":   len(roles),
			"roles":   rolesList,
		},
	)
}
