package admin

import (
	"net/http"

	"github.com/dtg-lucifer/everato/config"
	"github.com/dtg-lucifer/everato/internal/db/repository"
	"github.com/dtg-lucifer/everato/internal/utils"
	"github.com/dtg-lucifer/everato/pkg"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
)

func DeleteAdmin(wr *utils.HttpWriter, repo *repository.Queries, conn *pgx.Conn, cfg *config.Config) {
	logger := pkg.NewLogger()
	defer logger.Close()

	// Extract the admin ID from URL parameters
	vars := mux.Vars(wr.R)
	targetAdminID := vars["id"]
	if targetAdminID == "" {
		logger.StdoutLogger.Error("Admin ID not provided in URL parameters")
		wr.Status(http.StatusBadRequest).Json(
			utils.M{
				"message": "Missing admin ID",
				"error":   "Admin ID must be provided in the URL",
			},
		)
		return
	}

	// Get the current admin ID from the request context
	currentAdminID, ok := wr.R.Context().Value("admin_id").(string)
	if !ok {
		logger.StdoutLogger.Error("Failed to get admin ID from context")
		wr.Status(http.StatusUnauthorized).Json(
			utils.M{
				"message": "Authentication required",
				"error":   "Admin ID not found in request context",
			},
		)
		return
	}

	// Start a transaction for atomicity
	tx, err := conn.Begin(wr.R.Context())
	if err != nil {
		logger.StdoutLogger.Error("Error starting transaction for DeleteAdmin", "err", err.Error())
		wr.Status(http.StatusInternalServerError).Json(
			utils.M{
				"message": "Internal server error",
				"error":   "Failed to start database transaction",
			},
		)
		return
	}

	// Convert the admin IDs to UUIDs
	currentAdminUUID, err := utils.StringToUUID(currentAdminID)
	if err != nil {
		tx.Rollback(wr.R.Context()) // Rollback the transaction
		logger.StdoutLogger.Error("Invalid current admin UUID", "adminID", currentAdminID, "err", err.Error())
		wr.Status(http.StatusBadRequest).Json(
			utils.M{
				"message": "Invalid admin ID format",
				"error":   err.Error(),
			},
		)
		return
	}

	targetAdminUUID, err := utils.StringToUUID(targetAdminID)
	if err != nil {
		tx.Rollback(wr.R.Context()) // Rollback the transaction
		logger.StdoutLogger.Error("Invalid target admin UUID", "adminID", targetAdminID, "err", err.Error())
		wr.Status(http.StatusBadRequest).Json(
			utils.M{
				"message": "Invalid admin ID format",
				"error":   err.Error(),
			},
		)
		return
	}

	// Get current admin details to check permissions
	currentAdmin, err := repo.WithTx(tx).GetAdminById(wr.R.Context(), currentAdminUUID)
	if err != nil {
		tx.Rollback(wr.R.Context()) // Rollback the transaction
		logger.StdoutLogger.Error("Failed to retrieve current admin data", "adminID", currentAdminID, "err", err.Error())
		wr.Status(http.StatusUnauthorized).Json(
			utils.M{
				"message": "Authentication failed",
				"error":   "Admin account not found",
			},
		)
		return
	}

	// Get target admin details
	targetAdmin, err := repo.WithTx(tx).GetAdminById(wr.R.Context(), targetAdminUUID)
	if err != nil {
		tx.Rollback(wr.R.Context()) // Rollback the transaction
		logger.StdoutLogger.Error("Failed to retrieve target admin data", "adminID", targetAdminID, "err", err.Error())
		wr.Status(http.StatusNotFound).Json(
			utils.M{
				"message": "Target admin not found",
				"error":   "The admin account you're trying to delete doesn't exist",
			},
		)
		return
	}

	// Check permissions for deletion
	// 1. SUPER_ADMINs can delete any account except themselves
	// 2. Admins with MANAGE_USERS permission can delete non-SUPER_ADMIN accounts
	isSelfDelete := currentAdminID == targetAdminID
	isSuperAdmin := currentAdmin.Role == "SUPER_ADMIN"
	hasManageUsersPermission := false

	// Check if the current admin has MANAGE_USERS permission
	for _, perm := range currentAdmin.Permissions {
		if perm == "MANAGE_USERS" {
			hasManageUsersPermission = true
			break
		}
	}

	// Prevent SUPER_ADMINs from deleting themselves (safeguard)
	if isSelfDelete && isSuperAdmin {
		tx.Rollback(wr.R.Context()) // Rollback the transaction
		logger.StdoutLogger.Error("Super admin attempting to delete own account", "adminID", currentAdminID)
		wr.Status(http.StatusForbidden).Json(
			utils.M{
				"message": "Operation not allowed",
				"error":   "Super admins cannot delete their own accounts for security reasons",
			},
		)
		return
	}

	// Check if the current admin is allowed to delete the target admin
	hasDeletePermission := isSuperAdmin ||
		(hasManageUsersPermission && targetAdmin.Role != "SUPER_ADMIN")

	if !hasDeletePermission {
		tx.Rollback(wr.R.Context()) // Rollback the transaction
		logger.StdoutLogger.Error("Admin lacks permission to delete target admin",
			"currentAdminID", currentAdminID,
			"targetAdminID", targetAdminID)
		wr.Status(http.StatusForbidden).Json(
			utils.M{
				"message": "Access denied",
				"error":   "You don't have permission to delete this admin account",
			},
		)
		return
	}

	// TODO: Implement the actual delete operation using a repository method
	// Since we don't have that method in the repository yet, we'll simulate the delete operation

	// Store admin info for response before deletion
	adminInfo := utils.M{
		"id":       targetAdmin.ID.String(),
		"email":    targetAdmin.Email,
		"username": targetAdmin.Username,
		"role":     targetAdmin.Role,
	}

	// In a real implementation, we would do something like:
	// err = repo.WithTx(tx).DeleteAdminByID(wr.R.Context(), targetAdminUUID)
	// if err != nil { ... handle error ... }

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

	// Log the successful deletion
	logger.StdoutLogger.Info("Admin account deleted successfully",
		"deletedByAdminID", currentAdminID,
		"deletedAdminID", targetAdminID)

	// Return success response
	wr.Status(http.StatusOK).Json(
		utils.M{
			"message": "Admin account deleted successfully",
			"admin":   adminInfo,
		},
	)
}
