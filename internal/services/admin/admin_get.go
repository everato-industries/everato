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

func GetAllAdmins(wr *utils.HttpWriter, repo *repository.Queries, conn *pgx.Conn, cfg *config.Config) {
	logger := pkg.NewLogger()
	defer logger.Close()

	admins, err := repo.GetAllAdmins(wr.R.Context())
	if err != nil {
		logger.StdoutLogger.Error("Failed to retrieve admin accounts", "err", err.Error())
		wr.Status(http.StatusInternalServerError).Json(
			utils.M{
				"message": "Internal server error",
				"error":   "Failed to retrieve admin accounts",
			},
		)
		return
	}

	// Return the admin data
	wr.Status(http.StatusOK).Json(
		utils.M{
			"message": "Admin accounts retrieved successfully",
			"admins":  admins,
		},
	)
}

func GetAdminByID(wr *utils.HttpWriter, repo *repository.Queries, conn *pgx.Conn, cfg *config.Config) {
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

	// Start a transaction for consistency
	tx, err := conn.Begin(wr.R.Context())
	if err != nil {
		logger.StdoutLogger.Error("Error starting transaction for GetAdminByID", "err", err.Error())
		wr.Status(http.StatusInternalServerError).Json(
			utils.M{
				"message": "Internal server error",
				"error":   "Failed to start database transaction",
			},
		)
		return
	}

	// Get the current admin ID from the request context to ensure we have proper authorization
	currentAdminID, ok := wr.R.Context().Value("admin_id").(string)
	if !ok {
		tx.Rollback(wr.R.Context()) // Rollback the transaction
		logger.StdoutLogger.Error("Failed to get admin ID from context")
		wr.Status(http.StatusUnauthorized).Json(
			utils.M{
				"message": "Authentication required",
				"error":   "Admin ID not found in request context",
			},
		)
		return
	}

	// Convert the target admin ID to UUID
	targetAdminUUID, err := utils.StringToUUID(targetAdminID)
	if err != nil {
		tx.Rollback(wr.R.Context()) // Rollback the transaction
		logger.StdoutLogger.Error("Invalid target admin UUID", "targetAdminID", targetAdminID, "err", err.Error())
		wr.Status(http.StatusBadRequest).Json(
			utils.M{
				"message": "Invalid admin ID format",
				"error":   err.Error(),
			},
		)
		return
	}

	// Convert the current admin ID to UUID
	currentAdminUUID, err := utils.StringToUUID(currentAdminID)
	if err != nil {
		tx.Rollback(wr.R.Context()) // Rollback the transaction
		logger.StdoutLogger.Error("Invalid current admin UUID", "currentAdminID", currentAdminID, "err", err.Error())
		wr.Status(http.StatusBadRequest).Json(
			utils.M{
				"message": "Invalid admin ID format",
				"error":   err.Error(),
			},
		)
		return
	}

	// Check if the current admin exists and has appropriate permissions
	currentAdmin, err := repo.WithTx(tx).GetAdminById(wr.R.Context(), currentAdminUUID)
	if err != nil {
		tx.Rollback(wr.R.Context()) // Rollback the transaction
		logger.StdoutLogger.Error("Failed to retrieve current admin data", "currentAdminID", currentAdminID, "err", err.Error())
		wr.Status(http.StatusUnauthorized).Json(
			utils.M{
				"message": "Authentication failed",
				"error":   "Admin account not found",
			},
		)
		return
	}

	// Get the target admin from the database
	targetAdmin, err := repo.WithTx(tx).GetAdminById(wr.R.Context(), targetAdminUUID)
	if err != nil {
		tx.Rollback(wr.R.Context()) // Rollback the transaction
		logger.StdoutLogger.Error("Failed to retrieve target admin data", "targetAdminID", targetAdminID, "err", err.Error())
		wr.Status(http.StatusNotFound).Json(
			utils.M{
				"message": "Admin not found",
				"error":   "No admin account exists with the provided ID",
			},
		)
		return
	}

	// Check if the current admin is allowed to view the target admin
	// Either it's the same admin, or the current admin has MANAGE_USERS permission, or is a SUPER_ADMIN
	hasPermission := currentAdminID == targetAdminID // Self-viewing is always allowed

	if !hasPermission {
		// Check for MANAGE_USERS permission or SUPER_ADMIN role
		if currentAdmin.Role == "SUPER_ADMIN" {
			hasPermission = true
		} else {
			for _, permission := range currentAdmin.Permissions {
				if permission == "MANAGE_USERS" {
					hasPermission = true
					break
				}
			}
		}
	}

	if !hasPermission {
		tx.Rollback(wr.R.Context()) // Rollback the transaction
		logger.StdoutLogger.Error("Admin lacks permission to view other admin",
			"currentAdminID", currentAdminID,
			"targetAdminID", targetAdminID)
		wr.Status(http.StatusForbidden).Json(
			utils.M{
				"message": "Access denied",
				"error":   "Insufficient permissions to view this admin account",
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

	// Return the admin data (excluding the password for security)
	wr.Status(http.StatusOK).Json(
		utils.M{
			"message": "Admin account retrieved successfully",
			"admin": utils.M{
				"id":          targetAdmin.ID.String(),
				"email":       targetAdmin.Email,
				"username":    targetAdmin.Username,
				"name":        targetAdmin.Name,
				"role":        targetAdmin.Role,
				"permissions": targetAdmin.Permissions,
				"created_at":  targetAdmin.CreatedAt,
				"updated_at":  targetAdmin.UpdatedAt,
			},
		},
	)
}

func GetAdminByUserName(wr *utils.HttpWriter, repo *repository.Queries, conn *pgx.Conn, cfg *config.Config) {
	logger := pkg.NewLogger()
	defer logger.Close()

	// Extract the username from URL parameters
	vars := mux.Vars(wr.R)
	targetUsername := vars["username"]
	if targetUsername == "" {
		logger.StdoutLogger.Error("Username not provided in URL parameters")
		wr.Status(http.StatusBadRequest).Json(
			utils.M{
				"message": "Missing username",
				"error":   "Username must be provided in the URL",
			},
		)
		return
	}

	// Start a transaction for consistency
	tx, err := conn.Begin(wr.R.Context())
	if err != nil {
		logger.StdoutLogger.Error("Error starting transaction for GetAdminByUserName", "err", err.Error())
		wr.Status(http.StatusInternalServerError).Json(
			utils.M{
				"message": "Internal server error",
				"error":   "Failed to start database transaction",
			},
		)
		return
	}

	// Get the target admin from the database by username
	targetAdmin, err := repo.WithTx(tx).GetAdminByUsername(wr.R.Context(), targetUsername)
	if err != nil {
		tx.Rollback(wr.R.Context()) // Rollback the transaction
		logger.StdoutLogger.Error("Failed to retrieve admin by username", "username", targetUsername, "err", err.Error())
		wr.Status(http.StatusNotFound).Json(
			utils.M{
				"message": "Admin not found",
				"error":   "No admin account exists with the provided username",
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

	// Return the admin data (excluding the password for security)
	wr.Status(http.StatusOK).Json(
		utils.M{
			"message": "Admin account retrieved successfully",
			"admin": utils.M{
				"id":          targetAdmin.ID.String(),
				"email":       targetAdmin.Email,
				"username":    targetAdmin.Username,
				"name":        targetAdmin.Name,
				"role":        targetAdmin.Role,
				"permissions": targetAdmin.Permissions,
				"created_at":  targetAdmin.CreatedAt,
				"updated_at":  targetAdmin.UpdatedAt,
			},
		},
	)
}

func SearchAdminByQuery(wr *utils.HttpWriter, repo *repository.Queries, conn *pgx.Conn, cfg *config.Config) {
	logger := pkg.NewLogger()
	defer logger.Close()

	// Extract the search query from URL parameters
	vars := mux.Vars(wr.R)
	searchQuery := vars["query"]
	if searchQuery == "" {
		logger.StdoutLogger.Error("Search query not provided in URL parameters")
		wr.Status(http.StatusBadRequest).Json(
			utils.M{
				"message": "Missing search query",
				"error":   "Search query must be provided in the URL",
			},
		)
		return
	}

	// Start a transaction for consistency
	tx, err := conn.Begin(wr.R.Context())
	if err != nil {
		logger.StdoutLogger.Error("Error starting transaction for SearchAdminByQuery", "err", err.Error())
		wr.Status(http.StatusInternalServerError).Json(
			utils.M{
				"message": "Internal server error",
				"error":   "Failed to start database transaction",
			},
		)
		return
	}

	// Get the current admin ID from the request context to ensure we have proper authorization
	currentAdminID, ok := wr.R.Context().Value("admin_id").(string)
	if !ok {
		tx.Rollback(wr.R.Context()) // Rollback the transaction
		logger.StdoutLogger.Error("Failed to get admin ID from context")
		wr.Status(http.StatusUnauthorized).Json(
			utils.M{
				"message": "Authentication required",
				"error":   "Admin ID not found in request context",
			},
		)
		return
	}

	// Convert the admin ID to UUID
	adminUUID, err := utils.StringToUUID(currentAdminID)
	if err != nil {
		tx.Rollback(wr.R.Context()) // Rollback the transaction
		logger.StdoutLogger.Error("Invalid admin UUID", "adminID", currentAdminID, "err", err.Error())
		wr.Status(http.StatusBadRequest).Json(
			utils.M{
				"message": "Invalid admin ID format",
				"error":   err.Error(),
			},
		)
		return
	}

	// Check if the admin exists and has appropriate permissions
	currentAdmin, err := repo.WithTx(tx).GetAdminById(wr.R.Context(), adminUUID)
	if err != nil {
		tx.Rollback(wr.R.Context()) // Rollback the transaction
		logger.StdoutLogger.Error("Failed to retrieve admin data", "adminID", currentAdminID, "err", err.Error())
		wr.Status(http.StatusUnauthorized).Json(
			utils.M{
				"message": "Authentication failed",
				"error":   "Admin account not found",
			},
		)
		return
	}

	// Check if the admin has MANAGE_USERS permission or is a SUPER_ADMIN
	hasPermission := false
	if currentAdmin.Role == "SUPER_ADMIN" {
		hasPermission = true
	} else {
		for _, permission := range currentAdmin.Permissions {
			if permission == "MANAGE_USERS" {
				hasPermission = true
				break
			}
		}
	}

	if !hasPermission {
		tx.Rollback(wr.R.Context()) // Rollback the transaction
		logger.StdoutLogger.Error("Admin lacks permission", "adminID", currentAdminID, "role", currentAdmin.Role)
		wr.Status(http.StatusForbidden).Json(
			utils.M{
				"message": "Access denied",
				"error":   "Insufficient permissions to search admin accounts",
			},
		)
		return
	}

	// TODO: Implement a query to search admins by email or username
	// This would be something like:
	// SELECT * FROM super_users WHERE email ILIKE '%' || $1 || '%' OR username ILIKE '%' || $1 || '%'

	// For now, we'll just use the existing queries to try email and username separately
	var foundAdmins []utils.M

	// Try exact match by email first
	if adminByEmail, err := repo.WithTx(tx).GetAdminByEmail(wr.R.Context(), searchQuery); err == nil {
		foundAdmins = append(foundAdmins, utils.M{
			"id":          adminByEmail.ID.String(),
			"email":       adminByEmail.Email,
			"username":    adminByEmail.Username,
			"name":        adminByEmail.Name,
			"role":        adminByEmail.Role,
			"permissions": adminByEmail.Permissions,
			"created_at":  adminByEmail.CreatedAt,
			"updated_at":  adminByEmail.UpdatedAt,
		})
	}

	// Try exact match by username
	if adminByUsername, err := repo.WithTx(tx).GetAdminByUsername(wr.R.Context(), searchQuery); err == nil {
		// Check if we already added this admin (to avoid duplicates)
		alreadyAdded := false
		for _, admin := range foundAdmins {
			if admin["id"] == adminByUsername.ID.String() {
				alreadyAdded = true
				break
			}
		}

		if !alreadyAdded {
			foundAdmins = append(foundAdmins, utils.M{
				"id":          adminByUsername.ID.String(),
				"email":       adminByUsername.Email,
				"username":    adminByUsername.Username,
				"name":        adminByUsername.Name,
				"role":        adminByUsername.Role,
				"permissions": adminByUsername.Permissions,
				"created_at":  adminByUsername.CreatedAt,
				"updated_at":  adminByUsername.UpdatedAt,
			})
		}
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

	// Return the search results
	wr.Status(http.StatusOK).Json(
		utils.M{
			"message": "Search completed successfully",
			"query":   searchQuery,
			"results": foundAdmins,
			"count":   len(foundAdmins),
		},
	)
}
