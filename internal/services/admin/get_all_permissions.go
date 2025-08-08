package admin

import (
	"net/http"

	"github.com/dtg-lucifer/everato/internal/db/repository"
	"github.com/dtg-lucifer/everato/internal/utils"
	"github.com/dtg-lucifer/everato/pkg"
	"github.com/jackc/pgx/v5"
)

// GetAllPermissions retrieves all available admin permissions from the database.
// This function queries the database for all enum values of the PERMISSIONS type.
//
// Parameters:
//   - wr: HTTP response writer wrapper
//   - repo: Database repository for querying permission information
//   - conn: Database connection for transaction handling
func GetAllPermissions(wr *utils.HttpWriter, repo *repository.Queries, conn *pgx.Conn) {
	logger := pkg.NewLogger()
	defer logger.Close()

	// Start a transaction for consistency
	tx, err := conn.Begin(wr.R.Context())
	if err != nil {
		logger.StdoutLogger.Error("Error starting transaction for GetAllPermissions", "err", err.Error())
		wr.Status(http.StatusInternalServerError).Json(
			utils.M{
				"message": "Internal server error",
				"error":   "Failed to start database transaction",
			},
		)
		return
	}

	// Query the database for all admin permissions
	permissions, err := repo.WithTx(tx).GetAdminPermissions(wr.R.Context())
	if err != nil {
		tx.Rollback(wr.R.Context()) // Rollback the transaction
		logger.StdoutLogger.Error("Failed to retrieve permissions", "err", err.Error())
		wr.Status(http.StatusInternalServerError).Json(
			utils.M{
				"message": "Failed to retrieve permissions",
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

	// Transform the permission enum values into a structured format with descriptions
	permissionsList := make([]utils.M, len(permissions))
	for i, perm := range permissions {
		// Create a description for each permission
		var description string
		switch perm {
		case "MANAGE_EVENTS":
			description = "Full control over all events"
		case "CREATE_EVENT":
			description = "Ability to create new events"
		case "EDIT_EVENT":
			description = "Ability to modify existing events"
		case "DELETE_EVENT":
			description = "Ability to remove events from the system"
		case "VIEW_EVENT":
			description = "Ability to view event details"
		case "MANAGE_BOOKINGS":
			description = "Full control over all bookings"
		case "CREATE_BOOKING":
			description = "Ability to create new bookings"
		case "EDIT_BOOKING":
			description = "Ability to modify existing bookings"
		case "DELETE_BOOKING":
			description = "Ability to cancel bookings"
		case "VIEW_BOOKING":
			description = "Ability to view booking details"
		case "MANAGE_USERS":
			description = "Ability to manage user accounts and permissions"
		case "VIEW_REPORTS":
			description = "Ability to view system reports and analytics"
		default:
			description = "Custom permission"
		}

		permissionsList[i] = utils.M{
			"value":       perm,
			"description": description,
		}
	}

	// Return the permissions as a JSON response
	wr.Status(http.StatusOK).Json(
		utils.M{
			"message":     "Admin permissions retrieved successfully",
			"count":       len(permissions),
			"permissions": permissionsList,
		},
	)
}
